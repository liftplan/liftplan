resource "aws_apigatewayv2_api" "api" {
  name          = var.project_name
  protocol_type = "HTTP"
}


resource "aws_apigatewayv2_domain_name" "api" {
  domain_name = "api.${var.domain_name}"

  domain_name_configuration {
    certificate_arn = data.aws_acm_certificate.primary.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_stage" "api" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_apigatewayv2_api_mapping" "api" {
  api_id      = aws_apigatewayv2_api.api.id
  domain_name = aws_apigatewayv2_domain_name.api.id
  stage       = aws_apigatewayv2_stage.api.id
}

resource "aws_apigatewayv2_integration" "root_get" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "AWS_PROXY"
  connection_type        = "INTERNET"
  description            = ""
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.liftplan.arn
  passthrough_behavior   = "WHEN_NO_MATCH"
  payload_format_version = "2.0"
  request_parameters     = {}
  request_templates      = {}
}

resource "aws_apigatewayv2_integration" "plan_get" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "AWS_PROXY"
  connection_type        = "INTERNET"
  description            = ""
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.liftplan.arn
  passthrough_behavior   = "WHEN_NO_MATCH"
  payload_format_version = "2.0"
  request_parameters     = {}
  request_templates      = {}
}

resource "aws_apigatewayv2_integration" "plan_post" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "AWS_PROXY"
  connection_type        = "INTERNET"
  description            = ""
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.liftplan.arn
  passthrough_behavior   = "WHEN_NO_MATCH"
  payload_format_version = "2.0"
  request_parameters     = {}
  request_templates      = {}
}
