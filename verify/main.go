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
	Result   string `json:"result"`
	NumberId string `json:"numberId"`
}

func handler(ctx context.Context, event map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	if err := common.CheckSecret(event); err != nil {
		return common.RespondUnauthorized(err)
	}

	queryParams, paramsOk := event["queryStringParameters"].(map[string]interface{})
	number, numberOk := queryParams["number"].(string)
	if !paramsOk || !numberOk {
		err := errors.New("number parameter missing")
		return common.RespondBadRequest(err)
	}

	numberId := common.Hash(number)
	prefix := os.Getenv("SESSION_SECRETS_MANAGER_PREFIX")
	if _, err := common.GetSecret(prefix + numberId); err != nil {
		return common.RespondServerError(err)
	}

	res := &verifyResponse{
		Result:   "the number '" + number + "' was previously registered and can send messages",
		NumberId: numberId,
	}

	return common.RespondSuccess(res)
}

func main() {
	lambda.Start(handler)
}
