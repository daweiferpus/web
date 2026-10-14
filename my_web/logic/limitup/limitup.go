package limitup

import (
	"encoding/json"
	"fmt"
	"net/http"
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

/*
记录昨日涨停的二板及二板以上的相关行情信息
*/
type stockData struct {
	StockCode string   `json:"stock_code"`
	Fields    []string `json:"fields"`
}

var (
	curDay    string
	responses limitFormat
	codes     string
	hqChan    chan stockData
	errLog    log.LogInter
)

func init() {
	responses = limitFormat{}
	hqChan = make(chan stockData)
	errLog = log.GetLogInter()
	errLog.NewFileLogger("../log/error.log")
}

func changeStkCode(code string) string {
	if strings.HasPrefix(code, "6") {
		return "sh" + code
	}
	return "sz" + code
}

func MakeRouter() {
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

	router.GET("/yesterday/limit", func(ctx *gin.Context) {
		stkdatas := make([]limitFormat, 0)
		today := time.Now().Format("2006-01-02")
		if today != curDay {
			redisInter := redis.GetRedisInter()
			jsonStr := ""
			tmp := time.Now().AddDate(0, 0, -1)
			for jsonStr == "" {
				jsonStr, _ = redisInter.GetHashValue("stkdata", tmp.Format("2006-01-02"))
				tmp = tmp.AddDate(0, 0, -1)
			}

			json.Unmarshal([]byte(jsonStr), &stkdatas)
			strBuild := strings.Builder{}
			sort.Slice(stkdatas, func(i, j int) bool {
				return stkdatas[i].Num > stkdatas[j].Num
			})
			for _, v := range stkdatas {
				if v.Num > 1 || strings.HasPrefix(v.Code, "3") {
					strBuild.WriteString(changeStkCode(v.Code))
					strBuild.WriteByte(',')
				}
			}
			codes = strBuild.String()
			if len(codes) > 0 {
				codes = codes[:len(codes)-1]
			}
		}

		hqChan = make(chan stockData)
		url := fmt.Sprintf(public.SomeStkUrl, codes)
		// 创建回调函数
		callBack := func(url string, data []byte) {
			results := parseStockData(data)
			for _, v := range results {
				hqChan <- v
			}
			close(hqChan)
		}
		visit := visitsina.GetVisitInter()
		visit.InitSinaQuote(callBack, errLog)
		go visit.VisitSina(url)
		ctx.JSON(http.StatusOK, response{
			Code:    0,
			Message: "ok",
			Data:    GetYesterDayStk(stkdatas),
		})
	})

	router.Run(":8002")
}

/*
GetYesterDayStk http://qt.gtimg.cn/q=sz000037,sh605162
	response：

		第2位 ：证券名称
		第3位 ：证券代码
		第4位 ：最新价
		第33位：涨跌幅
		第38位：成交量 单位 万元
		第39位：换手率 单位 %
		第45位：流通市值 单位 /亿
*/
func GetYesterDayStk(stkdatas []limitFormat) (responses []limitFormat) {
	redisInter := redis.GetRedisInter()
MainLoop:
	for {
		select {
		case v, ok := <-hqChan:
			if !ok {
				break MainLoop
			}
			var response = limitFormat{}
			var reason = stkReason{}
			jsonStr, _ := redisInter.GetHashValue("stkreason", v.StockCode)
			json.Unmarshal([]byte(jsonStr), &reason)

			response.Reason = reason.Reason
			response.Limit = reason.Limit
			if len(v.Fields) > 44 {
				response.Code = v.Fields[2]
				response.Name = v.Fields[1]
				response.Percent = fmt.Sprintf("%s%%", v.Fields[32])       // 第33位：涨跌幅
				response.Trade = v.Fields[3]                               // 第4位：最新价
				response.Nmc = fmt.Sprintf("%s亿", v.Fields[44])            // 第45位：流通市值 单位 /亿
				response.Turnoverratio = fmt.Sprintf("%s%%", v.Fields[38]) // 第39位：换手率 单位 %
				// 第38位：成交量 单位 万元
				amount, _ := strconv.ParseFloat(v.Fields[37], 32)
				response.Amount = fmt.Sprintf("%.2f亿", amount/10000.0)
				for _, tmp := range stkdatas {
					if tmp.Code == response.Code {
						response.Num = tmp.Num
						break
					}
				}
			}
			responses = append(responses, response)
		case <-time.After(3 * time.Second):
			break MainLoop
		}
	}
	return
}

func parseStockData(data []byte) []stockData {
	// 转换为字符串并按行分割
	content := string(data)
	lines := strings.Split(content, "\n")

	var stocks []stockData

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 分割股票代码和数据部分
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}

		// 提取股票代码（去掉前缀"v_"）
		stockCode := strings.TrimPrefix(parts[0], "v_")

		// 处理数据部分：去掉首尾的引号和分号
		dataPart := strings.Trim(parts[1], `";`)

		// 分割字段
		fields := strings.Split(dataPart, "~")

		// 存入结构体
		stocks = append(stocks, stockData{
			StockCode: stockCode,
			Fields:    fields,
		})
	}

	return stocks
}
