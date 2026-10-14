package webhq

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"localhost.com/my_web/entity/public"
	"localhost.com/my_web/repo/log"
	"localhost.com/my_web/repo/redis"
	"localhost.com/my_web/repo/visitsina"
)

var (
	hqChan  chan []public.SinaQuoteData
	bidChan chan []public.SinaQuoteData
	errLog  log.LogInter
)

func init() {
	errLog = log.GetLogInter()
	errLog.NewFileLogger("../log/error.log")

	hqChan = make(chan []public.SinaQuoteData)
	bidChan = make(chan []public.SinaQuoteData)
}

func setresponse(ctx context.Context, c *gin.Context) {
	// 在函数内部初始化redis
	redisClient := redis.GetRedisInter()
	select { // 140  140 130 80 90 90
	case data := <-hqChan:
		result := make([]hqFormat, 0, len(data))
		for _, v := range data {
			if strings.HasPrefix(v.Symbol, "bj") {
				continue
			}
			temp := hqFormat{}
			temp.Code = v.Code
			temp.Name = v.Name
			temp.Percent = fmt.Sprintf("%.2f%%", v.Changepercent)
			temp.Price = v.Trade[:len(v.Trade)-1]
			reason, _ := redisClient.GetHashValue("stkinfo", v.Name)
			temp.Reason = reason
			limit, _ := redisClient.GetHashValue("stklimitup", v.Name)
			temp.Limit = limit
			temp.Nmc = fmt.Sprintf("%.2f亿", v.Nmc/10000.0)
			temp.Amount = fmt.Sprintf("%.2f亿", float64(v.Amount)/10000.0/10000.0)
			temp.Turnoverratio = fmt.Sprintf("%.2f%%", v.Turnoverratio)
			result = append(result, temp)
		}
		c.JSON(http.StatusOK, response{
			Code:    0,
			Message: "ok",
			Data:    result,
		})
	case <-ctx.Done():
		c.JSON(http.StatusOK, response{
			Code:    -1,
			Message: "time out",
			Data:    nil,
		})
		break
	}
}

func GethqInfo(url string, dealtime time.Time) {
	visit := visitsina.GetVisitInter()
	dealline, cancel := context.WithDeadline(context.Background(), dealtime)
	defer cancel()
	// 创建回调函数
	callBack := func(url string, data []byte) {
		quotes := make([]public.SinaQuoteData, 0)
		json.Unmarshal(data, &quotes)
		bidChan <- quotes
	}
	mapSinadata := make(map[string][]public.SinaQuoteData)

	for {
		visit.InitSinaQuote(callBack, errLog)
		go visit.VisitSina(url)
		select {
		case hq := <-bidChan:
			for _, v := range hq {
				mapSinadata[v.Symbol] = append(mapSinadata[v.Symbol], v)
			}
		case <-dealline.Done():
			MoniterBidInfo(mapSinadata)
			return
		case <-time.After(3 * time.Second):
		}
		time.Sleep(3 * time.Second)
	}
}

