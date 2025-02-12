package main

import (
	"fmt"
	"taskManager/router"
	"taskManager/templates"
)

func main() {
	err := templates.InitializeTemplates()
	if err != nil {
		fmt.Println(err)
	}
	router.StartServer()
}
