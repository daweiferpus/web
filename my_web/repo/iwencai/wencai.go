package iwencai

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// 涨停 沪深 非ST
func (wc *iWencaiReal) GetLimitup(info string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()
	url := "https://www.iwencai.com/customized/chart/get-robot-data"

	tmp := fmt.Sprintf(`{"source":"Ths_iwencai_Xuangu","version":"2.0","query_area":"","block_list":"","add_info":"{\"urp\":{\"scene\":1,\"company\":1,\"business\":1},\"contentType\":\"json\",\"searchInfo\":true}","question":"%s","perpage":"100","page":1,"secondary_intent":"stock","log_info":"{\"input_type\":\"typewrite\"}","rsh":"Ths_iwencai_Xuangu_91u8d1pyd682i75m322hjshyjwh8gwe6"}`, info)
	// 创建自定义请求体
	requestBody := []byte(tmp)

	// 创建请求对象
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))

	// 设置请求头
	req.Header = http.Header{
		"Accept":             []string{"application/json, text/plain, */*"},
		"Accept-Encoding":    []string{"gzip, deflate, br, zstd"},
		"Accept-Language":    []string{"zh-CN,zh;q=0.9"},
		"Cache-Control":      []string{"no-cache"},
		"Connection":         []string{"keep-alive"},
		"Content-Type":       []string{"application/json"},
		"Cookie":             []string{"other_uid=Ths_iwencai_Xuangu_91u8d1pyd682i75m322hjshyjwh8gwe6; ta_random_userid=c8b0ft7a1k; cid=4fd32da24bd5536dc9d0211e49b7f4111729211831; u_ukey=A10702B8689642C6BE607730E11E6E4A; u_uver=1.0.0; u_dpass=w7HNmrFjdWIdSn1FMYKcF7KuJXxp0dvW7w0kpinvBHZUq8xUSBHyp%2FQX8mK%2BvGRTHi80LrSsTFH9a%2B6rtRvqGg%3D%3D; u_did=7025FD0D9DDB4C01A2C08BE7FD827C3A; u_ttype=WEB; user_status=0; v=A005fPLGGSX_Hb0JcXrhlE1eXGLCKoPyC1_l0I_SisxTSWOcV3qRzJuu9bcc"},
		"hexin-v":            []string{"A005fPLGGSX_Hb0JcXrhlE1eXGLCKoPyC1_l0I_SisxTSWOcV3qRzJuu9bcc"},
		"Host":               []string{"www.iwencai.com"},
		"Origin":             []string{"https://www.iwencai.com"},
		"Pragma":             []string{"no-cache"},
		"Referer":            []string{"https://www.iwencai.com/unifiedwap/result?w=%E8%B7%8C%E5%81%9C%20%20%E6%B2%AA%E6%B7%B1%20%E9%9D%9EST&querytype=stock&addSign=1745369076507"},
		"Sec-Ch-Ua":          []string{`"Google Chrome";v="135", "Not-A.Brand";v="8", "Chromium";v="135"`},
		"Sec-Ch-Ua-Mobile":   []string{"?0"},
		"Sec-Ch-Ua-Platform": []string{`"Windows"`},
		"Sec-Fetch-Dest":     []string{"empty"},
		"Sec-Fetch-Mode":     []string{"cors"},
		"Sec-Fetch-Site":     []string{"same-origin"},
		"User-Agent":         []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"},
	}

	// 配置HTTP客户端
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic("请求失败: " + err.Error())
	}
	defer resp.Body.Close()

	// 根据Content-Encoding选择解压方式
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			panic(err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	// 读取响应
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		panic("读取响应失败: " + err.Error())
	}

	// 解析JSON响应
	response := ResponseData{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Unmarshal : %s\n", err.Error())
		return
	}

	// 输出结果
	fmt.Printf("响应状态: %s\n", resp.Status)
	// fmt.Printf("响应内容: %s\n", body)
}
