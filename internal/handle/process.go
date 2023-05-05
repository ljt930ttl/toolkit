package handle

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

// 心跳信息
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
	// p.count += 1
	logger.Debug("HeartBeat process..", p)
	return process("MS_HeartBeat_ACK", p.Station)
}

// ------------------end--------------------

// 主站请求全遥信以及闭锁信息
type DeviceStatus struct {
	DdeviceNo   string `json:"deviceNo"`
	Status      string `json:"status"`
	BsStatus    string `json:"bsStatus"`
	RelatedInfo string `json:"relatedInfo"`
	HookTIme    string `json:"HookTIme"`
}
type AskAllYXAndBS struct {
	BaseData
	DeviceStatus []*DeviceStatus `json:"deviceStatus,omitempty"`
}

// 发送请求
// func (p *AskAllYXAndBS) Request(content string) string {
// 	err := json.Unmarshal([]byte(content), p)
// 	if err != nil {
// 		logger.Error(err)
// 		return ""
// 	}
// 	content_byte, _ := json.Marshal(p)
// 	return string(content_byte)
// }

// 接收子站返回信息
func (p *AskAllYXAndBS) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("Recv AskAllYXAndBS ACK", p)
	return ""
}

// ------------------end--------------------

// 子站或者主站主动发送设备变位信息
type ChangedDeviceStatus struct {
	BaseData
	DeviceStatus []DeviceStatus `json:"deviceStatus"`
}

type ChangedDeviceStatusAck struct {
	MsgACK
}

// 发送请求
// func (p *ChangedDeviceStatus) Request(content string) string {
// 	err := json.Unmarshal([]byte(content), p)
// 	if err != nil {
// 		logger.Error(err)
// 		return ""
// 	}
// 	content_byte, _ := json.Marshal(p)
// 	return string(content_byte)
// }

// 接收子站返回信息
func (p *ChangedDeviceStatusAck) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("Recv ChangedDeviceStatus ACK", p)
	return ""
}

// 处理子站信息
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
type SendDeviceOperInfo struct {
	BaseData
	DeviceBSInfo []DeviceStatus `json:"deviceBSInfo"`
}
type SendDeviceOperInfoAck struct {
	MsgACK
}

// 发送请求
func (p *SendDeviceOperInfo) Request(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	content_byte, _ := json.Marshal(p)
	return string(content_byte)
}

// 接收子站返回信息
func (p *SendDeviceOperInfoAck) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("Recv SendDeviceOperInfo ACK", p)
	return ""
}

// 处理子站信息
func (p *SendDeviceOperInfo) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	// 标注一下
	p.direction = "SM"
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
	GraphData []GraphData `json:"graphData"`
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

// 子站上送设备信息到主站
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
	LogicalData LogicalData `json:"logicalData"`
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

// 子站上送五防设备资产详细信息（新疆）
type WFDeviceData struct {
	WFDeviceName string `json:"WFDeviceName"`
	DeviceType   string `json:"DeviceType"`
	ProductDate  string `json:"ProductDate"`
	StartupDate  string `json:"StartupDate"`
	ValidDate    string `json:"ValidDate"`
	IsGood       string `json:"IsGood"`
	DeviceModel  string `json:"DeviceModel"`
	Remarks      string `json:"Remarks"`
}
type WFlocks struct {
	LockName     string `json:"LockName"`
	LockType     string `json:"LockType"`
	RelateDevice string `json:"RelateDevice"`
	ProductDate  string `json:"ProductDate"`
	StartupDate  string `json:"StartupDate"`
	ValidDate    string `json:"ValidDate"`
	IsGood       string `json:"IsGood"`
	Remarks      string `json:"Remarks"`
}
type WFDeviceInfo struct {
	WFDeviceData []WFDeviceData `json:"WFDeviceData"`
	WFlocks      []WFlocks      `json:"WFlocks"`
}

type SendWFDeviceInfo struct {
	BaseData
	Data WFDeviceInfo `json:"data"`
}

func (p *SendWFDeviceInfo) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendWFDeviceInfo process..", p)
	return process("MS_SendWFDeviceInfo_ACK", p.Station)
}

// ------------------end--------------------

// 子站上送缺陷到主站（新疆）
type DefectData struct {
	ID           string `json:"ID"`
	Station      string `json:"Station"`
	DefectDescr  string `json:"DefectDescr"`
	DefectType   string `json:"DefectType"`
	WriteType    string `json:"WriteType"`
	DefectKind   string `json:"DefectKind"`
	Reason       string `json:"Reason"`
	Measure      string `json:"Measure"`
	HandPerson   string `json:"HandPerson"`
	RecordPerson string `json:"RecordPerson"`
	LogTime      string `json:"LogTime"`
	DefectObject string `json:"DefectObject"`
	ProductModel string `json:"ProductModel"`
	DefectStatus string `json:"DefectStatus"`
	SolvetTime   string `json:"SolvetTime"`
	Remark       string `json:"Remark"`
}
type WFDefectData struct {
	DefectData []DefectData `json:"DefectData"`
}

