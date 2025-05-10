import * as cdk from "aws-cdk-lib";
import { Match, Template } from "aws-cdk-lib/assertions";

import { RDSConfig, prefix } from "../lib/buildConfig";
import { RDSStack } from "../lib/rdsStack";
import { VPCStack } from "../lib/vpcStack";

import {
  useAlarmConfig,
  useConfig,
  useEC2Config,
  useECSConfig,
  useRDSConfig,
  useVPCConfig,
} from "./config.test";

const createTemplate = (rdsConfig: RDSConfig) => {
  const app = new cdk.App();

  const testConfig = useConfig(
    useVPCConfig({}),
    rdsConfig,
    useECSConfig({}),
    useEC2Config({}),
    useAlarmConfig({}),
  );
  const testPrefix = prefix(testConfig);

  const testVpcStack = new VPCStack(app, `${testPrefix}-vpc-stack`, {
    prefix: testPrefix,
    config: testConfig.vpc,
  });
  const testRdsStack = new RDSStack(app, `${testPrefix}-rds-stack`, {
    prefix: testPrefix,
    config: testConfig.rds,
    vpc: testVpcStack.vpc,
  });
  return Template.fromStack(testRdsStack);
};

describe("isAuroraEnabled: true", () => {
  const rdsConfig = useRDSConfig({ isAuroraEnabled: true });
  const template = createTemplate(rdsConfig);
  test("create rds stack: db cluster", () => {
    template.resourceCountIs("AWS::RDS::DBCluster", 1);
    template.hasResourceProperties("AWS::RDS::DBCluster", {
      Engine: "aurora-mysql",
      BackupRetentionPeriod: 20,
      DatabaseName: "testdatabase",
      DBClusterParameterGroupName: Match.anyValue(),
      DBSubnetGroupName: Match.anyValue(),
      MasterUsername: Match.anyValue(),
      MasterUserPassword: Match.anyValue(),
      Port: 3306,
      VpcSecurityGroupIds: Match.anyValue(),
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });
  });

  test("create rds stack: db instance", () => {
    template.resourceCountIs("AWS::RDS::DBInstance", 2);
    template.hasResourceProperties("AWS::RDS::DBInstance", {
      AllowMajorVersionUpgrade: false,
      AutoMinorVersionUpgrade: true,
      CACertificateIdentifier: "rds-ca-rsa2048-g1",
      DBClusterIdentifier: Match.anyValue(),
      DBInstanceClass: "db.t2.micro",
      DBInstanceIdentifier: "test-env-db-writer",
      Engine: "aurora-mysql",
      PromotionTier: 0,
      PubliclyAccessible: false,
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });

    template.hasResourceProperties("AWS::RDS::DBInstance", {
      CACertificateIdentifier: "rds-ca-rsa2048-g1",
      DBClusterIdentifier: Match.anyValue(),
      DBInstanceClass: "db.serverless",
      Engine: "aurora-mysql",
      PromotionTier: 2,
      PubliclyAccessible: false,
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });
  });

  test("create rds stack: parameter group", () => {
    template.hasResourceProperties("AWS::RDS::DBClusterParameterGroup", {
      Description: "Cluster parameter group for aurora-mysql8.0",
      Family: "aurora-mysql8.0",
      Parameters: {
        character_set_client: "utf8mb4",
        character_set_connection: "utf8mb4",
        character_set_database: "utf8mb4",
        character_set_results: "utf8mb4",
        character_set_server: "utf8mb4",
      },
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });
  });
});

describe("isAuroraEnabled: false", () => {
  const rdsConfig = useRDSConfig({ isAuroraEnabled: false });
  const template = createTemplate(rdsConfig);
  test("create rds stack: db instance", () => {
    template.resourceCountIs("AWS::RDS::DBInstance", 1);
    template.hasResourceProperties("AWS::RDS::DBInstance", {
      DBInstanceClass: "db.t2.micro",
      DBSubnetGroupName: Match.anyValue(),
      Engine: "mysql",
      PubliclyAccessible: false,
      AllocatedStorage: "100",
      CopyTagsToSnapshot: true,
      DBName: "testdatabase",
      DBParameterGroupName: Match.anyValue(),
      EngineVersion: "8.0",
      MasterUsername: Match.anyValue(),
      MasterUserPassword: Match.anyValue(),
      Port: "3306",
      StorageType: "gp2",
      VPCSecurityGroups: [Match.anyValue()],
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });
  });

  test("create rds stack: parameter group", () => {
    template.hasResourceProperties("AWS::RDS::DBParameterGroup", {
      Description: "Parameter group for mysql8.0",
      Family: "mysql8.0",
      Parameters: {
        character_set_client: "utf8mb4",
        character_set_connection: "utf8mb4",
        character_set_database: "utf8mb4",
        character_set_results: "utf8mb4",
        character_set_server: "utf8mb4",
      },
      Tags: [
        {
          Key: "project",
          Value: "test-env",
        },
      ],
    });
  });
});