func HqmoniterHttp() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	gRouter := router.Group("/monitor/hq")

	// 涨幅榜行情
	gRouter.GET("/increase/:num", func(c *gin.Context) {
		stknum := c.Param("num")
		visit := visitsina.GetVisitInter()
		url := fmt.Sprintf(public.IncPrcUrl, stknum)
		// 创建回调函数
		callBack := func(url string, data []byte) {
			quotes := make([]public.SinaQuoteData, 0)
			json.Unmarshal(data, &quotes)
			hqChan <- quotes
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		// 发起请求
		visit.InitSinaQuote(callBack, errLog)
		go visit.VisitSina(url)
		setresponse(ctx, c)
	})
	// 集合竞价行情，按涨幅来分析
	gRouter.GET("/bid/:num", func(c *gin.Context) {
		dealtime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 25, 10, 0, time.Local)
		beforetime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 15, 0, 0, time.Local)
		redisClient := redis.GetRedisInter()
		jsonBytes, _ := redisClient.GetHashValue("stkbidinfo", dealtime.Format("2006-01-02"))
		if time.Now().After(dealtime) || time.Now().Before(beforetime) { // 不在9.15~9.25之间直接查数据库里面的数据
			for i := 0; jsonBytes == "" && i < 15; i++ {
				dealtime = dealtime.AddDate(0, 0, -1)
				jsonBytes, _ = redisClient.GetHashValue("stkbidinfo", dealtime.Format("2006-01-02"))
			}
			curDate := dealtime.Format("2006-01-02")
			// 取前一日的成交量
			yesterBytes := ""
			for i := 0; yesterBytes == "" && i < 15; i++ { // 最多往前查15天
				dealtime = dealtime.AddDate(0, 0, -1)
				yesterBytes, _ = redisClient.GetHashValue("stkbidinfo", dealtime.Format("2006-01-02"))
			}
			// 取前一日的成交量 -- end
			bidHq := bidredis{}
			yesterHq := bidredis{}
			json.Unmarshal([]byte(jsonBytes), &bidHq)
			json.Unmarshal([]byte(yesterBytes), &yesterHq)
			bidHq.Date = curDate
			bidHq.YesterAmount = yesterHq.TotalAmount
			c.JSON(http.StatusOK, response{
				Code:    0,
				Message: "ok",
				Data:    bidHq,
			})
			return
		}

		stknum := c.Param("num")
		url := fmt.Sprintf(public.IncPrcUrl, stknum)
		if jsonBytes == "" { // 去查询竞价行情信息
			curtime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 25, 10, 0, time.Local)
			go GethqInfo(url, curtime)
		}
		c.JSON(http.StatusOK, response{
			Code:    0,
			Message: "ok",
			Data:    bidredis{},
		})
	})
	// 启动监听
	router.Run(":8000")
}

// SinaQuoteDatas slice排序
type SinaQuoteDatas []public.SinaQuoteData

func (p SinaQuoteDatas) Len() int {
	return len(p)
}

func (p SinaQuoteDatas) Less(i, j int) bool {
	if strings.Compare(p[i].Ticktime, p[j].Ticktime) < 0 {
		return true
	} else {
		return false
	}
}

func (p SinaQuoteDatas) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func setoneredisdata(temp *bidFormat, sinInfo *public.SinaQuoteData, redisClient redis.RedisInter) {
	reason, _ := redisClient.GetHashValue("stkinfo", sinInfo.Name)
	limitNum, _ := redisClient.GetHashValue("stklimitup", sinInfo.Name)

	temp.Code = sinInfo.Code
	temp.Name = sinInfo.Name
	temp.Percent = fmt.Sprintf("%.2f%%", sinInfo.Changepercent)
	temp.Open = sinInfo.Open[:len(sinInfo.Open)-1]
	temp.Reason = reason
	temp.Limit = limitNum
	temp.Nmc = fmt.Sprintf("%.2f亿", sinInfo.Nmc/10000.0)
	temp.Amount = fmt.Sprintf("%d万", sinInfo.Amount/10000)
}

