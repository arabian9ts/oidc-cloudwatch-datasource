package plugin

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

func newConfig(ctx context.Context, opts ...func(*config.LoadOptions) error) (aws.Config, error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("oidc", "oidc", ""))
	opts = append(opts, config.WithCredentialsProvider(creds))
	return config.LoadDefaultConfig(ctx, opts...)
}

func newCloudWatch(cfg aws.Config) *cloudwatch.Client {
	return cloudwatch.NewFromConfig(cfg)
}
