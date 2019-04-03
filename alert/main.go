package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/redaranj/whatsit-serverless/common"
)

type alertResponse struct {
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

	sender, senderOk := postParams["sender"].(string)
	if !senderOk {
		err = errors.New("'sender' parameter is required")
		return common.RespondBadRequest(err)
	}

	number, numberOk := postParams["number"].(string)
	if !numberOk {
		err = errors.New("'number' parameter is required")
		return common.RespondBadRequest(err)
	}

	message, messageOk := postParams["message"].(string)
	if !messageOk {
		err = errors.New("'message' parameter is required")
		return common.RespondBadRequest(err)
	}

	if err = common.SendMessage(sender, number, message); err != nil {
		return common.RespondServerError(err)
	}

	res := &alertResponse{
		Result: "ok",
	}

	return common.RespondSuccess(res)
}

func main() {
	lambda.Start(handler)
}
