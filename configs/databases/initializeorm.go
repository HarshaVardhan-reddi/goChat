package databases

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// "chatonetoone/configs/databases"

func IntializeGoOrm(config DatabaseConfig) (*gorm.DB) {
	// dsn - data source name
	pass := os.ExpandEnv(config.Password)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, pass, config.Host, config.Port, config.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if(err != nil){
		panic(err)
	}
	return db
}