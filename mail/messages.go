//The mail package interact with Google GMail through Google API
package mail

import (
	"encoding/base64"
	"errors"
	"strings"

	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/alexivanenko/egroupware_notifier_bot/config"
	"golang.org/x/net/context"
	"google.golang.org/api/gmail/v1"
)

type Message struct {
	From       string
	Subject    string
	TaskNumber int
	Body       *MessageBody
}

type MessageBody struct {
	TrackingSystem string
	Priority       string
	Summary        string
}

//Get Messages returns the unread messages from GMail inbox.
//Messages filtered by specific strings declared in constants
//and compacted to the slice of Message structures.
func GetMessages(ctx context.Context) ([]Message, error) {
	srv, err := Connect(ctx)
	if err != nil {
		return nil, err
	}

	var filteredMessages []Message
	var pageToken string

	for {
		response, err := loadUnreadMessages(srv, pageToken)
		if err != nil {
			return nil, err
		}

		messages, err := filterMessages(srv, response)
		if err != nil {
			return nil, err
		}

		filteredMessages = append(filteredMessages, messages...)

		pageToken = response.NextPageToken

		if pageToken == "" {
			break
		}
	}

	return filteredMessages, nil
}

func loadUnreadMessages(srv *gmail.Service, pageToken string) (*gmail.ListMessagesResponse, error) {
	result, err := srv.Users.Messages.List(User).Q("is:unread").PageToken(pageToken).Do()

	return result, err
}

func filterMessages(srv *gmail.Service, listResponse *gmail.ListMessagesResponse) ([]Message, error) {

	var result []Message

	if len(listResponse.Messages) == 0 {
		return nil, errors.New("Empty messages list")
	}

	for _, shortMessage := range listResponse.Messages {
		fullMessage, err := srv.Users.Messages.Get(User, shortMessage.Id).Format("full").Do()
		if err != nil {
			continue
		}

		var body string
		var decodedBody []byte

		if fullMessage.Payload.MimeType == "text/plain" {

			decodedBody, err = base64.StdEncoding.DecodeString(fullMessage.Payload.Body.Data)

			if err == nil {
				body = string(decodedBody)
			}

		} else if strings.HasPrefix(fullMessage.Payload.MimeType, "multipart") {

			for _, part := range fullMessage.Payload.Parts {
				decodedBody, err = base64.URLEncoding.DecodeString(part.Body.Data)
				if err == nil {
					body += "\n" + string(decodedBody)
				}
			}

		}

		if body != "" &&
			(strings.Contains(body, config.String("mail_filters", "new_ticket")) || strings.Contains(body, config.String("mail_filters", "modified_ticket"))) {

			var subject, from string

			for _, header := range fullMessage.Payload.Headers {
				if header.Name == "Subject" {
					subject = header.Value
				} else if header.Name == "From" {
					from = header.Value
				}
			}

			message := Message{
				From:       from,
				Subject:    subject,
				TaskNumber: extractTaskNumber(subject),
				Body:       parseBody(body),
			}

			result = append(result, message)
		}

	}

	return result, nil
}

func extractTaskNumber(subject string) int {
	re := regexp.MustCompile("#(.*?): ")
	rm := re.FindStringSubmatch(subject)

	result := 0

	if len(rm) >= 2 {
		num, err := strconv.Atoi(rm[1])
		if err == nil {
			result = num
		}
	}

	return result
}

func parseBody(sourceBody string) *MessageBody {
	body := new(MessageBody)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader((sourceBody)))

	if err == nil {
		doc.Find("td").Each(func(i int, td *goquery.Selection) {

			text := td.Text()
			text = strings.TrimSpace(text)

			if text == "Tracking System" {
				body.TrackingSystem = td.Next().Text()
			} else if text == "Priority" {
				body.Priority = td.Next().Text()
			}

			if colspan, found := td.Attr("colspan"); found == true && colspan == "2" {
				style, found := td.Attr("style")

				if found && !strings.Contains(style, "font-weight:bold;") {
					body.Summary = td.Text()
				}
			}
		})
	}

	return body
}
