package main

import (
	"bufio"
	"chatonetoone/configs/databases"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const WELCOME_MESSAGE = "<h1>welcome to the chat application<h1>"

var con *websocket.Conn
var err error

var MysqlDB *gorm.DB

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func main(){
	// r := mux.NewRouter()
	// r.HandleFunc("/",greetingMessage)
	// r.HandleFunc("/ws",upgradeToWebsocket)
	// log.Print("Listening on port 3000 ...")
	// if err := http.ListenAndServe(":3000",r); err != nil{
	// 	panic(err)
	// }

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
}

func greetingMessage(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(os.Stdout, "timestamp: %s | host: %s", time.Now().UTC().Format(time.RFC3339), req.Host)
	welcomemsg := "<h1>welcome to the chat application<h1>"
	w.WriteHeader(200)
	w.Write(([]byte)(welcomemsg))
}

func upgradeToWebsocket(w http.ResponseWriter, r *http.Request){
	w.Header().Add("status","101")
	con, err = upgrader.Upgrade(w,r,w.Header())
	if(err != nil){
		w.Write([]byte(err.Error()))
	}
	con.SetCloseHandler(
		func(code int, text string) error {
			log.Printf("connection closed: %d %s", code, text)
			return nil
		},
	)
	userinp := bufio.NewReader(os.Stdin)
	for{
		_, reader, commuerr := con.NextReader()
		if(commuerr != nil){
			log.Print(commuerr.Error())
			break
		}
		msg, _  := io.ReadAll(reader)
		log.Print((string)(msg))
		fmt.Println("Any thing you wanted to talk? If yes type a message:")
		fmt.Print("> ")
		inpmsg, _ := userinp.ReadString('\n')
		w,wrr := con.NextWriter(1)
		if(wrr != nil){
			log.Print(wrr.Error())
			break
		}
		if(inpmsg != ""){
			w.Write([]byte(inpmsg))
			w.Close()
		}
	}
	log.Print("Successfully upgraded to wesocket connecton!")
	w.WriteHeader(400)
}