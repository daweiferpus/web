package public

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/yanyiwu/gojieba"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	RDB      *redis.Client
	DB       *gorm.DB
	JB       *gojieba.Jieba
	VP       *viper.Viper
	LOG      *zap.SugaredLogger
	USECACHE bool
)
