package logger

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/***
在控制台打印：直接调用 Debug(***) Info(***) Warn(***) Error(***) Fatal(***)
可以设置打印格式：SetFormat(FORMAT_SHORTFILENAME|FORMAT_DATE|FORMAT_TIME)
	无其他格式，只打印日志内容
	FORMAT_NANO
	长文件名及行数
	FORMAT_LONGFILENAME
	短文件名及行数
	FORMAT_SHORTFILENAME
	精确到日期
	FORMAT_DATE
	精确到秒
	FORMAT_TIME
	精确到微秒
	FORMAT_MICROSECNDS
	—————————————————————————————————————————————————————————————————————
    写日志文件可以获取实例
    全局实例可以直接调用log := logging.GetStaticLogger()
    获取新实例可以调用log := logging.NewLogger()
	1. 按日期分割日志文件
    	log.SetRollingDaily("d://foldTest", "log.txt")
	2. 按文件大小分割日志文件
	log.SetRollingFile("d://foldTest", "log.txt", 300, MB)
	log.SetConsole(false)控制台不打日志,默认值true
    日志级别
***/

const (
	_VER string = "2.0.1"
)

type LEVEL int8
type UNIT int64
type MODE_TIME uint8
type ROLLTYPE int //dailyRolling ,rollingFile
type FORMAT int

const DATEFORMAT_DAY = "20060102"
const DATEFORMAT_HOUR = "2006010215"
const DATEFORMAT_MONTH = "200601"

var static_mu *sync.Mutex = new(sync.Mutex)

var static_lo *logger = NewLogger()

var TIME_DEVIATION time.Duration

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	MODE_HOUR  MODE_TIME = 1
	MODE_DAY   MODE_TIME = 2
	MODE_MONTH MODE_TIME = 3
)

const (
	/*无其他格式，只打印日志内容*/
	FORMAT_NANO FORMAT = 0
	/*长文件名(文件绝对路径)及行数*/
	FORMAT_LONGFILENAME = FORMAT(log.Llongfile)
	/*短文件名及行数*/
	FORMAT_SHORTFILENAME = FORMAT(log.Lshortfile)
	/*日期时间精确到天*/
	FORMAT_DATE = FORMAT(log.Ldate)
	/*时间精确到秒*/
	FORMAT_TIME = FORMAT(log.Ltime)
	/*时间精确到微秒*/
	FORMAT_MICROSECNDS = FORMAT(log.Lmicroseconds)
)

const (
	/*日志级别：ALL 最低级别*/
	LEVEL_ALL LEVEL = iota
	/*日志级别：DEBUG 小于INFO*/
	LEVEL_DEBUG
	/*日志级别：INFO 小于 WARN*/
	LEVEL_INFO
	/*日志级别：WARN 小于 ERROR*/
	LEVEL_WARN
	/*日志级别：ERROR 小于 FATAL*/
	LEVEL_ERROR
	/*日志级别：FATAL 小于 OFF*/
	LEVEL_FATAL
	/*日志级别：off 不打印任何日志*/
	LEVEL_OFF
)

const (
	DAYLY ROLLTYPE = iota
	ROLLFILE
)

var default_format FORMAT = FORMAT_SHORTFILENAME | FORMAT_DATE | FORMAT_TIME
var default_level = LEVEL_ALL

/*设置打印格式*/
func SetFormat(format FORMAT) *logger {
	default_format = format
	return static_lo.SetFormat(format)
}

/*设置控制台日志级别，默认ALL*/
func SetLevel(level LEVEL) *logger {
	default_level = level
	return static_lo.SetLevel(level)
}

func SetConsole(on bool) *logger {
	return static_lo.SetConsole(on)

}

/*获得全局Logger对象*/
func GetStaticLogger() *logger {
	return staticLogger()
}

func SetRollingFile(fileDir, fileName string, maxFileSize int64, unit UNIT) error {
	return SetRollingFileLoop(fileDir, fileName, maxFileSize, unit, 0)
}

