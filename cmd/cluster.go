// TODO: this command should be generated from the resource schema
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudcontrolapi"
	"github.com/spf13/cobra"
)

var cluster_type string
var dbname string
var masterusername string
var masterpassword string
var nodetype string

type RedshiftCluster struct {
	ClusterType        string
	DBName             string
	MasterUserPassword string
	MasterUsername     string
	NodeType           string
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Specifies a cluster",
	Long: `A cluster is a fully managed data warehouse that consists of a set of compute nodes.

To create a cluster in Virtual Private Cloud (VPC), you must provide a cluster subnet group name. The cluster subnet 
group identifies the subnets of your VPC that Amazon Redshift uses when creating the cluster. For more information 
about managing clusters, go to Amazon Redshift Clusters in the Amazon Redshift Cluster Management Guide.`,
	Run: func(cmd *cobra.Command, args []string) {
		desired_state := RedshiftCluster{
			ClusterType:        cluster_type,
			DBName:             dbname,
			MasterUserPassword: masterpassword,
			MasterUsername:     masterusername,
			NodeType:           nodetype,
		}
		fmt.Println("cluster called")
		CreateResource(desired_state)
	},
}

func CreateResource(cluster RedshiftCluster) {
	cc := cloudcontrolapi.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))
	jsonB, err := json.Marshal(cluster)
	fmt.Println(string(jsonB))
	if err != nil {
		panic(err)
	}
	input := cloudcontrolapi.CreateResourceInput{
		DesiredState: aws.String(string(jsonB)),
		TypeName:     aws.String("AWS::Redshift::Cluster"),
	}
	resp, err := cc.CreateResourceWithContext(context.TODO(), &input)
	if err != nil {
		fmt.Println(err)
	}
	if *resp.ProgressEvent.OperationStatus == "FAILED" {
		fmt.Println("Operation failed: " + *resp.ProgressEvent.StatusMessage)
	}
	progress := resp.ProgressEvent
	if *progress.OperationStatus == "IN_PROGRESS" {
		for progress.Identifier == nil && *progress.OperationStatus == "IN_PROGRESS" {
			time.Sleep(3 * time.Second)
			resp, err := cc.GetResourceRequestStatusWithContext(context.TODO(), &cloudcontrolapi.GetResourceRequestStatusInput{RequestToken: resp.ProgressEvent.RequestToken})
			if err != nil {
				fmt.Println(err)
			}
			progress = resp.ProgressEvent
		}
	}
	if progress.OperationStatus == nil || progress.Identifier == nil || progress.RequestToken == nil {
		fmt.Println("one of the needed properties is nil")
	} else {
		fmt.Printf("Cluster creation is %q with identifier %q and request token %q\n", *progress.OperationStatus, *progress.Identifier, *progress.RequestToken)
	}
}

func init() {
	redshiftCmd.AddCommand(clusterCmd)

	flags := clusterCmd.Flags()
	flags.StringVarP(
		&cluster_type,
		"type",
		"t",
		"single-node",
		"The type of the cluster. Valid Values: multi-node | single-node",
	)
	flags.StringVarP(
		&dbname,
		"dbname",
		"d",
		"dev",
		"The name of the first database to be created when the cluster is created",
	)
	flags.StringVarP(
		&masterusername,
		"master-username",
		"u",
		"",
		"The user name associated with the admin user account for the cluster that is being created",
	)
	flags.StringVarP(
		&masterpassword,
		"master-userpassword",
		"p",
		"",
		"The password associated with the admin user account for the cluster that is being created",
	)
	flags.StringVar(
		&nodetype,
		"node-type",
		"dc2.large",
		"The node type to be provisioned for the cluster",
	)
	clusterCmd.MarkFlagRequired("master-username")
	clusterCmd.MarkFlagRequired("master-userpassword")
}
