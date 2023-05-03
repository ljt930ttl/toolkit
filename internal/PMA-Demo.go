package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	// logger.SetConsole(true)
	// logger.SetLevel(logger.LEVEL_DEBUG)
	// service.ServerStart()

	jsonBuf := `
    {
    "company": "itcast",
    "subjects": [
        "Go",
        "C++",
        "Python",
        "Test"
    ],
    "isok": true,
    "price": 666.666
}`

	//创建一个map
	m := make(map[string]interface{}, 4)

	err := json.Unmarshal([]byte(jsonBuf), &m) //第二个参数要地址传递
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	fmt.Printf("m = %+v\n", m)
	fmt.Print(m["isok"].(bool))
}
