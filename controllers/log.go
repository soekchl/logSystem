package controllers

import (
	"fmt"
	"logSystem/libs"
	"logSystem/models"
	"time"

	"github.com/astaxie/beego"
	. "github.com/soekchl/myUtils"
)

type LogController struct {
	BaseController
}

func (this *LogController) List() {
	Notice("list")
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	logs, _ := (&models.Log{}).ReadAll(nil, 100)

	count := len(logs)

	list := make([]map[string]interface{}, count)
	for i, v := range logs {
		row := make(map[string]interface{})
		row["UserId"] = v.UserId
		row["Level"] = fmt.Sprint(v.Level)
		row["TimeStamp"] = time.Unix(v.TimeStamp, 0).Format("2006.01.02 15:04:05")
		row["FileName"] = v.FileName
		row["FuncName"] = v.FuncName
		row["FileNo"] = v.FileNo
		row["Info"] = v.LogInfo
		list[i] = row
	}

	this.Data["pageTitle"] = "日志预览"

	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(0, count, 1024, beego.URLFor("LogController.log", "groupid", 1), true).ToString()
	this.display()
}
