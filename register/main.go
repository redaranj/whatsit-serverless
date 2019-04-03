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
	Result   string `json:"result"`
	NumberId string `json:"numberId"`
	Secret   string `json:"secret"`
}

func handler(ctx context.Context, event map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	if err := common.CheckApiKey(event); err != nil {
		return common.RespondUnauthorized(err)
	}

	postParams, err := common.ParseJSONBody(event)
	if err != nil {
		return common.RespondBadRequest(err)
	}

	number, numberOk := postParams["number"].(string)
	if !numberOk {
		err = errors.New("'number' parameter is required")
		return common.RespondBadRequest(err)
	}

	email, emailOk := postParams["email"].(string)
	if !emailOk {
		err = errors.New("'email' parameter is required")
		return common.RespondBadRequest(err)
	}

	if err = common.SignIn(number, email); err != nil {
		common.RespondServerError(err)
	}

	secret, err := common.CreateRandomSecret()
	if err != nil {
		common.RespondServerError(err)
	}

	numberId := common.Hash(number)
	prefix := os.Getenv("SECRET_SECRETS_MANAGER_PREFIX")
	err = common.UpdateSecretString(prefix+numberId, secret)
	if err != nil {
		common.RespondServerError(err)
	}

	res := &registerResponse{
		Result:   "registration complete",
		NumberId: numberId,
		Secret:   secret,
	}

	return common.RespondSuccess(res)
}

func main() {
	lambda.Start(handler)
}
