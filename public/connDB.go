package public

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	user := VP.GetString("database.user")
	pwd := VP.GetString("database.pwd")
	host := VP.GetString("database.host")
	port := VP.GetString("database.port")
	dbName := VP.GetString("database.dbName")
	timeout := VP.GetString("database.timeout")

	if user == "" || pwd == "" || host == "" || port == "" || dbName == "" {
		LOG.Fatal("viper hasn't get config param")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", user, pwd, host, port, dbName, timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		LOG.Fatal("connect with db failed ():" + err.Error())
	}
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	if DB == nil {
		return InitDB()
	}
	return DB
}
