import {
  Stack,
  StackProps,
  aws_wafv2 as wafv2,
  aws_ecs_patterns as ecsp,
  Tags,
} from "aws-cdk-lib";
import { Construct } from "constructs";
import { PROJECT_TAG_KEY } from "./buildConfig";

export interface WafStackProps extends StackProps {
  readonly prefix: string;
  readonly fargate: ecsp.ApplicationLoadBalancedFargateService;
}

export class WafStack extends Stack {
  webACLId: string;

  constructor(scope: Construct, id: string, props: WafStackProps) {
    super(scope, id, props);

    const { prefix, fargate } = props;

    const webACL = new wafv2.CfnWebACL(this, `${prefix}-web-acl`, {
      name: `${prefix}-web-acl`,
      defaultAction: {
        allow: {},
      },
      scope: "REGIONAL",
      visibilityConfig: {
        cloudWatchMetricsEnabled: true,
        metricName: `${prefix}-webacl-rule-metric`,
        sampledRequestsEnabled: true,
      },
      rules: [
        // AWS Managed Rules
        // https://docs.aws.amazon.com/ja_jp/waf/latest/developerguide/aws-managed-rule-groups-baseline.html
        {
          name: "AWSManagedRulesCommonRuleSet",
          priority: 1,
          overrideAction: { none: {} },
          statement: {
            managedRuleGroupStatement: {
              vendorName: "AWS",
              name: "AWSManagedRulesCommonRuleSet",
              excludedRules: [
                {
                  name: "SizeRestrictions_BODY",
                },
              ],
            },
          },
          visibilityConfig: {
            cloudWatchMetricsEnabled: true,
            metricName: "AWSManagedRulesCommonRuleSet",
            sampledRequestsEnabled: true,
          },
        },
      ],
    });

    // ALBにWebACLを関連付け
    // eslint-disable-next-line no-new
    new wafv2.CfnWebACLAssociation(this, "WebAclAssociation", {
      resourceArn: fargate.loadBalancer.loadBalancerArn,
      webAclArn: webACL.attrArn,
    });

    Tags.of(this).add(PROJECT_TAG_KEY, prefix);
  }
}
