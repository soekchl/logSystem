package models

import (
	"time"
)

type User struct {
	Id         int64  `orm:"auto"`           // 每个帐号唯一ID
	UserName   string `orm:"size(32);index"` // 账号
	Password   string `orm:"size(32);"`      // 密码
	Salt       string `orm:"size(16);"`      // 加密项
	Email      string `orm:"size(32);index"` // 邮件
	Phone      string `orm:"size(11);index"` // 邮件
	LastLogin  time.Time
	LastIp     string    `orm:"size(32);"` // 最后登录的ip
	Status     int       // 账号状态  -1 禁用
	CreateTime time.Time `orm:"type(timestamp);null;auto_now_add"` // 创建时间
}

func (u *User) TableName() string {
	return TableName("user")
}

func (u *User) Update(fields ...string) error {
	if _, err := m_orm.Update(u, fields...); err != nil {
		return err
	}
	return nil
}

func UserAdd(user *User) (int64, error) {
	return m_orm.Insert(user)
}

func UserGetById(id int64) (*User, error) {
	u := new(User)

	err := m_orm.QueryTable(TableName("user")).Filter("id", id).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserGetByName(userName string) (*User, error) {
	u := new(User)

	err := m_orm.QueryTable(TableName("user")).Filter("user_name", userName).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func ReadAllUser() (users []User, err error) {
	qs := m_orm.QueryTable(&User{})
	_, err = qs.All(&users)
	return
}

func UserUpdate(user *User, fields ...string) error {
	_, err := m_orm.Update(user, fields...)
	return err
}

func UserGetMyList(page, pageSize int) ([]*User, int64) {
	offset := (page - 1) * pageSize
	list := make([]*User, 0)
	query := m_orm.QueryTable(TableName("user"))
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, int64(len(list))
}
