package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	code := req.QueryStringParameters["code"]
	clientID := os.Getenv("CLIENT_ID")
	clientKey := os.Getenv("CLIENT_SECRET")
	subject := "Message from Meetup Announcer Ouath API"
	toEmail := os.Getenv("EMAIL_TO_SEND")

	// Issue request to slack oath endpoint and get response:
	response, err := http.PostForm("https://slack.com/api/oauth.access", url.Values{
		"client_id":     {clientID},
		"client_secret": {clientKey},
		"code":          {code}})

	if err != nil {
		//handle postform error
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		fmt.Printf("Error issuing request: %v \n", err)
	}

	log.Println(string(body))

	// Send an email with oauth results
	emailClient := ses.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
	emailParams := &ses.SendEmailInput{
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String("client_id:" + clientID + "\n client_secret:" + clientKey + "\n code: " + code + "\n body: " + string(body)),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(toEmail)},
		},
		Source: aws.String(toEmail),
	}

	_, err = emailClient.SendEmail(emailParams)
	if err != nil {
		fmt.Printf("Error sending email with oauth key: %v \n", err)
	}

	// Return and redirect to Polyhack
	returnVal := events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers: map[string]string{
			"location": "http://polyhack.ca/",
		},
	}
	returnVal.Body = "Completed request successfully"
	return returnVal, nil
}

func main() {
	lambda.Start(Handler)
}
