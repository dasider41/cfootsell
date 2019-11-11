package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const footSellURL = "https://footsell.com"

type schCond struct {
	size    int
	keyword string
}

func main() {
	cond := schCond{
		size:    285,
		keyword: "고어",
	}

	baseURL, err := cond.getURL()
	errCheck(err)
	fmt.Println(baseURL)

	node, err := sendRequest(baseURL)
	errCheck(err)

	doc := goquery.NewDocumentFromNode(node)
	table := doc.Find("#list_table")
	table.Each(func(i int, item *goquery.Selection) {
		item.Find(".list_table_row").Each(func(j int, block *goquery.Selection) {
			title := getText(block, ".list_subject_a")
			condition := getText(block, "span.list_market_used")
			txtPrice := getText(block, ".list_market_price")
			price, _ := numberOnly(txtPrice)
			errCheck(err)
			member := getText(block, "span.member")
			txtDate := getText(block, "span.list_table_dates")
			date := sqlDateFormat(txtDate)
			fmt.Printf("%s, %s, %s, %d, %s\n", title, condition, member, price, date)
		})
	})
}

func (cond schCond) getURL() (string, error) {
	reqURL, err := url.Parse(footSellURL)
	if err != nil {
		return "", err
	}
	reqURL.Path = "g2/bbs/board.php"

	params := url.Values{}
	params.Add("size1", strconv.Itoa(cond.size))
	params.Add("stx", cond.keyword)
	params.Add("bo_table", "m51")
	params.Add("sfl", "wr_subject")

	reqURL.RawQuery = params.Encode()
	return reqURL.String(), nil
}

func sendRequest(url string) (*html.Node, error) {
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

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func numberOnly(text string) (int, error) {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(reg.ReplaceAllString(text, ""))
	if err != nil {
		return 0, err
	}
	return val, nil
}

func sqlDateFormat(tDate string) string {
	layoutIN := "06-01-02"
	layoutOUT := "2006-01-02"
	t, _ := time.Parse(layoutIN, tDate)
	return t.Format(layoutOUT)
}
