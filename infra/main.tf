terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "us-east-1"
}


resource "aws_kinesis_stream" "ingestion_stream" {
  name = "ingestion_stream"
  shard_count = 1
  retention_period = 48
  encryption_type = "KMS"
  kms_key_id = "alias/aws/kinesis"

  shard_level_metrics = [
    "IncomingBytes",
    "OutgoingBytes",
  ]

  tags = {
    Environment = "ingestion-test"
  }
}

resource "aws_lambda_function" "test_lambda" {
  filename = "../function.zip"
  function_name = "golang_kinesis_producer"
  role = aws_iam_role.iam_ingestion_lambda.arn
  handler = "dataConsumer"

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = filebase64sha256("../function.zip")

  runtime = "go1.x"

  environment {
    variables = {
      STREAM = aws_kinesis_stream.ingestion_stream.name
    }
  }
}

resource "aws_iam_role" "iam_ingestion_lambda" {
  name = "iam_ingestion_lambda"

  assume_role_policy = jsonencode(
  {
    Version: "2012-10-17",
    Statement: [

      {
        Action: "sts:AssumeRole",
        Principal: {
          Service: "lambda.amazonaws.com"
        },
        Effect: "Allow",
        Sid: ""
      }
    ]
  })
}

resource "aws_iam_role_policy" "iam_ingestion_kinesis_policy" {
  name        = "kinesis_ingestion_policy"
  role = aws_iam_role.iam_ingestion_lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "kinesis:PutRecord",
          "kinesis:PutRecords"
        ]
        Effect   = "Allow"
        Resource = aws_kinesis_stream.ingestion_stream.arn
      },
    ]
  })
}