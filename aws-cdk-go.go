package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsappsync"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type AwsLambdaCronStackProps struct {
	awscdk.StackProps
}

func NewLambdaCronStack(scope constructs.Construct, id string, props *AwsLambdaCronStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("AnimalTable"),
		&awsdynamodb.TableProps{
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String("UserId"),
				Type: "STRING",
			},
			BillingMode: "PAY_PER_REQUEST",
		})

	env := make(map[string]*string)
	env["DYNAMODB_TABLE"] = table.TableName()
	env["DYNAMODB_AWS_REGION"] = table.Env().Region

	// The code that defines your stack goes here

	createAnimalFunction := awslambda.NewFunction(stack, jsii.String("CreateAnimalFunction"), &awslambda.FunctionProps{
		Code:        awslambda.NewAssetCode(jsii.String("lambda/create"), nil),
		Handler:     jsii.String("main"),
		Timeout:     awscdk.Duration_Seconds(jsii.Number(300)),
		Runtime:     awslambda.Runtime_GO_1_X(),
		Environment: &env,
	})

	listAnimalFunction := awslambda.NewFunction(stack, jsii.String("ListAnimalFunction"), &awslambda.FunctionProps{
		Code:        awslambda.NewAssetCode(jsii.String("lambda/list"), nil),
		Handler:     jsii.String("main.go"),
		Timeout:     awscdk.Duration_Seconds(jsii.Number(300)),
		Runtime:     awslambda.Runtime_GO_1_X(),
		Environment: &env,
	})

	table.GrantReadData(listAnimalFunction)
	table.GrantReadWriteData(createAnimalFunction)

	api := awsappsync.NewGraphqlApi(stack, jsii.String("AnimalGrapghQL"), &awsappsync.GraphqlApiProps{
		Name:   jsii.String("animals-graphql-api"),
		Schema: awsappsync.Schema_FromAsset(jsii.String("./graphql/schema.graphql")),
	})

	api.AddLambdaDataSource(jsii.String("ListAnimalsLambdaResolver"), listAnimalFunction, &awsappsync.DataSourceOptions{
		Description: jsii.String("List Animals"),
		Name:        jsii.String("ListAnimal"),
	})

	api.AddLambdaDataSource(jsii.String("CreateAnimalsLambdaResolver"), createAnimalFunction, &awsappsync.DataSourceOptions{
		Description: jsii.String("Create Animal"),
		Name:        jsii.String("CreateAnimal"),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewLambdaCronStack(app, "AwsLambdaStack", &AwsLambdaCronStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
