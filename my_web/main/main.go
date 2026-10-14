package main

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"localhost.com/my_web/entity/public"
	"localhost.com/my_web/logic/active"
	"localhost.com/my_web/logic/limitup"
	"localhost.com/my_web/logic/task"
	"localhost.com/my_web/logic/webhq"
	"localhost.com/my_web/repo/log"
	"localhost.com/my_web/repo/redis"
	"localhost.com/my_web/repo/visitsina"

	"net/http"
	_ "net/http/pprof" // 自动注册pprof的handler
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	runtime.GOMAXPROCS(runtime.NumCPU())

	defer cancel()
	// 启动一个goroutine运行性能分析服务器
	go func() {
		errLog.WriteInfo(http.ListenAndServe("localhost:6060", nil))
	}()

	// 行情信息
	go webhq.HqmoniterHttp()

	// 统计活跃股
	go active.MakeRouter()

	go limitup.MakeRouter()

	go task.AutoTask(ctx)

	for {
		time.Sleep(time.Second * 10)
	}
	// makeredis()
}

var (
	pageCount int = 0
	hqChan        = make(chan public.SinaQuoteData)
	errLog    log.LogInter
)

func init() {
	errLog = log.GetLogInter()
	errLog.NewFileLogger("../log/error.log")
}

func getAllHqInfo(url string, data []byte) {
	if pageCount > 0 {
		visit := visitsina.GetVisitInter()
		pageIndex, _ := strconv.Atoi(url[strings.Index(url, "page=")+5 : strings.Index(url, "&num=")])

		quotes := make([]public.SinaQuoteData, 0)
		json.Unmarshal(data, &quotes)
		for _, v := range quotes {
			hqChan <- v
		}

		if pageIndex < pageCount {
			pageURL := fmt.Sprintf(public.PageUrl, pageIndex+1, "100")
			visit.InitSinaQuote(getAllHqInfo, errLog)
			visit.VisitSina(pageURL)
		} else {
			close(hqChan)
		}
	}
}

type stkReason struct {
	Symbol string `json:"symbol"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
	Limit  string `json:"limit"`
	Nmc    string `json:"nmc"`
}

func makeredis() {
	var errLog log.LogInter

	redisInter := redis.GetRedisInter()
	mStkinfo := redisInter.GetHashAllValue("stkinfo", false)
	mStklimit := redisInter.GetHashAllValue("stklimitup", false)

	visit := visitsina.GetVisitInter()

	// 创建回调函数
	callBack := func(url string, data []byte) {
		if url == public.CountUrl {
			str := string(data)
			str = str[1 : len(str)-1]
			pageCount, _ = strconv.Atoi(str)
			pageCount = (pageCount + 100 - 1) / 100

			pageURL := fmt.Sprintf(public.PageUrl, 1, "100") // 查第一页的行情
			visit.InitSinaQuote(getAllHqInfo, errLog)
			visit.VisitSina(pageURL)
		}
	}

	visit.InitSinaQuote(callBack, errLog)
	go visit.VisitSina(public.CountUrl)

	mapRedis := make(map[string]interface{})
	for v := range hqChan {
		if _, ok := mStkinfo[v.Name]; ok {
			if _, ok := mStklimit[v.Name]; ok {
				redisInfo := stkReason{}
				redisInfo.Symbol = v.Symbol
				redisInfo.Code = v.Code
				redisInfo.Name = v.Name
				redisInfo.Nmc = fmt.Sprintf("%.2f亿", v.Nmc/10000.0)
				redisInfo.Limit = mStklimit[v.Name]
				redisInfo.Reason = mStkinfo[v.Name]
				value, _ := json.Marshal(redisInfo)
				mapRedis[v.Symbol] = value
			}
		}
	}
	redisInter.SetHashInterface(mapRedis, "stkreason")
}
