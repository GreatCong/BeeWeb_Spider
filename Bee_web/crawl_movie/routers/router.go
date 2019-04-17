package routers

import (
	"../controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})                                 //默认路由
	beego.Router("/crawlmovie", &controllers.CrawlMovieController{}, "*:CrawlMovie") //新添加一个路由（CrawlMovie表示method）
	//*:表示所有回调都会调用,也可以改成get
}
