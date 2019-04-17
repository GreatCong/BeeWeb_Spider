package main

import (
	_ "./routers"
	"github.com/astaxie/beego"
)

/*
依赖于
"github.com/opentracing/opentracing-go/log" //OpenTracing是一个跨编程语言的标准，
"github.com/astaxie/goredis" //第三方redis包
"github.com/go-sql-driver/mysql" //mysql的第三方包
*/

func main() {
	beego.Run()
}
