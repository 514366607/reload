package reload

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

// S 服务
type s struct {
	L         net.Listener
	reLoad    uint // 初始化次数
	canReLoad uint // 允许重启次数，0为不限制
	SigHandle sigHandle
	sigChan   chan os.Signal
	sigs      []os.Signal
	isRun     bool
	stopChan  chan struct{}
	sync.WaitGroup
	sync.RWMutex
	logger Logger
}

// Start 开始监听信号量
func (s *s) Start() {
	s.Lock()
	s.isRun = true
	s.Unlock()

	signal.Notify(
		s.sigChan,
		s.sigs...,
	)

	go func() {
		var sig os.Signal
		pid := syscall.Getpid()
		for {
			sig = <-s.sigChan
			s.logger.Debug(pid, "Received SIG.", sig)
			if _, ok := s.SigHandle[sig]; ok {
				s.SigHandle[sig](s)
			}
		}
	}()

	<-s.stopChan
}

// CanReLoad 是否可以重启
func (s *s) CanReLoad() bool {
	s.RLock()
	defer s.RUnlock()
	if s.canReLoad == 0 {
		return true
	} else if s.reLoad >= s.canReLoad {
		return false
	}
	return true
}

// SetCanReLoad 设置可重启次数
func (s *s) SetCanReLoad(count uint) {
	s.Lock()
	defer s.Unlock()
	s.canReLoad = count
}

// IsChild 是否子进程
func (s *s) IsChild() bool {
	s.RLock()
	defer s.RUnlock()
	return s.reLoad > 0
}

// Shutdown 停止
func (s *s) Shutdown() {
	s.Lock()
	defer s.Unlock()

	s.Wait()

	s.stopChan <- struct{}{}
}

// Reload 重启
func (s *s) Reload() (err error) {
	if !s.CanReLoad() {
		return fmt.Errorf("forked count is %d More The %d ", s.reLoad, s.canReLoad)
	}

	s.Lock()
	defer s.Unlock()
	s.reLoad++

	s.logger.Debug("Restart: forked Start....")

	tl := s.L.(*net.TCPListener)
	fl, _ := tl.File()

	path := os.Args[0]
	var args []string
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			tag := strings.Split(arg, "=")
			if tag[0] == "-reLoad" {
				break
			}
			args = append(args, arg)
		}
	}
	args = append(args, fmt.Sprintf("-reLoad=%d", s.reLoad))

	s.logger.Debug(path, args)
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{fl}

	err = cmd.Start()
	if err != nil {
		s.logger.Error("Restart: Failed to launch, error: %v", err)
		return
	}
	return
}

// Logger Logger
func (s *s) Logger() Logger {
	return s.logger
}
