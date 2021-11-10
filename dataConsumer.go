package main

import (
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
	Teams []Team `json:"teams"`
}
type Team struct {
	Id int `json:"id"`
	Name string `json:"name"`
	FirstYearOfPlay string `json:"firstYearOfPlay"`
	Active bool `json:"active"`
}

var STREAM_NAME string = ""

func HandleRequest() (string, error) {
	nhlteams := getNHLTeamRecords()
	svc := connectToKinesis()
	records := buildPutRecordRequest(nhlteams)
	putRecordOutput, prErr :=  svc.PutRecords(records)
	if prErr != nil {
		fmt.Println("kinesis failed")
		fmt.Println(putRecordOutput.String())
		fmt.Println(prErr)
	} else {
		fmt.Println(putRecordOutput.String())
	}
	return fmt.Sprintf("Produced to stream: %s", STREAM_NAME), nil
}

func connectToKinesis() *kinesis.Kinesis {
	STREAM_NAME = os.Getenv("STREAM")
	if STREAM_NAME == "" {
		panic("No environment variable set. Set STREAM environment variable")
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	))

	return kinesis.New(sess)
}

func getNHLTeamRecords() *Teams {
	resp, err := http.Get("https://statsapi.web.nhl.com/api/v1/teams")
	if err != nil {
		panic(fmt.Sprintf("had an error %s", err.Error()))
	}
	b, processErr := ioutil.ReadAll(resp.Body)
	if processErr != nil {
		panic(fmt.Sprintf("had process error %s", processErr.Error()))
	}
	var nhlteams Teams
	err = json.Unmarshal(b, &nhlteams)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &nhlteams
}

func buildPutRecordRequest(nhlteams *Teams) *kinesis.PutRecordsInput {
	var records kinesis.PutRecordsInput
	var err error
	records.StreamName = &STREAM_NAME
	for id, data := range nhlteams.Teams {
		var record kinesis.PutRecordsRequestEntry
		record.Data, err = json.Marshal(data)
		if err != nil {
			panic(err)
		}
		idAsString := string(id)
		record.PartitionKey = &idAsString
		records.Records = append(records.Records, &record)
	}
	return &records
}

func main() {
	lambda.Start(HandleRequest)
}
