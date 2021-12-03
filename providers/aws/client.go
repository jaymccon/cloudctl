package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	typesCC "github.com/aws/aws-sdk-go-v2/service/cloudcontrol/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/cheggaaa/pb/v3"
	"log"
)

func FetchSchemas() (*map[string][]byte, error) {
	schemas := map[string][]byte{}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRetryer(func() aws.Retryer {
		return retry.AddWithMaxAttempts(retry.NewStandard(), 20)
	}))

	if err != nil {
		return nil, err
	}
	cfn := cloudformation.NewFromConfig(cfg)
	params := &cloudformation.ListTypesInput{
		Filters: &types.TypeFilters{
			Category: "AWS_TYPES",
		},
		ProvisioningType: "FULLY_MUTABLE",
		Type:             "RESOURCE",
		Visibility:       "PUBLIC",
	}

	paginator := cloudformation.NewListTypesPaginator(cfn, params, func(o *cloudformation.ListTypesPaginatorOptions) {})
	var typeArns []*string
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Printf("error: %v", err)
			return nil, err
		}
		for _, value := range output.TypeSummaries {
			typeArns = append(typeArns, value.TypeArn)
		}
	}
	bar := pb.StartNew(len(typeArns))
	for _, value := range typeArns {
		typeDesc, err := cfn.DescribeType(context.TODO(), &cloudformation.DescribeTypeInput{
			Arn: value,
		})
		if err != nil {
			log.Printf("error: %v", err)
			return nil, err
		}
		schemas[*typeDesc.TypeName] = []byte(*typeDesc.Schema)
		bar.Increment()
	}
	bar.Finish()
	return &schemas, nil
}

func ListResource(typeName string) (*[]typesCC.ResourceDescription, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRetryer(func() aws.Retryer {
		return retry.AddWithMaxAttempts(retry.NewStandard(), 20)
	}))
	if err != nil {
		fmt.Printf("ERROR: %q", err.Error())
		return nil, err
	}
	cc := cloudcontrol.NewFromConfig(cfg)
	resp, err := cc.ListResources(context.TODO(), &cloudcontrol.ListResourcesInput{TypeName: &typeName})
	if err != nil {
		fmt.Printf("ERROR: %q", err.Error())
		return nil, err
	}
	return &resp.ResourceDescriptions, nil
}
