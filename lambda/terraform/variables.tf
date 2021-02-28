variable "aws_region" {
  type    = string
  default = "us-east-2"
}

variable "project_name" {
  type    = string
  default = "liftplan"
}

variable "lambda_s3_bucket" {
  type = string
}
