package eye

import ()

var (
	//所有路由信息最终添加到这里汇总
	routerMux map[string]*RouterItem
)

type RouterItem struct {
	Controller *IController
	Method     string
}
