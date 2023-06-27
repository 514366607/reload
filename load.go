package reload

import (
	"flag"
	"net"
	"os"
	"syscall"
)

var (
	reLoad uint
)

func init() {
	flag.UintVar(&reLoad, "reLoad", 0, "重启次数")
}

// Service 服务
type Service interface {
	SetCanReLoad(count uint) // 设置可重启次数
	CanReLoad() bool         // 是否可以重启
	IsChild() bool           // 是否子进程
	Start()                  // 开启监控 注意：必须以这个来堵塞
	Reload() (err error)
	Shutdown()
	Wait()
	Add(delat int)
	Done()
	Logger() Logger
}

// NewService 初始化
func NewService(l net.Listener) Service {
	return NewServiceWith(l, WithDefaultHandle())
}

// NewService 初始化
func NewServiceWith(l net.Listener, opts ...Option) Service {
	if !flag.Parsed() {
		panic(" You Must Run Flag.Parse() at Main Pack ! ")
	}
	optsTmp := evaluateOptions(opts)

	var s = new(s)
	s.L = l
	s.reLoad = reLoad
	s.SigHandle = optsTmp.sigHandle
	s.sigs = make([]os.Signal, 0, len(optsTmp.sigHandle))
	for sig := range optsTmp.sigHandle {
		s.sigs = append(s.sigs, sig)
	}
	s.logger = optsTmp.logger
	s.sigChan = make(chan os.Signal)
	s.stopChan = make(chan struct{})
	return s
}

// GetListener 取得连接
func GetListener(laddr string) (l net.Listener, err error) {
	if reLoad > 0 {
		// 子进程用文字描述符来接收数据
		f := os.NewFile(3, "")
		l, err = net.FileListener(f)
		if err != nil {
			return
		}
		// 如果是子进程就杀掉父进程
		syscall.Kill(syscall.Getppid(), syscall.SIGTSTP) //干掉父进程
	} else {
		l, err = net.Listen("tcp", laddr)
		if err != nil {
			return
		}
	}
	return
}
