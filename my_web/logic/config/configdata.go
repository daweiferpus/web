package config

import "localhost.com/my_web/repo/log"

var (
	problemStock ProblemStock // 问题股信息

	sshqChan chan []SinaQuoteData // 实时行情

	bidChan chan SinaQuoteData

	svrConfig Config

	pageCount int = 0

	// 入口网址
	indexURL string = "https://vip.stock.finance.sina.com.cn/mkt/"

	// 获取涨幅榜前50的股票
	incPrcURL string = "https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?page=1&num=50&sort=changepercent&asc=0&node=hs_a&symbol=&_s_r_a=auto"

	// 获取某个市场板块的股票总数
	countURL string = "https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeStockCount?node=hs_a"

	// 分页
	pageURL string = "https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?page=%d&num=100&sort=changepercent&asc=0&node=hs_a&symbol=&_s_r_a=page"

	errLog log.LogInter

	problemTxtPath string = "../config/problem.txt"

	problemJSONPath string = "../config/problem.json"

	configJSONPath string = "../config/config.json"

	excelPath string = "../config/xuangu/隔日涨停股票分析.xlsx"

	bidfilePath string = "../config/bid_hq_calc.log"

	stockExcelPath string = "../config/xuangu/股票板块信息.xlsx"
)

func init() {
	errLog = log.GetLogInter()
	errLog.NewFileLogger("../log/error.log")

	sshqChan = make(chan []SinaQuoteData, 100)

	bidChan = make(chan SinaQuoteData, 6144)
}

type Config struct {
	OrderInfo  OrderInfo     `json:"OrderInfo"`
	ArrMoniter []MoniterInfo `json:"MoniterInfo"`
	FirStkInfo FirstStkInfo  `json:"FirstStkInfo"`
}

type FirstStkInfo struct {
	Reason      string     `json:"Reason"`   // 首板打板的板块信息
	ArrFirIsNot []FirIsNot `json:"FirIsNot"` // 首板不打的股票
}

type FirIsNot struct {
	Symbol  string `json:"Symbol"`
	Name    string `json:"Name"`
	IsOrder bool   `json:"IsOrder"`
}

type MoniterInfo struct {
	Symbol        string  `json:"Symbol"`
	Name          string  `json:"Name"`
	Changepercent float64 `json:"WarnPercent"`
	IsOrder       bool    `json:"IsOrder"`
	Amount        int     `json:"Amount"`
	Reason        string  `json:"Reason"`
	Activity      string  `json:"Activity"`
	AnchorName    string  `json:"AnchorName"`
	IsSealBack    bool    `json:"IsSealBack"` // 打回封
}

type OrderInfo struct {
	IsMoniter     bool    `json:"IsMoniter"`     // 监控是否弹出窗口
	IsFirstOrder  bool    `json:"IsFirstOrder"`  // 打首板委托开关
	FirstLimitNum int     `json:"FirstLimitNum"` // 首板数量 超过该数量之后的板不打
	IsSecondOrder bool    `json:"IsSecondOrder"` // 打二板委托开关
	IsFirst       bool    `json:"IsFirst"`       // 第一个上板的是否委托
	OrderMoney    float64 `json:"OrderMoney"`    // 委托金额
	IsMoreAmount  bool    `json:"IsMoreAmount"`  // 增量上板开关
	MaxAmount     float64 `json:"MaxAmount"`     // 最大成交量
	MaxNmc        float64 `json:"MaxNmc"`        // 最大流通市值
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

type StockPlm struct {
	Code string
	Name string
}

type ProblemStock struct {
	ArrStock []StockPlm
}