func SetRollingDaily(fileDir, fileName string) error {
	return SetRollingByTime(fileDir, fileName, MODE_DAY)
}

func SetRollingFileLoop(fileDir, fileName string, maxFileSize int64, unit UNIT, maxFileNum int) error {
	return static_lo.SetRollingFileLoop(fileDir, fileName, maxFileSize, unit, maxFileNum)
}

func SetRollingByTime(fileDir, fileName string, mode MODE_TIME) error {
	return static_lo.SetRollingByTime(fileDir, fileName, mode)
}

func staticLogger() *logger {
	return static_lo
}

func Debug(v ...interface{}) {
	print(default_format, LEVEL_DEBUG, default_level, 2, v...)
}
func Info(v ...interface{}) {
	print(default_format, LEVEL_INFO, default_level, 2, v...)
}
func Warn(v ...interface{}) {
	print(default_format, LEVEL_WARN, default_level, 2, v...)
}
func Error(v ...interface{}) {
	print(default_format, LEVEL_ERROR, default_level, 2, v...)
}
func Fatal(v ...interface{}) {
	print(default_format, LEVEL_FATAL, default_level, 2, v...)
}

func print(_format FORMAT, level, _default_level LEVEL, calldepth int, v ...interface{}) {
	if level < _default_level {
		return
	}
	staticLogger().println(level, k1(calldepth), v...)
}

func __print(_format FORMAT, level, _default_level LEVEL, calldepth int, v ...interface{}) {
	console(fmt.Sprint(v...), getlevelname(level, default_format), _format, k1(calldepth))
}

func getlevelname(level LEVEL, format FORMAT) (levelname string) {
	if format == FORMAT_NANO {
		return
	}
	switch level {
	case LEVEL_ALL:
		levelname = "[ALL]"
	case LEVEL_DEBUG:
		levelname = "[DEBUG]"
	case LEVEL_INFO:
		levelname = "[INFO]"
	case LEVEL_WARN:
		levelname = "[WARN]"
	case LEVEL_ERROR:
		levelname = "[ERROR]"
	case LEVEL_FATAL:
		levelname = "[FATAL]"
	default:
	}
	return
}

/*————————————————————————————————————————————————————————————————————————————*/
type logger struct {
	level      LEVEL
	format     FORMAT
	rwLock     *sync.RWMutex
	safe       bool
	fileDir    string
	fileName   string
	maxSize    int64
	unit       UNIT
	rolltype   ROLLTYPE
	mode       MODE_TIME
	fileObj    *fileObj
	maxFileNum int
	isConsole  bool
}

func NewLogger() (log *logger) {
	log = &logger{level: LEVEL_DEBUG, rolltype: DAYLY, rwLock: new(sync.RWMutex), format: FORMAT_SHORTFILENAME | FORMAT_DATE | FORMAT_TIME, isConsole: true}
	log.newfileObj()
	return
}

// 控制台日志是否打开
func (l *logger) SetConsole(_isConsole bool) *logger {
	l.isConsole = _isConsole
	return l
}
func (l *logger) Debug(v ...interface{}) {
	l.println(LEVEL_DEBUG, 2, v...)
}
func (l *logger) Info(v ...interface{}) {
	l.println(LEVEL_INFO, 2, v...)
}
func (l *logger) Warn(v ...interface{}) {
	l.println(LEVEL_WARN, 2, v...)
}
func (l *logger) Error(v ...interface{}) {
	l.println(LEVEL_ERROR, 2, v...)
}
func (l *logger) Fatal(v ...interface{}) {
	l.println(LEVEL_FATAL, 2, v...)
}
func (l *logger) SetFormat(format FORMAT) *logger {
	l.format = format
	return l
}
func (l *logger) SetLevel(level LEVEL) *logger {
	l.level = level
	return l
}

/*
按日志文件大小分割日志文件
fileDir 日志文件夹路径
fileName 日志文件名
maxFileSize  日志文件大小最大值
unit    日志文件大小单位
*/
func (l *logger) SetRollingFile(fileDir, fileName string, maxFileSize int64, unit UNIT) error {
	return l.SetRollingFileLoop(fileDir, fileName, maxFileSize, unit, 0)
}

