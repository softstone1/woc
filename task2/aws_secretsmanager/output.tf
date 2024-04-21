output "secret1_arn" {
  value = aws_secretsmanager_secret.secret1.arn
}

output "secret2_arn" {
  value = aws_secretsmanager_secret.secret2.arn
}