package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"class1/models"
	"strconv"
	"math"
	"github.com/gomodule/redigo/redis"
	"bytes"
	"encoding/gob"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowArticle() {

	o := orm.NewOrm()
	types := make([]models.ArticleType, 0)
	conn, _ := redis.Dial("tcp", "192.168.223.149:6379")
	defer conn.Close()
	rel, err02 := redis.Bytes(conn.Do("get", "types"))
	if err02 != nil {
		beego.Info("redis数据库读取操作错误", err02)
	}

	dec := gob.NewDecoder(bytes.NewReader(rel))
	dec.Decode(&types)
	beego.Info(types)

	//存入redis数据库
	if len(types) == 0 {
		o.QueryTable("ArticleType").All(&types)
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(&types)
		conn, _ := redis.Dial("tcp", "192.168.223.148:6379")
		defer conn.Close()
		_, err01 := conn.Do("set", "types", buffer.Bytes())
		if err01 != nil {
			beego.Info("redis数据库操作错误", err01)
		}

		beego.Info("从数据库读取")

	}

	typeName := this.GetString("select")

	qs := o.QueryTable("Article")

	pageIndex := this.GetString("pageIndex")
	pageIndex1, err := strconv.Atoi(pageIndex)
	if err != nil {
		pageIndex1 = 1 //自己
	}

	pageSize := 2
	start := pageSize * (pageIndex1 - 1)

	var count int
	var pageCount1 int
	var articleswithtype []models.Article

	if typeName == "" {
		beego.Info("下拉框传递数据失败")

		count01, err := qs.Count()
		if err != nil {
			beego.Info("Count", err)
		}
		count = int(count01)

		pageCount := float64(count) / float64(pageSize)
		pageCount1 = int(math.Ceil(pageCount))

		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articleswithtype)

	} else {
		count01, err := qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
		count = int(count01)
		if err != nil {
			beego.Info("Count", err)
		}
		pageCount := float64(count) / float64(pageSize)
		pageCount1 = int(math.Ceil(pageCount))

		//qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articleswithtype)
		qs.Filter("ArticleType__TypeName", typeName).RelatedSel().Limit(pageSize, start).All(&articleswithtype)
	}
	FirstPage := false
	//首页末页数据处理
	if pageIndex1 == 1 {
		FirstPage = true
	}
	LastPage := false
	//首页末页数据处理
	if pageIndex1 == int(pageCount1) {
		LastPage = true
	}
	userName := this.GetSession("userName")

	this.Data["types"] = types
	this.Data["typeName"] = typeName
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount1
	this.Data["pageIndex"] = pageIndex1
	this.Data["FirstPage"] = FirstPage
	this.Data["LastPage"] = LastPage
	this.Data["articles"] = articleswithtype
	this.Data["userName"] = userName

	this.TplName = "index.html"

}

func (this *ArticleController) ShowContent() {

	id := this.GetString(":id")
	id01, err := strconv.Atoi(id)
	if err != nil {
		beego.Info("转换错误", err)
	}
	o := orm.NewOrm()
	article := models.Article{Id: id01}
	o.Read(&article)
	article.Acount += 1
	o.Update(&article)

	m2m := o.QueryM2M(&article, "Users")
	userName := this.GetSession("userName")
	user := models.User{}
	user.UserName = userName.(string)
	o.Read(&user, "UserName")

	_, err = m2m.Add(&user)
	if err != nil {
		beego.Info("插入失败")
		return
	}
	o.Update(&article)

	users := []models.User{}
	o.QueryTable("User").Filter("Articles__Article__Id", id01).Distinct().All(&users)

	this.Data["userName"] = userName
	this.Data["users"] = users
	this.Data["article"] = article
	this.TplName = "content.html"
}

func (this *ArticleController) DeleteArticle() {
	id := this.GetString(":id")
	id01, err := strconv.Atoi(id)
	if err != nil {
		beego.Info("转换错误", err)
	}
	o := orm.NewOrm()
	article := models.Article{Id: id01}
	o.Delete(&article)
	this.Redirect("/ShowArticle", 302)
}
func (this *ArticleController) ShowUpdate() {
	id := this.GetString(":id")
	id01, err := strconv.Atoi(id)
	if err != nil {
		beego.Info("转换错误", err)
	}
	o := orm.NewOrm()
	article := models.Article{Id: id01}
	o.Read(&article)
	userName := this.GetSession("userName")
	this.Data["userName"] = userName
	this.Data["update"] = article
	this.TplName = "update.html"
}

