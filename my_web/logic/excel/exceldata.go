package excel

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"localhost.com/my_web/repo/log"
)

var (
	excelPath      = "../config/xuangu/交易日期_2026.xlsx"
	stockExcelPath = "../config/xuangu/隔日涨停股票分析.xlsx"

	currExcel     ExcelOrder // 当前涨停数据
	preExcel      ExcelOrder // 上一日涨停数据
	IsMonFirstDay bool       // 是否是一个月的第一天，需要特殊处理，涉及到两个表单
	redisChan     chan stockData
	monthInfo     []*OneWeek // 一个月的数据
	indexWeek     int        // 标记一个月的第几周
	errLog        log.LogInter
)

func init() {
	currExcel.stkInfo = make(map[int][]rowInfo)
	preExcel.stkInfo = make(map[int][]rowInfo)

	monthInfo = make([]*OneWeek, 0)

	indexWeek = 0

	IsMonFirstDay = false

	errLog = log.GetLogInter()
	errLog.NewFileLogger("../log/error.log")

	redisChan = make(chan stockData)
}

// 今日涨停和昨日涨停数据
// sheet1和sheet2里的一行数据
type rowInfo struct {
	Code          string  // 股票代码
	Name          string  // 股票名称
	Reason        string  // 涨停理由
	Money         string  // 封单金额
	changepercent float64 // T+1日涨幅
}

// 今日涨停和昨日涨停数据
type ExcelOrder struct {
	day     int
	mon     int
	year    int
	weekNum time.Weekday      // 记录星期几
	stkInfo map[int][]rowInfo // int : 第几板  []rowInfo每个板中的股票信息
}

// 一只股票的信息
type OneStk struct {
	Code    string `json:"Code"` // 股票代码
	Name    string `json:"Name"` // 股票名称
	Reason  string `json:"-"`    // 涨停理由
	Money   string `json:"-"`    // 封单金额
	Perform string `json:"-"`    // T+日表现  < -4：大面； < -4且小于0：小面； > 0 且没涨停红盘； 涨停：晋级
	PreRate string `json:"-"`    // 特殊处理，第一条数据记录板块表现 例："共4支 晋级率13.33%"
	Day     int    `json:"Day"`  // 连续涨停天数
}

// 一天的股票信息 map[int][]OneStk, int:表示第几板 []OneStk股票信息集合
type OneDay struct {
	Date    string           //日期
	WeekNum time.Weekday     // 星期数
	MapStk  map[int][]OneStk // 格式为 arrStk[1].[0] 连板数为1的第一支股票
}

// 一周的股票信息  Maxboard 表示每个连板占excel的行数数组，例如 [4 4 4 32] 分别代表 4,3,2,1连板的行数
// []OneDay每天股票信息的集合
type OneWeek struct {
	Maxboard []int    // 最大的板块数
	ArrDay   []OneDay // 每天的股票信息集合
}

// 接口用来操作切换每周的数据，最多只配置了五周，一个月的
type Weeker interface {
	Init()
	SetParam(x interface{})
	PushStk(int, int, OneStk)
	GetMaxBoard() []int
	GetArrDay() []OneDay
}

// func (w *OneWeek) Show() {
//     fmt.Printf("ArrDay : %p\n", w.ArrDay)
//     fmt.Printf("Maxboard : %p\n", w.Maxboard)
// }

// 初始化
func (w *OneWeek) Init() {
	w.Maxboard = make([]int, 0)
	w.ArrDay = make([]OneDay, 0)
}

// 设置每周的数据，包括股票信息，连板数组
func (w *OneWeek) SetParam(x interface{}) {
	switch v := x.(type) {
	case int:
		w.Maxboard = append(w.Maxboard, v)
	case []int:
		w.Maxboard = v
	case OneDay:
		w.ArrDay = append(w.ArrDay, v)
	}
}

// 写入一天的股票数据 day:第几天 k:第几板
func (w *OneWeek) PushStk(day, k int, oneStk OneStk) {
	if len(w.ArrDay) < day+1 {
		fmt.Println(day, k, oneStk, "写入失败")
		return
	}
	w.ArrDay[day].MapStk[k] = append(w.ArrDay[day].MapStk[k], oneStk)
}

// 获取每个连板占excel的行数数组，例如 [4 4 4 32] 分别代表 4,3,2,1连板的行数
func (w *OneWeek) GetMaxBoard() []int {
	return w.Maxboard
}

// 获取一周的股票数据
func (w *OneWeek) GetArrDay() []OneDay {
	return w.ArrDay
}

// []排序
type ByDescending []int

func (a ByDescending) Len() int           { return len(a) }
func (a ByDescending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescending) Less(i, j int) bool { return a[i] > a[j] } // 注意这里的比较符号 降序排序

