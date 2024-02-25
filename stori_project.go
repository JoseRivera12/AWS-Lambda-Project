package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type StoriProjectStackProps struct {
	awscdk.StackProps
}

func NewStoriProjectStack(scope constructs.Construct, id string, props *StoriProjectStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// S3 Bucket
	s3 := awss3.NewBucket(stack, jsii.String("StoriTransactions"), &awss3.BucketProps{})

	// Lambda function
	fn := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("FileProcessing"), &awscdklambdagoalpha.GoFunctionProps{
		Environment: &map[string]*string{
			"DATABASE_NAME":     jsii.String("db_name"),
			"DATABASE_USER":     jsii.String("db_user"),
			"DATABASE_HOST":     jsii.String("db_host"),
			"DATABASE_PORT":     jsii.String("db_port"),
			"DATABASE_PASSWORD": jsii.String("db_password"),
			"SES_EMAIL":         jsii.String("db_email"),
		},
		Runtime: awslambda.Runtime_GO_1_X(),
		Entry:   jsii.String("./lambda-handler/cmd/"),
		Timeout: awscdk.Duration_Seconds(jsii.Number(60)),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: jsii.Strings(`-ldflags "-s -w"`),
		},
	})

	// Policies
	fn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ses:SendEmail"),
		Resources: jsii.Strings("*"),
	}))

	fn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: jsii.Strings("s3:GetObject"),
		Resources: jsii.Strings(
			"arn:aws:s3:::" + *s3.BucketName() + "/*",
		),
	}))

	// Events
	notification := awss3notifications.NewLambdaDestination(fn)
	s3.AddEventNotification(awss3.EventType_OBJECT_CREATED, notification)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewStoriProjectStack(app, "StoriProjectStack", &StoriProjectStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
