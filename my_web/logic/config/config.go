package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"

	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	excelctl "localhost.com/my_web/repo/execelctl"
	"localhost.com/my_web/repo/filectl"
	"localhost.com/my_web/repo/redis"
	"localhost.com/my_web/repo/visitsina"
)

// SetConfigFile 生成监控config
func SetConfigFile(configFlag bool) {
	SetProbemJSON()
	SetMoniterConfig(configFlag)

	GetBidStkHqInfo()

	fmt.Println("写入 config.json文件完成")
	time.Sleep(3000 * time.Millisecond)

}

// SetProbemJSON 设置问题股json
func SetProbemJSON() {
	// 打开文件
	file, err := os.Open(problemTxtPath)
	if err != nil {
		fmt.Println("open file err:", err.Error())
		return
	}

	// 处理结束后关闭文件
	defer file.Close()
	// 使用bufio读取
	r := bufio.NewReader(file)

	var problemStock ProblemStock
	for {
		// 分行读取文件  ReadLine返回单个行，不包括行尾字节(\n  或 \r\n)
		filedata, _, err := r.ReadLine()

		// 读取到末尾退出
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("read err", err.Error())
			break
		}

		data := string(filedata)
		data = strings.TrimSpace(data)

		if len(data) > 0 && data[0] >= 65 && data[0] <= 122 {
			arrTmp := strings.Split(data, "\t")

			if len(arrTmp) > 1 {
				gbkBytes := []byte(arrTmp[1]) // 这些字节代表GBK编码
				// 创建一个GBK的解码器
				decoder := simplifiedchinese.GBK.NewDecoder()
				// 使用解码器将GBK字节转换为UTF-8
				utf8Bytes, _, err := transform.Bytes(decoder, gbkBytes)
				if err != nil {
					fmt.Println("Error decoding:", err)
					return
				}

				// utf8Bytes现在包含的是UTF-8编码的字节
				utf8String := string(utf8Bytes)

				problemStock.ArrStock = append(problemStock.ArrStock, StockPlm{strings.ToLower(arrTmp[0]), utf8String})
			}
		}
	}
	// 使用json.MarshalIndent美化输出的JSON（可选，增加可读性）
	jsonData, err := json.MarshalIndent(problemStock, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	// 打开一个文件用于写入
	fileW, errW := os.Create(problemJSONPath)
	if errW != nil {
		fmt.Println("Error creating file:", errW)
		return
	}
	defer fileW.Close()

	// 写入JSON数据到文件
	_, errW = fileW.Write(jsonData)
	if errW != nil {
		fmt.Println("Error writing to file:", errW)
		return
	}
}

func callIndexURL(url string, data []byte) {
	if url == indexURL {
		visit := visitsina.GetVisitInter()
		visit.InitSinaQuote(callIndexURL, errLog)
		visit.VisitSina(incPrcURL)
	} else {
		quotes := make([]SinaQuoteData, 0)
		json.Unmarshal(data, &quotes)

		sshqChan <- quotes
		close(sshqChan)
	}
}

