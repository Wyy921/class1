package models

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       int
	UserName string
	Password string
	Articles  []*Article `orm:"rel(m2m)"`
}
type Article struct {
	Id          int          `orm:"pk;auto"`
	ArtiName    string       `orm:"size(20)"`
	Atime       time.Time    `orm:"auto_now_add"`
	Acount      int          `orm:"default(0);null"`
	Acontent    string       `orm:"size(500)"`
	Aimg        string       `orm:"size(100)"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users        []*User      `orm:"reverse(many)"`
}
type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Article  []*Article `orm:"reverse(many)"`
}

func init() {
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(192.168.223.149:3306)/class1?charset=utf8&loc=Asia%2FShanghai", 30)
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	orm.RunSyncdb("default", false, true)
}