/*
按日志文件大小分割日志文件，指定保留的最大日志文件数
fileDir 日志文件夹路径
fileName 日志文件名
maxFileSize  日志文件大小最大值
unit    	日志文件大小单位
maxFileNum  留的日志文件数
*/
func (l *logger) SetRollingFileLoop(fileDir, fileName string, maxFileSize int64, unit UNIT, maxFileNum int) error {
	if fileDir == "" {
		fileDir, _ = os.Getwd()
	}
	if maxFileNum > 0 {
		maxFileNum--
	}
	l.fileDir, l.fileName, l.maxSize, l.maxFileNum, l.unit = fileDir, fileName, maxFileSize, maxFileNum, unit
	l.rolltype = ROLLFILE
	if l.fileObj != nil {
		l.fileObj.close()
	}
	l.newfileObj()
	err := l.fileObj.openFileHandler()
	return err
}

/*
按日期分割日志文件
fileDir 日志文件夹路径
fileName 日志文件名
*/
func (l *logger) SetRollingDaily(fileDir, fileName string) error {
	return l.SetRollingByTime(fileDir, fileName, MODE_DAY)
}

/*
指定按 小时，天，月 分割日志文件
fileDir 日志文件夹路径
fileName 日志文件名
mode   指定 小时，天，月
*/
func (l *logger) SetRollingByTime(fileDir, fileName string, mode MODE_TIME) error {
	if fileDir == "" {
		fileDir, _ = os.Getwd()
	}
	l.fileDir, l.fileName, l.mode = fileDir, fileName, mode
	l.rolltype = DAYLY
	if l.fileObj != nil {
		l.fileObj.close()
	}
	l.newfileObj()
	err := l.fileObj.openFileHandler()
	return err
}

func (l *logger) newfileObj() {
	l.fileObj = new(fileObj)
	l.fileObj.fileDir, l.fileObj.fileName, l.fileObj.maxSize, l.fileObj.rolltype, l.fileObj.unit, l.fileObj.maxFileNum, l.fileObj.mode = l.fileDir, l.fileName, l.maxSize, l.rolltype, l.unit, l.maxFileNum, l.mode
}

func (l *logger) backUp() (err, openFileErr error) {
	l.rwLock.Lock()
	defer l.rwLock.Unlock()
	if !l.fileObj.isMustBackUp() {
		return
	}
	err = l.fileObj.close()
	if err != nil {
		__print(l.format, LEVEL_ERROR, LEVEL_ERROR, 1, err.Error())
		return
	}
	err = l.fileObj.rename()
	if err != nil {
		__print(l.format, LEVEL_ERROR, LEVEL_ERROR, 1, err.Error())
		return
	}
	openFileErr = l.fileObj.openFileHandler()
	if openFileErr != nil {
		__print(l.format, LEVEL_ERROR, LEVEL_ERROR, 1, openFileErr.Error())
	}
	return
}

func (l *logger) println(_level LEVEL, calldepth int, v ...interface{}) {
	if l.level > _level {
		return
	}
	if l.fileObj.isFileWell {
		var openFileErr error
		if l.fileObj.isMustBackUp() {
			_, openFileErr = l.backUp()
		}
		if openFileErr == nil {
			func() {
				l.rwLock.RLock()
				defer l.rwLock.RUnlock()
				if l.format != FORMAT_NANO {
					s := fmt.Sprint(v...)
					buf := getOutBuffer(s, getlevelname(_level, l.format), l.format, k1(calldepth)+1)
					l.fileObj.write2file(buf.Bytes())
				} else {
					var bs []byte
					l.fileObj.write2file(fmt.Appendln(bs, v...))
				}
			}()
		}
	}
	if l.isConsole {
		__print(l.format, _level, l.level, k1(calldepth), v...)
	}
}

