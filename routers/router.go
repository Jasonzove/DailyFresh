package routers

import (
	"DailyFresh/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//方法名一定要大写

	//如果需要登陆之后才能访问的页面需要加过滤函数，即需要redirect("/user/......")前面需要加一个过序的/user/，并需要注册下面的代码beego.InsertFilter
	beego.InsertFilter("/user/*",beego.BeforeExec,filterFunc)
    //beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{},"get:ShowReg;post:HandleReg")
	beego.Router("/active", &controllers.UserController{},"get:ActiveUser")
	beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:Handlelogin")
	beego.Router("/", &controllers.GoodsController{},"get:ShowIndex")
	beego.Router("/user/logout", &controllers.UserController{},"get:Logout")
	//用户中心信息页
	beego.Router("/user/userCenterInfo",&controllers.UserController{},"get:ShowUserCenterInfo")
	//用户中心订单页
	beego.Router("/user/userCenterOrder",&controllers.UserController{},"get:ShowUserCenterOrder")
	//用户订单地址页
	beego.Router("/user/userCenterSite",&controllers.UserController{},"get:ShowUserCenterSite;post:HandleUserCenterSite")



}

//过滤器函数

 var filterFunc = func(ctx *context.Context){
	 userName := ctx.Input.Session("username")
	//userName := ctx.GetCookie("username")
	 if userName == "" {
		 ctx.Redirect(302,"/login")
		 return
	 }
}
