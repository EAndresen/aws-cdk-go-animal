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

	//Create DynamoDb table
	table := awsdynamodb.NewTable(stack, jsii.String("AnimalTable"),
		&awsdynamodb.TableProps{
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String("ID"),
				Type: "STRING",
			},
			BillingMode: "PAY_PER_REQUEST",
			TableName:   jsii.String("Animals"),
			SortKey: &awsdynamodb.Attribute{
				Name: jsii.String("Name"),
				Type: "STRING",
			},
		})

	env := make(map[string]*string)
	env["DYNAMODB_TABLE"] = table.TableName()
	env["DYNAMODB_AWS_REGION"] = table.Env().Region

	//Create Lambda functions
	createAnimalFunction := createLambdaFunction(stack,"CreateAnimal", "create", "CreateAnimalFunction", &env)
	listAnimalFunction := createLambdaFunction(stack,"ListAnimals", "list", "ListAnimalFunction", &env)

	table.GrantReadData(listAnimalFunction)
	table.GrantReadWriteData(createAnimalFunction)

	//Create GraphQL API - AppSync
	api := awsappsync.NewGraphqlApi(stack, jsii.String("AnimalGraphQL"), &awsappsync.GraphqlApiProps{
		Name:   jsii.String("animals-graphql-api"),
		Schema: awsappsync.Schema_FromAsset(jsii.String("graphql/schema.graphql")),
	})

	listAnimalDS := api.AddLambdaDataSource(jsii.String("ListAnimalsLambdaResolver"), listAnimalFunction, &awsappsync.DataSourceOptions{
		Description: jsii.String("List Animals"),
		Name:        jsii.String("ListAnimal"),
	})
	listAnimalDS.CreateResolver(&awsappsync.BaseResolverProps{
		FieldName: jsii.String("listAnimals"),
		TypeName:  jsii.String("Query"),
	})

	createAnimalDS := api.AddLambdaDataSource(jsii.String("CreateAnimalsLambdaResolver"), createAnimalFunction, &awsappsync.DataSourceOptions{
		Description: jsii.String("Create Animal"),
		Name:        jsii.String("CreateAnimal"),
	})
	createAnimalDS.CreateResolver(&awsappsync.BaseResolverProps{
		FieldName:              jsii.String("createAnimal"),
		TypeName:               jsii.String("Mutation"),
		RequestMappingTemplate: awsappsync.MappingTemplate_FromString(
			jsii.String("{\n                    " +
				"\"version\": \"2018-05-29\",\n" +
				"\"operation\": \"Invoke\",\n" +
				"\"payload\": $util.toJson($context.arguments)\n" +
				"}")),
	})

	return stack
}

// createLambdaFunction creates an Lambda Function.
// Takes in the stacked to were the function should be deployed, name, id and handler for the function,
// and also if the any specific environments.
// Returns an awslambda Function.
func createLambdaFunction( stack awscdk.Stack, name, handler, id string, env *map[string]*string) awslambda.Function {
	return awslambda.NewFunction(stack, jsii.String(id), &awslambda.FunctionProps{
		Environment:  env,
		FunctionName: jsii.String(name),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(300)),
		Code:         awslambda.NewAssetCode(jsii.String("bin"), nil),
		Handler:      jsii.String(handler),
		Runtime:      awslambda.Runtime_GO_1_X(),
	})
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
