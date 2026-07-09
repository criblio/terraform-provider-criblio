# User-supplied parameters
#
# Supply values for the declared variables in this file in terraform.tfvars.
#
# Required for the Worker Group:
# - worker_group_id: ID for the new Worker Group (must not already exist)
#
# Required for the S3 Destination:
# - aws_api_key: AWS Access Key ID
# - aws_secret_key: AWS Secret Access Key
# - aws_bucket_name: S3 bucket name
# - aws_region: S3 bucket region, such as us-east-2

variable "worker_group_id" {
  type        = string
  description = "ID for the new Worker Group. Must not already exist in your deployment."
}

variable "aws_api_key" {
  type        = string
  description = "AWS Access Key ID for the S3 Destination"
  sensitive   = true
}

variable "aws_secret_key" {
  type        = string
  description = "AWS Secret Access Key for the S3 Destination"
  sensitive   = true
}

variable "aws_bucket_name" {
  type        = string
  description = "S3 bucket name for the S3 Destination"
}

variable "aws_region" {
  type        = string
  description = "AWS region for the S3 bucket, such as us-east-2"
}
