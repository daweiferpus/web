package public

import (
	"github.com/yanyiwu/gojieba"
)

var IgnoreWordSet map[string]interface{}

func InitJieba() *gojieba.Jieba {
	IgnoreWordSet = make(map[string]interface{}, 0)
	JB = gojieba.NewJieba()

	return JB
}

func GetJB() *gojieba.Jieba {
	if JB == nil {
		return InitJieba()
	} else {
		return JB
	}
}

func IsStopWord(str string) bool {
	_, ok := IgnoreWordSet[str]
	return ok
}

func AddStopWord(str string) bool {
	_, ok := IgnoreWordSet[str]
	if !ok {
		IgnoreWordSet[str] = nil
		return true
	}
	return false
}

// 系统启动的时候初始化词典
func InitDictionary(stopWords []string, finWords []string) {

	for _, word := range finWords {
		JB.AddWord(word)
	}

	for _, word := range stopWords {
		IgnoreWordSet[word] = nil
	}
}

func DelStopWord(str string) bool {
	_, ok := IgnoreWordSet[str]
	if ok {
		delete(IgnoreWordSet, str)
		return true
	}
	return false
}
