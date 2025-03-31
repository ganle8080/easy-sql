package handler

// GetHandlerFactory 根据 handlerName 返回对应的处理器工厂
func GetHandlerFactory(name string) (func() (interface{}, error), bool) {
	if name == "DemoTestHandler" {
		return func() (interface{}, error) {
			return &DemoTestHandler{}, nil
		}, true
	}
	return nil, false
}
