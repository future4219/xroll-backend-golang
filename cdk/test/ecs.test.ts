import * as cdk from "aws-cdk-lib";
import { Match, Template } from "aws-cdk-lib/assertions";

import { BuildConfig, ECSConfig } from "../lib/buildConfig";
import * as Cdk from "../lib/ecsStack";
import { RDSStack } from "../lib/rdsStack";
import { SecretsManagerStack } from "../lib/smStack";
import { VPCStack } from "../lib/vpcStack";

import {
  useConfig,
  useVPCConfig,
  useRDSConfig,
  useAlarmConfig,
  useECSConfig,
  testPrefix,
  useEC2Config,
} from "./config.test";

const creteTemplate = (ecsConfig: ECSConfig) => {
  const app = new cdk.App();

  const testConfig: BuildConfig = useConfig(
    useVPCConfig({}),
    useRDSConfig({}),
    ecsConfig,
    useEC2Config({}),
    useAlarmConfig({}),
  );

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

  const template = Template.fromStack(testEcsStack);
  return template;
};

describe("create ecs stack: TaskDefinition", () => {
  test("タスク定義が3つあるか", () => {
    const ecsConfig = useECSConfig({});
    const template = creteTemplate(ecsConfig);
    template.resourceCountIs("AWS::ECS::TaskDefinition", 3);
  });
  test("echoサーバ用のタスク定義", () => {
    const ecsConfig = useECSConfig({});
    const template = creteTemplate(ecsConfig);
    template.hasResourceProperties("AWS::ECS::TaskDefinition", {
      ContainerDefinitions: [
        {
          Cpu: 256,
          Memory: 512,
          MemoryReservation: 512,
          Name: "test-env-app-container",
          Command: ["/go/src/app/main"],
          PortMappings: [
            {
              ContainerPort: 80,
              HostPort: 80,
              Protocol: "tcp",
            },
          ],
          LogConfiguration: {
            LogDriver: "awslogs",
            Options: {
              "awslogs-group": Match.anyValue(),
              "awslogs-stream-prefix": "test-env-ecs-log",
              "awslogs-region": {
                Ref: "AWS::Region",
              },
            },
          },
          Secrets: [
            { Name: "DB_USER" },
            { Name: "DB_PASSWORD" },
            { Name: "DB_HOST" },
            { Name: "DB_PORT" },
            { Name: "DB_NAME" },
            { Name: "ENV" },
            { Name: "SIG_KEY" },
            { Name: "AWS_ACCESS_KEY_ID" },
            { Name: "AWS_SECRET_ACCESS_KEY" },
            { Name: "AWS_REGION" },
            { Name: "S3_BUCKET" },
            { Name: "EMAIL_FROM" },
            { Name: "POST_CODE_JP_TOKEN" },
            { Name: "FRONTEND_URL" },
            { Name: "STRIPE_ENDPOINT_SECRET" },
            { Name: "STRIPE_API_KEY" },
            { Name: "VIDEO_CLOUD_FRONT_URL" },
            { Name: "VIDEO_CLOUD_FRONT_KEY_ID" },
            { Name: "VIDEO_CLOUD_FRONT_PRIVATE_KEY" },
          ],
        },
      ],
      Cpu: "256",
      Memory: "512",
      Family: "test-env-task-def-family",
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });
  });
  test("マイグレーション用のタスク定義", () => {
    const ecsConfig = useECSConfig({});
    const template = creteTemplate(ecsConfig);
    template.hasResourceProperties("AWS::ECS::TaskDefinition", {
      ContainerDefinitions: [
        {
          Cpu: 256,
          Memory: 512,
          MemoryReservation: 512,
          Command: [
            "sh",
            "-c",
            'migrate -path db/migrations -database "mysql://$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?multiStatements=true" up',
          ],
          Name: "test-env-app-migration",
          LogConfiguration: {
            LogDriver: "awslogs",
            Options: {
              "awslogs-stream-prefix": "test-env-migration",
            },
          },
          Secrets: [
            { Name: "DB_USER" },
            { Name: "DB_PASSWORD" },
            { Name: "DB_HOST" },
            { Name: "DB_PORT" },
            { Name: "DB_NAME" },
          ],
        },
      ],
      Cpu: "256",
      Memory: "512",
      Family: "test-env-task-migration-def-family",
    });
  });
  test("バッチ処理タスク定義", () => {
    const ecsConfig = useECSConfig({});
    const template = creteTemplate(ecsConfig);
    template.hasResourceProperties("AWS::ECS::TaskDefinition", {
      ContainerDefinitions: [
        {
          Command: ["sh", "-c", "go run cmd/initdb/initdb.go"],
          Cpu: 256,
          Essential: true,
          Image: Match.anyValue(),
          LogConfiguration: {
            LogDriver: "awslogs",
            Options: {
              "awslogs-group": {
                Ref: Match.anyValue(),
              },
              "awslogs-stream-prefix": "test-env-batch",
              "awslogs-region": {
                Ref: "AWS::Region",
              },
            },
          },
          Memory: 512,
          MemoryReservation: 512,
          Name: "test-env-batch-container",
          Secrets: [
            { Name: "DB_USER" },
            { Name: "DB_PASSWORD" },
            { Name: "DB_HOST" },
            { Name: "DB_PORT" },
            { Name: "DB_NAME" },
            { Name: "ENV" },
            { Name: "SIG_KEY" },
            { Name: "AWS_ACCESS_KEY_ID" },
            { Name: "AWS_SECRET_ACCESS_KEY" },
            { Name: "AWS_REGION" },
            { Name: "S3_BUCKET" },
            { Name: "EMAIL_FROM" },
            { Name: "POST_CODE_JP_TOKEN" },
            { Name: "FRONTEND_URL" },
            { Name: "STRIPE_ENDPOINT_SECRET" },
            { Name: "STRIPE_API_KEY" },
            { Name: "VIDEO_CLOUD_FRONT_URL" },
            { Name: "VIDEO_CLOUD_FRONT_KEY_ID" },
            { Name: "VIDEO_CLOUD_FRONT_PRIVATE_KEY" },
          ],
        },
      ],
      Cpu: "256",
      ExecutionRoleArn: {
        "Fn::GetAtt": [Match.anyValue(), "Arn"],
      },
      Family: "test-env-task-batch-def-family",
      Memory: "512",
      NetworkMode: "awsvpc",
      RequiresCompatibilities: ["FARGATE"],
      TaskRoleArn: Match.anyValue(),
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });
  });
  test("オートスケーリングポリシー", () => {
    const ecsConfig = useECSConfig({});
    const template = creteTemplate(ecsConfig);
    template.hasResourceProperties(
      "AWS::ApplicationAutoScaling::ScalingPolicy",
      {
        PolicyName: "test-env-ecs-scaling-policy",
        PolicyType: "TargetTrackingScaling",
        ScalingTargetId: Match.anyValue(),
        TargetTrackingScalingPolicyConfiguration: {
          PredefinedMetricSpecification: {
            PredefinedMetricType: "ECSServiceAverageCPUUtilization",
          },
          TargetValue: 56,
        },
      },
    );
  });
  test("オートスケーリングターゲット", () => {
    const ecsConfig = useECSConfig({});
    const template = creteTemplate(ecsConfig);
    template.hasResourceProperties(
      "AWS::ApplicationAutoScaling::ScalableTarget",
      {
        MaxCapacity: 34,
        MinCapacity: 12,
        ResourceId: Match.anyValue(),
        RoleARN: Match.anyValue(),
        ScalableDimension: "ecs:service:DesiredCount",
        ServiceNamespace: "ecs",
      },
    );
  });
});

