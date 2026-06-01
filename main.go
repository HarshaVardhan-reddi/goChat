package main

import (
	"chatonetoone/configs"
	"chatonetoone/configs/databases"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

var err error

var MysqlDB *gorm.DB

func main(){
	godotenv.Load()

	configloader := databases.MysqlConfigLoader{Environment: databases.Development}
	mysqlconfig, err := configloader.LoadMysqlConfiguration()

	if(err != nil){
		panic(err)
	}
	
	fmt.Println(mysqlconfig)

	MysqlDB = databases.IntializeGoOrm(*mysqlconfig)
	// fmt.Printf("%v",obj)
	fmt.Printf("%T\n",MysqlDB)
	fmt.Println(MysqlDB.Config)
	
	handler := configs.IntializeRoutes(MysqlDB)

	log.Print("Listening on port 3000 ...")
	if err := http.ListenAndServe(":3000", handler); err != nil{
		panic(err)
	}
}