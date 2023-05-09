package main

import "toolkit/took/logger"

func main() {
	logger.SetConsole(true)
	logger.SetLevel(logger.LEVEL_DEBUG)
	service.RunThriftServer()

}
