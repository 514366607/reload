package reload

import (
	"os"
	"syscall"
)

type options struct {
	logger    Logger
	sigHandle sigHandle
}

// Option 参数
type Option func(*options)

var defaultOptions = &options{
	logger:    &defaultLogger{},
	sigHandle: make(sigHandle),
}

// evaluateOptions 参数处理
func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// WithLogger 日志
func WithLogger(l Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

// WithHandleFunc 信号处理
func WithHandleFunc(s os.Signal, h HandleFunc) Option {
	return func(o *options) {
		o.sigHandle[s] = h
	}
}

// WithDefaultHandle 默认信号处理
func WithDefaultHandle() Option {
	return func(o *options) {
		o.sigHandle[syscall.SIGUSR1] = func(s Service) {
			if err := s.Reload(); err != nil {
				s.Logger().Error(err)
			}
		}
		// 设置会退出的信号量
		o.sigHandle[syscall.SIGINT] = func(s Service) { s.Shutdown() }
		o.sigHandle[syscall.SIGTERM] = func(s Service) { s.Shutdown() }
		o.sigHandle[syscall.SIGTSTP] = func(s Service) { s.Shutdown() }
	}
}
