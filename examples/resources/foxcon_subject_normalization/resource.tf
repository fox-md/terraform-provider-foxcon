resource "foxcon_subject_normalization" "test" {
  rest_endpoint         = "http://localhost:8081"
  subject_name          = "test"
  normalization_enabled = true
  credentials {
    key    = "admin"
    secret = "admin-secret"
  }
}