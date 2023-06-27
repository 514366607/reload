package reload

import "fmt"

// Logger 日志
type Logger interface {
	Debug(...interface{})
	Error(...interface{})
}

type defaultLogger struct{}

func (*defaultLogger) Debug(m ...interface{}) {
	fmt.Println(m...)
}

func (*defaultLogger) Error(m ...interface{}) {
	fmt.Println(m...)
}
