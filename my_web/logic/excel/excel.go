package excel

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	excelctl "localhost.com/my_web/repo/execelctl"
	"localhost.com/my_web/repo/redis"
	"localhost.com/my_web/repo/visitsina"
)

// SetStockExcel 写excel入口函数，将sheet1和sheet2的数据写入每月总结数据中
func SetStockExcel(isRedis bool) {
	// 获取当日涨停数据
	readSheet := excelctl.GetExcelInfo(excelPath, "Sheet1")
	if readSheet == nil {
		fmt.Println("Sheet1 is nil")
		return
	}
	// 将sheet1的格式转为ExcelOrder
	currExcel.GetOrderInfo(readSheet)

	// 获取前日涨停数据
	readSheet = excelctl.GetExcelInfo(excelPath, "Sheet2")
	if readSheet == nil {
		fmt.Println("Sheet2 is nil")
		return
	}
	// 将sheet2的格式转为ExcelOrder
	preExcel.GetOrderInfo(readSheet)

	if currExcel.mon <= 0 || preExcel.mon <= 0 || len(currExcel.stkInfo) == 0 || len(preExcel.stkInfo) == 0 {
		fmt.Println("写入失败...")
		return
	}

	// 根据当前交易日推算出是写在哪个sheet表单中
	// 最坏的情况是1号是周二
	// 最好的情况是1号是周一
	sheetName := ""
	sheetName = fmt.Sprintf("%d月", currExcel.mon)
	// 每月7号之前，如果对应的月份为空
	if currExcel.day < 7 && excelctl.GetExcelInfo(excelPath, sheetName) == nil {
		// 取前一交易日
		lastdate, err := time.Parse("2006.01.02", fmt.Sprintf("%d.%02d.%02d", preExcel.year, preExcel.mon, preExcel.day))
		if err != nil {
			fmt.Println(err)
			return
		}
		// 取当前交易日
		curdate, err := time.Parse("2006.01.02", fmt.Sprintf("%d.%02d.%02d", currExcel.year, currExcel.mon, currExcel.day))
		if err != nil {
			fmt.Println(err)
			return
		}
		diff := curdate.Sub(lastdate)
		// 如果间隔小于72小时(两天)
		if diff.Hours() < 72 {
			sheetName = fmt.Sprintf("%d月", currExcel.mon-1)
		}
	}

	lastMon := fmt.Sprintf("%d月", currExcel.mon-1)

	if isRedis {
		// 将当天涨停信息统计并写入redis，同时生成redis的execel文件
		SetRedisByExcelOrder(currExcel)
	}

	fmt.Printf("正在写入表格: %s ...\n", sheetName)
	// 获取写入单元格
	writeSheet := excelctl.GetExcelInfo(excelPath, sheetName)
	if writeSheet == nil { // 如果当前日期是一个月的第一天，取sheet取不到，就取上一个月的
		writeSheet = excelctl.GetExcelInfo(excelPath, lastMon)
		if writeSheet != nil {
			IsMonFirstDay = true
		}
	}

	// 1.将excel的二维数组转为OneWeek数据
	// 2.将curExcel、preExcel当日涨停和昨日涨停的数据插入OneWeek
	GetWriteSheet(writeSheet, lastMon)

	IsMonFirstDay = false
	// 按格式写入sheet表单
	WriteExcel(sheetName)

	fmt.Println("写入成功")
}

