package impl

import "encoding/json"

type S_HeartBeatData struct {
	Account string `json:"account"`
}

type HeartBeat struct {
	Method        string          `json:"method"`
	Station       string          `json:"station"`
	IsAck         string          `json:"isAck"`
	HeartBeatData S_HeartBeatData `json:"heartBeatData"`
}

type HeartBeatAckData struct {
	ReturnCode string `json:"returnCode"`
	Msg        string `json:"msg"`
}

type HeartBeatACK struct {
	Method  string           `json:"method"`
	Station string           `json:"Station"`
	IsAck   string           `json:"isAck"`
	Data    HeartBeatAckData `json:"data"`
}

func (hb *HeartBeat) Process(contnet string) string {
	hb_ack := &HeartBeatACK{
		Method: "MS_SendHeartBeat_ACK",
		IsAck:  "1",
	}
	json.Unmarshal([]byte(contnet), hb)
	hb_ack_data := HeartBeatAckData{
		ReturnCode: "1",
		Msg:        "",
	}

	hb_ack.Data = hb_ack_data
	hb_ack.Station = hb.Station
	contnet_byte, _ := json.Marshal(hb_ack)
	return string(contnet_byte)
}

type HandleFunction interface {
	Process(context string) string
}

var processMap = map[string]HandleFunction{
	"SM_SendHeartBeat": &HeartBeat{},
}
