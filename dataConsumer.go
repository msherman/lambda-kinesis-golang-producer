package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"io/ioutil"
	"net/http"
	"os"
)

type Teams struct {
	Copyright string `json:"copyright"`
	Teams []team `json:"teams"`
}
type team struct {
	Id int `json:"id"`
	Name string `json:"name"`
	FirstYearOfPlay string `json:"firstYearOfPlay"`
	Active bool `json:"active"`
}

func HandleRequest(ctx context.Context) (string, error) {
	streamName := os.Getenv("STREAM")
	if streamName == "" {
		panic("No environment variable set. Set STREAM environment variable")
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	))

	// Create kinesis service client
	svc := kinesis.New(sess)
	resp, err := http.Get("https://statsapi.web.nhl.com/api/v1/teams")
	if err != nil {
		return fmt.Sprintf("had an error %s", err.Error()), err
	}
	b, processErr := ioutil.ReadAll(resp.Body)
	if processErr != nil {
		return fmt.Sprintf("had process error %s", processErr.Error()), processErr
	}
	var nhlteams Teams
	err = json.Unmarshal(b, &nhlteams)
	if err != nil {
		fmt.Println(err.Error())
	}
	var records kinesis.PutRecordsInput
	records.StreamName = &streamName
	for id, data := range nhlteams.Teams {
		var record kinesis.PutRecordsRequestEntry
		record.Data, err = json.Marshal(data)
		fmt.Println(data)
		idAsString := string(id)
		record.PartitionKey = &idAsString
		records.Records = append(records.Records, &record)
	}
	putRecordOutput, prErr :=  svc.PutRecords(&records)
	if prErr != nil {
		fmt.Println("kinesis failed")
		fmt.Println(putRecordOutput.String())
		fmt.Println(prErr)
	} else {
		fmt.Println(putRecordOutput.String())
	}
	return fmt.Sprintf("Produced to stream: %s", streamName), nil
}
func main() {
	lambda.Start(HandleRequest)
}
