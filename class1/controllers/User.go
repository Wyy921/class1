package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"class1/models"
	"time"
)

type RegController struct {
	beego.Controller
}

func (this *RegController) ShowReg() {

	this.TplName = "register.html"
}

func (this *RegController) HandleReg() {
	//1.拿到浏览器传递的数据
	name := this.GetString("userName")
	passwd := this.GetString("password")
	//2.数据处理
	if name == "" || passwd == "" {
		beego.Info("用户名或者密码不能为空")
		this.TplName = "register.html"
		return
	}
	o := orm.NewOrm()
	user := models.User{}
	user.UserName = name
	user.Password = passwd
	_, err := o.Insert(&user)
	if err != nil {
		beego.Info("Inserterr", err)
		return
	}

	this.Redirect("/", 302)
}

type LoginController struct {
	beego.Controller
}

func (this *LoginController) ShowLogin() {
	name := this.Ctx.GetCookie("userName")
	beego.Info("AAAAAAAAAAAAA",name)

	if name != "" {
		this.Data["userName"] = name
		this.Data["check"] = "checked"
	}
	this.Data["data"] = "aaa"

	this.TplName = "login.html"
}

func (this *LoginController) HandleLogin() {
	name := this.GetString("userName")
	passwd := this.GetString("password")
	if name == "" || passwd == "" {
		beego.Info("用户名或者密码不能为空")
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	user := models.User{}
	user.UserName = name
	err := o.Read(&user, "Username")
	if err != nil {
		beego.Info("用户名失败")
		this.TplName = "login.html"
		return
	}
	if user.Password != passwd {
		beego.Info("密码错误")
		this.TplName = "login.html"
	}

	check := this.GetString("remember")
	beego.Info(check)
	if check == "on" {
		this.Ctx.SetCookie("userName", name, time.Second*3600)
	} else {
		this.Ctx.SetCookie("userName", "sss", -1)
	}

	this.SetSession("userName",name)

	this.Redirect("/Article/ShowArticle", 302)
}
