package xops

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/prometheus/prometheus/util/xfile"
)

type Manager struct {
	defaultPid string // prometheus.pid path
}

func NewManager() *Manager {
	pwd, _ := os.Getwd() //  获取prometheus可执行文件的运行目录
	return &Manager{
		defaultPid: pwd + "/prometheus.pid", // prometheus.pid 目录就是prometheus可执行文件同级目录
	}
}

// reload
func (m *Manager) Reload() error {
	return m.signal(syscall.SIGHUP)
}

// reload
func (m *Manager) Stop() error {
	return m.signal(syscall.SIGINT)
}

// 记录 Prometheus pid
func (m *Manager) CreateAndRecordPid() error {
	pid := os.Getpid() // 获取Prometheus pid
	if xfile.FileExist(m.defaultPid) {
		return errors.New("prometheus is running...")
	}
	return xfile.FileCreateAndWrite(
		m.defaultPid,
		func(fd *os.File) error {
			writer := bufio.NewWriter(fd)
			// 写入pid
			if _, err := writer.WriteString(fmt.Sprintf("%d", pid)); err != nil {
				return err
			}
			return writer.Flush()
		})
}

// 退出时 移除Prometheus pid
func (m *Manager) RemovePid() error {

	if !xfile.FileExist(m.defaultPid) {
		return errors.New(m.defaultPid + " not found")
	}
	return xfile.FileDelete(m.defaultPid)
}

// 读取Prometheus pid
func (m *Manager) getPID() (int, error) {
	b, err := xfile.FileRead(m.defaultPid)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(b))
}

// 发送信息
func (m *Manager) signal(s syscall.Signal) error {
	pid, err := m.getPID()
	if err != nil {
		return err
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(s)
}
