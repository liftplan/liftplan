resource "aws_lambda_function" "liftplan" {
    function_name = var.project_name
    role          = aws_iam_role.lambda.arn
    handler       = "main"
    runtime       = "go1.x"
}