// SetMoniterConfig 自动设置config文件
func SetMoniterConfig(bFlag bool) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for sshq := range sshqChan {
			for _, v := range sshq {
				if strings.Contains(strings.ToLower(v.Symbol), "sz") && !strings.Contains(strings.ToLower(v.Name), "st") { // 深圳 非ST
					if strings.HasPrefix(strings.ToLower(v.Name), "n") || strings.HasPrefix(strings.ToLower(v.Name), "c") {
						continue
					}
					if strings.Compare(v.Trade, v.High) == 0 && strings.Compare(v.Sell, "0.000") == 0 &&
						!lessOrEqual(v.Changepercent, 15.00, 0.00001) { // 20cm 涨停
						// 过滤问题股
						bProblem := false
						for _, pro := range problemStock.ArrStock {
							if strings.Compare(pro.Code, v.Symbol) == 0 {
								fmt.Printf("Code : %s, Name : %s 是问题股\n", pro.Code, pro.Name)
								bProblem = true
								break
							}
						}
						if bProblem {
							continue
						}
						dTrade, _ := strconv.ParseFloat(v.Trade, 64)
						// 市价小于6块 且 (市值 < 30 || 市盈率 < 0)
						if bFlag && lessOrEqual(dTrade, 6.00, 0.000001) && (lessOrEqual(v.Mktcap, 300000.0, 0.00001) ||
							lessOrEqual(v.Per, 0.00, 0.00001)) {
							continue
						} else if !bFlag || !lessOrEqual(dTrade, 3.00, 0.000001) {
							redisctl := redis.GetRedisInter()
							strName, _ := redisctl.GetHashValue("stkinfo", v.Name)
							strCnt, _ := redisctl.GetHashValue("stklimitup", v.Name)
							svrConfig.ArrMoniter = append(svrConfig.ArrMoniter, MoniterInfo{v.Symbol, v.Name, 15.00, true,
								v.Amount, strName, strCnt, "", false})
						}
					}
				}
			}
		}
	}()

	visit := visitsina.GetVisitInter()

	visit.InitSinaQuote(callIndexURL, errLog)
	// 发起调用
	visit.VisitSina(indexURL)

	wg.Wait()
	svrConfig.OrderInfo.IsMoniter = true
	svrConfig.OrderInfo.IsSecondOrder = false
	svrConfig.OrderInfo.IsFirstOrder = false
	svrConfig.OrderInfo.IsFirst = true
	svrConfig.OrderInfo.FirstLimitNum = 1
	svrConfig.OrderInfo.OrderMoney = 11000.00
	svrConfig.OrderInfo.IsMoreAmount = true
	svrConfig.OrderInfo.MaxAmount = 2.0
	svrConfig.OrderInfo.MaxNmc = 100.0

	if len(svrConfig.ArrMoniter) == 0 {
		fmt.Println("没有符合条件的创业板")
		svrConfig.ArrMoniter = append(svrConfig.ArrMoniter, MoniterInfo{"sz300750", "宁德时代", 15.00, true, 0, "", "", "", false})
	}

	svrConfig.FirStkInfo.ArrFirIsNot = append(svrConfig.FirStkInfo.ArrFirIsNot, FirIsNot{"sz301090", "华润材料", false})
	svrConfig.FirStkInfo.Reason = ""

	// 使用json.MarshalIndent美化输出的JSON（可选，增加可读性）
	jsonData, err := json.MarshalIndent(svrConfig, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	// 打开一个文件用于写入
	file, err := os.Create(configJSONPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// 写入JSON数据到文件
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func callCount(url string, data []byte) {
	if url == countURL {
		str := string(data)
		str = str[1 : len(str)-1]
		pageCount, _ = strconv.Atoi(str)
		pageCount = (pageCount + 100 - 1) / 100

		url := fmt.Sprintf(pageURL, 1)

		visit := visitsina.GetVisitInter()
		visit.InitSinaQuote(callCount, errLog)
		visit.VisitSina(url)
	} else if pageCount > 0 {
		pageIndex, _ := strconv.Atoi(url[strings.Index(url, "page=")+5 : strings.Index(url, "&num=")])

		quotes := make([]SinaQuoteData, 0)
		json.Unmarshal(data, &quotes)
		for _, v := range quotes {
			bidChan <- v
		}

		if pageIndex < pageCount {
			url := fmt.Sprintf(pageURL, pageIndex+1)

			visit := visitsina.GetVisitInter()
			visit.InitSinaQuote(callCount, errLog)
			visit.VisitSina(url)
		} else {
			close(bidChan)
		}
	}
}

func GetBidStkHqInfo() {
	now := time.Now()

	fileName := fmt.Sprintf("../log/bidinfo/bid_moniter_%04d_%02d_%02d.log", now.Year(), now.Month(), now.Day())
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	type bidHqInfo struct {
		Code          string
		Name          string
		Nmc           string
		Amount        string
		Open          string // 开盘价
		Changepercent string // 开盘涨幅
		Reason        string
	}
	arrHqInfo := make([]bidHqInfo, 0)
	writeSlice := make([]string, 0)

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		strLine := strings.TrimSpace(string(line))

		var stHqInfo = bidHqInfo{}

		bflag := false
		if strings.Contains(strLine, "成交量较大:") {
			strLine = strLine[strings.Index(strLine, "成交量较大:"):]
		} else if strings.Contains(strLine, "波动较大:") {
			// 两秒内涨幅波动超过3%:
			bflag = true
			strLine = strLine[strings.Index(strLine, "波动较大:"):]
		} else if strings.Contains(strLine, "涨幅榜前100成交量") {
			writeSlice = append(writeSlice, strLine)
			continue
		} else {
			continue
		}

		strLine = strings.ReplaceAll(strLine, " ", "")

		// Code
		stHqInfo.Code = strLine[strings.Index(strLine, "Code:")+5 : strings.Index(strLine, ",")]
		strLine = strLine[strings.Index(strLine, "Name:"):]

		// Name
		stHqInfo.Name = strLine[strings.Index(strLine, ":")+1 : strings.Index(strLine, ",")]
		strLine = strLine[strings.Index(strLine, "Reason:"):]
		if bflag {
			stHqInfo.Name = fmt.Sprintf("%%3-%s", stHqInfo.Name)
		}

		// Reason
		stHqInfo.Reason = strLine[strings.Index(strLine, ":")+1 : strings.Index(strLine, ",")]
		strLine = strLine[strings.Index(strLine, "Nmc:"):]

		// 流通盘
		stHqInfo.Nmc = strLine[strings.Index(strLine, ":")+1 : strings.Index(strLine, ",")]
		strLine = strLine[strings.Index(strLine, "Amount:"):]

		// 成交量
		stHqInfo.Amount = strLine[strings.Index(strLine, ":")+1 : strings.Index(strLine, ",")]
		strLine = strLine[strings.Index(strLine, "Percent:"):]

		// 开盘涨幅
		stHqInfo.Changepercent = strLine[strings.Index(strLine, ":")+1:]
		arrHqInfo = append(arrHqInfo, stHqInfo)
	}
	file.Close()

	ok, result := GetStkResonExcel("stkinfo")
	mapStkName := make(map[string]string)
	for _, v := range result {
		mapStkName[v] = ""
	}
	var wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		info := excelctl.GetExcelInfo(excelPath, "stock")
		suc, fail := 0, 0
		for stSin := range bidChan {
			for _, name := range result {
				if name == stSin.Name {
					mapStkName[name] = stSin.Code
				}
			}
			for _, v := range arrHqInfo {
				if stSin.Code == v.Code {
					strLine := fmt.Sprintf("%04d-%02d-%02d: Code:%s,Name:%s,开盘:%s%%,收盘:%.2f%%,流通盘:%s,集合竞价成交量:%s", now.Year(), now.Month(), now.Day(), v.Code,
						v.Name, v.Changepercent, stSin.Changepercent, v.Nmc, v.Amount)
					writeSlice = append(writeSlice, strLine)
					break
				}
			}

			for row := 1; row < len(info); row++ {
				tmp := info[row]
				if len(tmp) < 9 {
					if tmp[1] == stSin.Code {
						if len(tmp) < 4 {
							info[row] = append(info[row], make([]string, 6-len(tmp))...)

							info[row][3] = fmt.Sprintf("%.3f", stSin.Mktcap/10000.0)
							info[row][4] = fmt.Sprintf("%d", stSin.Amount/10000)
							info[row][5] = fmt.Sprintf("%.2f%%", (float64(stSin.Amount)/100)/stSin.Mktcap)
						} else if len(tmp) < 7 {
							info[row] = append(info[row], make([]string, 7-len(tmp))...)

							info[row][6] = stSin.Open
						} else {
							info[row] = append(info[row], make([]string, 9-len(tmp))...)

							price, _ := strconv.ParseFloat(stSin.Open, 64)
							price_last, _ := strconv.ParseFloat(info[row][6], 64)

							info[row][7] = stSin.Open

							rate := (price - price_last) / price_last * 100
							info[row][8] = fmt.Sprintf("%.2f", rate)
							if rate-0.00001 > 0.00 {
								suc++
							} else {
								fail++

							}
						}
					}
				}
			}
		}
		// info = append(info, []string{"合计", fmt.Sprintf("盈：%d", suc), fmt.Sprintf("亏：%d", fail)})
		WriteExecelInfo(excelPath, "stock", info)
	}()

	visit := visitsina.GetVisitInter()

	visit.InitSinaQuote(callCount, errLog)
	// 发起调用
	visit.VisitSina(countURL)

	wg.Wait()

	filectl.WriteFileInfo(bidfilePath, writeSlice)

	if ok {
		arrStkCode := []string{}
		for _, v := range result {
			arrStkCode = append(arrStkCode, mapStkName[v])
		}

		ModStkResonExcel(arrStkCode, "stkinfo")
	}
}

