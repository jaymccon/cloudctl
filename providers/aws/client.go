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
	"sync"
	"time"
)

var noWaitSleep = 3 * time.Second
var bar *pb.ProgressBar

func FetchSchemas() (*map[string][]byte, error) {
	schemas := map[string][]byte{}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRetryer(func() aws.Retryer {
		return retry.NewStandard(func(opts *retry.StandardOptions) {
			opts.MaxAttempts = 20
			opts.MaxBackoff = 60 * time.Second
			opts.RetryCost = 1
		})
	}))
	if err != nil {
		return nil, err
	}
	cfn := cloudformation.NewFromConfig(cfg)
	categories := []types.Category{"THIRD_PARTY", "AWS_TYPES"}
	var typeArns []*string
	for _, c := range categories {
		params := &cloudformation.ListTypesInput{
			Filters: &types.TypeFilters{
				Category: c,
			},
			ProvisioningType: "FULLY_MUTABLE",
			Type:             "RESOURCE",
			Visibility:       "PUBLIC",
		}

		paginator := cloudformation.NewListTypesPaginator(cfn, params, func(o *cloudformation.ListTypesPaginatorOptions) {})

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
	}
	bar = pb.StartNew(len(typeArns))
	typeDescriptions, err := asyncCfnDescribeType(*cfn, typeArns)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	for _, value := range typeDescriptions {
		schemas[*value.TypeName] = []byte(*value.Schema)
	}
	bar.Finish()
	return &schemas, nil
}

func NewCcClient() (*cloudcontrol.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRetryer(func() aws.Retryer {
		return retry.AddWithMaxAttempts(retry.NewStandard(), 20)
	}))
	if err != nil {
		fmt.Printf("ERROR: %q", err.Error())
		return nil, err
	}
	cc := cloudcontrol.NewFromConfig(cfg, func(options *cloudcontrol.Options) {
		options.APIOptions = append(options.APIOptions, attachCustomMiddleware)
	})
	return cc, nil
}

func ListResource(typeName string) (*[]typesCC.ResourceDescription, error) {
	cc, err := NewCcClient()
	if err != nil {
		fmt.Printf("ERROR: %q", err.Error())
		return nil, err
	}
	resp, err := cc.ListResources(context.TODO(), &cloudcontrol.ListResourcesInput{TypeName: &typeName})
	if err != nil {
		fmt.Printf("ERROR: %q", err.Error())
		return nil, err
	}
	return &resp.ResourceDescriptions, nil
}

func DeleteResource(cc *cloudcontrol.Client, typeName string, id string, async bool) error {
	resp, err := cc.DeleteResource(
		context.TODO(),
		&cloudcontrol.DeleteResourceInput{TypeName: &typeName, Identifier: &id},
	)
	if err != nil {
		fmt.Printf("ERROR: %q", err.Error())
		return err
	}
	// nil timeout will wait until operation completes
	var timeout *time.Time = nil
	timeoutTime := time.Now().Add(noWaitSleep)
	if async {
		// most errors come back really quickly, so we wait a bit even if wait is disabled
		timeout = &timeoutTime
	}
	pe, err := waitForComplete(*cc, *resp.ProgressEvent, timeout)
	if err != nil {
		fmt.Printf("ERROR: %q", err.Error())
		return err
	}
	if pe.OperationStatus == typesCC.OperationStatusFailed {
		fmt.Printf("ERROR: %s with identifier %q failed. [%s] %s\n", typeName, id, pe.ErrorCode, *pe.StatusMessage)
	}
	fmt.Printf("%s delete %s for resource with the identifier %q\n", typeName, pe.OperationStatus, id)
	if async && !isFinished(*pe) {
		fmt.Printf("Request token: %s\n", *pe.RequestToken)
	}
	return nil
}

func isFinished(pe typesCC.ProgressEvent) bool {
	var finalStatuses = []typesCC.OperationStatus{
		typesCC.OperationStatusSuccess,
		typesCC.OperationStatusFailed,
		typesCC.OperationStatusCancelComplete,
	}
	for _, s := range finalStatuses {
		if pe.OperationStatus == s {
			return true
		}
	}
	return false
}

func waitForComplete(cc cloudcontrol.Client, pe typesCC.ProgressEvent, timeout *time.Time) (*typesCC.ProgressEvent, error) {
	if !isFinished(pe) {
		if timeout != nil {
			if time.Now().After(*timeout) {
				return &pe, nil
			}
		}
		time.Sleep(2 * time.Second)
		resp, err := cc.GetResourceRequestStatus(
			context.TODO(),
			&cloudcontrol.GetResourceRequestStatusInput{RequestToken: pe.RequestToken},
		)
		if err != nil {
			return resp.ProgressEvent, err
		}
		return waitForComplete(cc, *resp.ProgressEvent, timeout)
	}
	return &pe, nil
}

const concurrentAwsCalls = 4

func streamInputs(done <-chan struct{}, inputs []*string) <-chan string {
	inputCh := make(chan string)
	go func() {
		defer close(inputCh)
		for _, input := range inputs {
			select {
			case inputCh <- *input:
			case <-done:
				break
			}
		}
	}()
	return inputCh
}

func describeType(client cloudformation.Client, arn string) (*cloudformation.DescribeTypeOutput, error) {
	typeDesc, err := client.DescribeType(context.TODO(), &cloudformation.DescribeTypeInput{
		Arn: &arn,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	bar.Increment()
	return typeDesc, nil
}

type result struct {
	response *cloudformation.DescribeTypeOutput
	err      error
}

func asyncCfnDescribeType(client cloudformation.Client, typeArns []*string) ([]*cloudformation.DescribeTypeOutput, error) {
	done := make(chan struct{})
	defer close(done)

	inputCh := streamInputs(done, typeArns)

	var wg sync.WaitGroup
	wg.Add(concurrentAwsCalls)

	resultCh := make(chan result)

	for i := 0; i < concurrentAwsCalls; i++ {
		go func() {
			for input := range inputCh {
				resp, err := describeType(client, input)
				resultCh <- result{resp, err}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results []*cloudformation.DescribeTypeOutput
	for result := range resultCh {
		if result.err != nil {
			return nil, result.err
		}
		results = append(results, result.response)
	}
	return results, nil
}

type deleteResourceErrors struct {
	id  string
	err error
}

func AsyncCcDeleteResource(client cloudcontrol.Client, typeName string, resourceIds []string, async bool) {
	done := make(chan struct{})
	defer close(done)

	var strSlice []*string

	for _, i := range resourceIds {
		strSlice = append(strSlice, &i)
	}
	inputCh := streamInputs(done, strSlice)

	var wg sync.WaitGroup
	wg.Add(concurrentAwsCalls)

	resultCh := make(chan deleteResourceErrors)

	for i := 0; i < concurrentAwsCalls; i++ {
		go func() {
			for input := range inputCh {
				err := DeleteResource(&client, input, typeName, async)
				resultCh <- deleteResourceErrors{input, err}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()
	for e := range resultCh {
		if e.err != nil {
			fmt.Printf("deleting %s with identifier %s failed: %s", typeName, e.id, e.err.Error())
		}
	}
}
