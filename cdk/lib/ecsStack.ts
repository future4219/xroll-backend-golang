/* eslint-disable no-new */

import {
  Stack,
  StackProps,
  aws_ec2 as ec2,
  aws_ecs_patterns as ecsp,
  aws_ecr as ecr,
  aws_secretsmanager as sm,
  aws_route53 as route53,
  aws_ecs as ecs,
  aws_certificatemanager as acm,
  aws_elasticloadbalancingv2 as elbv2,
  aws_logs as logs,
  RemovalPolicy,
  Duration,
  Tags,
} from "aws-cdk-lib";
import { Construct } from "constructs";

import { ECSConfig, PROJECT_TAG_KEY, isHTTPSAvailable } from "./buildConfig";

export interface ECSStackProps extends StackProps {
  readonly prefix: string;
  readonly config: ECSConfig;
  readonly vpc: ec2.IVpc;
  readonly rdsSecret: sm.Secret;
  readonly containerSecret: sm.Secret;
  readonly ecrTag: string;
}

export class EcsStack extends Stack {
  public readonly ecsFargateService: ecsp.ApplicationLoadBalancedFargateService;

  public readonly ecsServiceSecurityGroup: ec2.ISecurityGroup;

  public readonly ecsTargetGroup: elbv2.ApplicationTargetGroup;

  public readonly ecsApiLogGroup: logs.LogGroup;

  public readonly ecsBatchLogGroup: logs.LogGroup;

