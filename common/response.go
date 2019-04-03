package common

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type errorResponse struct {
	Result string `json:"result"`
}

func RespondSuccess(res interface{}) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(res)
	if err != nil {
		return RespondServerError(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func RespondBadRequest(err error) (events.APIGatewayProxyResponse, error) {
	log.Println("Bad request: ", err)

	return respondError(err, 400)
}

func RespondUnauthorized(err error) (events.APIGatewayProxyResponse, error) {
	log.Println("Unauthorized: ", err)

	return respondError(err, 401)
}

func RespondNotFound(err error) (events.APIGatewayProxyResponse, error) {
	log.Println("Not found: ", err)

	return respondError(err, 404)
}

func RespondServerError(err error) (events.APIGatewayProxyResponse, error) {
	log.Println("Server error: ", err)

	return respondError(err, 500)
}

func respondError(err error, code int) (events.APIGatewayProxyResponse, error) {
	res := &errorResponse{
		Result: err.Error(),
	}
	body, jsonErr := json.Marshal(res)
	if jsonErr != nil {
		return respondGenericError()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: code,
	}, err
}

func respondGenericError() (events.APIGatewayProxyResponse, error) {
	log.Println("Generic error")

	return events.APIGatewayProxyResponse{
		Body:       "{ \"result\": \"unknown error\" }",
		StatusCode: 500,
	}, nil
}
