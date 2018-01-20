package controllers

import (
	"logSystem/libs"
	slave "logSystem/models"
	"sync"

	"github.com/astaxie/beego"
	. "github.com/soekchl/myUtils"
)

type UserController struct {
	BaseController
}

var (
	userList      = make(map[string]*slave.User) // 用户帐号 -> 用户数据
	userListMutex sync.RWMutex                   // 用户操作锁
)

func init() {
	result, err := slave.ReadAllUser()
	if err != nil {
		Error(err)
		return
	}
	for i, _ := range result {
		SetUserList(&result[i])
	}
}

func (this *UserController) List() {
	if this.userName != "admin" {
		this.redirect("/")
		return
	}
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	userListMutex.RLock()
	defer userListMutex.RUnlock()

	count := len(userList)
	list := make([]map[string]interface{}, count)
	sort := make([]int, len(userList))
	n := 0
	for _, v := range userList {
		row := make(map[string]interface{})
		row["Id"] = v.Id
		row["account"] = v.UserName
		row["Email"] = v.Email
		if v.Status == -1 {
			row["status"] = "禁用"
		} else {
			row["status"] = "正常"
		}
		row["createTime"] = v.CreateTime.Format("2006/01/02 15:04:05")
		row["lastIp"] = v.LastIp
		row["lastLogin"] = v.LastLogin.Format("2006/01/02 15:04:05")
		list[n] = row
		sort[n] = int(v.Id)
		n++
	}
	for i := 0; i < count; i++ {
		for j := 0; j < count-i-1; j++ {
			if sort[j] < sort[j+1] {
				list[j+1], list[j] = list[j], list[j+1]
				sort[j], sort[j+1] = sort[j+1], sort[j]
			}
		}
	}

	this.Data["pageTitle"] = "帐号管理列表"

	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("UserController.List", "groupid", 1), true).ToString()
	this.display()
}

func (this *UserController) Status() {
	if this.userName != "admin" {
		this.redirect("/")
		return
	}

	id, _ := this.GetInt64("id")
	if id == 1 {
		this.redirect(beego.URLFor("UserController.List"))
	}
	Notice(id)
	u := GetUserListForId(id)
	if u == nil {
		Error("没有找到用户")
	} else {
		if u.Status == 0 {
			u.Status = -1
		} else {
			u.Status = 0
		}
		u.Update("status")
	}
	this.redirect(beego.URLFor("UserController.List"))
}

func SetUserList(user *slave.User) {
	Info("[SetUserList] name=", user.UserName)
	userListMutex.Lock()
	userList[user.UserName] = user
	userListMutex.Unlock()
}

func GetUserList(username string) *slave.User {
	Debug("[GetUserList] username=", username)
	if len(username) < 1 {
		return nil
	}

	userListMutex.RLock()         // rlock
	defer userListMutex.RUnlock() // runlock

	if c, exist := userList[username]; exist {
		return c
	}
	return nil
}

func GetUserListForId(id int64) *slave.User {
	Debug("[GetUserListForId] id=", id)
	if id < 1 {
		return nil
	}

	userListMutex.RLock()         // rlock
	defer userListMutex.RUnlock() // runlock
	for _, v := range userList {
		if v.Id == id {
			return v
		}
	}
	return nil
}
