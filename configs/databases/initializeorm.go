package databases

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// "chatonetoone/configs/databases"

func IntializeGoOrm(config DatabaseConfig) (*gorm.DB) {
	// dsn - data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if(err != nil){
		panic(err)
	}
	return db
}