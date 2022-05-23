package config

import (
	"bootcamp/logger"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func GetDb()*sql.DB{
	return Db

}

func Dbinit(){
	data,err:=sql.Open("mysql","root:123456789@tcp(localhost:3306)/employee?charset=utf8")
	if err!=nil{
		logger.GeneralLogger.Fatalln(err)
		return
	}
	// defer Db.Close()

	if err := data.Ping(); err != nil {
		logger.GeneralLogger.Fatal(err)
	  }

Db=data

fmt.Println("Database is Connected")

}