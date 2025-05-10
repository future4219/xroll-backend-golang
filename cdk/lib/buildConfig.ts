import { aws_ec2 as ec2 } from "aws-cdk-lib";

export interface BuildConfig {
  appName: string;
  envName: string;
  ecs: ECSConfig;
  vpc: VPCConfig;
  ec2: EC2Config;
  rds: RDSConfig;
  alarm: AlarmConfig;
}

interface TaskSpecConfig {
  cpu: number;
  memory: number;
}

interface AutoScalingConfig {
  minCapacity: number;
  maxCapacity: number;
  targetUtilizationPercent: number;
}

export interface ECSConfig {
  api: {
    spec: TaskSpecConfig;
    autoScaling: AutoScalingConfig;
    slowQueryThresholdMilliSecond?: string;
    idleTimeoutSeconds: number;
  };
  migration: {
    spec: TaskSpecConfig;
  };
  zoneName?: string;
  hostedZoneId?: string;
  apiDomainName?: string;
  certificateArn?: string;
}

export interface VPCConfig {
  cidr: string;
  cidrMask: number;
  availabilityZones: string[];
}

export interface EC2Config {
  useElasticIP: boolean;
  IpAddressToDBClient?: string;
}

export interface RDSConfig {
  instanceClass: ec2.InstanceClass;
  instanceSize: ec2.InstanceSize;
  port: number;
  defaultDatabaseName: string;
  backupRetentionDays: number;
  isAuroraEnabled: boolean;
  autoScaling: AutoScalingConfig;
}

export interface AlarmConfig {
  alertMonitoringEnabled: boolean;
  slackWorkspaceId: string;
  slackChannelId: string;
  rdsMetricFilter: {
    cpu: RDSAlarmOption; // 使用量(%)で管理
    freeMemory: RDSAlarmOption; // 使用量(残GB)で管理
    freeStorage: RDSAlarmOption; // 使用量(残GB)で管理
  };
  statusCode5xxFilter: AlarmOption;
  ecsUnhealthyFilter: AlarmOption;
  ecsApiErrorLogMetricFilter: AlarmOption;
  ecsBatchErrorLogMetricFilter: AlarmOption;
  ecsSlowQueryLogMetricFilter: AlarmOption;
}

interface RDSAlarmOption {
  threshold: number;
  period: number;
  evaluationPeriods: number;
}

interface AlarmOption {
  threshold: number;
  period: number; // 単位: 秒
  evaluationPeriods: number; // アラート監視の間隔は、period * evaluationPeriods
  dataPointsToAlarm: number;
}

export const prefix = (c: BuildConfig) => `${c.appName}-${c.envName}`;

export const isHTTPSAvailable = (c: ECSConfig): c is Required<ECSConfig> =>
  c.zoneName != null &&
  c.zoneName !== "" &&
  c.hostedZoneId != null &&
  c.hostedZoneId !== "" &&
  c.apiDomainName != null &&
  c.apiDomainName !== "" &&
  c.certificateArn != null &&
  c.certificateArn !== "";

export const isAlertMonitoringConfigured = (alarmConfig: AlarmConfig) =>
  alarmConfig.alertMonitoringEnabled &&
  alarmConfig.slackWorkspaceId !== "" &&
  alarmConfig.slackChannelId !== "";

/**
 * CDKで立ち上げるリソースに一律のタグをつける
 *
 * 請求情報のフィルターに用いる
 */
export const PROJECT_TAG_KEY = "project";