// SetRedisByExcelOrder 用于将当天涨停的股票的板块信息写入redis
func SetRedisByExcelOrder(currDay ExcelOrder) {
	// 用于统计涨停的板块
	redisMap := make(map[string]interface{}, 0)

	// 用于计算当日反包涨停的信息
	stkSlice := make([]OneStk, 0)

	visit := visitsina.GetVisitInter()
	// 创建回调函数
	callBack := func(url string, data []byte) {
		results := parseStockData(data)
		for _, v := range results {
			redisChan <- v
		}
	}

	redis := redis.GetRedisInter()
	for k, v := range currDay.stkInfo {
		for _, stk := range v {
			stkSlice = append(stkSlice, OneStk{stk.Code, stk.Name, "", "", "", "", k})

			jsonStr, _ := redis.GetHashValue("stkreason", changeStkCode(stk.Code))
			stkReason := stkReason{}
			json.Unmarshal([]byte(jsonStr), &stkReason)
			tmp, _ := strconv.Atoi(stkReason.Limit)
			stkReason.Symbol = changeStkCode(stk.Code)
			stkReason.Code = stk.Code
			stkReason.Name = stk.Name
			stkReason.Reason = stk.Reason
			stkReason.Limit = fmt.Sprintf("%d", tmp+1)

			// 发起调用
			visit.InitSinaQuote(callBack, errLog)
			go visit.VisitSina(fmt.Sprintf("http://qt.gtimg.cn/q=%s", changeStkCode(stk.Code)))

			v := <-redisChan
			stkReason.Name = v.Fields[1]
			stkReason.Nmc = fmt.Sprintf("%s亿", v.Fields[44])
			value, _ := json.Marshal(stkReason)
			redisMap[changeStkCode(stk.Code)] = value
		}
	}
	defer close(redisChan)

	redis.SetHashInterface(redisMap, "stkreason")
	// GetStkResonExcel("stkinfo")

	// 用于写入当天涨停的股票，且提示前日涨停，当日也涨停的股票
	// 将切片编码为 JSON
	value, err := json.Marshal(stkSlice)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		// 格式化为 YYYY-MM-DD 的字符串
		dateStr := fmt.Sprintf("%04d-%02d-%02d", currDay.year, currDay.mon, currDay.day)

		if err := redis.SetOneHashDate("stkdata", dateStr, value); err != nil {
			fmt.Println("SetRedisByExcelOrder SetOneHashDate :", err.Error())
		}

		dataTmp := ""
		currday, _ := time.Parse("2006-01-02", dateStr)
		for i := 0; i < 2; i++ { // 找T-2日的股票信息
			for cnt := 0; cnt < 15; cnt++ { // 因为不知道哪天是交易日最多往前找15天
				// 获取前一天的时间
				currday = currday.AddDate(0, 0, -1)
				// 取涨停股票信息
				dataTmp, _ = redis.GetHashValue("stkdata", currday.Format("2006-01-02"))
				if dataTmp != "" {
					break
				}
			}
		}

		if dataTmp != "" {
			stkSlice = make([]OneStk, 0)
			// 将 JSON 数据解码为数组
			if err := json.Unmarshal([]byte(dataTmp), &stkSlice); err != nil {
				// log.Fatalf("Failed to unmarshal person: %v", err)
				fmt.Printf("Failed to unmarshal person: %v", err)
			}
			for key, v := range currDay.stkInfo {
				if key == 1 {
					info := excelctl.GetExcelInfo(stockExcelPath, "stock")
					oldLen := len(info)
					for _, stk := range v {
						for _, obj := range stkSlice {
							if obj.Code == stk.Code || obj.Name == stk.Name {
								fmt.Println("隔日涨停的股票：", obj.Code, obj.Name)
								tmp := make([]string, 9)
								tmp[0] = currDay.ToDataString()
								tmp[1] = obj.Code
								tmp[2] = obj.Name
								info = append(info, tmp)
								break
							}
						}
					}
					if len(info) > oldLen {
						WriteExecelInfo(stockExcelPath, "stock", info)
					}
				}
			}
		}
	}
}

