package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"DailyFresh/models"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"math"
)

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
		return userName.(string)
	}
	return ""
}

func (this*GoodsController)ShowIndex()  {
	GetUser(&this.Controller)
	o := orm.NewOrm()
	//获取类型数据
	var goodTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goodTypes)
	this.Data["goodTypes"] = goodTypes
	//获取轮播数据
	var indexGoodsBanner []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&indexGoodsBanner)
	this.Data["indexGoodsBanner"]=indexGoodsBanner
	//获取促销商品数据
	var promotionGoods []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&promotionGoods)
	this.Data["promotionGoods"]=promotionGoods

	//首页展示商品数据
	//定义并初始化商品map切片(map是有key值的切片)[]map[string]interface{}是相当于二维数组
	//这个容器是一个切片，然后切片的每个元素是一个map，这个map的值类型是interface{},
	//每个切片元素中追加一个map，这个map包含三个key类型，分别是type，textGoods和imgGoods。
	goods := make([]map[string]interface{},len(goodTypes))

	//向切片interface中插入类型数据
	for index, value := range goodTypes{
		//获取对应类型的首页展示商品
		temp := make(map[string]interface{})
		temp["type"] = value
		goods[index] = temp
	}
	//商品数据，向切片中共加入与每个类型对应的text和商品数据
	for _,value := range goods{
		var textGoods []models.IndexTypeGoodsBanner
		var imgGoods []models.IndexTypeGoodsBanner
		//获取文字商品数据
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType","GoodsSKU").OrderBy("Index").Filter("GoodsType",value["type"]).Filter("DisplayType",0).All(&textGoods)
		//获取图片商品数据
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType","GoodsSKU").OrderBy("Index").Filter("GoodsType",value["type"]).Filter("DisplayType",1).All(&imgGoods)

		value["textGoods"] = textGoods
		value["imgGoods"] = imgGoods

		//var testTextGoods []models.IndexTypeGoodsBanner
		//var testTextGoods2 []models.IndexTypeGoodsBanner
		//o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType","GoodsSKU").OrderBy("Index").All(&testTextGoods)
		//beego.Info("[testTextGoods]len:",len(testTextGoods),"[ShowIndex]testTextGoods:",testTextGoods)
		//beego.Info("/n")
		//o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType").OrderBy("Index").All(&testTextGoods2)
		//beego.Info("[testTextGoods222]len:",len(testTextGoods2),"[ShowIndex]testTextGoods2222:",testTextGoods2)
	}
	this.Data["goods"] = goods

	//返回视图
	this.TplName="index.html"
}

func ShowLayout(this *beego.Controller)  {
	//查询所有的类型
	o:=orm.NewOrm()
	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)
	beego.Info("[ShowLayout]types :",types)
	this.Data["types"] = types
	//获取用户信息
	GetUser(this)
	//指定layout
	this.Layout="goodsLayout.html"
}

func (this*GoodsController)ShowGoodsDetail()  {
	GetUser(&this.Controller)
	id,err := this.GetInt("id")
	if err != nil {
		beego.Error("浏览器请求数据错误")
		this.Redirect("/",302)
		return
	}

	o:=orm.NewOrm()
	var goodsSKU models.GoodsSKU
	goodsSKU.Id = id
	//o.Read(&goodsSKU)
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType","Goods").Filter("Id",id).All(&goodsSKU)
	//获取同类型时间考前的两条商品数据
	var goodsNew []models.GoodsSKU

	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType",goodsSKU.GoodsType).OrderBy("Time").Limit(2,0).All(&goodsNew)
	this.Data["goodSku"] = goodsSKU
	this.Data["goodsNew"] = goodsNew
	//添加历史浏览记录
	//判断用户是否登陆
	userName := this.GetSession("username")
	if userName != nil {
		o:=orm.NewOrm()
		var user models.User
		user.Name = userName.(string)
		o.Read(&user,"Name")
		//添加历史记录用redis存储
		conn,err := redis.Dial("tcp","192.168.80.130:6379")
		defer conn.Close()
		if err!=nil {
			beego.Info("redis连接失败")
		}
		_,err = conn.Do("lrem","history_"+strconv.Itoa(user.Id),0,id)
		if err != nil {
			conn.Do("lpush","history_"+strconv.Itoa(user.Id),id)
		}
		conn.Do("lpush","history_"+strconv.Itoa(user.Id),id)


	}
	//添加历史记录


	//返回视图
	ShowLayout(&this.Controller)
	cartCount := GetCartCount(&this.Controller)
	this.Data["cartCount"] = cartCount
	this.TplName="detail.html"
}

