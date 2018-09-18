package main

import (
	_ "class1/routers"
	_ "class1/models"
	"github.com/astaxie/beego"
)

func main() {

	beego.AddFuncMap("ShowPrePage",HandlePrePage)
	beego.AddFuncMap("ShowNextPage",HandleNextPage)

	beego.Run()


}
func HandlePrePage(data int)(int){

	pageIndex := data - 1
	return pageIndex
}

func HandleNextPage(data int)(int){
	pageIndex := data + 1
	return pageIndex

}