// 将excel的二维数组转为OneWeek结构体数据
func GetWriteSheet(writeSheet [][]string, sheetName string) {
	var weeker Weeker

	maxBoard := 0 // 记录当前板最大板数

	changeRow := -1 // 切换到下一周时的行数，以防止反复切换下一周

	lastday := "" // 标记最后一天

	var curNum int64 = 0 // 记录当前是第几个

	actDay := 0 // 每周实际交易日天数，如周一，周二是节假日，那么这周就只有三个交易日
	for row, rows := range writeSheet {
		isNotStk := false // 不是股票信息的行
		var oneStk OneStk // 一条股票记录

		if row-1 >= 0 {
			if len(writeSheet[row-1]) > 0 && strings.TrimSpace(writeSheet[row-1][0]) == "板数" { // 板数下面一行也标为isNotStk
				isNotStk = true
			}
		}
		for col, cols := range rows {
			colStr := strings.TrimSpace(cols)
			if col == 0 && len(colStr) > 0 { // 第0列标记板数
				if colStr == "板数" {
					if maxBoard > 0 {
						weeker.SetParam(maxBoard)
					}
					isNotStk = true
				}
				curNum, _ = strconv.ParseInt(colStr, 10, 64)
				// 记录具体是第几板
				if curNum > 0 {
					isNotStk = true
					if maxBoard > 0 {
						weeker.SetParam(maxBoard)
					}
				}
				maxBoard = 0
			}
			// 记录日期
			if strings.HasPrefix(colStr, "202") && strings.Contains(colStr, "(周") {
				if changeRow != row {
					// 切换到下一周
					changeRow = row

					week := new(OneWeek)
					week.Init()

					monthInfo = append(monthInfo, week)
					indexWeek++

					weeker = week
					actDay = 0
				}
				strDate := colStr[:strings.Index(colStr, "(")]
				var oneDay OneDay
				oneDay.MapStk = make(map[int][]OneStk)
				oneDay.Date = strDate
				oneDay.WeekNum = time.Weekday(col/4 + 1)
				weeker.SetParam(oneDay)

				if actDay == 0 {
					actDay = 5 - col/4
				}

				lastday = strDate
			}
			if col > 0 && !isNotStk {
				if col%4 == 1 {
					if len(writeSheet[row-1]) > col {
						strTmp := strings.TrimSpace(writeSheet[row-1][col]) // 上一行
						if strings.HasPrefix(strTmp, "共") && strings.Contains(strTmp, "支") {
							oneStk.PreRate = strTmp
						}
					}
				}
				switch col % 4 {
				case 0:
					oneStk.Perform = colStr
					// col最少从4开始进入这里
					if len(oneStk.Name) > 0 && actDay > 0 {
						weeker.PushStk(int(col/4)-6+actDay, int(curNum), oneStk)
					}
				case 1:
					oneStk.Name = colStr
				case 2:
					oneStk.Money = convertChineseNumToDecimal(strings.TrimSpace(colStr))
				default:
					oneStk.Reason = colStr
					if col == len(rows)-1 { // 已经是最后一列了，push进去
						if len(oneStk.Name) > 0 && actDay > 0 {
							weeker.PushStk(int((col+1)/4)-6+actDay, int(curNum), oneStk)
						}
					}
				}
			}
		}
		if curNum > 0 {
			maxBoard++
		}
	}
	if maxBoard > 0 {
		weeker.SetParam(maxBoard)

		// 插入T+1日表现 -- begin
		tmpDate := weeker.GetArrDay()
		if len(tmpDate) > 0 && len(preExcel.stkInfo) > 0 {
			lastDate := tmpDate[len(tmpDate)-1]       // 取最后一天
			if preExcel.weekNum == lastDate.WeekNum { // 确保是同一天
				for k, v := range lastDate.MapStk { // 遍历第几板
					_, ok := preExcel.stkInfo[k]
					if ok {
						for i, value := range v { // 遍历板下面的股票
							for _, one := range preExcel.stkInfo[k] {
								if one.Name == value.Name {
									lastDate.MapStk[k][i].Code = one.Code
									if one.changepercent+4 < 0.00000 {
										lastDate.MapStk[k][i].Perform = "大面"
									} else if one.changepercent-0 < 0.00000 {
										lastDate.MapStk[k][i].Perform = "小面"
									} else {
										tmp, ok := currExcel.stkInfo[k+1]
										if ok {
											bFlag := false
											for _, ttmp := range tmp {
												if ttmp.Name == one.Name {
													lastDate.MapStk[k][i].Perform = "晋级"
													bFlag = true
												}
											}
											if !bFlag {
												lastDate.MapStk[k][i].Perform = "红盘"
											}
										} else {
											lastDate.MapStk[k][i].Perform = "红盘"
										}
									}
								}
							}
						}
					}

				}
			}
		}
		// 插入T+1日表现 -- end
	}

	if IsMonFirstDay {
		WriteExcel(sheetName)
		for _, week := range monthInfo {
			week.Init()
		}

		weeker = monthInfo[0]
		indexWeek = 0
	}

	if len(currExcel.stkInfo) == 0 {
		return
	}

	// 插入当日数据 -- begin
	var curDay OneDay // 当日数据

	// 1. 先把ExcelOrder 转换为 OneDay ////////////////////////////////////////////////////////////
	curDay.Date = fmt.Sprintf("%d.%02d.%02d", currExcel.year, currExcel.mon, currExcel.day)
	curDay.WeekNum = currExcel.weekNum
	curDay.MapStk = make(map[int][]OneStk)
	for k, v := range currExcel.stkInfo { // 遍历第几板
		for i, value := range v { // 遍历到单只股票
			var oneStk OneStk
			oneStk.Code = value.Code
			oneStk.Name = value.Name
			oneStk.Money = value.Money
			oneStk.Reason = value.Reason
			if i == 0 {
				if k == 1 {
					oneStk.PreRate = fmt.Sprintf("共%d支", len(v))
				} else {
					rate := float32(len(v)) / float32(len(preExcel.stkInfo[k-1])) * 100
					oneStk.PreRate = fmt.Sprintf("共%d支 晋级率%.2f%%", len(v), rate)
				}
			}
			curDay.MapStk[k] = append(curDay.MapStk[k], oneStk)
		}
	}

	// ExcelOrder 转换为 OneDay -- end ////////////////////////////////////////////////////////////

	// 2. 将转换成功的OneDay数据插入到对应的周中，此处的难点在于要知道是第几周，而且还包括节假日、切换到下一周、
	//    一个月的第一天、当日连板数与这一周之前的天数不一致(少于或多于之前的连板数)....

	// 定义日期格式
	dateLayout := "2006.01.02"
	// 解析日期字符串
	lastdate, err := time.Parse(dateLayout, lastday)
	curdate, err := time.Parse(dateLayout, curDay.Date)
	if err != nil {
		fmt.Println("Error parsing date:", curDay.Date)
		return
	}
	diff := 73 * time.Hour
	if lastday != "" {
		diff = curdate.Sub(lastdate)
	}

	// 如果是放一天假不跨周末的那种呢 72小时 >= 两天，应该没有放两天假不跨周的吧
	if diff.Hours() >= 72 || IsMonFirstDay {
		//if curDay.WeekNum == time.Monday || IsMonFirstDay { // 周一的时候切换到下一周 想想节假日
		indexWeek++
		for indexWeek > len(monthInfo) {
			week := new(OneWeek)
			week.Init()
			monthInfo = append(monthInfo, week)
		}
		weeker = monthInfo[indexWeek-1]

		keys := []int{}
		for k := range currExcel.stkInfo {
			keys = append(keys, k)
		}
		sort.Sort(ByDescending(keys)) // 4 2 1 如果是不连续的怎么办
		// for i, v := range currExcel.stkInfo { // map的遍历是随机的
		if len(keys) > 0 {
			for k := keys[0]; k > 0; k-- { // [4 2 1]
				value, ok := currExcel.stkInfo[k] // 当日有这个板的数据
				if ok {
					if len(value)+2 < 4 {
						weeker.SetParam(4)
						continue
					}
					weeker.SetParam(int(len(value) + 2))
				} else { // 当日没有这个板的数据，比如[4 2 1] 缺少3板的数据,但也要留出空行出来
					weeker.SetParam(4)
				}
			}
			// 打印出当日连板数组，方便确认当日写入有没有什么问题
			fmt.Println(weeker.GetMaxBoard())
		}
	} else { // 如果当天的数据很多，就可能需要改原来的maxBoard
		tmpSlice := make([]int, 0)
		arrBoard := weeker.GetMaxBoard() // [4 4 12 34]
		keys := []int{}
		for k := range currExcel.stkInfo {
			keys = append(keys, k)
		}
		sort.Sort(ByDescending(keys)) // 4 2 1 如果是不连续的怎么办
		// for i, v := range currExcel.stkInfo { // map的遍历是随机的
		fmt.Println(keys)
		if len(keys) > 0 {
			max := 0
			if keys[0] > len(arrBoard) { // 两边的板数可能不一样，不一定谁多谁少
				max = keys[0]
			} else {
				max = len(arrBoard)
			}
			for k := max; k > 0; k-- { // [3 2 1] | [4 2 1]
				v, ok := currExcel.stkInfo[k] // 当日有这个板的数据
				if ok {
					if len(arrBoard) < k { // 当日的板不在arrBoard里面 如 当日 [3 2 1] 之前的是 [2 1]
						if len(v)+2 < 4 {
							tmpSlice = append([]int{4}, tmpSlice...)
							continue
						}
						// 从前面插入
						tmpSlice = append([]int{len(v) + 2}, tmpSlice...)
						continue
					}
					if arrBoard[len(arrBoard)-k] < len(v)+2 { // 当日的板数量大于之前的数量，按当日的数量来
						tmpSlice = append(tmpSlice, len(v)+2)
						continue
					}
					tmpSlice = append(tmpSlice, arrBoard[len(arrBoard)-k]) // 延续之前的数量
				} else { // 当日没有这个板的数据，比如[4 2 1] 缺少3板的数据,但也要留出空行出来
					tmpSlice = append([]int{4}, tmpSlice...)
				}
			}
		}

		weeker.SetParam(tmpSlice)
		// 打印出当日连板数组，方便确认当日写入有没有什么问题
		fmt.Println(tmpSlice)
	}
	weeker.SetParam(curDay)
	// 插入当日数据 -- end
}

