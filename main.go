package main

import (
	"github.com/dasider41/cfootsell/email"
	"github.com/dasider41/cfootsell/market"
	"github.com/dasider41/cfootsell/models"
	"github.com/dasider41/cfootsell/util"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	util.ErrCheck(err)
	list, err := models.GetConditionList()
	util.ErrCheck(err)

	for _, v := range list {
		// fmt.Printf("%v\n", v)
		market.Update(v)
	}

	email.Notification()
}
