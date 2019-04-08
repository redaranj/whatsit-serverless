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

	postParams, err := ParseJSONBody(event)
	if err != nil {
		return err
	}

	var sender string
	sender, senderOk := postParams["sender"].(string)
	if !senderOk || sender == "" {
		number, numberOk := postParams["number"].(string)
		if !numberOk || number == "" {
			return errors.New("'number' or 'sender' parameter is required")
		} else {
			sender = Hash(number)
		}
	}

	prefix := os.Getenv("SECRET_SECRETS_MANAGER_PREFIX")
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
	svc := secretsmanager.New(session.New())
	getInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}

	secret, err := svc.GetSecretValue(getInput)
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

func UpdateSecretString(key string, value string) error {
	secret, err := GetSecret(key)
	if err != nil {
		return err
	} else if secret == nil {
		svc := secretsmanager.New(session.New())
		createInput := &secretsmanager.CreateSecretInput{
			Name:         aws.String(key),
			SecretString: aws.String(value),
			KmsKeyId:     aws.String(os.Getenv("KMS_KEY_ARN")),
		}
		_, err := svc.CreateSecret(createInput)
		return err
	} else {
		svc := secretsmanager.New(session.New())
		updateInput := &secretsmanager.UpdateSecretInput{
			SecretId:     aws.String(key),
			SecretString: aws.String(value),
		}

		if _, updateErr := svc.UpdateSecret(updateInput); updateErr != nil {
			return updateErr
		}

		return nil
	}
}

func UpdateSecretBinary(key string, data []byte) error {
	secret, err := GetSecret(key)

	svc := secretsmanager.New(session.New())

	if err != nil {
		return err
	} else if secret == nil {
		createInput := &secretsmanager.CreateSecretInput{
			Name:         aws.String(key),
			SecretBinary: data,
			KmsKeyId:     aws.String(os.Getenv("KMS_KEY_ARN")),
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

func CreateRandomSecret() (string, error) {
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

	return secret, nil
}

func DeleteSecret(key string) error {
	svc := secretsmanager.New(session.New())
	deleteInput := &secretsmanager.DeleteSecretInput{
		SecretId:                   aws.String(key),
		ForceDeleteWithoutRecovery: aws.Bool(true),
	}
	_, err := svc.DeleteSecret(deleteInput)
	if err != nil {
		return err
	}

	return nil
}
