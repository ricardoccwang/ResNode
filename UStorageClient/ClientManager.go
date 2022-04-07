/**
 * @Author: loyeller
 * @Description:
 * @File:  ClientManage
 * @Version: 1.0.0
 * @Date: 2021/11/10 10:29
 */
package UStorageClient

import (
	"UResNode/UTool"
	"UResNode/internal/Data"
	"UResNode/internal/config"
	"UResNode/internal/types"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ResUnit struct {
	Name          string // 文件名
	Filepath      string // 文件路径
	UrlPath       string // 文件检索路径，就是相对路径，为了提供给URL用的路径
	FileSignature string // 文件的标识值
	FileSize      int64  // 文件大小
	Op            string // 操作  add/delete/update 三种类型
}

type ClientManager struct {
	ResUnits        map[string]*ResUnit // map[filepath]*ResUnit
	Root            string
	SendAry         []*ResUnit
	mu              sync.Mutex
	serverUrl       string
	DownloadUrl     string
	node            config.NodeConf
	encrypt         string
	Token           *ClientToken
	IsConnecting    bool
	sickCount       int64
	watch           *fsnotify.Watcher
	pendingFilePath []string
	lock            sync.RWMutex
}

func NewClientManager(r string, serverUrl string, downloadUrl string, node config.NodeConf, encrypt string) *ClientManager {
	s := ClientManager{
		Root:         r,
		serverUrl:    serverUrl,
		DownloadUrl:  downloadUrl,
		node:         node,
		encrypt:      encrypt,
		IsConnecting: false,
	}
	return &s
}

func (c *ClientManager) notifyNodeToServer() error {
	if c.IsConnecting == false {
		go c._notifyNodeToServer()
	}
	return nil
}

// 为了token 连接服务器
func (c *ClientManager) _notifyNodeToServer() {
	c.IsConnecting = true
	for {
		if c.IsConnecting == false {
			break
		}
		token, err := SendNodeInfoToServer(c.serverUrl, c.node, c.encrypt)
		if err != nil {
			logx.Error(err.Error())
		} else {
			if token != nil {
				_ = c.ResendAll()
				if c.Token != nil {
					_ = c.Token.Close()
				}
				c.Token = NewClientToken(token, fmt.Sprintf("%s%s", c.serverUrl, Data.HTTP_URL_REFRESH_TOKEN))
				break
			}
		}
		c.Cure()                    // 这里使用cure是为了避免多次触发sick操作， 否则会触发Isconnect=false的代码，导致逻辑异常
		time.Sleep(1 * time.Minute) // 如果无法访问，则休眠一秒
	}
	c.IsConnecting = false
}

//monitor
func (c *ClientManager) monitor() {
	go c.monitorFileChange()
	go c.healthMonitor()
}

func (c *ClientManager) DelayTransfer() {
	ticker := time.NewTicker(time.Second * 3)
	go func() {
		for v := range ticker.C {
			_ = v
			c.notifyResToServer()
		}
	}()
}

// StartClient 开始client的任务，根据配置文件，见指定目录下的所有文件视为资源，隐藏文件也不例外
func (c *ClientManager) StartClient() error {
	err := c.recreateFileIndex()
	if err != nil {
		return err
	}
	err = c.ConnectToServer()
	if err != nil {
		return err
	}
	c.monitor()
	c.DelayTransfer()
	return nil
}

func (c *ClientManager) healthMonitor() {
	for {
		if c.sickFound() {
			c.Cure() // 避免多次触发sick操作
			_ = c.ReConnectToServer()
		}
		c.sick()
		time.Sleep(1 * time.Minute)
	}
}

func (c *ClientManager) sick() {
	c.sickCount += 1
}

func (c *ClientManager) sickFound() bool {
	return c.sickCount > Data.SickCount
}

func (c *ClientManager) Cure() {
	c.sickCount = 0
}

func (c *ClientManager) ConnectToServer() error {
	return c.notifyNodeToServer()
}

func (c *ClientManager) Close() {
	if c.Token != nil {
		c.Token.Close()
		c.Token = nil
	}
	c.IsConnecting = false
}

func (c *ClientManager) ReConnectToServer() error {
	c.Close()
	return c.ConnectToServer()
}

// recreateFileIndex 删除之前的索引，全新获取一边索引
func (c *ClientManager) recreateFileIndex() error {
	err := c.deleteAllFileIndex()
	if err != nil {
		logx.Info(err.Error())
		return err
	}

	files, err := c.scanDirectory()
	if err != nil {
		logx.Info(err)
		return err
	}
	for _, filePath := range files {
		_, err = c.AddNewFileIndex(filePath)
		if err != nil {
			logx.Info(err.Error())
		}
	}
	return nil
}

// deleteAllFileIndex 之前的索引
func (c *ClientManager) deleteAllFileIndex() error {
	c.ResUnits = make(map[string]*ResUnit)
	return nil
}

func (c *ClientManager) scanDirectory() ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(c.Root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

func (c *ClientManager) createOrGetWatch() (*fsnotify.Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		logx.Info(err.Error())
		return nil, err
	}
	if w == nil {
		return nil, errors.New("watch is nil")
	}
	dirs, err := UTool.GetAllDirs(c.Root)
	if err != nil {
		return nil, err
	}
	dirs = append(dirs, c.Root) // 将根目录也加入监视中

	for _, dir := range dirs {
		err = w.Add(dir)
		if err != nil {
			_ = w.Close()
			return nil, err
		}
	}

	c.watch = w
	return w, nil
}

//monitorFileChange 发现新的文件
func (c *ClientManager) monitorFileChange() {
	w, err := c.createOrGetWatch()
	if err != nil {
		logx.Info(err.Error())
		return
	}
	for {
		select {
		case ev := <-w.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					time.Sleep(50 * time.Millisecond)
					if ok, _ := UTool.IsFile(ev.Name); ok == false {
						c._addNewDir(ev.Name)
					} else {
						c.AddNewFileIndex(ev.Name)
					}

				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					c.deleteFileIndex(ev.Name)
				}

				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					c.RenameFileIndex(ev.Name)
				}

				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {

				}
			}
		case err := <-w.Errors:
			{
				logx.Error(err.Error())
				return
			}
		}
	}
	logx.Info("monitor file change over")
}

