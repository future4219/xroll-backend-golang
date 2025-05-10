import * as cdk from "aws-cdk-lib";
import { Template } from "aws-cdk-lib/assertions";

import { prefix } from "../lib/buildConfig";
import * as Cdk from "../lib/ecsStack";
import { RDSStack } from "../lib/rdsStack";
import { SecretsManagerStack } from "../lib/smStack";
import { VPCStack } from "../lib/vpcStack";
import { WafStack } from "../lib/wafStack";

import {
  useAlarmConfig,
  useConfig,
  useEC2Config,
  useECSConfig,
  useRDSConfig,
  useVPCConfig,
} from "./config.test";

const createTemplate = () => {
  const app = new cdk.App();

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

  const testWafStack = new WafStack(app, `${testPrefix}-waf-stack`, {
    stackName: `${testPrefix}-ecs-stack`,
    prefix: testPrefix,
    fargate: testEcsStack.ecsFargateService,
  });

  return Template.fromStack(testWafStack);
};

describe("Create WebACL", () => {
  const template = createTemplate();
  test("AWS::WAFv2::WebACL", () => {
    template.hasResourceProperties("AWS::WAFv2::WebACL", {
      DefaultAction: {
        Allow: {},
      },
      Scope: "REGIONAL",
      VisibilityConfig: {
        CloudWatchMetricsEnabled: true,
        MetricName: "test-env-webacl-rule-metric",
        SampledRequestsEnabled: true,
      },
      Name: "test-env-web-acl",
      Rules: [
        {
          Name: "AWSManagedRulesCommonRuleSet",
          OverrideAction: {
            None: {},
          },
          Priority: 1,
          Statement: {
            ManagedRuleGroupStatement: {
              Name: "AWSManagedRulesCommonRuleSet",
              VendorName: "AWS",
              ExcludedRules: [
                {
                  Name: "SizeRestrictions_BODY",
                },
              ],
            },
          },
          VisibilityConfig: {
            CloudWatchMetricsEnabled: true,
            MetricName: "AWSManagedRulesCommonRuleSet",
            SampledRequestsEnabled: true,
          },
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
});
