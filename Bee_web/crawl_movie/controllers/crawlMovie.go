package controllers

import (
	"../models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/opentracing/opentracing-go/log" //OpenTracing是一个跨编程语言的标准，这是第三方包
	"strings"
	"time"
)

type CrawlMovieController struct {
	beego.Controller
}

/*
func (c *CrawlMovieController) CrawlMovie() {
	//sUrl := "https://movie.douban.com/subject/25827935/" //七月与安生
	sUrl := "https://movie.douban.com/subject/26752088/" //我不是药神
	rsp := httplib.Get(sUrl)
	sMovieHtml, err := rsp.String()
	if err != nil {
		log.Error(err)
	}

	// c.Ctx.WriteString(models.GetMovieDirector(sMovieHtml) + "|")
	// c.Ctx.WriteString(models.GetMovieName(sMovieHtml) + "|")
	// c.Ctx.WriteString(models.GetMovieMainCharacters(sMovieHtml) + "|")
	// c.Ctx.WriteString(models.GetMovieGrade(sMovieHtml) + "|")
	// c.Ctx.WriteString(models.GetMovieGenre(sMovieHtml) + "|")
	// c.Ctx.WriteString(models.GetMovieOnTime(sMovieHtml) + "|")
	// c.Ctx.WriteString(models.GetMovieRunningTime(sMovieHtml))


		// var movieInfo models.MovieInfo

		// movieInfo.Movie_name = models.GetMovieName(sMovieHtml)
		// movieInfo.Movie_director = models.GetMovieDirector(sMovieHtml)
		// movieInfo.Movie_main_character = models.GetMovieMainCharacters(sMovieHtml)
		// movieInfo.Movie_type = models.GetMovieGenre(sMovieHtml)
		// movieInfo.Movie_on_time = models.GetMovieOnTime(sMovieHtml) // 上映时间：2016-09-14(中国大陆)
		// movieInfo.Movie_grade = models.GetMovieGrade(sMovieHtml)
		// movieInfo.Movie_span = models.GetMovieRunningTime(sMovieHtml)

		// movieInfo.Movie_on_time = movieInfo.Movie_on_time[0:strings.Index(models.GetMovieOnTime(sMovieHtml), "(")] //上映时间：2016-09-14
		// movieInfo.Create_time = time.Now().Format("2006-1-2 15:04:05")

		// id, _ := models.AddMovie(&movieInfo)

		// c.Ctx.WriteString(fmt.Sprintf("%v", id))


	c.Ctx.WriteString(fmt.Sprintf("%v", models.GetMovieUrls(sMovieHtml))) //打印获取的页面的超链接
}
*/

/*
主要思路：首先链接到redis上，然后将要爬的页面作为整个程序的主界面，所以将它加入到要爬取数据的队列中，
          然后在for循环中，每次都需要先判断队列的长度，如果为0表示已经爬取完所有的数据，否则获取队列中的元素，
          然后判断它是否已经被访问过，也就是爬取过。
          如果没有爬取过，那么我们就可以爬取数据了，爬取数据这里是之前已经实现过了，稍微修改一下就可以了，然后提取该页面的所有的子链接，
          加入到要爬取的队列中，并将当前页面添加到set中，表示已经被爬取过了。
*/
func (c *CrawlMovieController) CrawlMovie() { //添加了对redis的处理

	var movieInfo models.MovieInfo

	//连接到redis
	models.ConnectRedis("127.0.0.1:6379") //redis 的默认地址

	sUrl := "https://movie.douban.com/subject/25827935/" //七月与安生
	//sUrl := "https://movie.douban.com/subject/26752088/" //我不是药神

	//先添加到队列中
	models.PutinQueue(sUrl)

	for {
		length := models.GetQueueLength() //获取队列长度

		if length == 0 {
			break //如果url队列为空 则退出当前循环
		}

		sUrl = models.PopfromQueue()
		//我们应当判断sUrl是否应该被访问过
		if models.IsVisit(sUrl) {
			continue
		}
		rsp := httplib.Get(sUrl)
		//设置User-agent以及cookie是为了防止  豆瓣网的 403 !!!!!
		rsp.Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:50.0) Gecko/20100101 Firefox/50.0")
		rsp.Header("Cookie", `bid=gFP9qSgGTfA; __utma=30149280.1124851270.1482153600.1483055851.1483064193.8; __utmz=30149280.1482971588.4.2.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; ll="118221"; _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1483064193%2C%22https%3A%2F%2Fwww.douban.com%2F%22%5D; _pk_id.100001.4cf6=5afcf5e5496eab22.1482413017.7.1483066280.1483057909.; __utma=223695111.1636117731.1482413017.1483055857.1483064193.7; __utmz=223695111.1483055857.6.5.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; _vwo_uuid_v2=BDC2DBEDF8958EC838F9D9394CC5D9A0|2cc6ef7952be8c2d5408cb7c8cce2684; ap=1; viewed="1006073"; gr_user_id=e5c932fc-2af6-4861-8a4f-5d696f34570b; __utmc=30149280; __utmc=223695111; _pk_ses.100001.4cf6=*; __utmb=30149280.0.10.1483064193; __utmb=223695111.0.10.1483064193`)

		sMovieHtml, err := rsp.String()
		if err != nil {
			log.Error(err)
		}

		movieInfo.Movie_name = models.GetMovieName(sMovieHtml) //获取电影的名字
		if movieInfo.Movie_name != "" {

			movieInfo.Movie_director = models.GetMovieDirector(sMovieHtml)
			movieInfo.Movie_main_character = models.GetMovieMainCharacters(sMovieHtml)
			movieInfo.Movie_type = models.GetMovieGenre(sMovieHtml)
			movieInfo.Movie_on_time = models.GetMovieOnTime(sMovieHtml) // 上映时间：2016-09-14(中国大陆)
			movieInfo.Movie_grade = models.GetMovieGrade(sMovieHtml)
			movieInfo.Movie_span = models.GetMovieRunningTime(sMovieHtml)

			//找到第一次出现括号的索引 也就是把后面(中国大陆)字眼去除
			movieInfo.Movie_on_time = movieInfo.Movie_on_time[0:strings.Index(models.GetMovieOnTime(sMovieHtml), "(")] //上映时间：2016-09-14

			movieInfo.Create_time = time.Now().Format("2006-1-2 15:04:05")

			id, _ := models.AddMovie(&movieInfo) //添加到mySQL数据库中
			fmt.Println("id-->", id)

		}

		//提取该页面的所有连接
		urls := models.GetMovieUrls(sMovieHtml)

		for _, url := range urls {
			models.PutinQueue(url)
			c.Ctx.WriteString("<br>" + url + "</br>")
		}

		//sUrl 应当记录到 访问set中
		models.AddToSet(sUrl)

		time.Sleep(time.Second)
	}

	c.Ctx.WriteString("end of crawl!")

}