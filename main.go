package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/joho/godotenv"
	"golang.org/x/net/html"

	_ "github.com/go-sql-driver/mysql"
)

const footSellURL = "https://footsell.com"

type schCond struct {
	size    int
	keyword string
}

type product struct {
	title     string
	condition string
	price     int
	member    string
	size      int
	updated   string
}

const (
	// Subject : Email title
	Subject = "Footsell market notification"
	// CharSet : Email charactor set
	CharSet = "UTF-8"
)
const (
	// True : 1
	True = 1
	// False : 0
	False = 0
)

func main() {
	err := godotenv.Load()
	errCheck(err)
	list, err := getConditionList()
	errCheck(err)

	for _, v := range list {
		// fmt.Printf("%v\n", v)
		marketSearch(v)
	}

	alertNotification()
}

func getConditionList() ([]schCond, error) {
	conn := initDB()
	defer conn.Close()

	var list []schCond

	results, err := conn.Query("SELECT size, keyword  FROM conditions")
	if err != nil {
		return nil, err
	}

	for results.Next() {
		var cond schCond
		err = results.Scan(&cond.size, &cond.keyword)
		list = append(list, cond)
	}

	return list, nil
}

func marketSearch(cond schCond) {
	baseURL, err := cond.generateURL()
	errCheck(err)
	// fmt.Println(baseURL)

	node, err := footsellMarketRequest(baseURL)
	errCheck(err)

	doc := goquery.NewDocumentFromNode(node)

	table := doc.Find("#list_table")
	table.Each(func(i int, item *goquery.Selection) {
		item.Find(".list_table_row").Each(func(j int, block *goquery.Selection) {
			p := getProduct(cond, block)
			if p.isAvailable() {
				// fmt.Printf("%v\n", p)
				_, err := updateTransaction(p)
				errCheck(err)
			}
		})
	})
}

func getProduct(cond schCond, block *goquery.Selection) product {
	var p product

	p.title = getText(block, ".list_subject_a")
	p.condition = getText(block, "span.list_market_used")
	p.size = cond.size
	txtPrice := getText(block, ".list_market_price")
	price, _ := numberOnly(txtPrice)
	p.price = price
	p.member = getText(block, "span.member")
	txtDate := getText(block, "span.list_table_dates")
	p.updated = sqlDateFormat(txtDate)

	return p
}

func (cond schCond) generateURL() (string, error) {
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
	params.Add("sop", "and")
	params.Add("price1", "0")
	params.Add("price2", "0")

	reqURL.RawQuery = params.Encode()
	return reqURL.String(), nil
}

func (p product) isAvailable() bool {
	if len(p.member) > 0 {
		return true
	}
	return false
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

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
		// TODO :: Report an error by eamil
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
	t, err := time.Parse(layoutIN, tDate)

	if err != nil {
		return time.Now().Format(layoutOUT)
	}

	return t.Format(layoutOUT)
}

func initDB() *sql.DB {
	env, err := godotenv.Read()
	errCheck(err)

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		env["DB_USERNAME"],
		env["DB_PASSWORD"],
		env["DB_HOST"],
		env["DB_PORT"],
		env["DB_DATABASE"])
	db, err := sql.Open("mysql", dbConn)
	errCheck(err)
	return db
}

func updateTransaction(e product) (int64, error) {
	conn := initDB()
	defer conn.Close()

	stmInsertQuery, err := conn.Prepare("INSERT IGNORE INTO market SET " +
		"`title`=?," +
		"`condition`=?," +
		"`size`=?," +
		"`price`=?," +
		"`member`=?," +
		"`updated`=?")

	if err != nil {
		return 0, err
	}

	defer stmInsertQuery.Close()

	res, err := stmInsertQuery.Exec(e.title, e.condition, e.size, e.price, e.member, e.updated)
	errCheck(err)
	count, err := res.RowsAffected()
	return count, nil
}

func alertNotification() {
	conn := initDB()
	defer conn.Close()

	rows, err := conn.Query("SELECT `id`, "+
		"`title`, "+
		"`size`, "+
		"`condition`, "+
		"`price`, "+
		"`member`, "+
		"`updated` "+
		"FROM market "+
		"WHERE isNew = ?", True)
	errCheck(err)
	defer rows.Close()

	textBody := ""
	var arrNewID []int

	for rows.Next() {
		var id int
		var p product
		err = rows.Scan(
			&id,
			&p.title,
			&p.size,
			&p.condition,
			&p.price,
			&p.member,
			&p.updated)

		textBody += fmt.Sprintf("%s, %d, %s, %d, %s, %s\n",
			p.title,
			p.size,
			p.condition,
			p.price,
			p.member,
			p.updated)
		arrNewID = append(arrNewID, id)
	}

	if len(arrNewID) > 0 {
		sendEmail(textBody)
		updateItemStatus(arrNewID)
	}
}

func updateItemStatus(arr []int) {
	conn := initDB()
	defer conn.Close()

	stmUpdateQuery, err := conn.Prepare(
		"UPDATE market SET " +
			"`isNew` = ? " +
			" WHERE id = ?")
	errCheck(err)
	defer stmUpdateQuery.Close()

	for _, id := range arr {
		_, err := stmUpdateQuery.Exec(False, id)
		errCheck(err)
	}
}

func sendEmail(TextBody string) {
	env, err := godotenv.Read()
	errCheck(err)

	sess, err := session.NewSession()
	errCheck(err)
	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(env["RECIPIENT"]),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				// TODO:: change email template
				// Html: &ses.Content{
				// 	Charset: aws.String(CharSet),
				// 	Data:    aws.String(HtmlBody),
				// },
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(env["SENDER"]),
	}

	_, err = svc.SendEmail(input)
	errCheck(err)
}
