package models

import (
	"time"
)

type Log struct {
	Id         int64     `orm:"auto"`                              // 每个帐号唯一ID
	UserId     int64     `orm:"default(0);index"`                  // 操作用户id
	Level      int32     `orm:"default(0);index"`                  // 操作日志级别
	TimeStamp  int64     `orm:"default(0);index"`                  // 生成时间戳
	FileName   string    `orm:"size(32);index"`                    // 文件名
	FuncName   string    `orm:"size(64);index"`                    // 函数名
	FileNo     int32     `orm:"default(0)"`                        // 行号
	CreateTime time.Time `orm:"type(timestamp);null;auto_now_add"` // 创建时间
	LogInfo    string    `orm:"type(text)"`                        // 内容
}

const (
	Track = iota
	Debug
	Info
	Warr
	Error
)

func (u *Log) TableName() string {
	return TableName("Log")
}

func LogAdd(log *Log) (int64, error) {
	return m_orm.Insert(log)
}

func LogMultiAdd(logs []*Log) error {
	_, err := m_orm.InsertMulti(len(logs), logs)
	return err
}

func (this *Log) ReadAll(query map[string]string, limit int) (logs []Log, err error) {
	qs := slave_orm.QueryTable(this)
	if query != nil {
		for k, v := range query {
			qs = qs.Filter(k, v)
		}
	}

	qs = qs.OrderBy("-id").Limit(limit)

	_, err = qs.All(&logs)
	return
}
