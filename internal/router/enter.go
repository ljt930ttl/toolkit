package router

import pmas "toolkit/internal/router/PMAS"

type RouterGroup struct {
	PMAS pmas.PMAServiceRouter
}

var RouterGroupApp = new(RouterGroup)
