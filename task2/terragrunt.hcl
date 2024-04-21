# terragrunt.hcl

terraform {
  source = "./aws_secretsmanager"
}

inputs = {
  secret1_name = get_env("TF_VAR_secret1_name", "default-secret1-name")
  secret1_value = get_env("TF_VAR_secret1_value", "")
  secret2_name = get_env("TF_VAR_secret2_name", "default-secret2-name")
  secret2_value = get_env("TF_VAR_secret2_value", "")
}
