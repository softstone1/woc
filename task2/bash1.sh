#!/bin/bash

# Define the array of secret names
secrets=("testing/secret1" "testing/secret2")

# Loop through the secret names
for secret in "${secrets[@]}"; do
  echo "Fetching secret for: $secret"
  value=$(aws secretsmanager get-secret-value --secret-id "$secret" --query 'SecretString' --output text)
  echo "Value: $value"
done
