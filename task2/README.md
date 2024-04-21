# Task 2

This terragrunt.hcl configuration will load the secret values from the environment variables. If the environment variable is not set, it will use the default value specified as the second argument to get_env function. For secrets, we let it fail if the variable is not set, so we set the default to an empty string.

**Deployment with Terragrunt**

Before running Terragrunt, we'll set your environment variables:
```bash
export TF_VAR_secret1_name="testing/secret1"
export TF_VAR_secret1_value='{"name":"Example1Username","password":"Example1Pass"}'
export TF_VAR_secret2_name="testing/secret2"
export TF_VAR_secret2_value='{"key":"somestring"}'
```
Then, navigate to the directory where your terragrunt.hcl file is located and run:

```bash
terragrunt apply
```

**Testing**

Lists the names of all the secrets in the AWS Secrets Manager:
```bash
aws secretsmanager list-secrets | jq '.SecretList[] | .Name'
```
Gets the value of a specific secret using its ARN:
```bash
ARN1=$(aws secretsmanager list-secrets | jq '.SecretList[] | .ARN' | head -n1 | sed -e 's/"//g')
aws secretsmanager get-secret-value --secret-id $ARN1 | jq '.SecretString'
```
Fetch the ARN for testing/secret1:
```bash
ARN1=$(aws secretsmanager list-secrets | jq -r '.SecretList[] | select(.Name=="testing/secret1") | .ARN')
aws secretsmanager get-secret-value --secret-id $ARN1 | jq '.SecretString'
```

***bash1.sh***

Give it execute permissions with chmod +x bash1.sh, and then run it with ./bash1.sh