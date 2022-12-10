# KMS_PROVIDERS_PATH is a file containing KMS provider credentials.
# It should normally be kept secret. The sample file only includes a local test provider that is not secret.
export KMS_PROVIDERS_PATH="./sample-kms-providers.json"
go run -tags cse .
