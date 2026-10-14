package filectl

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func GetFileInfo(filePath string) []string {
	file, err := os.Open(filePath) // 打开文件
	if err != nil {
		fmt.Println("open file err:", err.Error())
		return nil
	}
	// 处理结束后关闭文件
	defer file.Close()

	// 使用bufio读取
	r := bufio.NewReader(file)
	fileSlice := make([]string, 0)
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
		// data = strings.TrimSpace(data)
		fileSlice = append(fileSlice, data)
	}
	return fileSlice
}

func WriteFileInfo(filePath string, writeSlice []string) {
	if writeSlice == nil {
		fmt.Println("writeSlice 没有初始化")
		return
	}
	// 创建或打开一个文件用于写入
	fileW, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer fileW.Close()
	// 创建一个bufio.Writer实例，提高写入效率
	writer := bufio.NewWriter(fileW)
	for _, v := range writeSlice {
		if _, err := writer.WriteString(v + "\n"); err != nil {
			panic(err)
		}
	}
	// 刷新缓冲区，确保所有数据被写入文件
	if err := writer.Flush(); err != nil {
		panic(err)
	}
}
