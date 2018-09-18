package routers

import (
	"class1/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/Article/*",beego.BeforeRouter,FilterFunc)
	beego.Router("/", &controllers.LoginController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/register", &controllers.RegController{}, "get:ShowReg;post:HandleReg")
	beego.Router("/Article/ShowArticle", &controllers.ArticleController{}, "get:ShowArticle")
	beego.Router("/Article/ShowContent/?:id", &controllers.ArticleController{}, "get:ShowContent")
	beego.Router("/Article/DeleteArticle/?:id", &controllers.ArticleController{}, "get:DeleteArticle")
	beego.Router("/Article/UpdateArticle/?:id", &controllers.ArticleController{}, "get:ShowUpdate;post:UpdateArticle")
	beego.Router("/Article/AddArticle", &controllers.AddArticleController{}, "get:ShowAddArticle;post:HandleArticle")
	beego.Router("/Article/AddType", &controllers.AddArticleController{}, "get:ShowAddType;post:HandleAddType")
	beego.Router("/Article/Logout",&controllers.ArticleController{},"get:Logout")
	beego.Router("/Article/deleteType",&controllers.ArticleController{},"get:DeleteType")
}

var FilterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil{
		ctx.Redirect(302,"/")
	}
}
