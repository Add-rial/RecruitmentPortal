package main

import (
	"fmt"
	"recruitmentportal/config"
	"recruitmentportal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Test")

	config.ConnectDB()

	router := gin.Default()

	routes.SetupRoutes(router)
	router.Run()
}
