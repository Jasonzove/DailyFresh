package controllers

import "github.com/astaxie/beego"

type GoodsController struct {
	beego.Controller
}

func GetUser(this *beego.Controller) string {
	userName := this.GetSession("username")
	beego.Info("[ShowIndex]userName : ",userName)
	if userName == nil{
		this.Data["userName"] = ""
	}else {
		this.Data["userName"] = userName.(string)
	}
	return userName.(string)
}

func (this*GoodsController)ShowIndex()  {
	GetUser(&this.Controller)
	this.TplName="index.html"
}