func  (this*GoodsController)ShowList(){
	//获取数据
	id, err := this.GetInt("typeId")

	//校验数据
	if err != nil{
		beego.Info("[ShowList]请求数据错误")
		this.Redirect("/",302)
		return
	}
	//处理数据
	ShowLayout(&this.Controller)
	//获取新品
	var goodsNew []models.GoodsSKU
	o:=orm.NewOrm()
	//o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).OrderBy("Time").Limit(2,0).All(&goodsNew)
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).OrderBy("Time").Limit(2,0).All(&goodsNew)
	this.Data["goodsNew"] = goodsNew
	//beego.Info("[ShowList]goodsNew:",goodsNew)
	var goods []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).All(&goods)
	this.Data["goods"] = goods

	//分页实现
	//count一个有多少个商品
	count,_ :=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).Count()
	pageSize := 3//每页有几个
	//每页可以显示多少页码？比如每页显示2个，总共有10个商品，那么总共可以显示5页面，pages是每个面总共几个页码数字可以选。
	pageCount := math.Ceil(float64(count)/float64(pageSize))//总共有多少页
	//pageIndex为当前页
	pageIndex,err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	pages := PageTool(int(pageCount),pageIndex)
	this.Data["pages"] = pages
	this.Data["typeId"] = id
	this.Data["pageIndex"] = pageIndex
	start := (pageIndex-1) * pageSize//每个页面是从第几个商品开始的


	prePage := pageIndex-1
	if prePage <= 1 {
		prePage = 1
	}
	nextPage := pageIndex + 1
	if nextPage > int(pageCount) {
		nextPage = int(pageCount)
	}
	this.Data["prePage"] = prePage
	this.Data["nextPage"] = nextPage

	sort := this.GetString("sort")
	if sort == "" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).Limit(pageSize,start).All(&goods)
		this.Data["goods"] = goods
		this.Data["sort"] = ""
	}else if sort == "price" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).OrderBy("Price").Limit(pageSize,start).All(&goods)
		this.Data["goods"] = goods
		this.Data["sort"] = "price"

	}else if sort == "sale" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).OrderBy("Sales").Limit(pageSize,start).All(&goods)
		this.Data["goods"] = goods
		this.Data["sort"] = "sale"
	}

	cartCount := GetCartCount(&this.Controller)
	this.Data["cartCount"] = cartCount
	//返回视图
	this.TplName = "list.html"
}

func PageTool(pageCount,pageIndex int) []int {
	var pages []int
	if pageCount <= 5 {
		pages = make([]int,pageCount)
		for i,_:= range pages{
			pages[i]=i+1
		}
	}else if pageIndex <= 3 {
		pages= []int{1,2,3,4,5}
	}else if pageIndex > pageCount-3 {
		pages=[]int{pageCount-4,pageCount-3,pageCount-2,pageCount-1,pageCount}
	}else {
		pages = []int{pageIndex-2,pageIndex -1,pageIndex,pageIndex+1,pageIndex+2}
	}
	return pages
}

func (this*GoodsController)HandleSearch() {
	//获取数据
	goodsName := this.GetString("goodsName")
	//校验数据
	o := orm.NewOrm()
	var goods []models.GoodsSKU
	if goodsName == "" {
		o.QueryTable("GoodsSKU").All(&goods)
		this.Data["goods"] = goods
		ShowLayout(&this.Controller)
		this.TplName = "search.html"
		return
	}
	//处理数据
	o.QueryTable("GoodsSKU").Filter("Name__icontains", goodsName).All(&goods)
	//返回视图
	this.Data["goods"] = goods
	ShowLayout(&this.Controller)
	this.TplName = "search.html"
}