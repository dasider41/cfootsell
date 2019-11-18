package market

import (
	"fmt"
	"github.com/dasider41/cfootsell/email"
	"github.com/dasider41/cfootsell/models"
	"github.com/dasider41/cfootsell/util"
)

// Scrap :
func Scrap() {

	list, err := models.GetConditionList()
	util.ErrCheck(err)
	fmt.Println("Start cron task")
	for _, v := range list {
		// fmt.Printf("%v\n", v)
		UpdateMarket(v)
	}

	email.Notification()
	fmt.Println("End cron task")
}
