package visitsina

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"localhost.com/my_web/entity/public"
	"localhost.com/my_web/repo/log"
)

type VisitInter interface {
	VisitSina(string)
	InitSinaQuote(func(string, []byte), log.LogInter) *colly.Collector
}

type visitReal struct {
	sinaClct *colly.Collector
}

func GetVisitInter() *visitReal {
	return &visitReal{}
}

// sinaQuoteDatas = make([]SinaQuoteData, 0)
// 使用colly库创建并配置一个网页爬虫。Colly 是一个强大的Go语言网络爬虫框架，用于抓取网页内容
func (vr *visitReal) InitSinaQuote(callbackFun func(string, []byte), errlog log.LogInter) *colly.Collector {
	sinaClct := colly.NewCollector(func(collector *colly.Collector) {
		// colly.AllowURLRevisit()
		// 随机设置用户代理(User-Agent)。用户代理字符串通常被浏览器用来向服务器标识自己，
		// 通过随机化用户代理，可以帮助爬虫避免因使用相同的User-Agent而被网站识别和屏蔽。
		extensions.RandomUserAgent(collector)
		colly.AllowURLRevisit()
	})
	// 设置HTTP客户端超时时间为3秒
	sinaClct.SetRequestTimeout(3 * time.Second)

	cookiestr := fmt.Sprintf("MONEY-FINANCE-SINA-COM-CN-WEB5=; UOR=,vip.stock.finance.sina.com.cn,; ULV=%d:1:1:1::", time.Now().UnixMilli())

	// 注册一个或多个函数，这些函数会在每个请求发出之前被调用
	// 这对于在发送HTTP请求前修改请求参数、记录日志或执行其他预处理操作非常有用。
	sinaClct.OnRequest(func(r *colly.Request) {
		if r.URL.String() == public.IndexUrl {
			// 设置HTTP请求头中的某个字段值
			r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
			r.Headers.Set("accept-encoding", "gzip, deflate, br")
			r.Headers.Set("accept-language", "zh-CN,zh;q=0.9")
			r.Headers.Set("sec-ch-ua", "\"Chromium\";v=\"112\", \"Google Chrome\";v=\"112\", \"Not:A-Brand\";v=\"99\"")
			r.Headers.Set("sec-ch-ua-mobile", "?0")
			r.Headers.Set("sec-ch-ua-platform", "\"Windows\"")
			r.Headers.Set("sec-fetch-dest", "document")
			r.Headers.Set("sec-fetch-mode", "navigate")
			r.Headers.Set("sec-fetch-site", "none")
			r.Headers.Set("sec-fetch-user", "?1")
			r.Headers.Set("upgrade-insecure-requests", "1")
			r.Headers.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
		} else {
			r.Headers.Set("accept", "*/*")
			r.Headers.Set("content-type", "application/x-www-form-urlencoded")
			r.Headers.Set("cookie", cookiestr)
			r.Headers.Set("referer", "https://vip.stock.finance.sina.com.cn/mkt/")
			r.Headers.Set("sec-fetch-dest", "empty")
			r.Headers.Set("sec-fetch-mode", "cors")
			r.Headers.Set("sec-fetch-site", "same-origin")
			r.Headers.Del("sec-fetch-user")
			r.Headers.Del("upgrade-insecure-requests")
		}
	})

	// 注册一个函数，该函数会在每个HTTP响应接收后被调用。
	// 定义一个处理响应的回调函数，该函数会在Colly成功从web服务器接收到响应数据后执行。
	sinaClct.OnResponse(func(r *colly.Response) {
		callbackFun(r.Request.URL.String(), r.Body)
	})

	// 错误处理
	sinaClct.OnError(func(r *colly.Response, err error) {
		errlog.WriteError(fmt.Sprintf("sinaClct.OnResponse err: %s", err))
	})
	vr.sinaClct = sinaClct
	return sinaClct
}
func (vr *visitReal) VisitSina(indexURL string) {
	if vr.sinaClct == nil {
		fmt.Println("visit is not init")
		return
	}
	vr.sinaClct.Visit(indexURL)
}
