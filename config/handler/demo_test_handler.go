package handler

import (
	"fmt"
)

// DemoTestHandler 示例处理器
type DemoTestHandler struct {
}

// AppendStr 示例方法
func (h *DemoTestHandler) AppendStr(args ...interface{}) interface{} {
	return "Hello, World!" + fmt.Sprintf("%v", args[0])
}

func (h *DemoTestHandler) AddMoney(args ...interface{}) interface{} {
	fmt.Printf("args: %v\n", args)
	return 0
}
