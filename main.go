package main

import (
	"github.com/dasider41/cfootsell/market"
	"github.com/dasider41/cfootsell/models"
	"github.com/dasider41/cfootsell/util"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	util.ErrCheck(err)

	c := cron.New()
	c.AddFunc("@hourly", func() { market.Scrap() })
	c.Start()

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/market", models.FetchAllProduct)
		v1.GET("/conditions", models.FetchAllConditions)
	}

	r.Run()
}
