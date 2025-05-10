import "source-map-support/register";
import { Stack, StackProps, Tags } from "aws-cdk-lib";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import { Construct } from "constructs";

import { EC2Config, PROJECT_TAG_KEY } from "./buildConfig";

export interface EC2BastionStackProps extends StackProps {
  readonly prefix: string;
  readonly vpc: ec2.IVpc;
  readonly config: EC2Config;
}

export class EC2BastionStack extends Stack {
  public readonly ec2SecurityGroup: ec2.ISecurityGroup;

  public readonly ec2Instance: ec2.Instance;

  private readonly ec2KeyPair: ec2.CfnKeyPair;

  constructor(scope: Construct, id: string, props: EC2BastionStackProps) {
    super(scope, id, props);
    const { prefix, vpc, config } = props;

    // EC2用のSSH鍵を作成
    // パラメータストアに/ec2/keypair/{key_pair_id}の名前で作成される。
    //! 現状はキーの名前から鍵を特定しにくいので作成後idを控えておく
    this.ec2KeyPair = new ec2.CfnKeyPair(this, `${prefix}-ec2-key-pair`, {
      keyName: `${prefix}-bastion-key-pair`,
    });

    this.ec2SecurityGroup = new ec2.SecurityGroup(this, `${prefix}-ec2-sg`, {
      vpc,
      securityGroupName: `${prefix}-ec2-sg`,
      allowAllOutbound: true,
    });

    const machineImage = ec2.MachineImage.fromSsmParameter(
      "/aws/service/canonical/ubuntu/server/focal/stable/current/amd64/hvm/ebs-gp2/ami-id",
    );
    const userData = ec2.UserData.forLinux({ shebang: "#!/bin/bash -xe" });
    userData.addCommands(
      "sudo apt-get update -y",
      "sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common",
      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
      'sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"',
      "sudo apt-get update -y",
      "sudo apt-get install -y docker-ce",
      "sudo usermod -aG docker ubuntu",
    );

    this.ec2Instance = new ec2.Instance(this, `${prefix}-ec2-instance`, {
      machineImage,
      vpc,
      vpcSubnets: {
        subnetType: ec2.SubnetType.PUBLIC,
      },
      securityGroup: this.ec2SecurityGroup,
      instanceType: new ec2.InstanceType("t2.micro"),
      keyName: this.ec2KeyPair.keyName,
      userData,
    });

    if (config.useElasticIP) {
      // ElasticIPの設定
      // eslint-disable-next-line no-new
      new ec2.CfnEIP(this, `${prefix}-ec2-eip`, {
        instanceId: this.ec2Instance.instanceId,
      });
    }

    Tags.of(this).add(PROJECT_TAG_KEY, prefix);
  }
}