func (this *ArticleController) UpdateArticle() {
	articleName := this.GetString("articleName")
	if articleName == "" {
		beego.Info("没取到数据")
	}
	content := this.GetString("content")
	f, h, err8 := this.GetFile("uploadname")
	if err8 != nil {
		beego.Info("文件上传失败")
	}
	defer f.Close()

	o := orm.NewOrm()
	id := this.GetString(":id")
	id01, err := strconv.Atoi(id)
	if err != nil {
		beego.Info("转换错误", err)
	}
	article := models.Article{Id: id01}
	o.Read(&article)
	if articleName != "" {
		article.ArtiName = articleName
	}
	if content != "" {
		article.Acontent = content
	}
	if h.Filename != "" {

		beego.Info("sjfsgi")
		ext := path.Ext(h.Filename)
		if ext != ".jpg" && ext != ".png" {
			beego.Info("上传文件格式不正确")
			return
		}

		if h.Size > 5000000 {
			beego.Info("文件太大，不允许上传")
			return
		}
		pathfile := time.Now().Format("2006-01-02 15-04-05")

		err01 := this.SaveToFile("uploadname", "/static/img/"+pathfile+ext)
		if err01 != nil {
			beego.Info("存储失败", err01)
			return
		}
		article.Aimg = "/static/img/" + pathfile + ext
	}
	o.Update(&article)
	this.Redirect("/Article/ShowArticle", 302)

}

func (this *ArticleController) Logout() {
	this.DelSession("userName")
	this.Redirect("/", 302)
}

func (this *ArticleController) DeleteType() {
	id := this.GetString("id")
	id2, _ := strconv.Atoi(id)
	if id2 == 0 {
		beego.Info("获取id错误")
		return
	}

	o := orm.NewOrm()
	artiType := models.ArticleType{Id: id2}
	o.Delete(&artiType)

	this.Redirect("/Article/AddType", 302)
}

type AddArticleController struct {
	beego.Controller
}

func (this *AddArticleController) ShowAddArticle() {
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	userName := this.GetSession("userName")
	this.Data["userName"] = userName
	this.Data["types"] = types
	this.TplName = "add.html"
}
func (this *AddArticleController) HandleArticle() {
	articleName := this.GetString("articleName")
	beego.Info(articleName)
	content := this.GetString("content")
	f, h, err := this.GetFile("uploadname")
	if err != nil {
		beego.Info("文件上传失败")
	}
	defer f.Close()
	ext := path.Ext(h.Filename)
	if ext != ".jpg" && ext != ".png" {
		beego.Info("上传文件格式不正确")
		return
	}

	if h.Size > 5000000 {
		beego.Info("文件太大，不允许上传")
		return
	}
	pathfile := time.Now().Format("2006-01-02-15-04-05")

	err01 := this.SaveToFile("uploadname", "./static/img/"+pathfile+ext)
	if err01 != nil {
		beego.Info("存储失败", err01)
		return
	}

	o := orm.NewOrm()
	article := models.Article{}
	article.ArtiName = articleName
	article.Acontent = content
	article.Aimg = "/static/img/" + pathfile + ext

	typeName := this.GetString("select")
	//类型判断
	if typeName == "" {
		beego.Info("下拉匡数据错误")
		return
	}
	//获取type对象
	var artiType models.ArticleType
	artiType.TypeName = typeName
	err = o.Read(&artiType, "TypeName")
	if err != nil {
		beego.Info("获取类型错误")
		return
	}
	article.ArticleType = &artiType

	_, err02 := o.Insert(&article)
	if err02 != nil {
		beego.Info("Inserterr", err02)
		return
	}
	beego.Info("储存成功")
	this.Redirect("/Article/ShowArticle", 302)
}
func (this *AddArticleController) ShowAddType() {
	o := orm.NewOrm()
	var artiTypes []models.ArticleType
	//查询
	_, err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil {
		beego.Info("查询类型错误")
	}
	userName := this.GetSession("userName")
	this.Data["userName"] = userName
	this.Data["types"] = artiTypes
	this.TplName = "addType.html"
}
func (this *AddArticleController) HandleAddType() {
	//1.获取数据
	typename := this.GetString("typeName")
	//2.判断数据
	if typename == "" {
		beego.Info("添加类型数据为空")
		this.TplName = "addType"
	}
	//3.执行插入操作
	o := orm.NewOrm()
	var artiType models.ArticleType
	artiType.TypeName = typename
	_, err := o.Insert(&artiType)
	if err != nil {
		beego.Info("插入失败")
		return
	}
	//4.展示视图？
	this.Redirect("/Article/AddType", 302)
}
