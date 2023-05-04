package impl

import (
	"encoding/json"
	"toolkit/internal/logger"
)

type BaseData struct {
	Method    string `json:"method"`
	Station   string `json:"station"`
	IsAck     string `json:"isAck"`
	direction string
}

type AckData struct {
	ReturnCode string `json:"returnCode"`
	Msg        string `json:"msg"`
}
type MsgACK struct {
	BaseData
	Data AckData `json:"data"`
}

type HeartBeatData struct {
	Account string `json:"account"`
}

type HeartBeat struct {
	BaseData
	HeartBeatData HeartBeatData `json:"heartBeatData"`
	count         int
}

func (p *HeartBeat) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	p.count += 1
	logger.Debug("HeartBeat process..", p)
	return process("MS_HeartBeat_ACK", p.Station)
}

// ------------------end--------------------

// 子站或者主站主动发送设备变位信息
type DeviceStatus struct {
	DdeviceNo   string `json:"deviceNo"`
	Status      string `json:"status"`
	RelatedInfo string `json:"relatedInfo"`
	HookTIme    string `json:"HookTIme"`
}
type ChangedDeviceStatus struct {
	BaseData
	DeviceStatus []DeviceStatus `json:"deviceStatus"`
}

func (p *ChangedDeviceStatus) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	// 标注一下
	p.direction = "SM"
	logger.Debug("ChangedDeviceStatus process..", p)
	return process("MS_SendChangedDeviceStatus_ACK", p.Station)
}

// ------------------end--------------------

// 子站或者主站主动发送设备闭锁信息
type DeviceBSInfo struct {
	DeviceNo string `json:"deviceNo"`
	BsStatus string `json:"bsStatus"`
}
type SendDeviceOperInfo struct {
	BaseData
	DeviceBSInfo []DeviceBSInfo `json:"deviceBSInfo"`
}

func (p *SendDeviceOperInfo) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendDeviceOperInfo process..", p)
	return process("MS_SendDeviceOperInfo_ACK", p.Station)
}

// ------------------end--------------------

// 主站或者子站请求对设备进行单点或者多点解锁/闭锁
type DeviceOperInfo struct {
	DeviceNo    string `json:"deviceNo"`
	OperateType string `json:"operateType"`
	LlockUnlock string `json:"lockUnlock"`
	Result      string `json:"result"`
}
type SendDeviceLockUnLockInfo struct {
	BaseData
	DeviceOperInfo []DeviceOperInfo `json:"deviceOperInfo"`
}

func (p *SendDeviceLockUnLockInfo) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendDeviceLockUnLockInfo process..", p)
	return process("MS_SendDeviceLockUnLockInfo_ACK", p.Station)
}

// ------------------end--------------------

// 主站或者子站请求对设备进行总解锁/总闭锁
type AllBSJSData struct {
	AlllockUnlock string `json:"alllockUnlock"`
	Result        string `json:"result"`
}
type SendAllLockUnLockInfo struct {
	BaseData
	AllBSJSData AllBSJSData `json:"allBSJSData"`
}
type AllLockUnLockACK struct {
	BaseData
	AllBSJSData AllBSJSData `json:"allBSJSData"`
}

func (p *SendAllLockUnLockInfo) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendAllLockUnLockInfo process..", p)

	ack := new(AllLockUnLockACK)
	ack.Method = "MS_SendAllLockUnLockInfo_ACK"
	ack.IsAck = "1"
	ack.AllBSJSData = AllBSJSData{AlllockUnlock: p.AllBSJSData.AlllockUnlock, Result: "1"}
	ack.Station = p.Station
	content_byte, _ := json.Marshal(ack)
	return string(content_byte)
}

// ------------------end--------------------

// 子站上送图形文件到主站
type GraphData struct {
	GraphName string `json:"graphName"`
	SvgData   string `json:"svgData"`
}
type SendGraphFileInfo struct {
	BaseData
	GraphData GraphData `json:"graphData"`
}

func (p *SendGraphFileInfo) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendGraphFileInfo process..", p)
	return process("MS_SendGraphFileInfo_ACK", p.Station)
}

// ------------------end--------------------

// 子站上送设备信息到主
type DeviceInfoData struct {
	DeviceName  string `json:"deviceName"`
	Description string `json:"description"`
	Voltage     string `json:"voltage"`
}
type SendDeviceInfo struct {
	BaseData
	Type           string           `json:"type"`
	DeviceInfoData []DeviceInfoData `json:"deviceInfoData"`
}

func (p *SendDeviceInfo) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendDeviceInfo process..", p)
	return process("MS_SendDeviceInfo_ACK", p.Station)
}

// ------------------end--------------------

// 子站上送地线桩的当前挂接情况(新疆)
type SendEarthHookInfo struct {
	BaseData
	// 在前面设备变位信息处定义了 DeviceStatus
	DeviceStatus []DeviceStatus `json:"deviceStatus"`
}

func (p *SendEarthHookInfo) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendEarthHookInfo process..", p)
	return process("MS_SendEarthHookInfo_ACK", p.Station)
}

// ------------------end--------------------

// 子站上送逻辑公式到主站
type LogicalData struct {
	LogicFormula string `json:"logicFormula"`
}
type SendLogicFormula struct {
	BaseData
	LogicalData []LogicalData `json:"logicalData"`
}

