package (
    "testing"

    awsSdk "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    rds "github.com/aws/aws-sdk-go/service/rds"
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

    assert.Equal(t, "10.0.0.0/16", *vpcOutput.Vpcs[0].CidrBlock)

    dbInstanceIDs := terraform.OutputList(t, terraformOptions, "db_instance_ids")
    assert.Greater(t, len(dbInstanceIDs), 0, "DB instances must be created")

    rdsClient := rds.New(sess)

    for _, dbInstanceID := range dbInstanceIDs {
        dbInstancesOutput, err := rdsClient.DescribeDBInstances(&rds.DescribeDBInstancesInput{
            DBInstanceIdentifier: awsSdk.String(dbInstanceID),
        })
        if err != nil {
            t.Fatal(err)
        }
        assert.False(t, *dbInstancesOutput.DBInstances[0].PubliclyAccessible, "Database must NOT be publicly accessible")
    }
}
