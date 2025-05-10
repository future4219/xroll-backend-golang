import { InstanceClass, InstanceSize } from "aws-cdk-lib/aws-ec2";

import {
  AlarmConfig,
  BuildConfig,
  EC2Config,
  ECSConfig,
  RDSConfig,
  VPCConfig,
  isAlertMonitoringConfigured,
} from "../lib/buildConfig";

export const testPrefix = "test-env";
export const useVPCConfig = (config: Partial<VPCConfig>): VPCConfig => ({
  cidr: "192.168.0.0/21",
  cidrMask: 24,
  availabilityZones: ["ap-northeast-1a", "ap-northeast-1c"],
  ...config,
});

export const useECSConfig = (config: Partial<ECSConfig>): ECSConfig => ({
  api: {
    spec: {
      cpu: 256,
      memory: 512,
    },
    autoScaling: {
      minCapacity: 12,
      maxCapacity: 34,
      targetUtilizationPercent: 56,
    },
    idleTimeoutSeconds: 400,
  },
  migration: {
    spec: {
      cpu: 256,
      memory: 512,
    },
  },
  zoneName: undefined,
  hostedZoneId: undefined,
  apiDomainName: undefined,
  certificateArn: undefined,
  ...config,
});

export const useRDSConfig = (config: Partial<RDSConfig>): RDSConfig => ({
  instanceClass: InstanceClass.BURSTABLE2,
  instanceSize: InstanceSize.MICRO,
  port: 3306,
  backupRetentionDays: 20,
  defaultDatabaseName: "testdatabase",
  isAuroraEnabled: false,
  autoScaling: {
    minCapacity: 12,
    maxCapacity: 34,
    targetUtilizationPercent: 56,
  },
  ...config,
});

export const useEC2Config = (config: Partial<EC2Config>): EC2Config => ({
  useElasticIP: false,
  ...config,
});

export const useAlarmConfig = (config: Partial<AlarmConfig>): AlarmConfig => ({
  alertMonitoringEnabled: true,
  slackWorkspaceId: "workworkwork",
  slackChannelId: "channelchannel",
  statusCode5xxFilter: {
    threshold: 10,
    period: 300,
    evaluationPeriods: 10,
    dataPointsToAlarm: 10,
  },
  rdsMetricFilter: {
    cpu: {
      threshold: 80,
      period: 300,
      evaluationPeriods: 5,
    },
    freeMemory: {
      threshold: 0.5,
      period: 300,
      evaluationPeriods: 5,
    },
    freeStorage: {
      threshold: 1,
      period: 300,
      evaluationPeriods: 5,
    },
  },
  ecsUnhealthyFilter: {
    threshold: 20,
    period: 300,
    evaluationPeriods: 20,
    dataPointsToAlarm: 20,
  },
  ecsApiErrorLogMetricFilter: {
    threshold: 0,
    evaluationPeriods: 5,
    dataPointsToAlarm: 1,
    period: 60,
  },
  ecsBatchErrorLogMetricFilter: {
    threshold: 0,
    evaluationPeriods: 5,
    dataPointsToAlarm: 1,
    period: 60,
  },
  ecsSlowQueryLogMetricFilter: {
    threshold: 0,
    evaluationPeriods: 5,
    dataPointsToAlarm: 1,
    period: 60,
  },
  ...config,
});

export const useConfig = (
  vpc: VPCConfig,
  rds: RDSConfig,
  ecs: ECSConfig,
  ec2: EC2Config,
  alarm: AlarmConfig,
): BuildConfig => ({
  appName: "test",
  envName: "env",
  vpc,
  rds,
  ec2,
  ecs,
  alarm,
});

describe("isAlertMonitoringConfigured:", () => {
  const targetConfigList: {
    pattern: string;
    config: AlarmConfig;
    expected: boolean;
  }[] = [
    {
      pattern: "alertMonitoringEnabledがtrueでslackの情報も揃っている",
      config: useAlarmConfig({
        alertMonitoringEnabled: true,
        slackWorkspaceId: "workworkwork",
        slackChannelId: "channelchannel",
      }),
      expected: true,
    },
    {
      pattern: "alertMonitoringEnabledがfalseでslackの情報もない",
      config: useAlarmConfig({
        alertMonitoringEnabled: false,
        slackWorkspaceId: "",
        slackChannelId: "",
      }),
      expected: false,
    },
    {
      pattern: "alertMonitoringEnabledがtrueだがslackの情報は無い",
      config: useAlarmConfig({
        alertMonitoringEnabled: true,
        slackWorkspaceId: "",
        slackChannelId: "",
      }),
      expected: false,
    },
    {
      pattern: "alertMonitoringEnabledがfalseだがslackの情報は揃っている",
      config: useAlarmConfig({
        alertMonitoringEnabled: false,
        slackWorkspaceId: "workworkwork",
        slackChannelId: "channelchannel",
      }),
      expected: false,
    },
    {
      pattern:
        "alertMonitoringEnabledがtrueでslackの情報はチャンネルだけ設定されている",
      config: useAlarmConfig({
        alertMonitoringEnabled: true,
        slackWorkspaceId: "",
        slackChannelId: "channelchannel",
      }),
      expected: false,
    },
    {
      pattern:
        "alertMonitoringEnabledがtrueでslackの情報はワークスペースだけ設定されている",
      config: useAlarmConfig({
        alertMonitoringEnabled: true,
        slackWorkspaceId: "workworkwork",
        slackChannelId: "",
      }),
      expected: false,
    },
    {
      pattern:
        "alertMonitoringEnabledがfalseでslackの情報はチャンネルだけ設定されている",
      config: useAlarmConfig({
        alertMonitoringEnabled: false,
        slackWorkspaceId: "",
        slackChannelId: "channelchannel",
      }),
      expected: false,
    },
    {
      pattern:
        "alertMonitoringEnabledがfalseでslackの情報はワークスペースだけ設定されている",
      config: useAlarmConfig({
        alertMonitoringEnabled: false,
        slackWorkspaceId: "workworkwork",
        slackChannelId: "",
      }),
      expected: false,
    },
  ];
  test.each(targetConfigList)("$patternの時$expected", (testCase) => {
    const actual = isAlertMonitoringConfigured(testCase.config);
    expect(actual).toBe(testCase.expected);
  });
});
