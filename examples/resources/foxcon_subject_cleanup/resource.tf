resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint  = "http://localhost:8081"
  subject_name   = "versioned"
  cleanup_method = "KEEP_LATEST_ONLY"
  credentials {
    key    = "admin"
    secret = "admin-secret"
  }
}

resource "foxcon_subject_cleanup" "active" {
  rest_endpoint  = "http://localhost:8081"
  subject_name   = "versioned"
  cleanup_method = "KEEP_ACTIVE_ONLY"
  credentials {
    key    = "admin"
    secret = "admin-secret"
  }
}
