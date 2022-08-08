package main

import (
	"log"
	"os"
	"market/controllers"
	"market/database"
	"market/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8043"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveFromCart())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.GET("/cartcheckout", app.BuyFromCart())
	log.Fatal(router.Run(":" + port))
}
