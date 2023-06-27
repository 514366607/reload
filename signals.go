package reload

import (
	"os"
)

type sigHandle map[os.Signal]HandleFunc

// HandleFunc 处理函数
type HandleFunc func(s Service)
