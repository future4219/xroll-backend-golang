import * as cdk from "aws-cdk-lib";
import { Template } from "aws-cdk-lib/assertions";

import { VPCStack } from "../lib/vpcStack";

import { useVPCConfig } from "./config.test";

const app = new cdk.App();

const vpcConfig = useVPCConfig({});

const testPrefix = "test-env";

const testVpcStack = new VPCStack(app, `${testPrefix}-vpc-stack`, {
  prefix: testPrefix,
  config: vpcConfig,
});

const template = Template.fromStack(testVpcStack);

test("create vpc stack: VPC", () => {
  template.hasResourceProperties("AWS::EC2::VPC", {
    CidrBlock: "192.168.0.0/21",
    Tags: [
      {
        Key: "Name",
        Value: "test-env-vpc",
      },
      {
        Key: "project",
        Value: "test-env",
      },
    ],
  });
});

test("create vpc stack: Subnet", () => {
  template.resourceCountIs("AWS::EC2::Subnet", 4);
  template.hasResourceProperties("AWS::EC2::Subnet", {
    CidrBlock: "192.168.0.0/24",
    AvailabilityZone: "ap-northeast-1a",
    Tags: [
      {
        Key: "aws-cdk:subnet-name",
        Value: "test-env-subnet-public",
      },
      {
        Key: "aws-cdk:subnet-type",
        Value: "Public",
      },
      {
        Key: "Name",
        Value: "test-env-vpc-stack/test-env-vpc/test-env-subnet-publicSubnet1",
      },
      {
        Key: "project",
        Value: "test-env",
      },
    ],
  });
  template.hasResourceProperties("AWS::EC2::Subnet", {
    CidrBlock: "192.168.2.0/24",
    AvailabilityZone: "ap-northeast-1a",
    Tags: [
      {
        Key: "aws-cdk:subnet-name",
        Value: "test-env-subnet-private",
      },
      {
        Key: "aws-cdk:subnet-type",
        Value: "Isolated",
      },
      {
        Key: "Name",
        Value: "test-env-vpc-stack/test-env-vpc/test-env-subnet-privateSubnet1",
      },
      {
        Key: "project",
        Value: "test-env",
      },
    ],
  });
  template.hasResourceProperties("AWS::EC2::Subnet", {
    CidrBlock: "192.168.1.0/24",
    AvailabilityZone: "ap-northeast-1c",
    Tags: [
      {
        Key: "aws-cdk:subnet-name",
        Value: "test-env-subnet-public",
      },
      {
        Key: "aws-cdk:subnet-type",
        Value: "Public",
      },
      {
        Key: "Name",
        Value: "test-env-vpc-stack/test-env-vpc/test-env-subnet-publicSubnet2",
      },
      {
        Key: "project",
        Value: "test-env",
      },
    ],
  });
  template.hasResourceProperties("AWS::EC2::Subnet", {
    CidrBlock: "192.168.3.0/24",
    AvailabilityZone: "ap-northeast-1c",
    Tags: [
      {
        Key: "aws-cdk:subnet-name",
        Value: "test-env-subnet-private",
      },
      {
        Key: "aws-cdk:subnet-type",
        Value: "Isolated",
      },
      {
        Key: "Name",
        Value: "test-env-vpc-stack/test-env-vpc/test-env-subnet-privateSubnet2",
      },
      {
        Key: "project",
        Value: "test-env",
      },
    ],
  });
});
