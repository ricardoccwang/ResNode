/**
 * @Author: loyeller
 * @Description:
 * @File:  common
 * @Version: 1.0.0
 * @Date: 2021/11/10 11:46
 */
package UTool

import (
	"UResNode/internal/Data"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"os"
	"runtime"
)

func IsFile(fPath string) (bool, error) {
	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(fPath); err != nil {
		if os.IsExist(err) == false {
			return false, err
		}

	}
	//上面判断是否存在文件/文件夹
	if fi.IsDir() {
		return false, errors.New("is directory")
	}
	return true, nil
}

// SaltSecret 加盐加密
func SaltSecret(secret string, salt string) string {
	newSecret := sha256.Sum256([]byte(secret))
	strSecret := string(newSecret[:])
	sSecret := sha256.Sum256([]byte(strSecret + salt))
	rst := fmt.Sprintf("%x", sSecret)
	return rst
}

func GetAllDirs(dirPth string) (dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			newDirs, err := GetAllDirs(dirPth + PthSep + fi.Name())
			if err != nil {
				return nil, err
			}
			dirs = append(dirs, newDirs...)
		}
	}
	return dirs, nil
}

func LogxFormatSuccess() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	return fmt.Sprintf(Data.LOGX_SUCCESS_FMT, funcName)
}

func LogxFormatFail(errorReason string) string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf(Data.LOGX_FAIL_FMT, funcName, errorReason)
}

func LogxBothFail(errorReason string) error {
	message := LogxFormatFail(errorReason)
	logx.Error(message)
	return errors.New(message)
}

func LogxLogFail(errorReason string){
	logx.Error(LogxFormatFail(errorReason))
}
