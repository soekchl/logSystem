// myUtils project myUtils.go
package log

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"logSystem/api"
	"logSystem/libs"
	"logSystem/net"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/gogo/protobuf/proto"
)

const (
	_ = iota
	LeaveDebug
	LeaveInfo
	LeaveWarning
	LeaveError
	LeaveNoShow

	max_buff_size = 65536
	max_file_size = 1024 * 1024 * 50 // 50M
)

var (
	file_log_name string
	dir_log_name  = "myLog"
	file_name     = ""
	file_log_flag = false
	show_leave    = LeaveDebug // 默认全输出
	out_put_leave = LeaveDebug // 默认全输出

	log_buff         = bytes.NewBuffer(make([]byte, max_buff_size))
	out_put_log_time = time.Second / 5
	out_put_log_chan = make(chan string, 1000)
	enter            = "\n"
	_file_format     string

	save_system   bool
	save_log_chan chan *api.LogInfo // 往系统发送的通道
	save_address  string            // 日志服务器地址
	save_log_list []*api.LogInfo    // 缓存日志数据
)

// 设定显示log等级
func SetShowLeave(leave int) {
	show_leave = getLeave(leave)
}

// 设定输出log等级
func SetOutPutLeave(leave int) {
	out_put_leave = getLeave(leave)
}

func getLeave(leave int) int {
	switch leave {
	case LeaveInfo:
		return LeaveInfo
	case LeaveWarning:
		return LeaveWarning
	case LeaveError:
		return LeaveError
	case LeaveNoShow:
		return LeaveNoShow
	}
	return LeaveDebug
}

func init() {
	if runtime.GOOS == "windows" {
		enter = "\r\n"
	} else {
		enter = "\n"
	}

	_file_format = "%s\\%s_%s_%d.log"
	if runtime.GOOS != "windows" {
		_file_format = "%s/%s_%s_%d.log"
	}
}

func SetOutputFileLog(log_file_name string) {
	file_name = log_file_name
	dir_log_name = fmt.Sprintf("%s_log", file_name)
	checkFileSize()
	file_log_flag = true
	log_buff.Reset()
	go outPutLogLoop()
}

func checkFileSize() {
	// 判断是否存在  判断大小
	var file os.FileInfo
	var name string
	var err error

	for i := 0; ; i++ {
		name = fmt.Sprintf(_file_format, dir_log_name, time.Now().Format("20060102"), file_name, i)
		file, err = os.Stat(name)
		if err != nil {
			break
		}
		if file.Size() < int64(max_file_size) {
			break
		}
	}
	file_log_name = name
}

func SetOutPutLogIntervalTime(interval int64) {
	if interval < 1 {
		return
	}
	out_put_log_time = time.Duration(interval)
}

func Track(v ...interface{}) {
	if show_leave <= LeaveDebug || (file_log_flag && out_put_leave <= LeaveDebug) || save_system {
		myLog(api.LOG_LEAVE_LEAVE_TRACK, "[T]", show_leave <= LeaveDebug, out_put_leave <= LeaveDebug, v...)
	}
}

func Debug(v ...interface{}) {
	if show_leave <= LeaveDebug || (file_log_flag && out_put_leave <= LeaveDebug) || save_system {
		myLog(api.LOG_LEAVE_LEAVE_DEBUG, "[D]", show_leave <= LeaveDebug, out_put_leave <= LeaveDebug, v...)
	}
}

func Info(v ...interface{}) {
	if show_leave <= LeaveInfo || (file_log_flag && out_put_leave <= LeaveInfo) || save_system {
		myLog(api.LOG_LEAVE_LEAVE_INFO, "[I]", show_leave <= LeaveInfo, out_put_leave <= LeaveInfo, v...)
	}
}

func Warn(v ...interface{}) {
	if show_leave <= LeaveWarning || (file_log_flag && out_put_leave <= LeaveWarning) || save_system {
		myLog(api.LOG_LEAVE_LEAVE_WARR, "[W]", show_leave <= LeaveWarning, out_put_leave <= LeaveWarning, v...)
	}
}

func Error(v ...interface{}) {
	if show_leave <= LeaveError || (file_log_flag && out_put_leave <= LeaveError) || save_system {
		myLog(api.LOG_LEAVE_LEAVE_ERROR, "【E】", show_leave <= LeaveError, out_put_leave <= LeaveError, v...)
	}
}

