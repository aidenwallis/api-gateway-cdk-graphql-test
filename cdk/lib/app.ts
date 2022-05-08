import * as cdk from "aws-cdk-lib";
import { InfraStack } from "./infra-stack";
import { VPCStack } from "./vpc-stack";

const app = new cdk.App();

const vpcStack = new VPCStack(app, "VPCStack", {
  stackName: "GraphQLVPCStack",
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
});

new InfraStack(app, "InfraStack", {
  vpc: vpcStack.vpc,
  stackName: "GraphQLInfraStack",
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
});
