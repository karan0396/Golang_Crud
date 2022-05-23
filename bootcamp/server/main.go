package main

import (
	// "bootcamp/controller"
	"bootcamp/controller"
	"bootcamp/route"
)

func main(){
	
	r:=route.Route()
	r.Use(controller.Cors)
	
}