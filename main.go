package main

import (
	_ "logSystem/models"
	_ "logSystem/routers"

	"github.com/astaxie/beego"
)

const (
	VERSION = "1.0.0"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}
