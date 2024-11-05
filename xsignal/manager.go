package xsignal

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/prometheus/prometheus/util/xfile"
)

type Manager struct {
	name       string
	ctx        context.Context
	cancel     context.CancelFunc
	defaultPid string
}

func NewManager(name string) *Manager {
	pwd, _ := os.Getwd()
	ctx, cancel := context.WithCancel(context.TODO())
	return &Manager{
		name:       name,
		defaultPid: pwd + "/prometheus.pid",
		cancel:     cancel,
		ctx:        ctx,
	}
}

func (m *Manager) RecordPid() error {
	pid := os.Getpid()
	if xfile.FileExist(m.defaultPid) {
		return errors.New("prometheus is running...")
	}
	return xfile.FileCreateAndWrite(
		m.defaultPid,
		func(fd *os.File) error {
			writer := bufio.NewWriter(fd)
			_, err := writer.WriteString(fmt.Sprintf("%d", pid))
			if err != nil {
				return err
			}
			writer.Flush()

			return nil
		})
}

func (m *Manager) RemovePid() error {

	if !xfile.FileExist(m.defaultPid) {
		return errors.New(m.defaultPid + " not found")
	}
	return xfile.FileDelete(m.defaultPid)
}

func (m *Manager) getPID() (int, error) {
	b, err := xfile.FileRead(m.defaultPid)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(b))
}

func (m *Manager) Signal(s syscall.Signal) error {
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
