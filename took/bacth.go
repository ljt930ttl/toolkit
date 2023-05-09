package main

import (
	"fmt"
	"os"
	"strings"
)

type lockDevice struct {
	number   string
	rfid     string
	lockType string
	name     string
}

func readfile() string {
	file, err := os.Open("device.txt")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Println("bytes read: ", bytesread)
	// fmt.Println("bytestream to string: ", string(buffer))
	return string(buffer)
}

func splitLine(lines string) {

	deviceList := make([]*lockDevice, 0)
	arrLines := strings.Split(lines, "\n")
	for i, line := range arrLines {

		if i == 0 {
			continue
		}
		device := new(lockDevice)
		block := strings.Split(line, ":")
		if len(block) == 8 {
			device.number = block[0]
			device.lockType = block[1]
			device.rfid = block[4]
			device.name = block[6]
			deviceList = append(deviceList, device)
		} else {
			fmt.Sprintln("line err", line)
		}
	}
	fmt.Print("end")
}
func main() {
	line := readfile()
	splitLine(line)

}
