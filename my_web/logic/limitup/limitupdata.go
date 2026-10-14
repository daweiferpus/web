package limitup

type response struct {
	Code    int         `json:"code"`    // 业务状态码（0表示成功）
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据（成功时为具体数据，失败时可忽略）
}

// 集合竞价表单
type limitFormat struct {
	Num           int    `json:"Day"`
	Code          string `json:"Code"`
	Name          string `json:"Name"`
	Percent       string `json:"Precent"`
	Trade         string `json:"Trade"`
	Reason        string `json:"Reason"`
	Limit         string `json:"Limit"`
	Nmc           string `json:"Nmc"`           // 流通市值
	Amount        string `json:"Amount"`        // 成交额
	Turnoverratio string `json:"Turnoverratio"` // 换手率
}

type stkReason struct {
	Symbol string `json:"symbol"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
	Limit  string `json:"limit"`
	Nmc    string `json:"nmc"`
}