/*————————————————————————————————————————————————————————————————————————————*/
type fileObj struct {
	fileDir     string
	fileName    string
	maxSize     int64
	fileSize    int64
	unit        UNIT
	fileHandler *os.File
	rolltype    ROLLTYPE
	tomorSecond int64
	isFileWell  bool
	maxFileNum  int
	mode        MODE_TIME
}

func (f *fileObj) openFileHandler() (e error) {
	if f.fileDir == "" || f.fileName == "" {
		e = errors.New("log filePath is null or error")
		return
	}
	e = mkdirDir(f.fileDir)
	if e != nil {
		f.isFileWell = false
		return
	}
	fname := fmt.Sprint(f.fileDir, "/", f.fileName)
	f.fileHandler, e = os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if e != nil {
		__print(default_format, LEVEL_ERROR, LEVEL_ERROR, 1, e.Error())
		f.isFileWell = false
		return
	}
	f.isFileWell = true
	f.tomorSecond = tomorSecond(f.mode)
	fs, err := f.fileHandler.Stat()
	if err == nil {
		f.fileSize = fs.Size()
	} else {
		e = err
	}
	return
}

func (f *fileObj) addFileSize(size int64) {
	atomic.AddInt64(&f.fileSize, size)
}

func (f *fileObj) write2file(bs []byte) (e error) {
	defer catchError()
	if bs != nil {
		f.addFileSize(int64(len(bs)))
		write2file(f.fileHandler, bs)
	}
	return
}

func (f *fileObj) isMustBackUp() bool {
	switch f.rolltype {
	case DAYLY:
		if _time().Unix() >= f.tomorSecond {
			return true
		}
	case ROLLFILE:
		return f.fileSize > 0 && f.fileSize >= f.maxSize*int64(f.unit)
	}
	return false
}

func (f *fileObj) rename() (err error) {
	bckupfilename := ""
	if f.rolltype == DAYLY {
		bckupfilename = getBackupDayliFileName(f.fileDir, f.fileName, f.mode)
	} else {
		bckupfilename, err = getBackupRollFileName(f.fileDir, f.fileName)
	}
	if bckupfilename != "" && err == nil {
		oldPath := fmt.Sprint(f.fileDir, "/", f.fileName)
		newPath := fmt.Sprint(f.fileDir, "/", bckupfilename)
		err = os.Rename(oldPath, newPath)
		if err == nil && f.rolltype == ROLLFILE && f.maxFileNum > 0 {
			go rmOverCountFile(f.fileDir, bckupfilename, f.maxFileNum)
		}
	}
	return
}

func (f *fileObj) close() (err error) {
	defer catchError()
	if f.fileHandler != nil {
		err = f.fileHandler.Close()
	}
	return
}

func tomorSecond(mode MODE_TIME) int64 {
	now := _time()
	switch mode {
	case MODE_DAY:
		return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()).Unix()
	case MODE_HOUR:
		return time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location()).Unix()
	case MODE_MONTH:
		return time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1).Unix()
	default:
		return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()).Unix()
	}
}

func yestStr(mode MODE_TIME) string {
	now := _time()
	switch mode {
	case MODE_DAY:
		return now.AddDate(0, 0, -1).Format(DATEFORMAT_DAY)
	case MODE_HOUR:
		return now.Add(-1 * time.Hour).Format(DATEFORMAT_HOUR)
	case MODE_MONTH:
		return now.AddDate(0, -1, 0).Format(DATEFORMAT_MONTH)
	default:
		return now.AddDate(0, 0, -1).Format(DATEFORMAT_DAY)
	}
}

/*————————————————————————————————————————————————————————————————————————————*/
func getBackupDayliFileName(dir, filename string, mode MODE_TIME) (bckupfilename string) {
	timeStr := yestStr(mode)
	index := strings.LastIndex(filename, ".")
	if index <= 0 {
		index = len(filename)
	}
	fname := filename[:index]
	suffix := filename[index:]
	bckupfilename = fmt.Sprint(fname, "_", timeStr, suffix)
	if isFileExist(fmt.Sprint(dir, "/", bckupfilename)) {
		bckupfilename = getBackupfilename(1, dir, fmt.Sprint(fname, "_", timeStr), suffix)
	}
	return
}

