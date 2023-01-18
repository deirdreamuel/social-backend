package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfrastructureStackProps struct {
	awscdk.StackProps
}

func NewInfrastructureStack(scope constructs.Construct, id string, props *InfrastructureStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// S3 bucket for image upload
	bucket := awss3.NewBucket(stack, jsii.String("profile_pic_s3_bucket"), &awss3.BucketProps{
		BucketName: jsii.String("profile.image.amuel.org"),
		BlockPublicAccess: awss3.NewBlockPublicAccess(
			&awss3.BlockPublicAccessOptions{
				BlockPublicAcls:       jsii.Bool(false),
				BlockPublicPolicy:     jsii.Bool(false),
				IgnorePublicAcls:      jsii.Bool(false),
				RestrictPublicBuckets: jsii.Bool(false),
			},
		),
		Cors: &[]*awss3.CorsRule{
			{
				AllowedHeaders: jsii.Strings("*"),
				AllowedMethods: &[]awss3.HttpMethods{
					awss3.HttpMethods_PUT,
					awss3.HttpMethods_HEAD,
					awss3.HttpMethods_GET,
				},
				AllowedOrigins: jsii.Strings("*"),
				ExposedHeaders: &[]*string{},
			},
		},
	})

	// Add read policy for bucket and objects
	bucket.AddToResourcePolicy(
		awsiam.NewPolicyStatement(
			&awsiam.PolicyStatementProps{
				Resources: &[]*string{
					bucket.ArnForObjects(jsii.String("*")),
					bucket.BucketArn(),
				},
				Actions: &[]*string{
					jsii.String("s3:GetObject"),
				},
				Principals: &[]awsiam.IPrincipal{
					awsiam.NewAnyPrincipal(),
				},
			},
		),
	)

	// Create new policy for storing to bucket
	awsiam.NewPolicy(stack, jsii.String("profile_pic_put_policy"),
		&awsiam.PolicyProps{
			PolicyName: jsii.String("profile.image.amuel.org-put-policy"),
			Statements: &[]awsiam.PolicyStatement{
				awsiam.NewPolicyStatement(
					&awsiam.PolicyStatementProps{
						Resources: &[]*string{
							bucket.ArnForObjects(jsii.String("*")),
							bucket.BucketArn(),
						},
						Actions: &[]*string{
							jsii.String("s3:PutObject"),
						},
					},
				),
			},
		},
	)

	awscloudfront.NewDistribution(stack, jsii.String("profile_pic_cdn"),
		&awscloudfront.DistributionProps{
			DefaultBehavior: &awscloudfront.BehaviorOptions{
				Origin: awscloudfrontorigins.NewS3Origin(bucket, &awscloudfrontorigins.S3OriginProps{}),
			},
		},
	)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewInfrastructureStack(app, "InfrastructureStack", &InfrastructureStackProps{
		awscdk.StackProps{
			Env: &awscdk.Environment{
				Account: jsii.String("582250362323"),
				Region:  jsii.String("us-east-1"),
			},
		},
	})

	app.Synth(nil)
}
