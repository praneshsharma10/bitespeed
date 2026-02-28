package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pranesh/bitespeed/api"
	"github.com/pranesh/bitespeed/db"
	"github.com/pranesh/bitespeed/home"
)

func main() {
	home.LoadConfig()
	db.InitDB()

	r := gin.Default()

	r.GET(home.RoutePing, api.HandlePing)          // GET  /ping
	r.POST(home.RouteIdentify, api.HandleIdentify) // POST /identify

	r.Run(":" + home.AppConfig.Port)
}
