data "foxcon_subject_versions" "test" {
  rest_endpoint = "http://localhost:8081"
  subject_name  = "data-source"
  credentials {
    key    = "admin"
    secret = "admin-secret"
  }
}