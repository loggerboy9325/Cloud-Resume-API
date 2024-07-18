resource "aws_dynamodb_table" "resume-api" {
  name           = "Resume-api"
  billing_mode   = "PROVISIONED"
  hash_key       = "id"
  read_capacity  = 5
  write_capacity = 5

  attribute {
    name = "id"
    type = "S"
  }
}
