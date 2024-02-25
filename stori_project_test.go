package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

func TestExampleAppStack(t *testing.T) {
	app := awscdk.NewApp(nil)
	stack := NewStoriProjectStack(app, "StoriProjectStack", nil)

	template := assertions.Template_FromStack(stack, nil)
	template.HasResourceProperties(jsii.String("AWS::S3::Bucket"), map[string]interface{}{})
	template.HasResourceProperties(jsii.String("AWS::Lambda::Function"), map[string]interface{}{
		"Runtime": awslambda.Runtime_GO_1_X().ToString(),
		"Timeout": 60,
		"Handler": "bootstrap",
	})

	template.HasResourceProperties(jsii.String("AWS::IAM::Policy"), map[string]interface{}{
		"PolicyDocument": map[string]interface{}{
			"Statement": []interface{}{
				// Policy statement for ses:SendEmail
				map[string]interface{}{
					"Action":   "ses:SendEmail",
					"Effect":   "Allow",
					"Resource": "*",
				},
				// Policy statement for s3:GetObject
				map[string]interface{}{
					"Action": "s3:GetObject",
					"Effect": "Allow",
				},
			},
		},
	})
}
