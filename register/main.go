package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/redaranj/whatsit-serverless/common"
)

type registerResponse struct {
	Message  string `json:"message"`
	NumberId string `json:"numberId"`
	Secret   string `json:"secret"`
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
	email, emailOk := postParams["email"].(string)
	if !numberOk || !emailOk {
		err = errors.New("missing required parameter")
		return common.RespondError(err)
	}

	err = common.SignIn(number, email)
	if err != nil {
		common.RespondError(err)
	}

	numberId := common.Hash(number)
	prefix := os.Getenv("SESSION_SECRETS_MANAGER_PREFIX")
	secret, err := common.CreateRandomSecret(prefix + numberId)

	res := &registerResponse{
		Message:  "successfully registered " + number,
		NumberId: numberId,
		Secret:   secret,
	}

	return common.RespondSuccess(res)
}

func main() {
	lambda.Start(handler)
}
