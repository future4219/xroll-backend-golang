import * as cdk from "aws-cdk-lib";
import { Match, Template } from "aws-cdk-lib/assertions";

import { BuildConfig, EC2Config } from "../lib/buildConfig";
import { EC2BastionStack } from "../lib/ec2BastionStack";
import { VPCStack } from "../lib/vpcStack";

import {
  testPrefix,
  useAlarmConfig,
  useConfig,
  useEC2Config,
  useECSConfig,
  useRDSConfig,
  useVPCConfig,
} from "./config.test";

const createTemplate = (config: EC2Config) => {
  const app = new cdk.App();

  const testConfig: BuildConfig = useConfig(
    useVPCConfig({}),
    useRDSConfig({}),
    useECSConfig({}),
    config,
    useAlarmConfig({}),
  );

  const testVpcStack = new VPCStack(app, `${testPrefix}-vpc-stack`, {
    prefix: testPrefix,
    config: testConfig.vpc,
  });

  const testEc2Stack = new EC2BastionStack(app, `${testPrefix}-ec2-stack`, {
    prefix: testPrefix,
    vpc: testVpcStack.vpc,
    config: testConfig.ec2,
  });

  return Template.fromStack(testEc2Stack);
};

test("create ec2 stack: instance", () => {
  const template = createTemplate(useEC2Config({}));
  template.hasResourceProperties("AWS::EC2::Instance", {
    InstanceType: "t2.micro",
    KeyName: Match.anyValue(),
    AvailabilityZone: "ap-northeast-1a",
    IamInstanceProfile: Match.anyValue(),
    ImageId: Match.anyValue(),
    SecurityGroupIds: [Match.anyValue()],
    SubnetId: Match.anyValue(),
    Tags: [
      {
        Key: "Name",
        Value: "test-env-ec2-stack/test-env-ec2-instance",
      },
      {
        Key: "project",
        Value: "test-env",
      },
    ],
    UserData: Match.anyValue(),
  });
});

test("create ec2 stack: key pair", () => {
  const template = createTemplate(useEC2Config({}));
  template.hasResourceProperties("AWS::EC2::KeyPair", {
    KeyName: "test-env-bastion-key-pair",
    Tags: [
      {
        Key: "project",
        Value: "test-env",
      },
    ],
  });
});

describe("Elastic IP", () => {
  test("useElasticIPがfalseのとき、ElasticIPは作成されない", () => {
    const template = createTemplate(useEC2Config({ useElasticIP: false }));
    template.resourceCountIs("AWS::EC2::EIP", 0);
  });
  test("useElasticIPがtrueのとき、ElasticIPは作成される", () => {
    const template = createTemplate(useEC2Config({ useElasticIP: true }));
    template.resourceCountIs("AWS::EC2::EIP", 1);
    template.hasResourceProperties("AWS::EC2::EIP", {
      InstanceId: {
        Ref: Match.anyValue(),
      },
    });
  });
});
