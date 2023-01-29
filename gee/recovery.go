package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	/*
		skip first 3 caller:
		第 0 个 Caller 是 Callers 本身，
		第 1 个是上一层 trace，
		第 2 个是再上一层的 defer func。
		是中间件的逻辑，为了日志简洁一点，我们跳过了前 3 个 Caller。
	*/

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")

	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		str.WriteString(fmt.Sprintf("\n\t%s:%d", frame.File, frame.Line))
		if !more {
			break
		}
	}
	return str.String()
}

// Recovery 使用 defer 挂载上错误恢复的函数，
//在这个函数中调用 *recover()*，捕获 panic，并且将堆栈信息打印在日志中，向用户返回 Internal Server Error。
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
