package configtask

import (
	"encoding/json"

	"localhost.com/my_web/repo/config"
	"localhost.com/my_web/repo/redis"
)

// 设置问题股json
func SetProbemJson() {
	// 获取问题股信息 proStk_s
	configStr, _ := redis.GetRedisInter().GetValue("problemstk")
	if configStr == "" {
		config.GetConfigInter().InitConfig(&proStk_s, "../config/problem.json")
		bytes, err := json.Marshal(proStk_s)
		if err == nil {
			redis.GetRedisInter().SetValue("problemstk", string(bytes))
		}
	} else {
		json.Unmarshal([]byte(configStr), &proStk_s)
	}
}