// MoniterBidInfo 分析集合竞价数据
func MoniterBidInfo(mapSinadata map[string][]public.SinaQuoteData) {
	var totalAmt int64 = 0

	redisClient := redis.GetRedisInter()
	// result := make([]bidFormat, 0, len(mapSinadata))
	result1 := make([]bidFormat, 0, len(mapSinadata))
	result2 := make([]bidFormat, 0, len(mapSinadata))
	result3 := make([]bidFormat, 0, len(mapSinadata))
	for _, v := range mapSinadata {
		// 排序
		sort.Sort(SinaQuoteDatas(v))                      // 对map里面排序
		var lastSina public.SinaQuoteData                 // 上一条行情
		var sliceMoniter = make([]string, 0)              // 用来判断重复的
		var sliceMoniterTmp = make([]string, 0)           // 用来判断重复的
		lastTime, _ := time.Parse("15:04:05", "08:30:00") // 上一条行情的时间

		for i, sinInfo := range v {
			if i == len(v)-1 { // 只取最后一条9:25分的数据
				totalAmt += int64(sinInfo.Amount)
			}
			hqTime, _ := time.Parse("15:04:05", sinInfo.Ticktime)
			duration := hqTime.Sub(lastTime)

			// 两秒内涨幅波动超过3%
			if duration <= 2*time.Second && duration != 0 {
				if sinInfo.Changepercent-lastSina.Changepercent-3.00 > 0.00001 {
					if !slices.Contains(sliceMoniter, sinInfo.Code) {
						var temp = bidFormat{}
						temp.Note = "波动较大"
						setoneredisdata(&temp, &sinInfo, redisClient)
						result1 = append(result1, temp)
						sliceMoniter = append(sliceMoniter, sinInfo.Code)
					}
				}
			}
			// 对最后一条9:25:00的集合竞价数据做分析
			if strings.Compare(sinInfo.Ticktime, "09:25:00") >= 0 {
				// 获取竞价涨停的信息
				if strings.HasPrefix(strings.ToUpper(sinInfo.Name), "N") || strings.HasPrefix(strings.ToUpper(sinInfo.Name), "*ST") ||
					strings.HasPrefix(strings.ToUpper(sinInfo.Name), "ST") || strings.HasPrefix(strings.ToUpper(sinInfo.Name), "C") ||
					strings.HasPrefix(strings.ToUpper(sinInfo.Name), "BJ") {
					continue
				}
				// 昨收盘
				price_sett, _ := strconv.ParseFloat(sinInfo.Settlement, 64)
				// 市场价
				price_trade, _ := strconv.ParseFloat(sinInfo.Trade, 64)

				price_max := 0.00 // 涨停价
				// && !strings.HasPrefix(v.Symbol, "sh688")
				if strings.HasPrefix(sinInfo.Code, "30") || strings.HasPrefix(sinInfo.Code, "68") {
					// 涨停价
					price_max = math.Round(price_sett*120) / 100 // 四舍五入到两位小数
				} else {
					price_max = math.Round(price_sett*110) / 100 // 四舍五入到两位小数
				}
				if math.Abs(price_max-price_trade) <= 0.0001 && strings.Compare(sinInfo.Sell, "0.000") == 0 {
					if !slices.Contains(sliceMoniterTmp, sinInfo.Code) {
						var temp = bidFormat{}
						temp.Note = "竞价涨停"
						setoneredisdata(&temp, &sinInfo, redisClient)
						result2 = append(result2, temp)
						sliceMoniterTmp = append(sliceMoniterTmp, sinInfo.Code)
						continue
					}
				}

				// 成交量判断
				if !slices.Contains(sliceMoniterTmp, sinInfo.Code) && sinInfo.Amount-30000000.0 > 0.00 &&
					float64(sinInfo.Amount/10000)/sinInfo.Nmc > 0.0075 {
					var temp = bidFormat{}
					temp.Note = "成交量较大"
					setoneredisdata(&temp, &sinInfo, redisClient)
					result3 = append(result3, temp)
					sliceMoniterTmp = append(sliceMoniterTmp, sinInfo.Code)
				}
			}
			lastTime = hqTime
			lastSina = sinInfo
		}
	}
	if totalAmt == 0 {
		return
	}
	var bidhq = bidredis{}
	bidhq.TotalAmount = totalAmt
	bidhq.Bidhq = append(bidhq.Bidhq, result1...)
	bidhq.Bidhq = append(bidhq.Bidhq, result2...)
	bidhq.Bidhq = append(bidhq.Bidhq, result3...)
	jsonBytes, err := json.Marshal(bidhq)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	redisClient.SetOneHashData("stkbidinfo", time.Now().Format("2006-01-02"), jsonBytes)
}
