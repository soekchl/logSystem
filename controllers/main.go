package controllers

import (
	"fmt"
	"logSystem/libs"
	"logSystem/modelsMaster"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	. "github.com/soekchl/myUtils"
)

type MainController struct {
	BaseController
}

// 首页
func (this *MainController) Index() {

	this.Data["pageTitle"] = "待定"
	this.Data["month"] = "0"
	this.Data["leave"] = "0"
	this.display()
}

//个人信息
func (this *MainController) Profile() {
	beego.ReadFromRequest(&this.Controller)
	user := GetUserListForId(this.userId)

	if this.isPost() && user != nil {
		user.Email = this.GetString("email")
		user.Update()
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")
		if password1 != "" {
			if len(password1) < 6 || len(password1) > 32 {
				this.ajaxMsg("密码长度必须大于6位 小于32位", MSG_ERR)
			} else if password2 != password1 {
				this.ajaxMsg("两次输入的密码不一致", MSG_ERR)
			} else {
				user.Salt = string(utils.RandomCreateBytes(10))
				user.Password = libs.Md5([]byte(password1 + user.Salt))
				user.Update()
			}
		}
		this.ajaxMsg("", MSG_OK)
	}

	this.Data["pageTitle"] = "资料修改"
	this.Data["user"] = user
	this.display()
}

// 登录
func (this *MainController) Login() {
	if this.userId > 0 {
		this.redirect("/")
	}
	beego.ReadFromRequest(&this.Controller)
	if this.isPost() {

		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")

		if username != "" && password != "" {
			user := GetUserList(username)
			flash := beego.NewFlash()
			errorMsg := ""
			if user == nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "帐号或密码错误"
			} else if user.Status == -1 {
				errorMsg = "该帐号已禁用"
			} else {
				user.LastIp = this.getClientIp()
				user.LastLogin = time.Now()
				models.UserUpdate(user)

				authkey := libs.Md5([]byte(this.getClientIp() + "|" + user.Password + user.Salt))
				if remember == "yes" {
					this.Ctx.SetCookie("auth", fmt.Sprint(user.Id)+"|"+authkey, 7*86400)
				} else {
					this.Ctx.SetCookie("auth", fmt.Sprint(user.Id)+"|"+authkey, 86400)
				}
				this.redirect(beego.URLFor("MainController.Index"))
			}
			flash.Error(errorMsg)
			flash.Store(&this.Controller)
			this.redirect(beego.URLFor("MainController.Login"))
		}
	}
	this.TplName = "public/login.html"
}

// 重置密码
func (this *MainController) ResetPasswd() {
	Notice("ResetPasswd")
	this.TplName = "public/resetpasswd.html"
}

func (this *MainController) Agreement() {
	Notice("Agreement")
	this.TplName = "public/agreement.html"
}

// 退出登录
func (this *MainController) Logout() {
	this.Ctx.SetCookie("auth", "")
	this.redirect(beego.URLFor("MainController.Login"))
}

// 获取系统时间
func (this *MainController) GetTime() {
	out := make(map[string]interface{})
	out["time"] = time.Now().UnixNano() / int64(time.Millisecond)
	this.jsonResult(out)
}
