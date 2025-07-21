package test

import (
	"testing"

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
	assert.Equal(t, "10.0.0.0/16", aws.GetCidrBlockFromVpc(t, vpc), "VPC CIDR should be 10.0.0.0/16")

	
	dbInstanceIDs := terraform.OutputList(t, terraformOptions, "db_instance_ids")
	assert.Greater(t, len(dbInstanceIDs), 0, "DB instances must be created")

	
	for _, dbInstanceID := range dbInstanceIDs {
		dbInfo := aws.GetRdsInstanceDetails(t, "us-east-1", dbInstanceID)
		assert.False(t, dbInfo.PubliclyAccessible, "Database must NOT be publicly accessible")
	}
}
