
resource "aws_api_gateway_rest_api" "Cloud-Resume-API-gateway" {
  name        = "Cloud-Resume-API"
  description = "This is my API for demonstration purposes"
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_resource" "Cloud-Resume-API-gatewayresource" {
  rest_api_id = aws_api_gateway_rest_api.Cloud-Resume-API-gateway.id
  parent_id   = aws_api_gateway_rest_api.Cloud-Resume-API-gateway.root_resource_id
  path_part   = "cloud-resume"
}

resource "aws_api_gateway_method" "Cloud-Resume-API-gateway" {
  rest_api_id   = aws_api_gateway_rest_api.Cloud-Resume-API-gateway.id
  resource_id   = aws_api_gateway_resource.Cloud-Resume-API-gatewayresource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "Cloud-Resume-API" {
  rest_api_id             = aws_api_gateway_rest_api.Cloud-Resume-API-gateway.id
  resource_id             = aws_api_gateway_resource.Cloud-Resume-API-gatewayresource.id
  http_method             = aws_api_gateway_method.Cloud-Resume-API-gateway.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.cloud-resume-api.invoke_arn

}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.cloud-resume-api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.Cloud-Resume-API-gateway.execution_arn}/*/*"

}

resource "aws_api_gateway_deployment" "cloud-resume-deply" {
  depends_on  = [aws_api_gateway_integration.Cloud-Resume-API]
  rest_api_id = aws_api_gateway_rest_api.Cloud-Resume-API-gateway.id

}

resource "aws_api_gateway_stage" "Cloud-Resume-API-stage" {
  deployment_id = aws_api_gateway_deployment.cloud-resume-deply.id
  rest_api_id   = aws_api_gateway_rest_api.Cloud-Resume-API-gateway.id
  stage_name    = "prod"
}




output "api_gateway_stage_details" {
  value = {
    "stage_name" = "prod",
    "stage_url"  = "${aws_api_gateway_stage.Cloud-Resume-API-stage.invoke_url}/${aws_api_gateway_resource.Cloud-Resume-API-gatewayresource.path_part}"
  }
}
