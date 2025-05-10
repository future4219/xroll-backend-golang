import {
  Stack,
  StackProps,
  aws_chatbot as chatbot,
  aws_cloudwatch as cloudwatch,
  aws_sns as sns,
  aws_cloudwatch_actions as cwActions,
  aws_rds as rds,
  aws_logs as logs,
  aws_elasticloadbalancingv2 as elbv2,
  Duration,
  Tags,
} from "aws-cdk-lib";
import { Construct } from "constructs";

import { AlarmConfig, PROJECT_TAG_KEY } from "./buildConfig";

export interface AlertStackProps extends StackProps {
  readonly prefix: string;
  readonly config: AlarmConfig;
  readonly database: rds.IDatabaseCluster | rds.IDatabaseInstance;
  readonly isRdsMetricsAlertEnabled: boolean;
  readonly ecsApiLogGroup: logs.LogGroup;
  readonly ecsBatchLogGroup: logs.LogGroup;
  readonly ecsTargetGroup: elbv2.ApplicationTargetGroup;
}

export class AlertStack extends Stack {
  constructor(scope: Construct, id: string, props: AlertStackProps) {
    super(scope, id, props);

    const {
      config,
      prefix,
      database,
      isRdsMetricsAlertEnabled,
      ecsApiLogGroup,
      ecsBatchLogGroup,
      ecsTargetGroup,
    } = props;

    // SNSトピックの作成
    const notificationTopic = new sns.Topic(this, `${prefix}-sns-topic`, {
      displayName: `${prefix}-sns-topic`,
      topicName: `${prefix}-sns-topic`,
    });
    const action = new cwActions.SnsAction(notificationTopic);

    // eslint-disable-next-line no-new
    new chatbot.SlackChannelConfiguration(this, `${prefix}-slack-channel`, {
      slackChannelConfigurationName: `${prefix}-slack-channel`,
      slackWorkspaceId: config.slackWorkspaceId,
      slackChannelId: config.slackChannelId,
      notificationTopics: [notificationTopic],
    });

    const statusCode5xxMetricAlarm = new cloudwatch.Alarm(
      this,
      `${prefix}-status-code-5xx-metric-alert`,
      {
        alarmName: `${prefix}-STATUS_CODE_5XX_ALERT`,
        metric: ecsTargetGroup.metrics.httpCodeTarget(
          elbv2.HttpCodeTarget.TARGET_5XX_COUNT,
          {
            period: Duration.seconds(config.statusCode5xxFilter.period),
          },
        ),
        actionsEnabled: true,
        threshold: config.statusCode5xxFilter.threshold,
        evaluationPeriods: config.statusCode5xxFilter.evaluationPeriods,
        datapointsToAlarm: config.statusCode5xxFilter.dataPointsToAlarm,
        comparisonOperator:
          cloudwatch.ComparisonOperator.GREATER_THAN_THRESHOLD,
      },
    );
    statusCode5xxMetricAlarm.addAlarmAction(action);

    const ecsUnhealthyMetricAlarm = new cloudwatch.Alarm(
      this,
      `${prefix}-ecs-unhealthy-metric-alert`,
      {
        alarmName: `${prefix}-ECS_UNHEALTHY_ALERT`,
        metric: ecsTargetGroup.metrics.unhealthyHostCount({
          period: Duration.seconds(config.ecsUnhealthyFilter.period),
        }),
        actionsEnabled: true,
        threshold: config.ecsUnhealthyFilter.threshold,
        evaluationPeriods: config.ecsUnhealthyFilter.evaluationPeriods,
        datapointsToAlarm: config.ecsUnhealthyFilter.dataPointsToAlarm,
        comparisonOperator:
          cloudwatch.ComparisonOperator.GREATER_THAN_THRESHOLD,
      },
    );
    ecsUnhealthyMetricAlarm.addAlarmAction(action);
    ecsUnhealthyMetricAlarm.addOkAction(action);

    // level = error となっているログを拾うフィルター
    const ecsApiErrorLogMetricFilter = new logs.MetricFilter(
      this,
      `${prefix}-ecs-api-error-log-metric-filter`,
      {
        logGroup: ecsApiLogGroup,
        metricNamespace: `${prefix}-ecs-api-error-log-metric-filter-name-space`,
        metricName: `${prefix}-ecs-api-error-log-metric-filter`,
        filterPattern: logs.FilterPattern.stringValue("$.level", "=", "error"),
        metricValue: "1",
      },
    );

    const ecsApiErrorLogMetricAlarm = new cloudwatch.Alarm(
      this,
      `${prefix}-ecs-api-error-log-metric-alarm`,
      {
        alarmName: `${prefix}-ECS_API_ERROR_LOG_ALARM`,
        alarmDescription: "APPLICATIONでエラーが発生しました",
        metric: ecsApiErrorLogMetricFilter.metric({
          period: Duration.seconds(config.ecsApiErrorLogMetricFilter.period),
        }),
        actionsEnabled: true,
        threshold: config.ecsApiErrorLogMetricFilter.threshold,
        evaluationPeriods: config.ecsApiErrorLogMetricFilter.evaluationPeriods,
        datapointsToAlarm: config.ecsApiErrorLogMetricFilter.dataPointsToAlarm,
        comparisonOperator:
          cloudwatch.ComparisonOperator.GREATER_THAN_THRESHOLD,
      },
    );
    ecsApiErrorLogMetricAlarm.addAlarmAction(action);

    // level = error となっているログを拾うフィルター
    const ecsBatchErrorLogMetricFilter = new logs.MetricFilter(
      this,
      `${prefix}-ecs-batch-error-log-metric-filter`,
      {
        logGroup: ecsBatchLogGroup,
        metricNamespace: `${prefix}-ecs-batch-error-log-metric-filter-name-space`,
        metricName: `${prefix}-ecs-batch-error-log-metric-filter`,
        filterPattern: logs.FilterPattern.stringValue("$.level", "=", "error"),
        metricValue: "1",
      },
    );

    const ecsBatchErrorLogMetricAlarm = new cloudwatch.Alarm(
      this,
      `${prefix}-ecs-batch-error-log-metric-alarm`,
      {
        alarmName: `${prefix}-ECS_BATCH_ERROR_LOG_ALARM`,
        alarmDescription: "バッチ処理でエラーが発生しました",
        metric: ecsBatchErrorLogMetricFilter.metric({
          period: Duration.seconds(config.ecsBatchErrorLogMetricFilter.period),
        }),
        actionsEnabled: true,
        threshold: config.ecsBatchErrorLogMetricFilter.threshold,
        evaluationPeriods:
          config.ecsBatchErrorLogMetricFilter.evaluationPeriods,
        datapointsToAlarm:
          config.ecsBatchErrorLogMetricFilter.dataPointsToAlarm,
        comparisonOperator:
          cloudwatch.ComparisonOperator.GREATER_THAN_THRESHOLD,
      },
    );
    ecsBatchErrorLogMetricAlarm.addAlarmAction(action);
    ecsBatchErrorLogMetricAlarm.addOkAction(action);

    // スロークエリのログを拾うフィルター
    const ecsSlowQueryLogMetricFilter = new logs.MetricFilter(
      this,
      `${prefix}-ecs-slow-query-log-metric-filter`,
      {
        logGroup: ecsApiLogGroup,
        metricNamespace: `${prefix}-ecs-slow-query-log-metric-filter-name-space`,
        metricName: `${prefix}-ecs-slow-query-log-metric-filter`,
        filterPattern: logs.FilterPattern.literal("SLOW SQL"),
        metricValue: "1",
      },
    );

    const ecsSlowQueryLogMetricAlarm = new cloudwatch.Alarm(
      this,
      `${prefix}-ecs-slow-query-log-metric-alarm`,
      {
        alarmName: `${prefix}-ECS_SLOW_QUERY_LOG_ALARM`,
        alarmDescription: "スロークエリが実行されました",
        metric: ecsSlowQueryLogMetricFilter.metric({
          period: Duration.seconds(config.ecsSlowQueryLogMetricFilter.period),
        }),
        actionsEnabled: true,
        threshold: config.ecsSlowQueryLogMetricFilter.threshold,
        evaluationPeriods: config.ecsSlowQueryLogMetricFilter.evaluationPeriods,
        datapointsToAlarm: config.ecsSlowQueryLogMetricFilter.dataPointsToAlarm,
        comparisonOperator:
          cloudwatch.ComparisonOperator.GREATER_THAN_THRESHOLD,
      },
    );
    ecsSlowQueryLogMetricAlarm.addAlarmAction(action);

    // RDSのメトリクス監視が必要な場合Alarmの設定を追加
    // CPU, Memory, Storageを監視
    // RDSStack内でDBがAuroraかどうかで決まっている。
    if (isRdsMetricsAlertEnabled) {
      // RDSのCPU利用量の監視
      const rdsCPUUtilizationAlarm = new cloudwatch.Alarm(
        this,
        `${prefix}-rds-cpu-alarm`,
        {
          alarmName: `${prefix}-RDS_CPU_UTILIZATION_ALERT`,
          actionsEnabled: true,
          metric: database.metricCPUUtilization({
            period: Duration.seconds(config.rdsMetricFilter.cpu.period),
          }),
          threshold: config.rdsMetricFilter.cpu.threshold, // configの単位 %
          evaluationPeriods: config.rdsMetricFilter.cpu.evaluationPeriods,
          comparisonOperator:
            cloudwatch.ComparisonOperator.GREATER_THAN_OR_EQUAL_TO_THRESHOLD,
        },
      );
      rdsCPUUtilizationAlarm.addAlarmAction(action);
      rdsCPUUtilizationAlarm.addOkAction(action);

      // RDSのメモリ使用量の監視
      const rdsFreeableMemoryAlarm = new cloudwatch.Alarm(
        this,
        `${prefix}-rds-memory-alarm`,
        {
          alarmName: `${prefix}-RDS_FREE_MEMORY_ALERT`,
          actionsEnabled: true,
          metric: database.metricFreeableMemory({
            period: Duration.seconds(config.rdsMetricFilter.freeMemory.period),
          }),
          threshold:
            1024 * 1024 * 1024 * config.rdsMetricFilter.freeMemory.threshold, // configの単位 GB。それをB単位に変換して代入
          evaluationPeriods:
            config.rdsMetricFilter.freeMemory.evaluationPeriods,
          comparisonOperator: cloudwatch.ComparisonOperator.LESS_THAN_THRESHOLD,
        },
      );
      rdsFreeableMemoryAlarm.addAlarmAction(action);
      rdsFreeableMemoryAlarm.addOkAction(action);

      // RDSのストレージ容量の監視
      const rdsFreeLocalStorageAlarm = new cloudwatch.Alarm(
        this,
        `${prefix}-rds-storage-alarm`,
        {
          alarmName: `${prefix}-RDS_FREE_STORAGE_ALERT`,
          actionsEnabled: true,
          // IDatabaseInstanceストレージ監視はそれ用の関数がないのでカスタムメトリクスを作成
          metric: database.metric("FreeStorageSpace", {
            statistic: "Average",
            period: Duration.seconds(config.rdsMetricFilter.freeStorage.period),
          }),
          threshold:
            1024 * 1024 * 1024 * config.rdsMetricFilter.freeStorage.threshold, // configの単位 GB。それをB単位に変換して代入
          evaluationPeriods:
            config.rdsMetricFilter.freeStorage.evaluationPeriods,
          comparisonOperator: cloudwatch.ComparisonOperator.LESS_THAN_THRESHOLD,
        },
      );
      rdsFreeLocalStorageAlarm.addAlarmAction(action);
      rdsFreeLocalStorageAlarm.addOkAction(action);
    }

    Tags.of(this).add(PROJECT_TAG_KEY, prefix);
  }
}
