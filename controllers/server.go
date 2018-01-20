package controllers

// 开启服务器  接收数据并且 保存到数据库里

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"logSystem/api"
	"logSystem/libs"
	. "logSystem/log"
	"logSystem/modelsMaster"
	"logSystem/net"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/gogo/protobuf/proto"
)

var (
	log_chan         = make(chan *net.FormatData, 5000)
	buffLogList      []*models.Log // 多条插入log缓存
	buffLogListMutex sync.Mutex    // 锁
)

func init() {
	go processSaveLog()
	go autoSaveLog()

	log_chan, _ := beego.AppConfig.Int("log.chan_num")
	if log_chan < 2000 {
		log_chan = 2000
	}
	go startLogServer(
		beego.AppConfig.String("log.api"),
		log_chan,
	)
}

func startLogServer(address string, chan_num int) {
	Warn("启动日志系统 add=", address)
	cert, err := tls.LoadX509KeyPair("key/server.pem", "key/server.key")
	if err != nil {
		Error(err)
		return
	}
	certBytes, err := ioutil.ReadFile("key/client.pem")
	if err != nil {
		panic("Unable to read cert.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}
	server, err := net.Listen("tcp", address, config, chan_num, net.HandlerFunc(outServerSessionLoop))
	if err != nil {
		Error(err)
		return
	}
	server.Serve()
}

func outServerSessionLoop(session *net.Session) {
	Warn("[outServerSessionLoop] 新用户 ", session.RemoteAddr())
	defer session.Close() // 关闭连接

	var buff []byte
	var ok bool

	for {
		// 服务器 接收
		select {
		case buff, ok = <-session.ByteRecvChan:
			if !ok {
				Error("Session is Closed! ", session.RemoteAddr())
				return
			}
		}

		log_chan <- &net.FormatData{
			Id:   int32(libs.DecodeUint32(buff)),
			Body: buff[net.HeadSize-net.FirstReadSize:],
		}
	}
}

func processSaveLog() {
	Warn("[processSaveLog] 启动")
	for d := range log_chan {

		if d.Id == int32(api.API_add_log) {
			li := &api.LogInfoReq{}
			err := proto.Unmarshal(d.Body, li)
			if err != nil {
				Error(err)
				continue
			}
			for _, v := range li.GetLogDataList() {
				addLogBuff(&models.Log{
					UserId:    v.GetUserId(),
					Level:     int32(v.GetLogLeave()),
					TimeStamp: v.GetTimestamp(),
					FileName:  v.GetFileName(),
					FuncName:  v.GetFuncName(),
					FileNo:    v.GetFileNo(),
					LogInfo:   v.GetLogInfo(),
				})
			}
		}
	}
}

// 每秒 循环一次 判断是否需要插入
func autoSaveLog() {
	Warn("[autoSaveLog] 启动")
	for {
		time.Sleep(time.Second)
		if len(buffLogList) > 0 {
			saveLogBuff()
		}
	}
}

// 增加缓存
func addLogBuff(log *models.Log) {
	buffLogListMutex.Lock()
	defer buffLogListMutex.Unlock()

	buffLogList = append(buffLogList, log)
}

// 插入到数据库
func saveLogBuff() {
	buffLogListMutex.Lock()
	defer buffLogListMutex.Unlock()
	Info("[saveLogBuff] len=", len(buffLogList))

	err := models.LogMultiAdd(buffLogList)
	if err != nil {
		Error(err)
	} else {
		buffLogList = []*models.Log{} // 清空数据
	}

}