func getDirList(dir string) ([]os.DirEntry, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.ReadDir(-1)
}

func getBackupRollFileName(dir, filename string) (bckupfilename string, er error) {
	list, err := getDirList(dir)
	if err != nil {
		er = err
		return
	}
	index := strings.LastIndex(filename, ".")
	if index <= 0 {
		index = len(filename)
	}
	fname := filename[:index]
	suffix := filename[index:]
	length := len(list)
	bckupfilename = getBackupfilename(length, dir, fname, suffix)
	return
}

func getBackupfilename(count int, dir, filename, suffix string) (bckupfilename string) {
	bckupfilename = fmt.Sprint(filename, "_", count, suffix)
	if isFileExist(fmt.Sprint(dir, "/", bckupfilename)) {
		return getBackupfilename(count+1, dir, filename, suffix)
	}
	return
}

func write2file(f *os.File, bs []byte) (e error) {
	_, e = f.Write(bs)
	return
}

func console(s string, levelname string, flag FORMAT, calldepth int) {
	if flag != FORMAT_NANO {
		buf := getOutBuffer(s, levelname, flag, k1(calldepth))
		fmt.Print(&buf)
	} else {
		fmt.Println(s)
	}
}

func outwriter(out io.Writer, prefix string, flag FORMAT, calldepth int, s string) {
	l := log.New(out, prefix, int(flag))
	l.Output(k1(calldepth), s)
}

func k1(calldepth int) int {
	return calldepth + 1
}

func getOutBuffer(s string, levelname string, flag FORMAT, calldepth int) (buf bytes.Buffer) {
	outwriter(&buf, levelname, flag, k1(calldepth), s)
	return
}

func mkdirDir(dir string) (e error) {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, 0666); err != nil {
			if os.IsPermission(err) {
				e = err
			}
		}
	}
	return
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func catchError() {
	if err := recover(); err != nil {
		Fatal(string(debug.Stack()))
	}
}

func rmOverCountFile(dir, backupfileName string, maxFileNum int) {
	static_mu.Lock()
	defer static_mu.Unlock()
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	dirs, _ := f.ReadDir(-1)
	f.Close()
	if len(dirs) <= maxFileNum {
		return
	}
	sort.Slice(dirs, func(i, j int) bool {
		f1, _ := dirs[i].Info()
		f2, _ := dirs[j].Info()
		return f1.ModTime().Unix() > f2.ModTime().Unix()
	})
	index := strings.LastIndex(backupfileName, "_")
	indexSuffix := strings.LastIndex(backupfileName, ".")
	if indexSuffix == 0 {
		indexSuffix = len(backupfileName)
	}
	prefixname := backupfileName[:index+1]
	suffix := backupfileName[indexSuffix:]
	suffixlen := len(suffix)
	rmfiles := make([]string, 0)
	i := 0
	for _, f := range dirs {
		if len(f.Name()) > len(prefixname) && f.Name()[:len(prefixname)] == prefixname && matchString("^[0-9]+$", f.Name()[len(prefixname):len(f.Name())-suffixlen]) {
			finfo, err := f.Info()
			if err == nil && !finfo.IsDir() {
				i++
				if i > maxFileNum {
					rmfiles = append(rmfiles, fmt.Sprint(dir, "/", f.Name()))
				}
			}
		}
	}
	if len(rmfiles) > 0 {
		for _, k := range rmfiles {
			os.Remove(k)
		}
	}
}

func matchString(pattern string, s string) bool {
	b, err := regexp.MatchString(pattern, s)
	if err != nil {
		b = false
	}
	return b
}

func _time() time.Time {
	if TIME_DEVIATION != 0 {
		return time.Now().Add(TIME_DEVIATION)
	} else {
		return time.Now()
	}
}
