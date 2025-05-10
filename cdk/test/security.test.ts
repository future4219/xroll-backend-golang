import * as cdk from "aws-cdk-lib";
import { Template } from "aws-cdk-lib/assertions";

import { BuildConfig, EC2Config } from "../lib/buildConfig";
import { EC2BastionStack } from "../lib/ec2BastionStack";
import * as Cdk from "../lib/ecsStack";
import { RDSStack } from "../lib/rdsStack";
import { AssignSecurityGroup } from "../lib/securityGroups";
import { SecretsManagerStack } from "../lib/smStack";
import { VPCStack } from "../lib/vpcStack";

import {
  useVPCConfig,
  useRDSConfig,
  useECSConfig,
  testPrefix,
  useAlarmConfig,
  useConfig,
  useEC2Config,
} from "./config.test";

/** テストで用いる２つのテンプレート */
type Templates = {
  /** ECSスタックのテンプレート */
  ecsTemp: Template;
  /** EC2スタックのテンプレート */
  ec2Temp: Template;
};

/**
 * テスト対象であるEC2とECSのテンプレートを作成する。
 * 引数のEC2Config以外のconfigはデフォルトの値を用いる
 * @param ec2Config DBクライアントツール(PhpMyAdmin)への接続IPを指定するために用いる
 * @returns
 */
const createTemplate = (ec2Config: EC2Config): Templates => {
  const testConfig: BuildConfig = useConfig(
    useVPCConfig({}),
    useRDSConfig({}),
    useECSConfig({}),
    ec2Config,
    useAlarmConfig({}),
  );

  const app = new cdk.App();

  const testEcrTag = "92f5d00b61d5c4fcedeecb08ed928bb3bddd7d4d";

  const testVpcStack = new VPCStack(app, `${testPrefix}-vpc-stack`, {
    prefix: testPrefix,
    config: testConfig.vpc,
  });

  const testRdsStack = new RDSStack(app, `${testPrefix}-rds-stack`, {
    prefix: testPrefix,
    config: testConfig.rds,
    vpc: testVpcStack.vpc,
  });

  const testEc2Stack = new EC2BastionStack(app, `${testPrefix}-ec2-stack`, {
    prefix: testPrefix,
    vpc: testVpcStack.vpc,
    config: testConfig.ec2,
  });

  const testSmStack = new SecretsManagerStack(app, `${testPrefix}-sm-stack`, {
    prefix: testPrefix,
  });

  const testEcsStack = new Cdk.EcsStack(app, `${testPrefix}-ecs-stack`, {
    stackName: `${testPrefix}-ecs-stack`,
    vpc: testVpcStack.vpc,
    prefix: testPrefix,
    config: testConfig.ecs,
    rdsSecret: testRdsStack.rdsSecret,
    ecrTag: testEcrTag,
    containerSecret: testSmStack.secret,
  });

  new AssignSecurityGroup(
    testEcsStack.ecsServiceSecurityGroup,
    testRdsStack.rdsSecurityGroup,
    testEc2Stack.ec2SecurityGroup,
  ).assign(testConfig.ec2);

  return {
    ec2Temp: Template.fromStack(testEc2Stack),
    ecsTemp: Template.fromStack(testEcsStack),
  };
};

test("ecs security group", () => {
  const { ecsTemp } = createTemplate(useEC2Config({}));
  ecsTemp.resourceCountIs("AWS::EC2::SecurityGroup", 2);
  ecsTemp.hasResourceProperties("AWS::EC2::SecurityGroup", {
    GroupName: "test-env-ecs-sg",
    SecurityGroupEgress: [
      {
        CidrIp: "0.0.0.0/0",
        Description: "Allow all outbound traffic by default",
        IpProtocol: "-1",
      },
    ],
    Tags: [
      {
        Key: "project",
        Value: "test-env",
      },
    ],
  });
  ecsTemp.hasResourceProperties("AWS::EC2::SecurityGroup", {
    SecurityGroupIngress: [
      {
        CidrIp: "0.0.0.0/0",
        Description: "Allow from anyone on port 80",
        FromPort: 80,
        IpProtocol: "tcp",
        ToPort: 80,
      },
    ],
    Tags: [
      {
        Key: "project",
        Value: "test-env",
      },
    ],
  });
});

test("ecs security group from elb", () => {
  const { ecsTemp } = createTemplate(useEC2Config({}));
  ecsTemp.hasResourceProperties("AWS::EC2::SecurityGroupIngress", {
    IpProtocol: "tcp",
    Description: "Load balancer to target",
    FromPort: 80,
    ToPort: 80,
  });
});

test("ecs security group to rds", () => {
  const { ecsTemp } = createTemplate(useEC2Config({}));
  ecsTemp.hasResourceProperties("AWS::EC2::SecurityGroupIngress", {
    Description: "ECS-RDS",
    FromPort: 3306,
    IpProtocol: "tcp",
    ToPort: 3306,
  });
});

test("bastion security group to rds", () => {
  const { ec2Temp } = createTemplate(useEC2Config({}));
  ec2Temp.hasResourceProperties("AWS::EC2::SecurityGroupIngress", {
    Description: "EC2-RDS",
    FromPort: 3306,
    IpProtocol: "tcp",
    ToPort: 3306,
  });
});

describe("踏み台EC2のセキュリテイグループ", () => {
  test.each([
    {
      IpAddressToDBClient: "2.2.2.2/32",
      exceptedPhpMyAdminRule: {
        CidrIp: "2.2.2.2/32",
        Description: "Allow access to PhpMyAdmin from specific IP",
        FromPort: 8080,
        IpProtocol: "tcp",
        ToPort: 8080,
      },
    },
    {
      IpAddressToDBClient: undefined,
      exceptedPhpMyAdminRule: {
        CidrIp: "0.0.0.0/0",
        Description: "Allow access to PhpMyAdmin from all users",
        FromPort: 8080,
        IpProtocol: "tcp",
        ToPort: 8080,
      },
    },
  ])(
    "ec2 security group",
    ({ IpAddressToDBClient, exceptedPhpMyAdminRule }) => {
      const { ec2Temp } = createTemplate(
        useEC2Config({
          IpAddressToDBClient,
        }),
      );
      ec2Temp.hasResourceProperties("AWS::EC2::SecurityGroup", {
        GroupDescription: "test-env-ec2-stack/test-env-ec2-sg",
        GroupName: "test-env-ec2-sg",
        SecurityGroupEgress: [
          {
            CidrIp: "0.0.0.0/0",
            Description: "Allow all outbound traffic by default",
            IpProtocol: "-1",
          },
        ],
        SecurityGroupIngress: [
          {
            CidrIp: "0.0.0.0/0",
            Description: "SSH connection to EC2",
            FromPort: 22,
            IpProtocol: "tcp",
            ToPort: 22,
          },
          exceptedPhpMyAdminRule,
        ],
        Tags: [
          {
            Key: "project",
            Value: "test-env",
          },
        ],
      });
    },
  );
});
