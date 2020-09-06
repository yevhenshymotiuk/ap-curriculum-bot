package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yevhenshymotiuk/ap-curriculum-bot/curriculum"
	"github.com/yevhenshymotiuk/telegram-lambda-helpers/apigateway"
)

func getObjectFromS3Bucket(
	bucketName string,
	objectName string,
) *s3.GetObjectOutput {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String("eu-north-1")})

	client := s3.New(sess)

	resp, err := client.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectName),
		},
	)

	if err != nil {
		log.Fatalf("Unable to get file %q, %v", objectName, err)
	}

	return resp
}

func handler(
	request events.APIGatewayProxyRequest,
) (apigateway.Response, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		return apigateway.Response404, err
	}

	update := tgbotapi.Update{}

	err = json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		return apigateway.Response404, err
	}

	assetsBucket := os.Getenv("ASSETS_BUCKET")
	curriculumFile := os.Getenv("CURRICULUM_FILE")

	message := update.Message
	var responseMessageText string

	switch message.Command() {
	case "today":
		resp := getObjectFromS3Bucket(assetsBucket, curriculumFile)
		w, err := curriculum.NewWeek(io.Reader(resp.Body))
		if err != nil {
			return apigateway.Response404, err
		}

		responseMessageText = curriculum.Today(*w).Format()
		log.Println(responseMessageText)
	default:
		responseMessageText = `¯\_(ツ)_/¯`
	}

	responseMessage := tgbotapi.NewMessage(message.Chat.ID, responseMessageText)
	bot.Send(responseMessage)

	return apigateway.Response200, nil
}

func main() {
	lambda.Start(handler)
}