// lessOrEqual 浮点数比较 小于或等于
func lessOrEqual(a, b, epsilon float64) bool {
	return a < b || math.Abs(a-b) <= epsilon
}

// GetStkResonExcel 用于生成股票板块信息的excel，每天的config操作中写一个新的excel出来
func GetStkResonExcel(sheetName string) (bool, []string) {
	redisctl := redis.GetRedisInter()
	// 股票的板块信息
	mStkInfo := redisctl.GetHashAllValue("stkinfo", true)
	// 股票的涨停次数--这里用它来类比股票活跃度
	mStkLimit := redisctl.GetHashAllValue("stklimitup", false)

	result := []string{}

	// 打开一个已存在的 Excel 文件
	f, err := excelize.OpenFile(stockExcelPath)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return false, nil
	}
	defer f.Close()

	// 检查工作表是否存在，存在就删除
	_, err = f.GetSheetIndex(sheetName)
	if err == nil {
		f.DeleteSheet(sheetName)
		f.SaveAs(stockExcelPath)
	}

	// 创建新Sheet
	f.NewSheet(sheetName)
	// 冻结第一行的单元格
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		YSplit:      1,
		TopLeftCell: "A1",
		ActivePane:  "bottomLeft"}) // bottomLeft

	// 设置列宽
	// 参数分别为：工作表名，起始列，结束列（如果只设置一个列，则起始列和结束列相同），宽度值
	f.SetColWidth(sheetName, "A", "B", 20)
	f.SetColWidth(sheetName, "C", "C", 100)
	f.SetColWidth(sheetName, "D", "D", 20)

	row := 1
	// 写第一行
	f.SetCellValue(sheetName, GetExcelCol(row, 1, time.Sunday), "股票名称")
	f.SetCellValue(sheetName, GetExcelCol(row, 2, time.Sunday), "涨停次数")
	f.SetCellValue(sheetName, GetExcelCol(row, 3, time.Sunday), "板块信息")
	f.SetCellValue(sheetName, GetExcelCol(row, 4, time.Sunday), "股票代码")
	row++

	stkInfo := []struct {
		Name     string
		LimitNum int
		Reason   string
	}{}
	for name, v := range mStkLimit {
		reason, ok := mStkInfo[name]
		if ok {
			num, err := strconv.Atoi(v)
			if err == nil {
				stkInfo = append(stkInfo, struct {
					Name     string
					LimitNum int
					Reason   string
				}{name, num, reason})
			}
		}
	}

	sort.Slice(stkInfo, func(i, j int) bool {
		return stkInfo[i].LimitNum > stkInfo[j].LimitNum
	})

	for _, v := range stkInfo {
		f.SetCellValue(sheetName, GetExcelCol(row, 1, time.Sunday), v.Name)
		f.SetCellValue(sheetName, GetExcelCol(row, 2, time.Sunday), v.LimitNum)
		f.SetCellValue(sheetName, GetExcelCol(row, 3, time.Sunday), v.Reason)
		row++

		result = append(result, v.Name)
	}

	f.SaveAs(stockExcelPath)
	return true, result
}

