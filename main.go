package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type event struct {
	Action      string `json:"action"`
	Environment string `json:"environment"`
}

func filter(e event) map[string]bool {
	svc := ec2.New(session.New())
	include, _ := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Environment"),
				Values: []*string{aws.String(e.Environment)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	})

	exclude, _ := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("platform"),
				Values: []*string{aws.String("windows")},
			},
		},
	})

	exclude2, _ := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:AutoPatch"),
				Values: []*string{aws.String("ignore")},
			},
		},
	})

	excludeMap := make(map[string]bool)
	excludeMap = filterMapper(exclude, excludeMap)
	excludeMap = filterMapper(exclude2, excludeMap)

	includeMap := make(map[string]bool)
	includeMap = filterMapper(include, includeMap)

	for key := range includeMap {
		if excludeMap[key] {
			delete(includeMap, key)
		}
	}

	return includeMap
}

func filterMapper(ec2 *ec2.DescribeInstancesOutput, instanceMap map[string]bool) map[string]bool {
	for _, reservation := range ec2.Reservations {
		for _, instance := range reservation.Instances {
			instanceMap[*instance.InstanceId] = true
		}
	}
	return instanceMap
}

func handler(e event) {
	instances := filter(e)
	fmt.Print(instances)
}

func main() {
	lambda.Start(handler)
}
