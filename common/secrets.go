package common

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func CheckApiKey(event map[string]interface{}) error {
	queryParams, paramsOk := event["queryStringParameters"].(map[string]interface{})
	apiKey, apiKeyOk := queryParams["api_key"].(string)

	if !paramsOk || !apiKeyOk {
		return errors.New("api key missing")
	}

	key := os.Getenv("API_KEY_SECRETS_MANAGER_KEY")
	secretOutput, err := GetSecret(key)
	if err != nil {
		return err
	}

	if *secretOutput.SecretString != apiKey {
		return errors.New("incorrect api key")
	}

	return nil
}

func CheckSecret(event map[string]interface{}) error {
	queryParams, paramsOk := event["queryStringParameters"].(map[string]interface{})
	secret, secretOk := queryParams["secret"].(string)

	if !paramsOk || !secretOk {
		return errors.New("secret parameter missing")
	}

	number, numberOk := queryParams["number"].(string)
	sender, senderOk := queryParams["sender"].(string)

	if !numberOk && !senderOk {
		return errors.New("number or sender parameter is required")
	}

	if numberOk {
		sender = Hash(number)
	}

	prefix := os.Getenv("SESSION_SECRETS_MANAGER_PREFIX")
	secretOutput, err := GetSecret(prefix + sender)
	if err != nil {
		return err
	}

	if *secretOutput.SecretString != secret {
		return errors.New("incorrect secret")
	}

	return nil
}

func GetSecret(key string) (*secretsmanager.GetSecretValueOutput, error) {
	var secret *secretsmanager.GetSecretValueOutput
	var err error

	svc := secretsmanager.New(session.New())
	getInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}

	secret, err = svc.GetSecretValue(getInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == secretsmanager.ErrCodeResourceNotFoundException {
				return nil, nil
			}
		}

		return nil, err
	}

	return secret, nil
}

func UpdateSecret(key string, data []byte) error {
	svc := secretsmanager.New(session.New())
	getInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}

	secret, err := svc.GetSecretValue(getInput)
	if err != nil {
		return err
	} else if secret == nil {
		createInput := &secretsmanager.CreateSecretInput{
			Name:         aws.String(key),
			SecretBinary: data,
		}
		_, err := svc.CreateSecret(createInput)
		return err
	} else {
		updateInput := &secretsmanager.UpdateSecretInput{
			SecretId:     aws.String(key),
			SecretBinary: data,
		}
		_, updateErr := svc.UpdateSecret(updateInput)
		if updateErr != nil {
			return updateErr
		}

		return nil
	}
}

func CreateRandomSecret(key string) (string, error) {
	var secret string
	svc := secretsmanager.New(session.New())
	input := &secretsmanager.GetRandomPasswordInput{
		ExcludePunctuation: aws.Bool(true),
		IncludeSpace:       aws.Bool(false),
		PasswordLength:     aws.Int64(20),
	}

	generatedSecret, err := svc.GetRandomPassword(input)
	if err != nil {
		return secret, err
	}

	secret = *generatedSecret.RandomPassword
	createInput := &secretsmanager.CreateSecretInput{
		Name:         aws.String(key),
		SecretString: aws.String(secret),
	}

	_, err = svc.CreateSecret(createInput)

	return secret, err
}

func DeleteSecret(key string) error {
	svc := secretsmanager.New(session.New())
	deleteInput := &secretsmanager.DeleteSecretInput{
		SecretId: aws.String(key),
	}
	_, err := svc.DeleteSecret(deleteInput)
	if err != nil {
		return err
	}

	return nil
}