func (p *SendLogicFormula) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendLogicFormula process..", p)
	return process("MS_SendLogicFormula_ACK", p.Station)
}

// ------------------end--------------------

// 子站上送钥匙管理机操作记录到主站（新疆）
type KeyOperateData struct {
	UserName          string `json:"UserName"`
	DeviceDescription string `json:"DeviceDescription"`
	DeviceType        string `json:"DeviceType"`
	LogedTime         string `json:"LogedTime"`
	OperationType     string `json:"OperationType"`
	LogReason         string `json:"LogReason"`
	AUUsrName         string `json:"AUUsrName"`
	AuthorizeCatalog  string `json:"AuthorizeCatalog"`
	HasWFTask         string `json:"HasWFTask"`
	HasJXTask         string `json:"HasJXTask"`
}
type SendKeyOperateRecord struct {
	BaseData
	Data []KeyOperateData `json:"data"`
}

func (p *SendKeyOperateRecord) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendKeyOperateRecord process..", p)
	return process("MS_SendKeyOperateRecord_ACK", p.Station)
}

// ------------------end--------------------

// 子站上五防设备台账到主站（新疆）
type WFSystemInfo struct {
	WFSystemProducer     string `json:"WFSystemProducer"`
	WFSystemModel        string `json:"WFSystemModel"`
	WFSystemRunDate      string `json:"WFSystemRunDate"`
	HasSmartCabinetKey   string `json:"HasSmartCabinetKey"`
	CabinetKeyNetApprove string `json:"CabinetKeyNetApprove"`
	CabinetKeysmsApprove string `json:"CabinetKeysmsApprove"`
}
type MonitorSystemInfo struct {
	MonitorManufacturer string `json:"MonitorManufacturer"`
	MonitorStartYear    string `json:"MonitorStartYear"`
	MonitorModel        string `json:"MonitorModel"`
	MonitorHaveFW       string `json:"MonitorHaveFW"`
	MonitorWFEnable     string `json:"MonitorWFEnable"`
	WFToMonitorProtocol string `json:"WFToMonitorProtocol"`
}
type GYEQInfo struct {
	GYEQCount750kV string `json:"GYEQCount750kV"`
	GYEQCount500kV string `json:"GYEQCount500kV"`
	GYEQCount220kV string `json:"GYEQCount220kV"`
	GYEQCount110kV string `json:"GYEQCount110kV"`
	GYEQCount35kV  string `json:"GYEQCount35kV"`
	GYEQCount10kV  string `json:"GYEQCount10kV"`
}
type WFThreeRate struct {
	NeedLockCount    string `json:"NeedLockCount"`
	RealLockCount    string `json:"RealLockCount"`
	PutIntoLockCount string `json:"PutIntoLockCount"`
	IntactLockCount  string `json:"IntactLockCount"`
	InstallRate      string `json:"installRate"`
	InputRate        string `json:"InputRate"`
	IntactRate       string `json:"IntactRate"`
}
type StationBaseData struct {
	StationName       string            `json:"StationName"`
	IsSmartStation    string            `json:"IsSmartStation"`
	Voltage           string            `json:"Voltage"`
	StartRunTime      string            `json:"StartRunTime"`
	WFSystemType      string            `json:"WFSystemType"`
	WFSystemInfo      WFSystemInfo      `json:"WFSystemInfo"`
	MonitorSystemInfo MonitorSystemInfo `json:"MonitorSystemInfo"`
	GYEQInfo          GYEQInfo          `json:"GYEQInfo"`
	WFThreeRate       WFThreeRate       `json:"WFThreeRate"`
}
type SendStationBaseInfo struct {
	BaseData
	Data StationBaseData `json:"data"`
}

func (p *SendStationBaseInfo) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendStationBaseInfo process..", p)
	return process("MS_SendStationBaseInfo_ACK", p.Station)
}

// ------------------end--------------------

func process(method, station string) string {
	ack := new(MsgACK)
	ack.Method = method
	ack.IsAck = "1"
	ack.Data = AckData{ReturnCode: "1", Msg: ""}
	ack.Station = station
	content_byte, _ := json.Marshal(ack)
	return string(content_byte)

}

type HandleFunction interface {
	Process(content string) string
}

// var processMap = map[string]HandleFunction{
// 	"SM_SendHeartBeat":           new(HeartBeat),
// 	"SM_SendChangedDeviceStatus": &ChangedDeviceStatus{},
// 	"SM_SendDeviceOperInfo":      &SendDeviceOperInfo{},
// 	"SM_SendAllLockUnLockInfo":   &SendAllLockUnLockInfo{},
// }

var processMap = map[string]HandleFunction{
	"SM_SendHeartBeat":           new(HeartBeat),
	"SM_SendChangedDeviceStatus": new(ChangedDeviceStatus),
	"SM_SendDeviceOperInfo":      new(SendDeviceOperInfo),
	"SM_SendAllLockUnLockInfo":   new(SendAllLockUnLockInfo),
	"SM_SendGraphFileInfo":       new(SendGraphFileInfo),
	"SM_SendDeviceInfo":          new(SendDeviceInfo),
	"SM_SendEarthHookInfo":       new(SendEarthHookInfo),
	"SM_SendLogicFormula":        new(SendLogicFormula),
	"SM_SendKeyOperateRecord":    new(SendKeyOperateRecord),
	"SM_SendStationBaseInfo":     new(SendStationBaseInfo),
}
