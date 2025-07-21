package test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestCheckInfrastructure(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../infrastructure",
	}

	terraform.Init(t, terraformOptions)

	instanceIDs := terraform.OutputList(t, terraformOptions, "instance_ids")
	assert.Greater(t, len(instanceIDs), 0, "EC2 instances must be created")

	vpcID := terraform.Output(t, terraformOptions, "vpc_id")
	vpc := aws.GetVpcById(t, vpcID, "us-east-1")
	assert.Equal(t, "10.0.0.0/16", *vpc.CidrBlock) // исправлено

	dbInstanceIDs := terraform.OutputList(t, terraformOptions, "db_instance_ids")
	assert.Greater(t, len(dbInstanceIDs), 0, "DB instances must be created")

	rdsClient := aws.NewRdsClient(t, "us-east-1")

	for _, dbInstanceID := range dbInstanceIDs {
		dbInstancesOutput, err := rdsClient.DescribeDBInstances(&rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(dbInstanceID),
		})

		if err != nil {
			t.Fatal(err)
		}

		assert.False(t, *dbInstancesOutput.DBInstances[0].PubliclyAccessible, "Database must NOT be publicly accessible")
	}
}
