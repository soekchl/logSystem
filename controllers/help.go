package controllers

type HelpController struct {
	BaseController
}

func (this *HelpController) Index() {

	this.Data["pageTitle"] = "使用帮助"
	this.Data["adminFlag"] = this.userName
	this.display()
}
