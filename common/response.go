package common

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type errorResponse struct {
	Error string `json:"error"`
}

func RespondSuccess(res interface{}) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(res)

	if err != nil {
		return RespondError(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func RespondUnauthorized(err error) (events.APIGatewayProxyResponse, error) {
	log.Println("Unauthorized", err)

	res := &errorResponse{
		Error: err.Error(),
	}
	body, jsonErr := json.Marshal(res)

	if jsonErr != nil {
		return RespondError(jsonErr)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 401,
	}, err
}

func RespondError(err error) (events.APIGatewayProxyResponse, error) {
	log.Println("Server error", err)

	res := &errorResponse{
		Error: err.Error(),
	}
	body, jsonErr := json.Marshal(res)

	if jsonErr != nil {
		return respondGenericError()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 500,
	}, err
}

func respondGenericError() (events.APIGatewayProxyResponse, error) {
	log.Println("Generic error")

	return events.APIGatewayProxyResponse{
		Body:       "unknown error",
		StatusCode: 500,
	}, nil
}
