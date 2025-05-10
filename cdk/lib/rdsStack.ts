import {
  Duration,
  Stack,
  StackProps,
  aws_ec2 as ec2,
  aws_rds as rds,
  aws_secretsmanager as sm,
  aws_applicationautoscaling as appscaling,
  Tags,
} from "aws-cdk-lib";
import { RetentionDays } from "aws-cdk-lib/aws-logs";
import { DatabaseCluster, DatabaseClusterEngine } from "aws-cdk-lib/aws-rds";
import { Construct } from "constructs";

import { PROJECT_TAG_KEY, RDSConfig } from "./buildConfig";

export interface RDSStackProps extends StackProps {
  readonly prefix: string;
  readonly vpc: ec2.IVpc;
  readonly config: RDSConfig;
}

export class RDSStack extends Stack {
  public readonly database: rds.IDatabaseCluster | rds.IDatabaseInstance;

  public readonly isRdsMetricsAlertEnabled: boolean;

  public readonly rdsSecurityGroup: ec2.ISecurityGroup;

  public readonly rdsSecret: sm.Secret;

  constructor(scope: Construct, id: string, props: RDSStackProps) {
    super(scope, id, props);
    const { vpc, prefix, config } = props;

    this.rdsSecurityGroup = new ec2.SecurityGroup(
      this,
      `${props.prefix}-rds-sg`,
      {
        vpc,
        allowAllOutbound: true,
        securityGroupName: `${prefix}-rds-sg`,
      },
    );

    // RDSを起動させるサブネット
    const privateSubnet = vpc.selectSubnets({
      subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
    }).subnets;

    // RDSのSubnetGroup作成
    const subnetGroup = new rds.SubnetGroup(
      this,
      `${prefix}-rds-subnet-group`,
      {
        description: `aurora subnet group for ${prefix}`,
        vpc,
        vpcSubnets: {
          subnets: privateSubnet,
        },
      },
    );
    // RDS for PostgreSQL
    this.rdsSecret = new sm.Secret(this, `${prefix}-rds-secret`, {
      secretName: `${prefix}-rds-secret`,
      generateSecretString: {
        excludePunctuation: true,
        includeSpace: false,
        secretStringTemplate: JSON.stringify({
          username: "dbuser", // TODO: configを使う
        }),
        generateStringKey: "password",
      },
    });

    if (config.isAuroraEnabled) {
      this.isRdsMetricsAlertEnabled = false; // Auroraはメトリクス監視アラートは不要
      this.database = new DatabaseCluster(this, `${prefix}-db-cluster`, {
        engine: DatabaseClusterEngine.auroraMysql({
          version: rds.AuroraMysqlEngineVersion.VER_3_02_0,
        }),
        cloudwatchLogsRetention: RetentionDays.ONE_MONTH,
        parameters: {
          character_set_client: "utf8mb4",
          character_set_connection: "utf8mb4",
          character_set_database: "utf8mb4",
          character_set_results: "utf8mb4",
          character_set_server: "utf8mb4",
          slow_query_log: "1",
          long_query_time: "1",
          log_output: "FILE",
        },
        vpc,
        vpcSubnets: {
          subnets: privateSubnet,
        },
        writer: rds.ClusterInstance.provisioned(`${prefix}-db-writer`, {
          instanceType: ec2.InstanceType.of(
            config.instanceClass,
            config.instanceSize,
          ),
          allowMajorVersionUpgrade: false,
          autoMinorVersionUpgrade: true,
          publiclyAccessible: false,
          instanceIdentifier: `${prefix}-db-writer`,
          caCertificate: rds.CaCertificate.RDS_CA_RDS2048_G1,
        }),
        readers: [
          rds.ClusterInstance.serverlessV2(`${prefix}-db-reader`, {
            caCertificate: rds.CaCertificate.RDS_CA_RDS2048_G1,
          }),
        ],
        securityGroups: [this.rdsSecurityGroup],
        backup: {
          retention: Duration.days(config.backupRetentionDays),
        },
        port: config.port,
        defaultDatabaseName: config.defaultDatabaseName,
        subnetGroup,
        credentials: rds.Credentials.fromSecret(this.rdsSecret),
      });

      // オートスケーリングの設定
      const target = new appscaling.ScalableTarget(
        this,
        `${prefix}-db-cluster-scalable-target`,
        {
          serviceNamespace: appscaling.ServiceNamespace.RDS,
          maxCapacity: config.autoScaling.maxCapacity,
          minCapacity: config.autoScaling.minCapacity,
          resourceId: `cluster:${this.database.clusterIdentifier}`,
          scalableDimension: "rds:cluster:ReadReplicaCount",
        },
      );
      target.scaleToTrackMetric(`${prefix}-db-autoscaling`, {
        targetValue: config.autoScaling.targetUtilizationPercent,
        predefinedMetric:
          appscaling.PredefinedMetric.RDS_READER_AVERAGE_CPU_UTILIZATION,
      });
    } else {
      this.isRdsMetricsAlertEnabled = true; // 非Auroraはオートスケーリング設定が無いのでメトリクスを監視するようにする
      this.database = new rds.DatabaseInstance(this, `${prefix}-db-instance`, {
        engine: rds.DatabaseInstanceEngine.mysql({
          version: rds.MysqlEngineVersion.VER_8_0,
        }),
        instanceType: ec2.InstanceType.of(
          config.instanceClass,
          config.instanceSize,
        ),
        vpc,
        cloudwatchLogsExports: ["error", "slowquery"],
        cloudwatchLogsRetention: RetentionDays.ONE_MONTH,
        parameters: {
          character_set_client: "utf8mb4",
          character_set_connection: "utf8mb4",
          character_set_database: "utf8mb4",
          character_set_results: "utf8mb4",
          character_set_server: "utf8mb4",
          slow_query_log: "1",
          long_query_time: "1",
          log_output: "FILE",
        },
        vpcSubnets: {
          subnets: privateSubnet,
        },
        securityGroups: [this.rdsSecurityGroup],
        port: config.port,
        databaseName: config.defaultDatabaseName,
        subnetGroup,
        credentials: rds.Credentials.fromSecret(this.rdsSecret),
      });
    }

    Tags.of(this).add(PROJECT_TAG_KEY, prefix);
  }
}
