

output "lambda_function_arn" {
  value = aws_lambda_function.liftplan.arn
}

output "lambda_iam_role_arn" {
  value = aws_iam_role.lambda.arn
}