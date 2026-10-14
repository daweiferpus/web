package excelctl

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

// GetExcelInfo 取表格数据，结果为二维数组，表示行列
func GetExcelInfo(filepath, sheetName string) [][]string {
	// 打开一个现有的 Excel 文件
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	// 获取工作表名称
	// sheetName = f.GetSheetName(f.GetActiveSheetIndex())

	// 读取整个工作表的数据
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return rows
}

// GetSheetName 获取execel所有Sheet名称
func GetSheetName(filepath string) []string {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	sheetName := make([]string, 0)
	for i := 0; i < 20; i++ {
		name := f.GetSheetName(i)
		if name != "" {
			sheetName = append(sheetName, name)
		}
	}
	return sheetName
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
