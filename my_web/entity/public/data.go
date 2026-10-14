/*
贯穿整个服务的实体数据结构，类似常量，错误码
贫血模型仅包含基础的数据结构读写方法
*/
package public

import (
	"math"
)

const (
	LOGGER_LEVEL_ERROR = iota
	LOGGER_LEVEL_INFO
	LOGGER_LEVEL_NOPRE
)

var (
	// 新浪入口网址
	IndexUrl string = "https://vip.stock.finance.sina.com.cn/mkt/"

	// 分页
	PageUrl string = "https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?page=%d&num=%s&sort=changepercent&asc=0&node=hs_a&_s_r_a=page"

	// 获取A股市场板块的股票总数
	CountUrl string = "https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeStockCount?node=hs_a"

	// 按涨幅顺序获取行情
	IncPrcUrl string = "https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?page=1&num=%s&sort=changepercent&asc=0&node=hs_a&symbol=&_s_r_a=auto"

	// 查询指定股票数据（单支或多支,','分隔）
	SomeStkUrl string = "http://qt.gtimg.cn/q=%s"
)

// 浮点数比较 小于或等于
func LessOrEqual(a, b, epsilon float64) bool {
	return a < b || math.Abs(a-b) <= epsilon
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