func (c *ClientManager) notifyResToServer() {
	if c.Token == nil {
		return
	}

	if len(c.SendAry) == 0 {
		return
	}
	c.mu.Lock()
	// 每次上传上线是1000个， json的body有长度限制，size也有限制
	sendCount := math.Min(float64(len(c.SendAry)), Data.MaxSendFileInfoCount)
	sData := make([]*ResUnit, int64(sendCount))
	copy(sData, c.SendAry[:int64(sendCount)])
	c.SendAry = c.SendAry[int64(sendCount):]
	c.mu.Unlock()

	err := SendFileInfoToServer(c.serverUrl, c.DownloadUrl, sData, c.Token)
	if err != nil {
		UTool.LogxLogFail(err.Error())
	}
}

//AddNewFileIndex 添加新的文件到索引中
func (c *ClientManager) AddNewFileIndex(fPath string) (*ResUnit, error) {
	unit, err := c._addNewFileIndex(fPath)
	if err != nil {
		c.pendingFilePath = append(c.pendingFilePath, fPath)
		logx.Error(UTool.LogxFormatFail(fmt.Sprintf("Add New Path(%s), but fail", fPath)))
		return nil, err
	}
	return unit, nil
}

func (c *ClientManager) _addNewDir(dirPath string) error {
	logx.Info(fmt.Sprintf("add new dir, %s", dirPath))

	if c.watch != nil {
		// 新增文件前要对文件进行删除（主要是对inode的复用进行处理）
		// 因为重命名时，事件是文件已经修改后再传过来的，导致删除事件无法真正删除watch里面的监控对象
		_ = c.watch.Remove(dirPath)
		err := c.watch.Add(dirPath)
		if err != nil {
			logx.Error(UTool.LogxFormatFail(fmt.Sprintf("error when add dir %s", dirPath)))
		}
	}
	return nil
}

