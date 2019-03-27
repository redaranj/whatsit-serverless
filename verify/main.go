package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/redaranj/whatsit-serverless/common"
)

type verifyResponse struct {
	Message  string `json:"message"`
	NumberId string `json:"numberId"`
}

func handler(ctx context.Context, event map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	var err error
	err = common.CheckSecret(event)
	if err != nil {
		return common.RespondUnauthorized(err)
	}

	queryParams, paramsOk := event["queryStringParameters"].(map[string]interface{})
	number, numberOk := queryParams["number"].(string)
	if !paramsOk || !numberOk {
		err = errors.New("number parameter missing")
		return common.RespondError(err)
	}

	prefix := os.Getenv("SESSION_SECRETS_MANAGER_PREFIX")
	_, err = common.GetSecret(prefix + sender)
	if err != nil {
		return common.RespondError(err)
	}

	res := &verifyResponse{
		Message:  "number is registered",
		NumberId: numberId,
	}

	return common.RespondSuccess(res)
}

func main() {
	lambda.Start(handler)
}