// GetExcelCol 将行列数转为excel里面特有的表格字符，如：第1行，第1列 ==> A1
// GetExcelCol 特殊处理，按照星期几增加列数，这样就可以把每天的数据按照星期数固定在特定的列上
func GetExcelCol(row, col int, weekday time.Weekday) string {
	addcol := 0
	if weekday != time.Sunday && weekday != time.Saturday { // 周一到周五按照固定的列写入，防止节假日格式错乱的情况
		addcol = 4 * int(weekday-1)
	}
	str, _ := excelize.CoordinatesToCellName(col+addcol, row)
	// return fmt.Sprintf("%c%d", byte('A'+(col-1)%26), row)
	return str
}

// WriteExecelInfo 写入execel文件
func WriteExecelInfo(fileName, sheetName string, info [][]string) bool {
	// 打开一个已存在的 Excel 文件
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return false
	}
	defer f.Close()

	// 检查工作表是否存在，存在就删除
	_, err = f.GetSheetIndex(sheetName)
	if err == nil {
		f.DeleteSheet(sheetName)
		// f.SaveAs("fileName")
	}

	// 创建新Sheet
	f.NewSheet(sheetName)
	// 冻结第一行的单元格
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		YSplit:      1,
		TopLeftCell: "A1",
		ActivePane:  "bottomLeft"}) // bottomLeft

	for row, rows := range info {
		for col, value := range rows {
			f.SetCellValue(sheetName, GetExcelCol(row+1, col+1, time.Sunday), value)
			if row == 0 {
				f.SetColWidth(sheetName, getExcelColumnName(col), getExcelColumnName(col), 20)
			}
		}
	}

	f.SaveAs(fileName)
	return true
}

func ModStkResonExcel(arrName []string, sheetName string) bool {
	// 打开一个已存在的 Excel 文件
	f, err := excelize.OpenFile(stockExcelPath)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return false
	}
	defer f.Close()

	// 检查工作表是否存在
	_, err = f.GetSheetIndex(sheetName)
	if err != nil {
		return false
	}
	f.SetColWidth(sheetName, "D", "D", 20)

	for row, v := range arrName {
		f.SetCellValue(sheetName, GetExcelCol(row+2, 4, time.Sunday), v)
	}
	f.SaveAs(stockExcelPath)
	return true
}

func getExcelColumnName(col int) string {
	var result string
	for col >= 0 {
		remainder := col % 26
		result = string('A'+remainder) + result
		col = col/26 - 1
		if col < 0 {
			break
		}
	}
	return result
}