describe("create ecs stack: Fargate HTTP", () => {
  const ecsConfig = useECSConfig({});
  const template = creteTemplate(ecsConfig);
  test("ドメイン情報未設定時にHTTP版のサービスが作成される", () => {
    template.hasResourceProperties("AWS::ECS::Service", {
      ServiceName: "test-env-alb-fargate",
      LoadBalancers: [
        {
          ContainerName: "test-env-app-container",
          ContainerPort: 80,
        },
      ],
      NetworkConfiguration: {
        AwsvpcConfiguration: {
          AssignPublicIp: "ENABLED",
        },
      },
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });
  });
  test("ドメイン情報未設定時にHTTP版のリスナーが作成される", () => {
    template.hasResourceProperties("AWS::ElasticLoadBalancingV2::Listener", {
      DefaultActions: [
        {
          TargetGroupArn: {
            Ref: Match.anyValue(),
          },
          Type: "forward",
        },
      ],
      LoadBalancerArn: {
        Ref: Match.anyValue(),
      },
      Port: 80,
      Protocol: "HTTP",
    });
  });

  test("ターゲットグループに対してヘルスチェック監視の設定が作成される", () => {
    template.hasResourceProperties("AWS::ElasticLoadBalancingV2::TargetGroup", {
      HealthCheckIntervalSeconds: 30,
      HealthCheckPath: "/health",
      Port: 80,
      Protocol: "HTTP",
      TargetGroupAttributes: [
        {
          Key: "stickiness.enabled",
          Value: "false",
        },
      ],
      TargetType: "ip",
      VpcId: Match.anyValue(),
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });
  });
});

