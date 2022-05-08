import * as cdk from "aws-cdk-lib";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as path from "path";

interface InfraStackProps extends cdk.StackProps {
  vpc: ec2.Vpc;
}

export class InfraStack extends cdk.Stack {
  constructor(scope: cdk.App, id: string, props: InfraStackProps) {
    super(scope, id, props);

    // setup the lambda for api gateway
    const lambdaService = new lambda.Function(this, "GQLLambda", {
      vpc: props.vpc,
      allowPublicSubnet: true,
      runtime: lambda.Runtime.GO_1_X,
      handler: "main",
      code: lambda.Code.fromAsset(path.join(__dirname, "../.."), {
        bundling: {
          image: lambda.Runtime.GO_1_X.bundlingImage,
          user: "root",
          environment: {
            CGO_ENABLED: "0",
            GOOS: "linux",
            GOARCH: "amd64",
          },
          command: [
            "bash",
            "-c",
            [
              "cd /asset-input",
              "go build -o /asset-output/main cmd/app/*.go",
            ].join(" && "),
          ],
        },
      }),
    });

    const api = new apigateway.LambdaRestApi(this, "GQLAPI", {
      handler: lambdaService,
    });

    // this is the /graphql endpoint proxy handler, all requests to /graphql are now handled by the lambda
    const gqlEndpoint = api.root.addResource("graphql", {
      // make it invoke the lambda when /graphql is hit
      defaultIntegration: new apigateway.LambdaIntegration(lambdaService),

      // handle cors requests
      defaultCorsPreflightOptions: {
        allowOrigins: apigateway.Cors.ALL_ORIGINS,
        allowHeaders: ["*"],
        allowMethods: ["METHOD", "POST", "OPTIONS"],
      },
    });

    // gql only accepts get or post, we can drill down our handler a little by only allowing api gateway to accept those.
    gqlEndpoint.addMethod("POST");
    gqlEndpoint.addMethod("GET");
  }
}
