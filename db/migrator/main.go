package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConfiguration struct{
	Development map[string]any `yaml:"development"`
	Test map[string]any `yaml:"test"`
	Production map[string]any `yaml:"production"`
}

// required flags for environemt
const(
	DEVELOPMENT int = iota + 1
	TEST
	PRODUCTION
)

func main(){
	dbconfig := readDatabaseConfiguration()
	flags := takeInputFlags()
	envVal := *flags["environment"].(*int)
	currentenv := getAppropriateDbConfig(dbconfig, envVal)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	currentenv["username"],currentenv["password"], currentenv["host"], currentenv["port"],currentenv["database"])
	db, err := sql.Open("mysql",dsn)
	if(err != nil){
		panic(err)
	}
	mysqldriver, errmysql := mysql.WithInstance(db,&mysql.Config{})
	if (errmysql != nil){
		panic(errmysql)
	}
	migrator, migerror := migrate.NewWithDatabaseInstance("file://db/migrator/migrations","mysql",mysqldriver)
	if(migerror != nil){
		panic(migerror)
	}

	dir := *flags["direction"].(*string)
	if(dir == "up"){
		migrator.Up()
	}
	if(dir == "down"){
		migrator.Down()
	}
}

func readDatabaseConfiguration() *DatabaseConfiguration {
	rawDbConfig, err := os.ReadFile("configs/databases/database.yml")
	dbconfig := DatabaseConfiguration{}
	if(err != nil){
		panic(err)
	}
	if errun := yaml.Unmarshal(rawDbConfig, &dbconfig); errun != nil{
		panic(errun)
	}
	return &dbconfig
}

func takeInputFlags() map[string]any {
	result := make(map[string]any)
	env := flag.Int("env",1,"--env=1 or --env 1 is for selecting the right database environment for performing migrations. possible environemnts are 1 -> development, 2 -> test and 3 -> production")
	direction := flag.String("direction","up","direction of the migration wether it has to be up or down. defaults to up if non provided")
	flag.Parse()
	result["environment"] = env
	result["direction"] = direction
	return result
}

func getAppropriateDbConfig(DBConfig *DatabaseConfiguration, environment int) map[string]any {
	var dbconfig map[string]any
	switch environment {
	case DEVELOPMENT:
		dbconfig = DBConfig.Development
	case TEST:
		dbconfig = DBConfig.Test
	case PRODUCTION:
		dbconfig = DBConfig.Production
	}
	return dbconfig
}
