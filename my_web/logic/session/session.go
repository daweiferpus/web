package session

import (
	"strings"

	"github.com/gin-gonic/gin"
	"localhost.com/my_web/repo/redis"
)

func MakeRouter() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	redisClient := redis.GetRedisInter()
	router.GET("session/hq/", func(ctx *gin.Context) {
		reason := ctx.DefaultQuery("reason", "")
		mstkinfo := redisClient.GetHashAllValue("stkinfo", false)

		formats := make([]sessionFormat, 0)
		for k, v := range mstkinfo {
			if reason != "" && !strings.Contains(v, reason) {
				delete(mstkinfo, k)
			}
			format := sessionFormat{}
			format.Name = k
			format.Reason = v
			formats = append(formats, format)
		}

	})
}
