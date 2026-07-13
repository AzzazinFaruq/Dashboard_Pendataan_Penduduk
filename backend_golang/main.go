package main

import (
	"backend_golang/controllers"
	"backend_golang/middlewares"
	"backend_golang/setup"

	"github.com/gin-gonic/gin"
)

func main() {
	//Declare New Gin Route System
	router := gin.New()
	// Logger mencatat request; Recovery mencegah satu panic mematikan seluruh server.
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	//Run Database Setup
	setup.ConnectDatabase()

	//Auth
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	//User
	protected.POST("/logout", controllers.Logout)
	protected.GET("/user", controllers.GetCurrentUser)
	protected.GET("/user/all", controllers.GetAllUser)
	protected.PUT("/update/:id", controllers.UpdateUser)
	protected.PUT("/update/password/:id", controllers.PasswordUpdate)

	//Keluarga
	protected.GET("/keluarga", controllers.Index)
	protected.GET("/latestkel", controllers.Latest)
	protected.GET("/latestkelinput", controllers.LatestForInput)
	protected.GET("/keluarga/:id", controllers.GetKeluargaByID)
	protected.POST("/addkeluarga", controllers.AddKeluarga)
	protected.PUT("/editkeluarga/:id", controllers.UpdateKeluarga)
	protected.DELETE("/deletekeluarga/:id", controllers.DeleteKeluarga)

	//Penduduk
	protected.GET("/penduduk", controllers.GetPenduduk)
	protected.GET("/latestpend", controllers.GetLatestPenduduk)
	protected.GET("/penduduk/:id", controllers.GetPendudukByID)
	protected.POST("/addpenduduk", controllers.AddPenduduk)
	protected.PUT("/updatependuduk/:id", controllers.UpdatePenduduk)
	protected.DELETE("/deletependuduk/:id", controllers.DeletePenduduk)

	//Universal
	protected.GET("/alldata", controllers.AllData)
	protected.GET("/jumlah", controllers.DataCount)
	protected.GET("/alive", controllers.AliveCount)
	protected.GET("/marry", controllers.MarryCount)
	protected.GET("/gender", controllers.GenderCount)
	protected.GET("/data", controllers.RangeData)

	router.Static("/public", "./public")

	router.Run(":8080")
}
