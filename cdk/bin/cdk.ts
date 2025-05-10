#!/usr/bin/env node
/* eslint-disable no-new */
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";

import { AlertStack } from "../lib/alertStack";
import {
  BuildConfig,
  isAlertMonitoringConfigured,
  prefix,
} from "../lib/buildConfig";
import { EC2BastionStack } from "../lib/ec2BastionStack";
import { EcsStack } from "../lib/ecsStack";
import { RDSStack } from "../lib/rdsStack";
import { AssignSecurityGroup } from "../lib/securityGroups";
import { SecretsManagerStack } from "../lib/smStack";
import { VPCStack } from "../lib/vpcStack";
import { WafStack } from "../lib/wafStack";

const app = new cdk.App();
const appName: string = app.node.tryGetContext("appName");
const stage: string = app.node.tryGetContext("stage"); // stg or prod
const ecrTag: string = app.node.tryGetContext("tag");
const config: BuildConfig = { ...app.node.tryGetContext(stage), appName };
const appPrefix = prefix(config);

const vpcStack = new VPCStack(app, `${appPrefix}-vpc-stack`, {
  prefix: appPrefix,
  config: config.vpc,
  terminationProtection: true,
});

const ec2BastionStack = new EC2BastionStack(app, `${appPrefix}-ec2-stack`, {
  prefix: appPrefix,
  vpc: vpcStack.vpc,
  config: config.ec2,
});

const rdsStack = new RDSStack(app, `${appPrefix}-rds-stack`, {
  prefix: appPrefix,
  config: config.rds,
  vpc: vpcStack.vpc,
  terminationProtection: true,
});

const smStack = new SecretsManagerStack(app, `${appPrefix}-sm-stack`, {
  prefix: appPrefix,
});

const ecsStack = new EcsStack(app, `${appPrefix}-ecs-stack`, {
  stackName: `${appPrefix}-ecs-stack`,
  vpc: vpcStack.vpc,
  prefix: appPrefix,
  config: config.ecs,
  rdsSecret: rdsStack.rdsSecret,
  containerSecret: smStack.secret,
  ecrTag,
});

new AssignSecurityGroup(
  ecsStack.ecsServiceSecurityGroup,
  rdsStack.rdsSecurityGroup,
  ec2BastionStack.ec2SecurityGroup,
).assign(config.ec2);

if (isAlertMonitoringConfigured(config.alarm)) {
  new AlertStack(app, `${appPrefix}-alarm-stack`, {
    prefix: appPrefix,
    config: config.alarm,
    database: rdsStack.database,
    isRdsMetricsAlertEnabled: rdsStack.isRdsMetricsAlertEnabled,
    ecsApiLogGroup: ecsStack.ecsApiLogGroup,
    ecsBatchLogGroup: ecsStack.ecsBatchLogGroup,
    ecsTargetGroup: ecsStack.ecsTargetGroup,
  });
}

new WafStack(app, `${appPrefix}-waf-stack`, {
  prefix: appPrefix,
  fargate: ecsStack.ecsFargateService,
});
