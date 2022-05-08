import * as cdk from "aws-cdk-lib";
import * as ec2 from "aws-cdk-lib/aws-ec2";

export class VPCStack extends cdk.Stack {
  public readonly vpc: ec2.Vpc;

  public constructor(scope: cdk.App, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    this.vpc = new ec2.Vpc(this, "LambdaGQLVPC", {
      cidr: "10.0.0.0/16",
      natGateways: 0,
      maxAzs: 3, // add to each AZ within a given region
      subnetConfiguration: [
        {
          name: "public-subnet-1",
          subnetType: ec2.SubnetType.PUBLIC,
          cidrMask: 24,
        },
      ],
    });

    cdk.Tags.of(this.vpc).add(
      "source",
      "github.com/aidenwallis/api-gateway-cdk-graphql-test",
    );
  }
}
