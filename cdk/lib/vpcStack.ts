import { Stack, StackProps, Tags, aws_ec2 as ec2 } from "aws-cdk-lib";
import { IpAddresses } from "aws-cdk-lib/aws-ec2";
import { Construct } from "constructs";

import { PROJECT_TAG_KEY, VPCConfig } from "./buildConfig";

export interface VPCStackProps extends StackProps {
  readonly prefix: string;
  readonly config: VPCConfig;
}

export class VPCStack extends Stack {
  public readonly vpc: ec2.IVpc;

  constructor(scope: Construct, id: string, props: VPCStackProps) {
    super(scope, id, props);
    const { config, prefix } = props;
    this.vpc = new ec2.Vpc(this, `${prefix}-vpc`, {
      vpcName: `${prefix}-vpc`,
      ipAddresses: IpAddresses.cidr(config.cidr),
      availabilityZones: config.availabilityZones,
      subnetConfiguration: [
        {
          cidrMask: config.cidrMask,
          name: `${prefix}-subnet-public`,
          subnetType: ec2.SubnetType.PUBLIC,
        },
        {
          cidrMask: config.cidrMask,
          name: `${prefix}-subnet-private`,
          subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
        },
      ],
    });

    Tags.of(this).add(PROJECT_TAG_KEY, prefix);
  }
}
