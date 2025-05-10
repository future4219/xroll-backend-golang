import { App } from "aws-cdk-lib";
import { Match, Template } from "aws-cdk-lib/assertions";

import { AlertStack } from "../lib/alertStack";
import { prefix } from "../lib/buildConfig";
import * as Cdk from "../lib/ecsStack";
import { RDSStack } from "../lib/rdsStack";
import { SecretsManagerStack } from "../lib/smStack";
import { VPCStack } from "../lib/vpcStack";

import {
  useAlarmConfig,
  useConfig,
  useEC2Config,
  useECSConfig,
  useRDSConfig,
  useVPCConfig,
} from "./config.test";

const createTemplate = () => {
  const app = new App();
  const testConfig = useConfig(
    useVPCConfig({}),
    useRDSConfig({}),
    useECSConfig({}),
    useEC2Config({}),
    useAlarmConfig({}),
  );

  const testPrefix = prefix(testConfig);
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

  const testAlertStack = new AlertStack(app, `${testPrefix}-alarm-stack`, {
    config: testConfig.alarm,
    prefix: testPrefix,
    database: testRdsStack.database,
    isRdsMetricsAlertEnabled: testRdsStack.isRdsMetricsAlertEnabled,
    ecsApiLogGroup: testEcsStack.ecsApiLogGroup,
    ecsBatchLogGroup: testEcsStack.ecsBatchLogGroup,
    ecsTargetGroup: testEcsStack.ecsTargetGroup,
  });

  const template = Template.fromStack(testAlertStack);
  return template;
};

test("create alarm stack: SlackChannelConfiguration", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::Chatbot::SlackChannelConfiguration", {
    ConfigurationName: "test-env-slack-channel",
    SlackChannelId: "channelchannel",
    SlackWorkspaceId: "workworkwork",
    SnsTopicArns: [Match.anyValue()],
  });
});

test("create alarm stack: SNS Topic", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::SNS::Topic", {
    DisplayName: "test-env-sns-topic",
    TopicName: "test-env-sns-topic",
    Tags: [
      {
        Key: "project",
        Value: "test-env",
      },
    ],
  });
});

test("create status code alarm stack: Alarm", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::CloudWatch::Alarm", {
    ComparisonOperator: "GreaterThanThreshold",
    EvaluationPeriods: 10,
    AlarmActions: [Match.anyValue()],
    DatapointsToAlarm: 10,
    MetricName: "HTTPCode_Target_5XX_Count",
    Namespace: "AWS/ApplicationELB",
    Period: 300,
    Statistic: "Sum",
    Threshold: 10,
    ActionsEnabled: true,
  });
});

test("create alarm stack: ECS API Error Log Metric", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::CloudWatch::Alarm", {
    ComparisonOperator: "GreaterThanThreshold",
    EvaluationPeriods: 5,
    AlarmActions: [Match.anyValue()],
    AlarmDescription: "APPLICATIONでエラーが発生しました",
    AlarmName: "test-env-ECS_API_ERROR_LOG_ALARM",
    DatapointsToAlarm: 1,
    Namespace: "test-env-ecs-api-error-log-metric-filter-name-space",
    MetricName: "test-env-ecs-api-error-log-metric-filter",
    Period: 60,
    Statistic: "Average",
    Threshold: 0,
    ActionsEnabled: true,
  });
});

test("create alarm stack: ECS Batch Error Log Metric", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::CloudWatch::Alarm", {
    ComparisonOperator: "GreaterThanThreshold",
    EvaluationPeriods: 5,
    AlarmActions: [Match.anyValue()],
    AlarmDescription: "バッチ処理でエラーが発生しました",
    AlarmName: "test-env-ECS_BATCH_ERROR_LOG_ALARM",
    DatapointsToAlarm: 1,
    Namespace: "test-env-ecs-batch-error-log-metric-filter-name-space",
    MetricName: "test-env-ecs-batch-error-log-metric-filter",
    Period: 60,
    Statistic: "Average",
    Threshold: 0,
    ActionsEnabled: true,
  });
});

test("create alarm stack: ECS Slow Query Log Metric Filter", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::Logs::MetricFilter", {
    FilterPattern: "SLOW SQL",
    LogGroupName: Match.anyValue(),
    MetricTransformations: [
      {
        MetricName: "test-env-ecs-slow-query-log-metric-filter",
        MetricNamespace: "test-env-ecs-slow-query-log-metric-filter-name-space",
        MetricValue: "1",
      },
    ],
  });
});

