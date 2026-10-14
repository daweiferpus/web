package configtask

var (
	proStk_s ProblemStock // 问题股信息
)

func init() {
	proStk_s = ProblemStock{}
}

type StockPlm struct {
	Code string
	Name string
}
type ProblemStock struct {
	ArrStock []StockPlm
}
