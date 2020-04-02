package controllers

import(
	"github.com/astaxie/beego"
	"regexp"
	"github.com/astaxie/beego/orm"
	"DailyFresh/models"
	"github.com/astaxie/beego/utils"
	"strconv"
	"encoding/base64"
)

type  UserController struct {
	beego.Controller
}

func (this*UserController)ShowReg(){
	this.TplName="register.html"
}
func (this*UserController) HandleReg()  {
	//1获取数据
	userName := this.GetString("user_name")
	pwd := this.GetString("pwd")
	cpwd := this.GetString("cpwd")
	email := this.GetString("email")
	//2校验数据
	if userName == "" || pwd == "" || cpwd == "" || email == "" {
		this.Data["errmsg"] = "数据不完整，清重新注册～"
		this.TplName="register.html"
		return
	}
	if pwd != cpwd {
		this.Data["errmsg"] = "两次输入密码不一致，清重新注册～"
		this.TplName="register.html"
		return
	}
	reg,_ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res := reg.FindString(email)
	if res == "" {
		this.Data["errmsg"] = "邮箱格式不正确，清重新注册～"
		this.TplName="register.html"
		return
	}
	//3处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.PassWord = pwd
	user.Email = email
	_,err := o.Insert(&user)
	if err != nil {
		this.Data["errmsg"] = "注册失败，清重新注册～"
		this.TplName="register.html"
		return
	}
	emailConfig := `{"username":"anlz729@163.com","password":"XBKQFTLJHHPJMNAH","host":"smtp.163.com","port":25}`
	emailConn := utils.NewEMail(emailConfig)
	emailConn.From ="天天生鲜系统注册服务"//"anlz729@163.com" //"天天生鲜系统注册服务"
	emailConn.To = []string{email}
	emailConn.Subject = "天天生鲜用户注册"
	//注意这里我们发送给用户的是激活请求地址
	emailConn.Text="192.168.1.106:8080/Active?id="+strconv.Itoa(user.Id)

	err = emailConn.Send()
	if err != nil {
		beego.Info(err)
		this.Ctx.WriteString("发送到邮箱失败")
	}
	//4返回视图
	this.Ctx.WriteString("注册成功,请去相应的邮箱激活")
}

func (this*UserController)ActiveUser()  {
	id,err:= this.GetInt("id")
	if err != nil {
		this.Data["errmsg"] = "要激活的用户不存在"
		this.TplName="register.html"
		return
	}
	o := orm.NewOrm()
	user := &models.User{Id:id}
	err = o.Read(user)
	if err != nil {
		this.Data["errmsg"] = "要激活的用户不存在"
		this.TplName="register.html"
		return
	}

	user.Active =true
	o.Update(user)
	this.Redirect("/login",302)

}

func (this*UserController)ShowLogin()  {
	userName := this.Ctx.GetCookie("username")
	beego.Info("userName:",userName)
	temp,_ := base64.StdEncoding.DecodeString(userName)
	beego.Info(string(temp))
	if string(temp)== "" {
		this.Data["username"]= ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"]= string(temp)
		//要选中该控件需要设置属性checked,否则为unchecked
		this.Data["checked"] = "checked"
	}
	this.TplName="login.html"
}

func (this*UserController) Handlelogin()  {
	//获取数据
	username := this.GetString("username")
	pwd := this.GetString("pwd")
	remember := this.GetString("remember")
	//校验数据
	if username ==""||pwd == "" {
		this.Data["errmsg"] = "用户名或密码不存在，清重新输入"
		this.TplName="login.html"
		return
	}
	//处理数据
	user := &models.User{}
	user.Name = username
	o:=orm.NewOrm()
	err := o.Read(user,"Name")
	if err != nil {
		this.Data["errmsg"] = "用户名错误，清重新输入"
		this.TplName="login.html"
		return
	}
	if user.PassWord != pwd {
		this.Data["errmsg"] = "用户密码错误，清重新输入"
		this.TplName="login.html"
		return
	}
	if user.Active != true {
		this.Data["errmsg"] = "用户未激活，请返回邮箱激活"
		this.TplName="login.html"
		return
	}
	beego.Info("remember:",remember)
	//checkbox控件表明一个特定的状态（即选项）是选定 (on，值为true) 还是清除
	//(off，值为false)
	if  remember == "on"{
		temp := base64.StdEncoding.EncodeToString([]byte(username))
		this.Ctx.SetCookie("username",temp,24*3600*30)
	} else {
		this.Ctx.SetCookie("username",username,-1)
	}
	//返回视图
	this.Ctx.WriteString("登陆成功")
	//this.Redirect("index.html",302)
}