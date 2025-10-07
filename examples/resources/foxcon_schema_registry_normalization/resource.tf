resource "foxcon_schema_registry_normalization" "test" {
  rest_endpoint         = "http://localhost:8081"
  normalization_enabled = true
  credentials {
    key    = "admin"
    secret = "admin-secret"
  }
}