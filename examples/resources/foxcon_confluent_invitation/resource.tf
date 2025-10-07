resource "foxcon_confluent_invitation" "this" {
  email     = "user@company.com"
  auth_type = "AUTH_TYPE_SSO"
}