  constructor(scope: Construct, id: string, props: ECSStackProps) {
    super(scope, id, props);
    const { config, prefix, vpc, rdsSecret, ecrTag, containerSecret } = props;

    // CloudWatch Logsグループの作成
    const apiLogGroup = new logs.LogGroup(this, `${prefix}-ecs-log`, {
      logGroupName: `/ecs/${prefix}`,
      removalPolicy: RemovalPolicy.DESTROY, // スタックの削除時にロググループも削除
    });

    const batchLogGroup = new logs.LogGroup(this, `${prefix}-ecs-batch-log`, {
      logGroupName: `/batch/${prefix}`,
      removalPolicy: RemovalPolicy.DESTROY, // スタックの削除時にロググループも削除
    });

    // ECR
    const repo = ecr.Repository.fromRepositoryName(
      this,
      `${prefix}-repo`,
      `${prefix}-repo`,
    );

    const taskDefinition = new ecs.FargateTaskDefinition(
      this,
      `${prefix}-task-def`,
      {
        family: `${prefix}-task-def-family`,
      },
    );

    taskDefinition.addContainer(`${prefix}-app-container`, {
      containerName: `${prefix}-app-container`,
      image: ecs.ContainerImage.fromEcrRepository(repo, ecrTag),
      cpu: config.api.spec.cpu,
      essential: true,
      command: ["/go/src/app/main"],
      logging: new ecs.AwsLogDriver({
        logGroup: apiLogGroup,
        streamPrefix: `${prefix}-ecs-log`,
      }),
      memoryLimitMiB: config.api.spec.memory,
      memoryReservationMiB: config.api.spec.memory,
      portMappings: [
        { hostPort: 80, protocol: ecs.Protocol.TCP, containerPort: 80 },
      ],
      secrets: {
        DB_USER: ecs.Secret.fromSecretsManager(rdsSecret, "username"),
        DB_PASSWORD: ecs.Secret.fromSecretsManager(rdsSecret, "password"),
        DB_HOST: ecs.Secret.fromSecretsManager(rdsSecret, "host"),
        DB_PORT: ecs.Secret.fromSecretsManager(rdsSecret, "port"),
        DB_NAME: ecs.Secret.fromSecretsManager(rdsSecret, "dbname"),
        ENV: ecs.Secret.fromSecretsManager(containerSecret, "ENV"),
        SIG_KEY: ecs.Secret.fromSecretsManager(containerSecret, "SIG_KEY"),
        AWS_ACCESS_KEY_ID: ecs.Secret.fromSecretsManager(
          containerSecret,
          "AWS_ACCESS_KEY_ID",
        ),
        AWS_SECRET_ACCESS_KEY: ecs.Secret.fromSecretsManager(
          containerSecret,
          "AWS_SECRET_ACCESS_KEY",
        ),
        AWS_REGION: ecs.Secret.fromSecretsManager(
          containerSecret,
          "AWS_REGION",
        ),
        S3_BUCKET: ecs.Secret.fromSecretsManager(containerSecret, "S3_BUCKET"),
        EMAIL_FROM: ecs.Secret.fromSecretsManager(
          containerSecret,
          "EMAIL_FROM",
        ),
        POST_CODE_JP_TOKEN: ecs.Secret.fromSecretsManager(
          containerSecret,
          "POST_CODE_JP_TOKEN",
        ),
        FRONTEND_URL: ecs.Secret.fromSecretsManager(
          containerSecret,
          "FRONTEND_URL",
        ),
        STRIPE_ENDPOINT_SECRET: ecs.Secret.fromSecretsManager(
          containerSecret,
          "STRIPE_ENDPOINT_SECRET",
        ),
        STRIPE_API_KEY: ecs.Secret.fromSecretsManager(
          containerSecret,
          "STRIPE_API_KEY",
        ),
        VIDEO_CLOUD_FRONT_URL: ecs.Secret.fromSecretsManager(
          containerSecret,
          "VIDEO_CLOUD_FRONT_URL",
        ),
        VIDEO_CLOUD_FRONT_KEY_ID: ecs.Secret.fromSecretsManager(
          containerSecret,
          "VIDEO_CLOUD_FRONT_KEY_ID",
        ),
        VIDEO_CLOUD_FRONT_PRIVATE_KEY: ecs.Secret.fromSecretsManager(
          containerSecret,
          "VIDEO_CLOUD_FRONT_PRIVATE_KEY",
        ),
      },
      environment: {
        SLOW_QUERY_THRESHOLD_MILLISECOND:
          config.api.slowQueryThresholdMilliSecond ?? "1000",
      },
    });

    const migrationTaskDefinition = new ecs.FargateTaskDefinition(
      this,
      `${prefix}-task-migration-def`,
      {
        family: `${prefix}-task-migration-def-family`,
      },
    );

    migrationTaskDefinition.addContainer(`${prefix}-app-migration`, {
      containerName: `${prefix}-app-migration`,
      image: ecs.ContainerImage.fromEcrRepository(repo, ecrTag),
      cpu: config.migration.spec.cpu,
      essential: true,
      command: [
        "sh",
        "-c",
        'migrate -path db/migrations -database "mysql://$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?multiStatements=true" up',
      ],
      logging: new ecs.AwsLogDriver({
        logGroup: batchLogGroup,
        streamPrefix: `${prefix}-migration`,
      }),
      memoryLimitMiB: config.migration.spec.memory,
      memoryReservationMiB: config.migration.spec.memory,
      secrets: {
        DB_USER: ecs.Secret.fromSecretsManager(rdsSecret, "username"),
        DB_PASSWORD: ecs.Secret.fromSecretsManager(rdsSecret, "password"),
        DB_HOST: ecs.Secret.fromSecretsManager(rdsSecret, "host"),
        DB_PORT: ecs.Secret.fromSecretsManager(rdsSecret, "port"),
        DB_NAME: ecs.Secret.fromSecretsManager(rdsSecret, "dbname"),
      },
    });

    const batchTaskDefinition = new ecs.FargateTaskDefinition(
      this,
      `${prefix}-task-batch-def`,
      {
        family: `${prefix}-task-batch-def-family`,
      },
    );

    batchTaskDefinition.addContainer(`${prefix}-batch-container`, {
      containerName: `${prefix}-batch-container`,
      image: ecs.ContainerImage.fromEcrRepository(repo, ecrTag),
      cpu: config.api.spec.cpu,
      essential: true,
      command: ["sh", "-c", "go run cmd/initdb/initdb.go"],
      logging: new ecs.AwsLogDriver({
        logGroup: this.ecsBatchLogGroup,
        streamPrefix: `${prefix}-batch`,
      }),
      memoryLimitMiB: config.api.spec.memory,
      memoryReservationMiB: config.api.spec.memory,
      secrets: {
        DB_USER: ecs.Secret.fromSecretsManager(rdsSecret, "username"),
        DB_PASSWORD: ecs.Secret.fromSecretsManager(rdsSecret, "password"),
        DB_HOST: ecs.Secret.fromSecretsManager(rdsSecret, "host"),
        DB_PORT: ecs.Secret.fromSecretsManager(rdsSecret, "port"),
        DB_NAME: ecs.Secret.fromSecretsManager(rdsSecret, "dbname"),
        ENV: ecs.Secret.fromSecretsManager(containerSecret, "ENV"),
        SIG_KEY: ecs.Secret.fromSecretsManager(containerSecret, "SIG_KEY"),
        AWS_ACCESS_KEY_ID: ecs.Secret.fromSecretsManager(
          containerSecret,
          "AWS_ACCESS_KEY_ID",
        ),
        AWS_SECRET_ACCESS_KEY: ecs.Secret.fromSecretsManager(
          containerSecret,
          "AWS_SECRET_ACCESS_KEY",
        ),
        AWS_REGION: ecs.Secret.fromSecretsManager(
          containerSecret,
          "AWS_REGION",
        ),
        S3_BUCKET: ecs.Secret.fromSecretsManager(containerSecret, "S3_BUCKET"),
        EMAIL_FROM: ecs.Secret.fromSecretsManager(
          containerSecret,
          "EMAIL_FROM",
        ),
        POST_CODE_JP_TOKEN: ecs.Secret.fromSecretsManager(
          containerSecret,
          "POST_CODE_JP_TOKEN",
        ),
        FRONTEND_URL: ecs.Secret.fromSecretsManager(
          containerSecret,
          "FRONTEND_URL",
        ),
        STRIPE_ENDPOINT_SECRET: ecs.Secret.fromSecretsManager(
          containerSecret,
          "STRIPE_ENDPOINT_SECRET",
        ),
        STRIPE_API_KEY: ecs.Secret.fromSecretsManager(
          containerSecret,
          "STRIPE_API_KEY",
        ),
        VIDEO_CLOUD_FRONT_URL: ecs.Secret.fromSecretsManager(
          containerSecret,
          "VIDEO_CLOUD_FRONT_URL",
        ),
        VIDEO_CLOUD_FRONT_KEY_ID: ecs.Secret.fromSecretsManager(
          containerSecret,
          "VIDEO_CLOUD_FRONT_KEY_ID",
        ),
        VIDEO_CLOUD_FRONT_PRIVATE_KEY: ecs.Secret.fromSecretsManager(
          containerSecret,
          "VIDEO_CLOUD_FRONT_PRIVATE_KEY",
        ),
      },
    });

    const cluster = new ecs.Cluster(this, `${prefix}-ecs-cluster`, {
      vpc,
    });

    this.ecsServiceSecurityGroup = new ec2.SecurityGroup(
      this,
      `${prefix}-ecs-service-sg`,
      { vpc, securityGroupName: `${prefix}-ecs-sg` },
    );

    let albHTTPSConfig: ecsp.ApplicationLoadBalancedFargateServiceProps = {};
    if (isHTTPSAvailable(config)) {
      const domainZone = route53.HostedZone.fromHostedZoneAttributes(
        this,
        `${prefix}-hosted-zone`,
        {
          zoneName: config.zoneName,
          hostedZoneId: config.hostedZoneId,
        },
      );
      const certificate = acm.Certificate.fromCertificateArn(
        this,
        `${prefix}-cert`,
        config.certificateArn,
      );
      albHTTPSConfig = {
        certificate,
        sslPolicy: elbv2.SslPolicy.RECOMMENDED_TLS,
        domainName: config.apiDomainName,
        domainZone,
        redirectHTTP: true,
      };
    }

    this.ecsFargateService = new ecsp.ApplicationLoadBalancedFargateService(
      this,
      `${prefix}-alb-fargate`,
      {
        serviceName: `${prefix}-alb-fargate`,
        cluster,
        taskSubnets: {
          subnetType: ec2.SubnetType.PUBLIC, // パブリックサブネットにECSをデプロイする
        },
        assignPublicIp: true,
        cpu: config.api.spec.cpu,
        memoryLimitMiB: config.api.spec.memory,
        taskDefinition,
        publicLoadBalancer: true,
        circuitBreaker: { rollback: true }, // deploy失敗時にロールバックを行うためのパラメータ
        deploymentController: {
          type: ecs.DeploymentControllerType.ECS,
        },
        securityGroups: [this.ecsServiceSecurityGroup],
        platformVersion: ecs.FargatePlatformVersion.VERSION1_4,
        // for HTTPS
        ...albHTTPSConfig,
      },
    );

    // LBが504を返すtimeout時間を設定する
    this.ecsFargateService.loadBalancer.setAttribute(
      "idle_timeout.timeout_seconds",
      `${config.api.idleTimeoutSeconds}`,
    );

    // apiのオートスケーリング設定
    const scalingPolicy = this.ecsFargateService.service.autoScaleTaskCount({
      minCapacity: config.api.autoScaling.minCapacity,
      maxCapacity: config.api.autoScaling.maxCapacity,
    });
    scalingPolicy.scaleOnCpuUtilization(`${prefix}-ecs-scaling-policy`, {
      policyName: `${prefix}-ecs-scaling-policy`,
      targetUtilizationPercent: config.api.autoScaling.targetUtilizationPercent,
    });

    // ヘルスチェックの設定
    this.ecsTargetGroup = this.ecsFargateService.targetGroup;
    this.ecsTargetGroup.configureHealthCheck({
      path: "/health",
      interval: Duration.seconds(30), // ヘルスチェックの感覚はアプリケーションに依存せず30秒固定させる
    });

    this.ecsApiLogGroup = apiLogGroup;
    this.ecsBatchLogGroup = batchLogGroup;

    Tags.of(this).add(PROJECT_TAG_KEY, prefix);
  }
}
