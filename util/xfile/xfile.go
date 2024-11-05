package xfile

import (
	"errors"
	"os"
)

// 文件是否存在
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 创建文件并写入内容
func FileCreateAndWrite(path string, writeHandler func(fd *os.File) error) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	return writeHandler(fd)
}

// 删除文件
func FileDelete(path string) error {
	return os.Remove(path)
}

// 读取文件内容
func FileRead(path string) ([]byte, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	content := make([]byte, 1024)
	siz, err := fd.Read(content)
	if err != nil {
		return nil, err
	}
	if siz >= 1024 {
		return nil, errors.New("")
	}
	if content[siz] == '\n' {
		siz -= 1
	}
	return content[:siz], nil

}
