# Welcome to AWS CDK Go project-animal!

This is an example project for the aws-cdk-go IoC.

**NOTICE**: Go support is still in Developer Preview. This implies that APIs may
change while we address early feedback from the community. We would love to hear
about your experience through GitHub issues.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests


## NewApp
Initializes a CDK application. Returns an App interface:

A construct which represents an entire CDK app. This construct is normally the root of the construct tree.

You would normally define an `App` instance in your program's entrypoint, then define constructs where the app is used 
as the parent scope.

After all the child constructs are defined within the app, you should call `app.synth()` which will emit a 
"cloud assembly" from this app into the directory specified by `outdir`. Cloud assemblies includes artifacts such as
CloudFormation templates and assets that are needed to deploy this app into the AWS cloud. 
See: https://docs.aws.amazon.com/cdk/latest/guide/apps.html

## NewStack
Creating a new stack that is a root construct which represents a single CloudFormation stack.
You can add many stacks to the same App, though you have to specify which stack to deploy when running 
cdk deploy.

## NewTable - DynamoDB
Provides a DynamoDB table.
Creating a new DynamoDSB table to hold the animals.
Here is specified the `Partition` key that is mandatory, and also added a sort key for the name.
Also, the `TableName` and the `BillingType`

## LambdaFunctions - Lambda
### `createLambdaFunction`:
The `createLambdaFunction` is a helper function that's return an interface of a Lambda Function, called Function.
It takes the `stack`, `name`, `hanler`, `id` and a key-value map of environment variables needed for the lambda
functions. Then it takes the premade binary lambda file corresponding to the name prompt in the header.
Deploys a file from inside the construct library as a function.
The supplied file is subject to the 4096 bytes limit of being embedded in a CloudFormation template.
The construct includes an associated role with the lambda.
This construct does not yet reproduce all features from the underlying resource library.

##NewGraphQL - AppSync
An AppSync GraphQL API.
Here we specify the name, and the GraphQL Schema located in a separate folder.
### Lambdas
Adding the Lambda data sources for the GraphQL API.
### Resolvers
Here we are adding a resolver to the data source. This is connecting the API to the database.
And the `FieldNames` here need to match the names in the `GraphQL schema`.

## Makefile
This helper file is to create and binary of the `go` lambda function and place in a separate bin folder. 