type SendWFDefect struct {
	BaseData
	Data WFDefectData `json:"data"`
}

func (p *SendWFDefect) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendWFDefect process..", p)
	return process("MS_SendWFDefect_ACK", p.Station)
}

// ------------------end--------------------

// 子站上送操作票信息（新疆）
type OpSheet struct {
	Sequence    string `json:"Sequence"`
	Hint        string `json:"Hint"`
	Sbbh        string `json:"SBBH"`
	From        string `json:"from"`
	To          string `json:"to"`
	RelatedInfo string `json:"relatedInfo"`
	ToKey       string `json:"ToKey"`
	ToPrinter   string `json:"ToPrinter"`
	IsHintItem  string `json:"isHintItem"`
	FinishTime  string `json:"FinishTime"`
}
type HistoryTicketData struct {
	AreaName        string    `json:"AreaName"`
	Station         string    `json:"Station"`
	TaskType        string    `json:"TaskType"`
	TaskName        string    `json:"TaskName"`
	HistoryID       string    `json:"HistoryID"`
	ClassName       string    `json:"ClassName"`
	ToKeyTime       string    `json:"ToKeyTime"`
	FinishTime      string    `json:"FinishTime"`
	Writer          string    `json:"Writer"`
	WriterStartTime string    `json:"WriterStartTime"`
	WriterEndTime   string    `json:"WriterEndTime"`
	Operator        string    `json:"Operator"`
	User4           string    `json:"User4"`
	OpSheet         []OpSheet `json:"OpSheet"`
}

type SendHistoryTicket struct {
	BaseData
	Data HistoryTicketData `json:"data"`
}

func (p *SendHistoryTicket) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendHistoryTicket process..", p)
	return process("MS_SendHistoryTicket_ACK", p.Station)
}

// 上送操作票追忆(黑匣子)信息（新疆）
type OpOpRecords struct {
	Sequence     string `json:"Sequence"`
	Index        string `json:"Index"`
	Hint         string `json:"Hint"`
	OpMode       string `json:"OpMode"`
	OpTime       string `json:"OpTime"`
	IsStandardOp string `json:"IsStandardOp"`
	OpDevice     string `json:"OpDevice"`
}
type RecallInfoData struct {
	AreaName    string        `json:"AreaName"`
	Station     string        `json:"Station"`
	HistoryID   string        `json:"HistoryID"`
	OpOpRecords []OpOpRecords `json:"OpOpRecords"`
}
type SendRecallInfo struct {
	BaseData
	Data RecallInfoData `json:"data"`
}

func (p *SendRecallInfo) Process(content string) string {
	err := json.Unmarshal([]byte(content), p)
	if err != nil {
		logger.Error(err)
		return ""
	}
	logger.Debug("SendRecallInfo process..", p)
	return process("MS_SendRecallInfo_ACK", p.Station)
}

// ------------------end-------------------

func process(method, station string) string {
	ack := new(MsgACK)
	ack.Method = method
	ack.IsAck = "1"
	ack.Data = AckData{ReturnCode: "1", Msg: ""}
	ack.Station = station
	content_byte, _ := json.Marshal(ack)
	return string(content_byte)

}

type ServerHandleFunc interface {
	Process(content string) string
}

var processMap = map[string]ServerHandleFunc{
	"SM_SendHeartBeat": new(HeartBeat),
	// 子站返回
	"SM_AskAllYXAndBS_ACK":           new(AskAllYXAndBS),
	"SM_SendChangedDeviceStatus":     new(ChangedDeviceStatus),
	"SM_SendChangedDeviceStatus_ACK": new(ChangedDeviceStatusAck),
	"SM_SendDeviceOperInfo":          new(SendDeviceOperInfo),
	"SM_SendDeviceOperInfo_ACK":      new(SendDeviceOperInfoAck),
	"SM_SendDeviceLockUnLockInfo":    new(SendDeviceLockUnLockInfo),
	"SM_SendAllLockUnLockInfo":       new(SendAllLockUnLockInfo),
	"SM_SendGraphFileInfo":           new(SendGraphFileInfo),
	"SM_SendDeviceInfo":              new(SendDeviceInfo),
	"SM_SendEarthHookInfo":           new(SendEarthHookInfo),
	"SM_SendLogicFormula":            new(SendLogicFormula),
	"SM_SendKeyOperateRecord":        new(SendKeyOperateRecord),
	"SM_SendStationBaseInfo":         new(SendStationBaseInfo),
	"SM_SendWFDeviceInfo":            new(SendWFDeviceInfo),
	"SM_SendWFDefect":                new(SendWFDefect),
	"SM_SendHistoryTicket":           new(SendHistoryTicket),
	"SM_SendRecallInfo":              new(SendRecallInfo),
}
