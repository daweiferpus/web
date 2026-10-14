package iwencai

func GetIWencaiInter() *iWencaiReal {
	return &iWencaiReal{}
}

type IWencaiInter struct {
	GetLimitup (string)
}

type iWencaiReal struct {
}

// 定义返回数据的结构体
type ResponseData struct {
	Status_code int    `json:"status_code"`
	Data        Data   `json:"data"`
	Status_msg  string `json:"status_msg"`
	Logid       string `json:"logid"`
}

type Data struct {
	Status_msg string   `json:"status_msg"`
	Logid      string   `json:"logid"`
	Question   string   `json:"question"`
	Answers    []answer `json:"answer"`
}

type answer struct {
	Txts []txt `json:"txt"`
}
type txt struct {
	Contents content `json:"content"`
}
type content struct {
	Status_code string       `json:"string"`
	Components  []components `json:"components"`
}
type components struct {
	Data data `json:"data"`
}

type data struct {
	Datas []datas `json:"datas"`
}
type datas struct {
	Code    string `json:"code"`
	Amount  string `json:"涨停封单量[20250422]"`
	Price   string `json:"最新价"`
	Percent string `json:"最新涨跌幅"`
	Num     int    `json:"连续涨停天数[20250422]"`
	Reason  string `json:"涨停原因类别[20250422]"`
}