type stkReason struct {
	Symbol string `json:"symbol"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
	Limit  string `json:"limit"`
	Nmc    string `json:"nmc"`
}

// convertChineseNumToDecimal 转换字符串格式 例："1300万" ==> "0.13亿"
func convertChineseNumToDecimal(inputStr string) string {
	// 定义单位和对应的数值
	units := map[string]float64{"万": 10000, "亿": 100000000}

	// 使用正则表达式去除中文逗号，并获取数字和单位
	re := regexp.MustCompile(`(\d+,?\d*\.\d*)\s*([万亿]?)`)
	matches := re.FindStringSubmatch(inputStr)

	if len(matches) < 3 {
		return ""
	}

	numStr := strings.Replace(matches[1], ",", "", -1)
	unit := matches[2]

	// 将数字字符串转换为浮点数
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return ""
	}

	// 根据单位进行转换
	if unit == "" {
		num *= units["万"]
	} else {
		num *= units[unit]
	}

	// 格式化输出
	return fmt.Sprintf("%.2f亿", num/100000000)
	// return fmt.Sprintf("%.2f万", num/10000)
}

// GetOrderInfo 取昨日、今日涨停数据，将sheet1、sheet2的表格二维数组转换为ExcelOrder格式
func (e *ExcelOrder) GetOrderInfo(readSheet [][]string) {
	if e.stkInfo == nil {
		e.stkInfo = make(map[int][]rowInfo)
	}

	changepercentCol, reasonCol, indexCol, moneyCol := 5, 7, 6, 9

	for rowNum, row := range readSheet {
		indexKey := 0
		var stRowInfo rowInfo
		for sheetNum, sheeCol := range row {
			// if changepercentCol == 0 && strings.Contains(sheeCol, "涨跌幅(%)") {
			// 	changepercentCol = rowNum - 38
			// }
			// if reasonCol == 0 && strings.Contains(sheeCol, "涨停原因类别") {
			// 	reasonCol = rowNum - 38
			// }
			// if indexCol == 0 && strings.Contains(sheeCol, "连续涨停天数(天)") {
			// 	indexCol = rowNum - 38
			// }
			// if moneyCol == 0 && strings.Contains(sheeCol, "涨停封单额(元)") {
			// 	moneyCol = rowNum - 38
			// }
			if sheetNum == 0 && e.year <= 0 {
				// 2025.07.09 涨停 沪深 非ST
				// 涨停 2024.07.31
				dateStr := ""
				idx := strings.Index(sheeCol, "2")
				if idx < 0 || idx+10 > len(sheeCol) {
					continue
				}
				dateStr = sheeCol[idx : idx+10]

				// 指定布局字符串，与dateStr的格式相匹配
				layout := "2006.01.02"
				// 使用Parse函数解析日期字符串
				t, err := time.Parse(layout, dateStr)
				if err != nil {
					fmt.Println("Error parsing date:", err)
					return
				}
				e.weekNum = t.Weekday()
				// 2006.01.02
				arrTemp := strings.Split(dateStr, ".")
				if len(arrTemp) > 2 {
					e.year, _ = strconv.Atoi(arrTemp[0])
					e.mon, _ = strconv.Atoi(arrTemp[1])
					e.day, _ = strconv.Atoi(arrTemp[2])
				}
			}
			if rowNum > 35 && len(row) > 15 {
				switch sheetNum {
				case 2:
					// 不足6位的用0补齐
					stRowInfo.Code = fmt.Sprintf("%06s", strings.TrimSpace(sheeCol))
				case 3:
					stRowInfo.Name = strings.TrimSpace(sheeCol)
				case changepercentCol:
					stRowInfo.changepercent, _ = strconv.ParseFloat(strings.TrimSpace(sheeCol), 64)
				case indexCol:
					indexKey, _ = strconv.Atoi(strings.TrimSpace(sheeCol))
				case reasonCol:
					stRowInfo.Reason = strings.TrimSpace(sheeCol)
				case moneyCol:
					// 2,096.27万
					stRowInfo.Money = convertChineseNumToDecimal(strings.TrimSpace(sheeCol))
				default:
					break
				}
			}
		}
		if indexKey > 0 {
			e.stkInfo[indexKey] = append(e.stkInfo[indexKey], stRowInfo)
		}
	}
}

func (el *ExcelOrder) ToDataString() string {
	return fmt.Sprintf("%04d%02d%02d", el.year, el.mon, el.day)
}

type stockData struct {
	StockCode string   `json:"stock_code"`
	Fields    []string `json:"fields"`
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

type SinaQuoteData struct {
	Symbol        string  `json:"symbol"`        // 证券内码 含市场
	Code          string  `json:"code"`          // 证券代码
	Name          string  `json:"name"`          // 股票名称
	Trade         string  `json:"trade"`         // 市场价格
	Pricechange   float64 `json:"pricechange"`   // 涨跌价格 涨了几块或跌了几块
	Changepercent float64 `json:"changepercent"` // 涨跌幅
	Buy           string  `json:"buy"`           // 买一价
	Sell          string  `json:"sell"`          // 卖一价
	Settlement    string  `json:"settlement"`    // 昨收盘
	Open          string  `json:"open"`          // 开盘价
	High          string  `json:"high"`          // 最高价
	Low           string  `json:"low"`           // 最低价
	Volume        int     `json:"volume"`        // 成交量  单位 股
	Amount        int     `json:"amount"`        // 成交额  单位 元
	Ticktime      string  `json:"ticktime"`      // 时间，行情时间
	Per           float64 `json:"per"`           // 市盈率
	Pb            float64 `json:"pb"`            // 市净率
	Mktcap        float64 `json:"mktcap"`        // 总市值 单位 万
	Nmc           float64 `json:"nmc"`           // 流通市值 单位 万
	Turnoverratio float64 `json:"turnoverratio"` // 换手率
}
