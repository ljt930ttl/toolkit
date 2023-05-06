package main

import (
	"toolkit/internal/logger"
	"toolkit/internal/service"
)

func main() {

	logger.SetConsole(true)
	logger.SetLevel(logger.LEVEL_DEBUG)
	service.RunThriftServer()

}
