package session

type response struct {
	Code    int         `json:"code"`    // 业务状态码（0表示成功）
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据（成功时为具体数据，失败时可忽略）
}

type sessionFormat struct {
	Code          string // 股票代码
	Name          string // 股票名称
	Percent       string // 涨幅
	Price         string // 当前价
	Reason        string // 涨停原因
	Limit         string // 累计涨停次数
	Nmc           string // 流通市值
	Amount        string // 成交额
	Turnoverratio string // 换手率
}
