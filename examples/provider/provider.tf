# Confluent cloud connection
provider "foxcon" {
  cloud_api_key    = "test" # optionally use CONFLUENT_CLOUD_API_KEY env var
  cloud_api_secret = "test" # optionally use CONFLUENT_CLOUD_API_SECRET env var
}

# Assumes values are configured using environment variables or schema registry credentials are set individually on a resource level
provider "foxcon" {}

# Manages a single Schema Registry cluster in the same Terraform workspace
provider "foxcon" {
  schema_registry_rest_endpoint = "https://psrc-abcde.uksouth.azure.confluent.cloud" # optionally use SCHEMA_REGISTRY_REST_ENDPOINT env var
  schema_registry_api_key = "test"                                                   # optionally use SCHEMA_REGISTRY_API_KEY env var
  schema_registry_api_secret = "test"                                                # optionally use SCHEMA_REGISTRY_API_SECRET env var
}

# Manages a single Schema Registry cluster in the same Terraform workspace and Confluent cloud connection
provider "foxcon" {
  cloud_api_key    = "test"                                                          # optionally use CONFLUENT_CLOUD_API_KEY env var
  cloud_api_secret = "test"                                                          # optionally use CONFLUENT_CLOUD_API_SECRET env var
  schema_registry_rest_endpoint = "https://psrc-abcde.uksouth.azure.confluent.cloud" # optionally use SCHEMA_REGISTRY_REST_ENDPOINT env var
  schema_registry_api_key = "test"                                                   # optionally use SCHEMA_REGISTRY_API_KEY env var
  schema_registry_api_secret = "test"                                                # optionally use SCHEMA_REGISTRY_API_SECRET env var
}