func myLog(leave api.LOG_LEAVE, mark string, show bool, out_put bool, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	temp := fmt.Sprint(v...)
	outstring := fmt.Sprintf("%s %s %-16s %v%s",
		time.Now().Format("2006/01/02 15:04:05"), mark, fmt.Sprintf("%s:%d", filename, line), temp, enter)

	if show {
		fmt.Print(outstring)
	}
	if file_log_flag && out_put {
		out_put_log_chan <- outstring
	}
	if save_system {
		save_log_chan <- &api.LogInfo{
			LogLeave:  leave,
			Timestamp: time.Now().Unix(),
			FileName:  filename,
			FileNo:    int32(line),
			LogInfo:   temp,
		}
	}
}

func outPutLogLoop() {
	t := time.Now().UnixNano() // 最后一次输出log时间
	for file_log_flag {
		select {
		case <-time.After(out_put_log_time):
			if log_buff.Len() > 0 { //	等待后续log到一定时间 以后输出log
				outputLog()
				t = time.Now().UnixNano()
			}
		case buff, ok := <-out_put_log_chan:
			if ok {
				if log_buff.Len()+len(buff) > max_buff_size { // 当缓存 超过限定的时候 提前输出
					outputLog()
					t = time.Now().UnixNano()
				}
				log_buff.Write([]byte(buff)) // 写入到缓冲区
			}
		}

		// 当log 一定时间段内没有输出就输出一次log
		if log_buff.Len() > 0 && (time.Now().UnixNano()-t) > int64(out_put_log_time) {
			outputLog()
		}
	}
}

func outputLog() {
	if _, err := os.Stat(dir_log_name); err != nil {
		if err := os.Mkdir(dir_log_name, 0755); err != nil {
			fmt.Println(err, "Mkdir")
			return
		}
	}

	file, err := os.OpenFile(file_log_name, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		file, err = os.Create(file_log_name)
		if err != nil {
			fmt.Println("Error!!! file", err)
			return
		}
	}

	file.Write(log_buff.Bytes())
	log_buff.Reset()
	file.Close()
	checkFileSize()
}

func sendData(conn *net.Session) {
	ld := &api.LogInfoReq{
		LogDataList: save_log_list,
	}
	data := &net.FormatData{
		Id: int32(api.API_add_log),
	}
	var err error

	data.Body, err = proto.Marshal(ld)
	if err != nil {
		Error(err)
		return
	}

	buff := make([]byte, len(data.Body)+net.HeadSize)
	copy(buff[net.HeadSize:], data.Body)
	data.Size = int32(len(data.Body) + net.HeadSize - net.FirstReadSize)
	libs.EncodeUint32(uint32(data.Size), buff)
	libs.EncodeUint32(uint32(data.Id), buff[4:])

	if !conn.IsClosed() {
		conn.ByteSendChan <- buff
		save_log_list = []*api.LogInfo{} // 发送完毕初始化
	} else {
		Error("日志系统连接已关闭！！！")
	}
}

// 设置往服务器发送日志数据
func SetSaveLog(address string, chan_num int) (err error) {
	// 创建tcp连接
	conn := getConn(address)
	if conn == nil {
		return errors.New("连接不了...")
	}
	save_system = true

	save_log_chan = make(chan *api.LogInfo, chan_num)
	go func() {
		for {
			select {
			case <-time.After(time.Second / 2):
				// 往服务器发送日志数据
				sendData(conn)
			case li := <-save_log_chan:
				save_log_list = append(save_log_list, li)
			}
		}
	}()
	return nil
}

func getConn(address string) *net.Session {
	cert, err := tls.LoadX509KeyPair("key/client.pem", "key/client.key")
	if err != nil {
		Error(err)
		return nil
	}
	certBytes, err := ioutil.ReadFile("key/client.pem")
	if err != nil {
		Error("Unable to read cert.pem ", err)
		return nil
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		Error("failed to parse root certificate")
		return nil
	}
	conf := &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	client, err := net.Dial("tcp", address, conf, 2000)
	if err != nil {
		Error(err)
		return nil
	}

	go connProcess(client)
	return client
}

func connProcess(session *net.Session) {
	for {
		// 服务器 接收
		select {
		// 只有异常情况才会返回
		case <-session.ByteRecvChan:
			Error("logSystem 异常！ 数据发送错误！")
		}
	}
}
