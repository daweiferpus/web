package active

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"localhost.com/my_web/repo/redis"
)

type response struct {
	Code    int         `json:"code"`    // 业务状态码（0表示成功）
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据（成功时为具体数据，失败时可忽略）
}

// 一只股票的信息
type oneStk struct {
	Code    string `json:"Code"` // 股票代码
	Name    string `json:"Name"` // 股票名称
	Reason  string `json:"-"`    // 涨停理由
	Money   string `json:"-"`    // 封单金额
	Perform string `json:"-"`    // T+日表现  < -4：大面； < -4且小于0：小面； > 0 且没涨停红盘； 涨停：晋级
	PreRate string `json:"-"`    // 特殊处理，第一条数据记录板块表现 例："共4支 晋级率13.33%"
}

type paihangData struct {
	Code   string `json:"Code"` // 股票代码
	Name   string `json:"Name"` // 股票名称
	Num    int    `json:"Num"`  // 涨停次数
	Reason string `json:Reason` // 涨停原因
}

func GetpaihangData(stkNum int) (result []paihangData) {
	redisInter := redis.GetRedisInter()
	curDay := time.Now()
	dayNum := 0
	mapCnt := make(map[string]int)
	for {
		// 计算前一天的时间
		stkdata, err := redisInter.GetHashValue("stkdata", curDay.Format("2006-01-02"))
		if err != nil || stkdata == "" {
			curDay = curDay.AddDate(0, 0, -1)
			continue
		}
		dayNum++
		arrStk := make([]oneStk, 0)
		json.Unmarshal([]byte(stkdata), &arrStk)
		for _, v := range arrStk {
			var builder = strings.Builder{}
			builder.WriteString(v.Code)
			builder.WriteByte('|')
			builder.WriteString(v.Name)
			mapCnt[builder.String()]++
		}
		// 统计50天的数据
		if dayNum >= 50 || curDay.Year() <= 2024 {
			break
		}
		curDay = curDay.AddDate(0, 0, -1)
	}
	for k, v := range mapCnt {
		arr := strings.Split(k, "|")
		if len(arr) > 1 {
			reason, _ := redisInter.GetHashValue("stkinfo", arr[1])
			result = append(result, paihangData{arr[0], arr[1], v, reason})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Num > result[j].Num
	})
	return result[:stkNum]
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

	router.GET("/active/limit/:num", func(ctx *gin.Context) {
		stknum := ctx.Param("num")
		stknumInt, _ := strconv.Atoi(stknum)
		ctx.JSON(http.StatusOK, response{
			Code:    0,
			Message: "ok",
			Data:    GetpaihangData(stknumInt),
		})
	})

	router.Run(":8001")
}
