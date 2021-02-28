resource "aws_lambda_function" "liftplan" {
  function_name = var.project_name
  role          = aws_iam_role.lambda.arn
  handler       = var.project_name
  publish       = true
  s3_bucket     = var.lambda_s3_bucket
  s3_key        = "${var.project_name}/function.zip"

  memory_size = 128
  timeout     = 3
  runtime     = "go1.x"

  lifecycle {
    ignore_changes = [last_modified, qualified_arn, version]
  }
}