func (c *ClientManager) _addNewFileIndex(fPath string) (*ResUnit, error) {
	logx.Info(fmt.Sprintf("add new file, %s", fPath))
	fPath, _ = filepath.Abs(fPath)
	data, err := ioutil.ReadFile(fPath)

	if err != nil {
		return nil, err
	}
	fileSignature := fmt.Sprintf("%x", sha256.Sum256(data))
	absFPath, err := filepath.Abs(fPath)
	if err != nil {
		return nil, err
	}

	relPath, _ := filepath.Rel(c.Root, fPath)

	fi, err := os.Stat(absFPath)
	if err != nil {
		return nil, err
	}

	ru := &ResUnit{
		Name:          filepath.Base(fPath),
		Filepath:      absFPath,
		UrlPath:       relPath,
		FileSignature: fileSignature,
		FileSize:      fi.Size(),
		Op:            Data.SyncNodeOperatorAdd,
	}
	c.lock.Lock()
	c.ResUnits[fPath] = ru
	c.lock.Unlock()

	c.mu.Lock()
	c.SendAry = append(c.SendAry, ru)
	defer c.mu.Unlock()
	return ru, err
}

// deleteFileIndex 删除文件/目录的事件，watch不管被删除的，只管新增的。
func (c *ClientManager) deleteFileIndex(fPath string) error {
	logx.Info(fmt.Sprintf("delete path = %s", fPath))
	c.lock.Lock()
	defer c.lock.Unlock()
	if res, ok := c.ResUnits[fPath]; ok {
		ru := ResUnit{
			Name:          res.Name,
			Filepath:      res.Filepath,
			UrlPath:       res.UrlPath,
			FileSignature: res.FileSignature,
			FileSize:      res.FileSize,
			Op:            Data.SyncNodeOperatorDelete,
		}
		delete(c.ResUnits, fPath)
		c.mu.Lock()
		c.SendAry = append(c.SendAry, &ru)
		c.mu.Unlock()
	}

	return nil
}

func (c *ClientManager) RenameFileIndex(fPath string) error {
	return c.deleteFileIndex(fPath)
}

func (c *ClientManager) GetResUnitBySignature(sig string) (bool, *ResUnit, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for _, value := range c.ResUnits {
		if value.FileSignature == sig {
			return true, value, nil
		}
	}
	return false, nil, nil
}

func (c *ClientManager) GetFileByFilePath(fpath string) (*ResUnit, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if res, ok := c.ResUnits[fpath]; ok {
		return res, nil
	} else {
		return nil, errors.New(fmt.Sprintf("can't find resource by fpath = %s", fpath))
	}
}

func (c *ClientManager) SaveFileFromRequest(req types.UploadResReq, r *http.Request) (string, error) {
	f, fh, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	data := make([]byte, fh.Size)
	_, err = f.Read(data)
	if err != nil {
		return "", err
	}

	saveFilePath := fmt.Sprintf("%s/%s", c.Root, req.FileName)
	err = ioutil.WriteFile(saveFilePath, data, 0755)
	if err != nil {
		return "", err
	}
	return saveFilePath, nil
}

// ResendAll 将所有的res的信息都上传
func (c *ClientManager) ResendAll() error {
	c.mu.Lock()
	c.SendAry = make([]*ResUnit, len(c.ResUnits))
	count := 0
	c.lock.RLock()
	for _, v := range c.ResUnits {
		c.SendAry[count] = v
		count += 1
	}
	c.lock.RUnlock()
	defer c.mu.Unlock()
	return nil
}