// WriteExcel 按照firWeek,secWeek...写入excel表
func WriteExcel(sheetName string) bool {
	// 打开一个已存在的 Excel 文件
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return false
	}
	defer f.Close()

	// 总的每个板的数量，用于计算总体晋级率的
	advanceRate := make(map[int]int, 0)

	// 获取前日数据
	readSheet := excelctl.GetExcelInfo(excelPath, "涨跌停复盘")
	replayRow := len(readSheet)
	preDate := fmt.Sprintf("%d.%02d.%02d", preExcel.year, preExcel.mon, preExcel.day)
	curDate := fmt.Sprintf("%d.%02d.%02d", currExcel.year, currExcel.mon, currExcel.day)

	// 检查工作表是否存在，存在就删除
	_, err = f.GetSheetIndex(sheetName)
	if err == nil {
		f.DeleteSheet(sheetName)
		f.SaveAs(excelPath)
	}

	// 创建新Sheet
	f.NewSheet(sheetName)
	// 冻结第一列的单元格
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      1,
		TopLeftCell: "A1",
		ActivePane:  "bottomLeft"}) // bottomLeft

	styleId, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "center", // 水平对齐
			Indent:          0,        // 缩进
			JustifyLastLine: true,     // 两端对齐
			ReadingOrder:    0,        // 文字方向
			RelativeIndent:  0,        // 相对缩进
			ShrinkToFit:     false,    // 缩小字体
			TextRotation:    0,        // 文字旋转
			Vertical:        "center", // 垂直对齐
			WrapText:        false,    // 自动换行
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	// 左边对齐
	styleLeft, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "left",   // 左边对齐
			Indent:          0,        // 缩进
			JustifyLastLine: true,     // 两端对齐
			ReadingOrder:    0,        // 文字方向
			RelativeIndent:  0,        // 相对缩进
			ShrinkToFit:     false,    // 缩小字体
			TextRotation:    0,        // 文字旋转
			Vertical:        "center", // 垂直对齐
			WrapText:        false,    // 自动换行
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	// 设置列宽
	// 参数分别为：工作表名，起始列，结束列（如果只设置一个列，则起始列和结束列相同），宽度值
	f.SetColWidth(sheetName, "D", "D", 40)
	f.SetColWidth(sheetName, "H", "H", 40)
	f.SetColWidth(sheetName, "L", "L", 40)
	f.SetColWidth(sheetName, "P", "P", 40)
	f.SetColWidth(sheetName, "T", "T", 40)

	f.SetColWidth(sheetName, "E", "E", 13)
	f.SetColWidth(sheetName, "I", "I", 13)
	f.SetColWidth(sheetName, "M", "M", 13)
	f.SetColWidth(sheetName, "Q", "Q", 13)
	f.SetColWidth(sheetName, "U", "U", 13)

	indexWeek := 1 // 记录是第几周
	var weeker Weeker
	if len(monthInfo) > 0 {
		weeker = monthInfo[indexWeek-1]
	}

	row := 1 // excel的行数
