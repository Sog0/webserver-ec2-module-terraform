package test

import (
    "testing"

    awsSdk "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    ec2 "github.com/aws/aws-sdk-go/service/ec2"

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

    sess := session.Must(session.NewSession(&awsSdk.Config{Region: awsSdk.String("us-east-1")}))
    ec2Client := ec2.New(sess)

    vpcOutput, err := ec2Client.DescribeVpcs(&ec2.DescribeVpcsInput{
        VpcIds: []*string{awsSdk.String(vpcID)},
    })
    if err != nil {
        t.Fatal(err)
    }
    if len(vpcOutput.Vpcs) == 0 {
        t.Fatal("VPC not found")
    }

    assert.Equal(t, "192.168.0.0/16", *vpcOutput.Vpcs[0].CidrBlock)

    dbInstanceIDs := terraform.OutputList(t, terraformOptions, "db_instance_ids")
    assert.Greater(t, len(dbInstanceIDs), 0, "DB instances must be created")

    for _, dbInstanceID := range dbInstanceIDs {
        ec2InstanceOutput, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
            InstanceIds: []*string{awsSdk.String(dbInstanceID)},
        })
        if err != nil {
            t.Fatal(err)
        }

        instances := ec2InstanceOutput.Reservations[0].Instances
        if len(instances) == 0 {
            t.Fatalf("EC2 Instance %s not found", dbInstanceID)
        }

        assert.Nil(t, instances[0].PublicIpAddress, "Database EC2 instance must NOT have a public IP")
    }
}
