package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

type DatabaseConfiguration struct{
	Development map[string]any `yaml:"development"`
	Test map[string]any `yaml:"test"`
	Production map[string]any `yaml:"production"`
}

// TODO: constants should be written for
// 1. env
// 2. drivr

// TODO: flags should be taken for
// 1. env
// 2. up / down

func main(){
	dbconfig := readDatabaseConfiguration()
	currentenv := dbconfig.Development // if flag is development
	dsn := fmt.Sprintf("%s:%stcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	currentenv["username"],currentenv["password"], currentenv["host"], currentenv["port"],currentenv["database"])
	db, err := sql.Open("mysql",dsn)
	if(err != nil){
		panic(err)
	}
	mysqldriver, errmysql := mysql.WithInstance(db,&mysql.Config{})
	if (errmysql != nil){
		panic(errmysql)
	}
	migrator, migerror := migrate.NewWithDatabaseInstance("file:///db/migrator/migrations","mysql",mysqldriver)
	if(migerror != nil){
		panic(migerror)
	}
	migrator.Up() // if flag is up
	migrator.Down() // if flag is down
}

func readDatabaseConfiguration() *DatabaseConfiguration {
	rawDbConfig, err := os.ReadFile("configs/databases/database.yml")
	dbconfig := DatabaseConfiguration{}
	if(err != nil){
		panic(err)
	}
	if errun := yaml.Unmarshal(rawDbConfig, dbconfig); errun != nil{
		panic(errun)
	}
	return &dbconfig
}
