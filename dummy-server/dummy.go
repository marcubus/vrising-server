package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

func main() {
	logger := log.New(os.Stdout, "[dummy-server] ", log.LstdFlags)

	udpServer(logger)
}

func udpServer(logger *log.Logger) {
	host := "0.0.0.0"
	port := 9876 // Steam connection port
	blockSize := 1024

	ip := net.ParseIP(host)
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	logger.Printf("listening on addr=%s with block size=%d", listener.LocalAddr(), blockSize)

	data := make([]byte, blockSize)
	for {
		// UDP Listener
		_, remoteAddr, err := listener.ReadFrom(data)
		if err != nil {
			logger.Fatalf("error during read: %s", err)
		}

		logger.Printf("Connection request from <%s>", remoteAddr)
		enableRealServer()
		os.Exit(0)
	}
}

func enableRealServer() {
	targetService := os.Getenv("TARGET_SERVICE")
	dummyService := os.Getenv("DUMMY_SERVICE")
	targetCluster := os.Getenv("TARGET_CLUSTER")
	dummyCluster := os.Getenv("DUMMY_CLUSTER")
	targetAsg := os.Getenv("TARGET_ASG")
	dummyAsg := os.Getenv("DUMMY_ASG")

	changeAsgCount(targetAsg, 1)
	changeServiceCount(targetCluster, targetService, 1)
	changeServiceCount(dummyCluster, dummyService, 0)
	changeAsgCount(dummyAsg, 0)
}

func changeServiceCount(cluster string, service string, count int) {
	region := os.Getenv("AWS_REGION")
	svc := ecs.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	input := &ecs.UpdateServiceInput{
			DesiredCount: aws.Int64(int64(count)),
			Service:      aws.String(service),
			Cluster:			aws.String(cluster),
	}

	result, err := svc.UpdateService(input)
	if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case ecs.ErrCodeServerException:
							fmt.Println(ecs.ErrCodeServerException, aerr.Error())
					case ecs.ErrCodeClientException:
							fmt.Println(ecs.ErrCodeClientException, aerr.Error())
					case ecs.ErrCodeInvalidParameterException:
							fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
					case ecs.ErrCodeClusterNotFoundException:
							fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
					case ecs.ErrCodeServiceNotFoundException:
							fmt.Println(ecs.ErrCodeServiceNotFoundException, aerr.Error())
					case ecs.ErrCodeServiceNotActiveException:
							fmt.Println(ecs.ErrCodeServiceNotActiveException, aerr.Error())
					case ecs.ErrCodePlatformUnknownException:
							fmt.Println(ecs.ErrCodePlatformUnknownException, aerr.Error())
					case ecs.ErrCodePlatformTaskDefinitionIncompatibilityException:
							fmt.Println(ecs.ErrCodePlatformTaskDefinitionIncompatibilityException, aerr.Error())
					case ecs.ErrCodeAccessDeniedException:
							fmt.Println(ecs.ErrCodeAccessDeniedException, aerr.Error())
					default:
							fmt.Println(aerr.Error())
					}
			} else {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
			}
			return
	}

	fmt.Println(result)
}

func changeAsgCount(asg string, count int) {
	region := os.Getenv("AWS_REGION")
	svc := autoscaling.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	input := &autoscaling.SetDesiredCapacityInput{
			AutoScalingGroupName: aws.String(asg),
			DesiredCapacity:      aws.Int64(int64(count)),
	}

	result, err := svc.SetDesiredCapacity(input)
	if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case autoscaling.ErrCodeScalingActivityInProgressFault:
							fmt.Println(autoscaling.ErrCodeScalingActivityInProgressFault, aerr.Error())
					case autoscaling.ErrCodeResourceContentionFault:
							fmt.Println(autoscaling.ErrCodeResourceContentionFault, aerr.Error())
					case autoscaling.ErrCodeServiceLinkedRoleFailure:
							fmt.Println(autoscaling.ErrCodeServiceLinkedRoleFailure, aerr.Error())
					default:
							fmt.Println(aerr.Error())
					}
			} else {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
			}
			return
	}

	fmt.Println(result)
}
