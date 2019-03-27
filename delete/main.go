package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/redaranj/whatsit-serverless/common"
)

type deleteResponse struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, event map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	var err error
	err = common.CheckSecret(event)
	if err != nil {
		return common.RespondUnauthorized(err)
	}

	queryParams, paramsOk := event["queryStringParameters"].(map[string]interface{})
	sender, senderOk := queryParams["sender"].(string)
	if !paramsOk || !senderOk {
		err = errors.New("sender parameter missing")
		return common.RespondError(err)
	}

	prefix := os.Getenv("SESSION_SECRETS_MANAGER_PREFIX")
	err = common.DeleteSecret(prefix + sender)
	if err != nil {
		return common.RespondError(err)
	}

	res := &deleteResponse{
		Message: "number deleted",
	}

	return common.RespondSuccess(res)
}

func main() {
	lambda.Start(handler)
}
