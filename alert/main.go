package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/redaranj/whatsit-serverless/common"
)

type alertResponse struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, event map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	var err error
	err = common.CheckSecret(event)
	if err != nil {
		return common.RespondUnauthorized(err)
	}

	postParams, err := common.ParseJSONBody(event)
	if err != nil {
		return common.RespondError(err)
	}

	number, numberOk := postParams["number"].(string)
	message, messageOk := postParams["message"].(string)
	if !numberOk || !messageOk {
		err = errors.New("missing required parameter")
		return common.RespondError(err)
	}

	err = common.SendMessage(number, message)
	if err != nil {
		return common.RespondError(err)
	}

	res := &alertResponse{
		Message: "message sent",
	}

	return common.RespondSuccess(res)
}

func main() {
	lambda.Start(handler)
}
