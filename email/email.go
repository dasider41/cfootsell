package email

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/dasider41/cfootsell/db"
	"github.com/dasider41/cfootsell/models"
	"github.com/dasider41/cfootsell/util"
	"github.com/joho/godotenv"
)

const (
	// Subject : Email title
	Subject = "Footsell market notification"
	// CharSet : Email charactor set
	CharSet = "UTF-8"
)

// Notification :
func Notification() {
	conn := db.InitDB()
	defer conn.Close()

	rows, err := conn.Query("SELECT `id`, "+
		"`title`, "+
		"`size`, "+
		"`condition`, "+
		"`price`, "+
		"`member`, "+
		"`updated` "+
		"FROM market "+
		"WHERE isNew = ?", true)
	util.ErrCheck(err)
	defer rows.Close()

	textBody := ""
	var arrNewID []int

	for rows.Next() {
		var id int
		var p models.Product
		err = rows.Scan(
			&id,
			&p.Title,
			&p.Size,
			&p.Condition,
			&p.Price,
			&p.Member,
			&p.Updated)

		textBody += fmt.Sprintf("%s, %d, %s, %d, %s, %s\n",
			p.Title,
			p.Size,
			p.Condition,
			p.Price,
			p.Member,
			p.Updated)
		arrNewID = append(arrNewID, id)
	}

	if len(arrNewID) > 0 {
		sendEmail(textBody)
		models.UpdateItemStatus(arrNewID)
	}
}

func sendEmail(TextBody string) {
	env, err := godotenv.Read()
	util.ErrCheck(err)

	sess, err := session.NewSession()
	util.ErrCheck(err)
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
	util.ErrCheck(err)
}
