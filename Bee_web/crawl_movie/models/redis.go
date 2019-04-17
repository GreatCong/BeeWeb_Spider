package models

/*
我们需要在程序里面弄清楚，那些Url已经访问，那些没有访问，这个可以用redis实现
用set可以避免重复的元素
*/

import (
	"github.com/astaxie/goredis" //第三方redis包
)

const ( //这样类似于define 或者 enum
	URL_QUEUE     = "url_queue"
	URL_VISIT_SET = "url_visit_set"
)

var (
	client goredis.Client
)

//链接redis
func ConnectRedis(addr string) {
	client.Addr = addr
}

// 存入到queue
func PutinQueue(url string) {
	client.Lpush(URL_QUEUE, []byte(url))
}

//从queue中获取
func PopfromQueue() string {
	res, err := client.Rpop(URL_QUEUE)
	if err != nil {
		panic(err) //假如函数F中书写了panic语句，会终止其后要执行的代码，在panic所在函数F内如果存在要执行的defer函数列表，按照defer的逆序执行
	}

	return string(res)
}

//获取队列的长度
func GetQueueLength() int {
	length, err := client.Llen(URL_QUEUE)
	if err != nil {
		return 0
	}

	return length
}

//添加到已经访问过的set
func AddToSet(url string) {
	client.Sadd(URL_VISIT_SET, []byte(url))
}

//是否已经被访问过
func IsVisit(url string) bool {
	bIsVisit, err := client.Sismember(URL_VISIT_SET, []byte(url))
	if err != nil {
		return false
	}

	return bIsVisit
}
