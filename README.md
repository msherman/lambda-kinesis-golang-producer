# Using golang to produce data to kinesis data stream

## What is this
The idea behind this repo was to quickly determine how easy it would be to add a serverless function
to consume an API and add the data to a kinesis data stream.

## To build
1. Install golang: https://golang.org/doc/install
2. Compile the program: `GOOS=linux go build dataConsumer.go`
3. Zip up the program: `zip function.zip dataConsumer`

## To Deploy
1. Install terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli
2. Navigate to infra folder: `cd infra`
3. Apply the infrastructure: `terraform apply --auto-approve`

## Whats next?
Up next would be to either consume in a lambda function and perform some sort of ETL work
prior to storing in to a database.  
The other option is to send it to AWS GLUE to perform the ETL to continue learning the data
ingestion then storing in to a database.