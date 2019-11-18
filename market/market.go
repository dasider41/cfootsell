package market

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dasider41/cfootsell/db"
	"github.com/dasider41/cfootsell/models"
	"github.com/dasider41/cfootsell/util"
	"golang.org/x/net/html"
)

// UpdateMarket :
func UpdateMarket(cond models.SchCond) {
	baseURL, err := cond.GenerateURL()
	util.ErrCheck(err)
	// fmt.Println(baseURL)

	node, err := footsellMarketRequest(baseURL)
	util.ErrCheck(err)

	doc := goquery.NewDocumentFromNode(node)

	table := doc.Find("#list_table")
	table.Each(func(i int, item *goquery.Selection) {
		item.Find(".list_table_row").Each(func(j int, block *goquery.Selection) {
			p := getProduct(cond, block)
			if p.IsAvailable() {
				// fmt.Printf("%v\n", p)
				_, err := models.UpdateTransaction(p)
				util.ErrCheck(err)
			}
		})
	})
}

func footsellMarketRequest(url string) (*html.Node, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36")
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	parseBody, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	return parseBody, err
}

func getText(block *goquery.Selection, selector string) string {
	return strings.TrimSpace(block.Find(selector).Text())
}

func getProduct(cond models.SchCond, block *goquery.Selection) models.Product {
	var p models.Product

	p.Title = getText(block, ".list_subject_a")
	p.Condition = getText(block, "span.list_market_used")
	p.Size = cond.Size
	txtPrice := getText(block, ".list_market_price")
	price, _ := util.NumberOnly(txtPrice)
	p.Price = price
	p.Member = getText(block, "span.member")
	txtDate := getText(block, "span.list_table_dates")
	p.Updated = db.DateFormat(txtDate)

	return p
}