test("create alarm stack: ECS Slow Query Log Metric", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::CloudWatch::Alarm", {
    ComparisonOperator: "GreaterThanThreshold",
    EvaluationPeriods: 5,
    AlarmActions: [Match.anyValue()],
    AlarmDescription: "スロークエリが実行されました",
    AlarmName: "test-env-ECS_SLOW_QUERY_LOG_ALARM",
    DatapointsToAlarm: 1,
    Namespace: "test-env-ecs-slow-query-log-metric-filter-name-space",
    MetricName: "test-env-ecs-slow-query-log-metric-filter",
    Period: 60,
    Statistic: "Average",
    Threshold: 0,
    ActionsEnabled: true,
  });
});

test("create alarm stack: Unhealthy Metric", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::CloudWatch::Alarm", {
    ComparisonOperator: "GreaterThanThreshold",
    EvaluationPeriods: 20,
    ActionsEnabled: true,
    AlarmActions: [Match.anyValue()],
    OKActions: [Match.anyValue()],
    AlarmName: "test-env-ECS_UNHEALTHY_ALERT",
    DatapointsToAlarm: 20,
    Dimensions: [
      {
        Name: "LoadBalancer",
        Value: {
          "Fn::Join": [
            "",
            [
              {
                "Fn::Select": [
                  1,
                  {
                    "Fn::Split": ["/", Match.anyValue()],
                  },
                ],
              },
              "/",
              {
                "Fn::Select": [
                  2,
                  {
                    "Fn::Split": ["/", Match.anyValue()],
                  },
                ],
              },
              "/",
              {
                "Fn::Select": [
                  3,
                  {
                    "Fn::Split": ["/", Match.anyValue()],
                  },
                ],
              },
            ],
          ],
        },
      },
      {
        Name: "TargetGroup",
        Value: Match.anyValue(),
      },
    ],
    MetricName: "UnHealthyHostCount",
    Namespace: "AWS/ApplicationELB",
    Period: 300,
    Statistic: "Average",
    Threshold: 20,
  });
});

test("create alarm stack: RDS CPU Utilization", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::CloudWatch::Alarm", {
    ComparisonOperator: "GreaterThanOrEqualToThreshold",
    EvaluationPeriods: 5,
    ActionsEnabled: true,
    AlarmActions: [Match.anyValue()],
    OKActions: [Match.anyValue()],
    AlarmName: "test-env-RDS_CPU_UTILIZATION_ALERT",
    Dimensions: [
      {
        Name: "DBInstanceIdentifier",
        Value: Match.anyValue(),
      },
    ],
    MetricName: "CPUUtilization",
    Namespace: "AWS/RDS",
    Period: 300,
    Statistic: "Average",
    Threshold: 80,
  });
});

test("create alarm stack: RDS Free Memory", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::CloudWatch::Alarm", {
    ComparisonOperator: "LessThanThreshold",
    EvaluationPeriods: 5,
    ActionsEnabled: true,
    AlarmActions: [Match.anyValue()],
    OKActions: [Match.anyValue()],
    AlarmName: "test-env-RDS_FREE_MEMORY_ALERT",
    Dimensions: [
      {
        Name: "DBInstanceIdentifier",
        Value: Match.anyValue(),
      },
    ],
    MetricName: "FreeableMemory",
    Namespace: "AWS/RDS",
    Period: 300,
    Statistic: "Average",
    Threshold: 536870912,
  });
});

test("create alarm stack: RDS Free Storage", () => {
  const template = createTemplate();
  template.hasResourceProperties("AWS::CloudWatch::Alarm", {
    ComparisonOperator: "LessThanThreshold",
    EvaluationPeriods: 5,
    ActionsEnabled: true,
    AlarmActions: [Match.anyValue()],
    OKActions: [Match.anyValue()],
    AlarmName: "test-env-RDS_FREE_STORAGE_ALERT",
    Dimensions: [
      {
        Name: "DBInstanceIdentifier",
        Value: Match.anyValue(),
      },
    ],
    MetricName: "FreeStorageSpace",
    Namespace: "AWS/RDS",
    Period: 300,
    Statistic: "Average",
    Threshold: 1073741824,
  });
});
