package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type event struct {
	Action string `json:"action"`
	App    string `json:"app"`
}

func handler(e event) {
	switch e.Action {
	case "Update":
		update(e)
	}
}

func update(e event) {
	svc := ec2.New(session.New())
	filters := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Environment"),
				Values: []*string{aws.String(e.App)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
			{
				Name:   aws.String("platform"),
				Values: nil,
			},
		},
	}
	resp, err := svc.DescribeInstances(filters)
	check(err)

	for idx, res := range resp.Reservations {
		fmt.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("    - Instance ID: ", *inst.InstanceId)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	log.Println("Trigger")
	lambda.Start(handler)
}
