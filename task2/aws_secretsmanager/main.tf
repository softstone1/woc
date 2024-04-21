resource "aws_secretsmanager_secret" "secret1" {
  name = var.secret1_name
}

resource "aws_secretsmanager_secret_version" "secret1_version" {
  secret_id     = aws_secretsmanager_secret.secret1.id
  secret_string = var.secret1_value
}

resource "aws_secretsmanager_secret" "secret2" {
  name = var.secret2_name
}

resource "aws_secretsmanager_secret_version" "secret2_version" {
  secret_id     = aws_secretsmanager_secret.secret2.id
  secret_string = var.secret2_value
}