MainLoop:
	for {
		if weeker == nil {
			break
		}
		// 取一周的数据
		arrDay := weeker.GetArrDay()
		boardSlice := weeker.GetMaxBoard()
		if len(boardSlice) == 0 { // 出口，maxBoard记录每个板的数量切片为0
			break
		}

		weekRow := row // 一周开始的行
		col := 1       // 一周开始的列

		rowTmp := weekRow
		// 写第一列
		f.SetCellValue(sheetName, GetExcelCol(rowTmp, col, time.Sunday), "板数")
		// 合并单元格
		f.MergeCell(sheetName, GetExcelCol(rowTmp, col, time.Sunday), GetExcelCol(rowTmp+1, col, time.Sunday))
		f.SetCellStyle(sheetName, GetExcelCol(rowTmp, col, time.Sunday), GetExcelCol(rowTmp+1, col, time.Sunday), styleId)

		rowTmp = rowTmp + 2
		for i, v := range boardSlice {
			f.SetCellValue(sheetName, GetExcelCol(rowTmp, col, time.Sunday), fmt.Sprintf("%d", len(boardSlice)-i))
			f.MergeCell(sheetName, GetExcelCol(rowTmp, col, time.Sunday), GetExcelCol(rowTmp+v-1, col, time.Sunday))
			f.SetCellStyle(sheetName, GetExcelCol(rowTmp, col, time.Sunday), GetExcelCol(rowTmp+v-1, col, time.Sunday), styleId)
			rowTmp += v
		}

		col++
		// 多想一下，遇到节假日怎么办
		rowTmp = weekRow                // 调回到一周开始的地方
		for _, oneDay := range arrDay { // 按天写入
			// 2024.07.01(周一)
			strTmp := fmt.Sprintf("%s(周%d)", oneDay.Date, oneDay.WeekNum)
			f.SetCellValue(sheetName, GetExcelCol(rowTmp, col, oneDay.WeekNum), strTmp)
			f.MergeCell(sheetName, GetExcelCol(rowTmp, col, oneDay.WeekNum), GetExcelCol(rowTmp, col+2, oneDay.WeekNum))

			f.SetCellValue(sheetName, GetExcelCol(rowTmp, col+3, oneDay.WeekNum), "T+1日表现")
			f.SetCellStyle(sheetName, GetExcelCol(rowTmp, col, oneDay.WeekNum), GetExcelCol(rowTmp, col+3, oneDay.WeekNum), styleId)

			rowTmp++

			f.SetCellValue(sheetName, GetExcelCol(rowTmp, col, oneDay.WeekNum), "名称")
			f.SetCellValue(sheetName, GetExcelCol(rowTmp, col+1, oneDay.WeekNum), "封单额")
			f.SetCellValue(sheetName, GetExcelCol(rowTmp, col+2, oneDay.WeekNum), "涨停原因")
			f.SetCellStyle(sheetName, GetExcelCol(rowTmp, col, oneDay.WeekNum), GetExcelCol(rowTmp, col+2, oneDay.WeekNum), styleId)

			// 标记特殊的行，之后写入大面概率用的
			specialRow := rowTmp
			tmpNum := 0   // 记录大面的股票数量
			totalNum := 0 // 记录二板及二板以上的股票总数

			rowTmp++

			// 有点复杂啊
			for pos, v := range boardSlice { // 按板数遍历
				arrStk, ok := oneDay.MapStk[len(boardSlice)-pos]
				index := rowTmp
				if ok {
					if curDate == oneDay.Date && !IsMonFirstDay {
						f.SetCellValue("涨跌停复盘", GetExcelCol(replayRow+1, 12+len(boardSlice)-pos, time.Sunday), fmt.Sprintf("%d", len(arrStk)))
						// 算总晋级率不能包括当天吧
						// v, ok := advanceRate[len(boardSlice)-pos]
						// if !ok {
						// 	v = 0
						// }
						// advanceRate[len(boardSlice)-pos] = v + len(arrStk)
					}
					// 遍历股票
					for i, v := range arrStk {
						if len(boardSlice)-pos > 1 {
							if v.Perform == "大面" {
								tmpNum++
							}
							totalNum++
						}
						if i == 0 {
							// 共4支 晋级率16.00%
							f.SetCellValue(sheetName, GetExcelCol(index, col, oneDay.WeekNum), v.PreRate)
							f.MergeCell(sheetName, GetExcelCol(index, col, oneDay.WeekNum), GetExcelCol(index, col+2, oneDay.WeekNum))
							f.SetCellStyle(sheetName, GetExcelCol(index, col, oneDay.WeekNum), GetExcelCol(index, col+2, oneDay.WeekNum), styleId)
						}
						f.SetCellValue(sheetName, GetExcelCol(index+1, col, oneDay.WeekNum), v.Name)
						f.SetCellValue(sheetName, GetExcelCol(index+1, col+1, oneDay.WeekNum), v.Money)
						f.SetCellValue(sheetName, GetExcelCol(index+1, col+2, oneDay.WeekNum), v.Reason)
						if len(v.Perform) > 0 {
							f.SetCellValue(sheetName, GetExcelCol(index+1, col+3, oneDay.WeekNum), v.Perform)
						} else {
							if strings.HasPrefix(v.Code, "3") {
								f.SetCellValue(sheetName, GetExcelCol(index+1, col+3, oneDay.WeekNum), "创")
							} else if strings.HasPrefix(v.Code, "688") {
								f.SetCellValue(sheetName, GetExcelCol(index+1, col+3, oneDay.WeekNum), "科")
							}
						}

						f.SetCellStyle(sheetName, GetExcelCol(index+1, col, oneDay.WeekNum), GetExcelCol(index+1, col+3, oneDay.WeekNum), styleLeft)
						index++
					}
				}
				rowTmp += v
			}

			f.SetCellValue(sheetName, GetExcelCol(specialRow, col+3, oneDay.WeekNum), fmt.Sprintf("大面概率:%.f%%", float32(tmpNum)/float32(totalNum)*100))
			f.SetCellStyle(sheetName, GetExcelCol(specialRow, col+3, oneDay.WeekNum), GetExcelCol(specialRow, col+3, oneDay.WeekNum), styleLeft)

			if preDate == oneDay.Date && !IsMonFirstDay && totalNum > 0 { // 写前一日的大面概率
				str := fmt.Sprintf("%.2f%%", float32(tmpNum)/float32(totalNum)*100)
				f.SetCellValue("涨跌停复盘", GetExcelCol(replayRow+1, 12, time.Sunday), str)
			}
			if curDate == oneDay.Date && !IsMonFirstDay { // 写当日的连板数量
				f.SetCellValue("涨跌停复盘", GetExcelCol(replayRow+1, 8, time.Sunday), fmt.Sprintf("%d", totalNum))
				_, ok := currExcel.stkInfo[1]
				if ok {
					f.SetCellValue("涨跌停复盘", GetExcelCol(replayRow+1, 6, time.Sunday), fmt.Sprintf("%d", totalNum+len(currExcel.stkInfo[1])))
				}
			}

			rowTmp = weekRow // 调回一周开始的位置
		}

		// 调到下一周的行
		// [4 4 4 4 4 4 4 5 7 14 45]
		for _, v := range boardSlice {
			row += v
		}
		row += 3

		indexWeek++
		if indexWeek > len(monthInfo) {
			break MainLoop
		}
		// 切换到下一周
		weeker = monthInfo[indexWeek-1]
	}
	if !IsMonFirstDay {
		f.DeleteSheet("Sheet1")
		f.DeleteSheet("Sheet2")
	}

	// 计算1->2,2->3这种的总成功率
	// 获取前日数据
	startCol := 11
	for _, rows := range readSheet {
		for col, sheet := range rows {
			if col-startCol > 0 {
				v := advanceRate[col-startCol]
				temp, _ := strconv.Atoi(sheet)
				if temp > 0 {
					advanceRate[col-startCol] = temp + v
				}
			}
		}
	}

	for k, v := range advanceRate {
		num, ok := advanceRate[k+1]
		if ok {
			f.SetCellValue("涨跌停复盘", GetExcelCol(2, 12+k, time.Sunday), fmt.Sprintf("%d%%", int(math.Floor(float64(num)/float64(v)*100))))
		} else {
			f.SetCellValue("涨跌停复盘", GetExcelCol(2, 12+k, time.Sunday), "0.00%")
		}
	}
	if !IsMonFirstDay {
		f.SetCellValue("涨跌停复盘", GetExcelCol(replayRow+1, 1, time.Sunday), curDate)
	}

	f.SaveAs(excelPath)
	return true
}

func changeStkCode(code string) string {
	if strings.HasPrefix(code, "6") {
		return "sh" + code
	}
	return "sz" + code
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

// 将行列数转为excel里面特有的表格字符，如：第1行，第1列 ==> A1
// 特殊处理，按照星期几增加列数，这样就可以把每天的数据按照星期数固定在特定的列上
func GetExcelCol(row, col int, weekday time.Weekday) string {
	addcol := 0
	if weekday != time.Sunday && weekday != time.Saturday { // 周一到周五按照固定的列写入，防止节假日格式错乱的情况
		addcol = 4 * int(weekday-1)
	}
	str, _ := excelize.CoordinatesToCellName(col+addcol, row)
	// return fmt.Sprintf("%c%d", byte('A'+(col-1)%26), row)
	return str
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
