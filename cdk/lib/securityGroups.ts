import { aws_ec2 as ec2 } from "aws-cdk-lib";

import { EC2Config } from "./buildConfig";

export class AssignSecurityGroup {
  private ecsSecurityGroup: ec2.ISecurityGroup;

  private rdsSecurityGroup: ec2.ISecurityGroup;

  private ec2SecurityGroup: ec2.ISecurityGroup;

  constructor(
    ecsSecurityGroup: ec2.ISecurityGroup,
    rdsSecurityGroup: ec2.ISecurityGroup,
    ec2SecurityGroup: ec2.ISecurityGroup,
  ) {
    this.ecsSecurityGroup = ecsSecurityGroup;
    this.rdsSecurityGroup = rdsSecurityGroup;
    this.ec2SecurityGroup = ec2SecurityGroup;
  }

  /**
   * セキュリティグループの設定
   * @param ec2Config allowFromが指定されているとき、PhpMyAdminはそのIPアドレスからのみしか接続できなくなる
   */
  assign(ec2Config: EC2Config) {
    this.ecsSecurityGroup.connections.allowTo(
      this.rdsSecurityGroup,
      ec2.Port.tcp(3306),
      "ECS-RDS",
    );
    this.ec2SecurityGroup.connections.allowTo(
      this.rdsSecurityGroup,
      ec2.Port.tcp(3306),
      "EC2-RDS",
    );
    this.ec2SecurityGroup.addIngressRule(
      ec2.Peer.anyIpv4(),
      ec2.Port.tcp(22),
      "SSH connection to EC2",
    );
    // 制限させたいIPアドレスの有無でphpMyAdminの接続設定を変える
    if (ec2Config.IpAddressToDBClient != null) {
      this.ec2SecurityGroup.addIngressRule(
        ec2.Peer.ipv4(ec2Config.IpAddressToDBClient),
        ec2.Port.tcp(8080),
        "Allow access to PhpMyAdmin from specific IP",
      );
    } else {
      this.ec2SecurityGroup.addIngressRule(
        ec2.Peer.anyIpv4(),
        ec2.Port.tcp(8080),
        "Allow access to PhpMyAdmin from all users",
      );
    }
  }
}
