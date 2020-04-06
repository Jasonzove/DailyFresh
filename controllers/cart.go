package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"DailyFresh/models"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type CartController struct {
	beego.Controller
}


func (this *CartController)HandleAddCart(){
	//获取数据
	skuid,err1 := this.GetInt("skuid")
	count,err2 := this.GetInt("count")
	resp := make(map[string]interface{})
	defer this.ServeJSON()//发送json格式的数据
	//检验数据
	if err1!= nil || err2 != nil {
		resp["code"] = 1
		resp["msg"] ="传递的数据正确"
		this.Data["json"] = resp
		return
	}
	//处理数据
	userName := this.GetSession("username")
	if userName == nil {
		resp["code"] = 2
		resp["msg"] ="当前用户未登陆"
		this.Data["json"] = resp
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	o.Read(&user,"Name")

	//购物车数据存在redis中用hash
	conn,err := redis.Dial("tcp","192.168.1.102:6379")
	if err != nil {
		resp["code"] = 3
		resp["msg"] ="redis数据库连接错误"
		this.Data["json"] = resp
		return
	}
	//先获取原来的数量，然后给数量加起来
	preRep,err :=conn.Do("hget","cart_"+strconv.Itoa(user.Id),skuid)
	preCount,err:= redis.Int(preRep,err)
	conn.Do("hset","cart_"+strconv.Itoa(user.Id),skuid,count+preCount)
	rep,err := conn.Do("hlen","cart_"+strconv.Itoa(user.Id))
	//回复助手函数
	cartCount,_ := redis.Int(rep,err)
	//返回json数据
	resp["code"] = 5
	resp["msg"] ="OK"
	resp["cartCount"] = cartCount
	this.Data["json"] = resp
}

//获取购物车数量函数
func GetCartCount(this *beego.Controller) int {
	userName := this.GetSession("username")
	if userName == nil {
		beego.Info("[GetCartCount]userName:",userName)
		return 0
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	o.Read(&user,"Name")
	beego.Info("[GetCartCount]o.Read,user:",user.Id)
	conn,err := redis.Dial("tcp","192.168.1.102:6379")
	if err != nil {
		return 0
	}
	defer conn.Close()
	rep,err := conn.Do("hlen","cart_"+strconv.Itoa(user.Id))
	cartCount,_:=redis.Int(rep,err)
	beego.Info("[GetCartCount]cartCount:",cartCount)
	return cartCount
}

func (this *CartController)ShowCart() {
	//获取用户信息
	userName := GetUser(&this.Controller)
	o:=orm.NewOrm()
	var user models.User
	user.Name = userName
	o.Read(&user,"Name")
	//获取redis数据库信息
	conn,_ := redis.Dial("tcp","192.168.1.102:6379")
	rep,err := conn.Do("hgetall","cart_"+strconv.Itoa(user.Id)) //返回的是map[string]int
	goodsMap,_:=redis.IntMap(rep,err)
	goods := make([]map[string]interface{},len(goodsMap))
	i:=0
	totalPrice := 0
	totalCount := 0
	for index,value := range goodsMap {
		skuid, _ := strconv.Atoi(index)
		var goodsSku models.GoodsSKU
		goodsSku.Id=skuid
		o.Read(&goodsSku)
		temp := make(map[string]interface{})
		temp["goodsSku"] = goodsSku
		temp["count"] = value
		totalPrice += goodsSku.Price*value
		totalCount += value
		temp["addPrice"] = goodsSku.Price * value
		goods[i] = temp
		i++
	}
	this.Data["goods"] = goods
	this.Data["totalPrice"] = totalPrice
	this.Data["totalCount"] = totalCount


	this.TplName="cart.html"
}

func (this *CartController)HandleUpdateCart()  {
	skuid,err1 := this.GetInt("skuid")
	count,err2 := this.GetInt("count")
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	if err1 != nil || err2 !=nil {
		resp["code"] = 1
		resp["errmsg"] = "请求数据不正确"
		this.Data["json"] = resp
		return
	}
	userName := this.GetSession("username")
	if userName == nil {
		resp["code"] = 3
		resp["errmsg"] = "当前用户未登陆"
		this.Data["json"] = resp
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Name=userName.(string)
	o.Read(&user,"Name")


	conn,err := redis.Dial("tcp","192.168.1.102:6379")
	if err!= nil {
		resp["code"] = 2
		resp["errmsg"] = "redis连接失败"
		this.Data["json"] = resp
		return
	}
	defer conn.Close()
	conn.Do("hset","cart_"+strconv.Itoa(user.Id),skuid,count)
	resp["code"]= 5
	resp["errmsg"] = "OK"
	this.Data["json"] = resp
}