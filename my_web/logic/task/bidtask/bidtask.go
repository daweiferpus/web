package bidtask

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"localhost.com/my_web/entity/public"
	"localhost.com/my_web/logic/webhq"
	"localhost.com/my_web/repo/redis"
)

// slice排序
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

// 集合竞价表单
type bidFormat struct {
	Note    string `json:"note"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Percent string `json:"percent"`
	Open    string `json:"open"`
	Reason  string `json:"reason"`
	Limit   string `json:"limit"`
	Nmc     string `json:"nmc"`
	Amount  string `json:"amount"`
}

func GethqInfo() {
	dealtime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 25, 10, 0, time.Local)
	redisClient := redis.GetRedisInter()
	jsonBytes, _ := redisClient.GetHashValue("stkbidinfo", dealtime.Format("2006-01-02"))
	if jsonBytes == "" { // 去查询竞价行情信息
		url := fmt.Sprintf(public.IncPrcUrl, "100")
		webhq.GethqInfo(url, dealtime)
	}
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

type bidredis struct {
	TotalAmount int64       `json:"totalamount"`
	Bidhq       []bidFormat `json:"bidhq"`
}

// 分析集合竞价数据
func moniterBidInfo(mapSinadata map[string][]public.SinaQuoteData) {
	var totalAmt int64 = 0

	redisClient := redis.GetRedisInter()
	result := make([]bidFormat, 0, len(mapSinadata))
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
						result = append(result, temp)
						sliceMoniter = append(sliceMoniter, sinInfo.Code)
					}
				}
			}
			// 对最后一条9:25:00的集合竞价数据做分析
			if strings.Compare(sinInfo.Ticktime, "09:25:00") == 0 {
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
						result = append(result, temp)
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
					result = append(result, temp)
					sliceMoniterTmp = append(sliceMoniterTmp, sinInfo.Code)
				}
			}
			lastTime = hqTime
			lastSina = sinInfo
		}
	}
	var bidhq = bidredis{}
	bidhq.TotalAmount = totalAmt
	bidhq.Bidhq = result
	jsonBytes, err := json.Marshal(bidhq)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	redisClient.SetOneHashData("stkbidinfo", time.Now().Format("2006-01-02"), jsonBytes)
}
