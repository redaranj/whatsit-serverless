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
	Result string `json:"result"`
}

func handler(ctx context.Context, event map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	if err := common.CheckSecret(event); err != nil {
		return common.RespondUnauthorized(err)
	}

	queryParams, paramsOk := event["queryStringParameters"].(map[string]interface{})
	number, numberOk := queryParams["number"].(string)
	if !paramsOk || !numberOk {
		err := errors.New("'number' parameter is required")
		return common.RespondBadRequest(err)
	}

	prefix := os.Getenv("SESSION_SECRETS_MANAGER_PREFIX")
	sender := common.Hash(number)
	if err := common.DeleteSecret(prefix + sender); err != nil {
		return common.RespondServerError(err)
	}

	res := &deleteResponse{
		Result: "the number '" + number + "' is deleted",
	}

	return common.RespondSuccess(res)
}

func main() {
	lambda.Start(handler)
}
