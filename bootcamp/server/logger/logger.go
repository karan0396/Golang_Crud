package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var GeneralLogger *log.Logger

var ErrorLogger *log.Logger

func init(){
	absPath,err := filepath.Abs("../")
	if err!=nil{
		fmt.Println("Error reading given path:",err)
	}

	generalLog,err:=os.OpenFile(absPath+"/log.log",os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
	if err!=nil{
		fmt.Println("Error Opening file:",err)
		os.Exit(1)
}

GeneralLogger=log.New(generalLog,"General Logger:\t",log.Ldate|log.Ltime|log.Lshortfile)
ErrorLogger=log.New(generalLog,"ErrorLogger:\t",log.Ldate|log.Ltime|log.Lshortfile)
}