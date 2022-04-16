package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetSecret(secretName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("could not load aws configurtation: %w", err)
	}

	client := ssm.NewFromConfig(cfg, func(o *ssm.Options){
		o.Region = Region
	})
	param, err := client.GetParameter(context.TODO(),&ssm.GetParameterInput{Name: &secretName,WithDecryption: true})
	if err != nil {
		return "", fmt.Errorf("could not get the parameter from AWS SSM store: %w", err)
	}

	return *param.Parameter.Value, nil
}

