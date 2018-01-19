package log

import (
	"logSystem/api"
	"logSystem/libs"
	"logSystem/net"

	"github.com/gogo/protobuf/proto"
)

type SaveLog struct {
	UserId    int64  // 操作用户id
	Level     int32  // 操作日志级别
	TimeStamp int64  // 生成时间戳
	FileName  string // 文件名
	FuncName  string // 函数名
	FileNo    int32  // 行号
	Info      string // 内容
}

func SendSaveLog(conn *net.Session, sl []*SaveLog) {
	if len(sl) < 1 || conn.IsClosed() {
		return
	}

	ld := &api.LogInfoReq{}

	for _, v := range sl {
		ld.LogDataList = append(ld.LogDataList, &api.LogInfo{
			UserId:    v.UserId,
			LogLeave:  getSaveLogLeave(v.Level),
			Timestamp: v.TimeStamp,
			FileName:  v.FileName,
			FuncName:  v.FuncName,
			FileNo:    v.FileNo,
			LogInfo:   v.Info,
		})
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

func getSaveLogLeave(leave int32) api.LOG_LEAVE {
	switch leave {
	case int32(api.LOG_LEAVE_LEAVE_DEBUG):
		return api.LOG_LEAVE_LEAVE_DEBUG
	case int32(api.LOG_LEAVE_LEAVE_TRACK):
		return api.LOG_LEAVE_LEAVE_TRACK
	case int32(api.LOG_LEAVE_LEAVE_ERROR):
		return api.LOG_LEAVE_LEAVE_ERROR
	case int32(api.LOG_LEAVE_LEAVE_INFO):
		return api.LOG_LEAVE_LEAVE_INFO
	case int32(api.LOG_LEAVE_LEAVE_WARR):
		return api.LOG_LEAVE_LEAVE_WARR
	}
	return api.LOG_LEAVE_LEAVE_TRACK
}
