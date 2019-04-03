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

	postParams, err := common.ParseJSONBody(event)
	if err != nil {
		return common.RespondBadRequest(err)
	}

	number, numberOk := postParams["number"].(string)
	if !numberOk {
		err := errors.New("'number' parameter is required")
		return common.RespondBadRequest(err)
	}

	sender := common.Hash(number)

	prefix := os.Getenv("SESSION_SECRETS_MANAGER_PREFIX")
	if err := common.DeleteSecret(prefix + sender); err != nil {
		return common.RespondServerError(err)
	}

	prefix = os.Getenv("SECRET_SECRETS_MANAGER_PREFIX")
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
