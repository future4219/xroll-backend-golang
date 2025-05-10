/* eslint-disable no-new */
import { Stack, aws_secretsmanager as sm, StackProps, Tags } from "aws-cdk-lib";
import { Construct } from "constructs";
import { PROJECT_TAG_KEY } from "./buildConfig";

export interface SecretsManagerStackProps extends StackProps {
  readonly prefix: string;
}

export class SecretsManagerStack extends Stack {
  public readonly secret: sm.Secret;

  constructor(scope: Construct, id: string, props: SecretsManagerStackProps) {
    super(scope, id);
    const { prefix } = props;
    this.secret = new sm.Secret(this, `${prefix}-secret`, {
      secretName: `${prefix}-secret`,
      secretObjectValue: {},
    });

    Tags.of(this).add(PROJECT_TAG_KEY, prefix);
  }
}
