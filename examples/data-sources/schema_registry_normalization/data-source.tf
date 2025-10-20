data "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "http://localhost:8081"
  credentials {
    key    = "admin"
    secret = "admin-secret"
  }
}