test("LBが生成される", () => {
  const template = creteTemplate(useECSConfig({}));
  template.hasResourceProperties("AWS::ElasticLoadBalancingV2::LoadBalancer", {
    LoadBalancerAttributes: [
      {
        Key: "deletion_protection.enabled",
        Value: "false",
      },
      {
        Key: "idle_timeout.timeout_seconds",
        Value: "400",
      },
    ],
    Scheme: "internet-facing",
    SecurityGroups: Match.anyValue(),
    Subnets: Match.anyValue(),
    Tags: [
      {
        Key: "project",
        Value: "test-env",
      },
    ],
    Type: "application",
  });
});

describe("create ecs stack: Fargate HTTPS", () => {
  const ecsConfig = useECSConfig({
    zoneName: "example.com",
    hostedZoneId: "1234567890",
    apiDomainName: "api.example.com",
    certificateArn:
      "arn:aws:acm:ap-northeast-1:1234567890:certificate/12345678-1234-1234-1234-123456789012",
  });
  const template = creteTemplate(ecsConfig);
  test("ドメイン情報設定時にHTTPS版のサービスが作成される", () => {
    template.hasResourceProperties("AWS::ECS::Service", {
      ServiceName: "test-env-alb-fargate",
      LoadBalancers: [
        {
          ContainerName: "test-env-app-container",
          ContainerPort: 80,
        },
      ],
      NetworkConfiguration: {
        AwsvpcConfiguration: {
          AssignPublicIp: "ENABLED",
        },
      },
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });
  });
  test("ドメイン情報設定時にHTTPS版のリスナーが作成される", () => {
    template.hasResourceProperties("AWS::ElasticLoadBalancingV2::Listener", {
      DefaultActions: [
        {
          TargetGroupArn: {
            Ref: Match.anyValue(),
          },
          Type: "forward",
        },
      ],
      LoadBalancerArn: {
        Ref: Match.anyValue(),
      },
      Certificates: [
        {
          CertificateArn:
            "arn:aws:acm:ap-northeast-1:1234567890:certificate/12345678-1234-1234-1234-123456789012",
        },
      ],
      Port: 443,
      Protocol: "HTTPS",
      SslPolicy: "ELBSecurityPolicy-TLS13-1-2-2021-06",
    });
  });
});
