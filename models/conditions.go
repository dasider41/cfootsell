package models

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/dasider41/cfootsell/db"
	"github.com/gin-gonic/gin"
)

const footSellURL = "https://footsell.com"

// SchCond :
type SchCond struct {
	Size    int
	Keyword string
}

// GenerateURL :
func (cond SchCond) GenerateURL() (string, error) {
	reqURL, err := url.Parse(footSellURL)
	if err != nil {
		return "", err
	}
	reqURL.Path = "g2/bbs/board.php"

	params := url.Values{}
	params.Add("size1", strconv.Itoa(cond.Size))
	params.Add("stx", cond.Keyword)
	params.Add("bo_table", "m51")
	params.Add("sfl", "wr_subject")
	params.Add("sop", "and")
	params.Add("price1", "0")
	params.Add("price2", "0")

	reqURL.RawQuery = params.Encode()
	return reqURL.String(), nil
}

// GetConditionList :
func GetConditionList() ([]SchCond, error) {
	conn := db.InitDB()
	defer conn.Close()

	var list []SchCond

	results, err := conn.Query("SELECT size, keyword  FROM conditions")
	if err != nil {
		return nil, err
	}

	for results.Next() {
		var cond SchCond
		err = results.Scan(&cond.Size, &cond.Keyword)
		list = append(list, cond)
	}

	return list, nil
}

// FetchAllConditions :
func FetchAllConditions(c *gin.Context) {
	conn := db.InitDB()
	defer conn.Close()

	list, err := GetConditionList()

	if err != nil {
		c.JSON(200, gin.H{
			"status":  http.StatusNotFound,
			"message": "Not found.",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  http.StatusOK,
		"message": list,
	})
}
