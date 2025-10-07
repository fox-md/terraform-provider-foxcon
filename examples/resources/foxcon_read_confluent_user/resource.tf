resource "foxcon_read_confluent_user" "iid" {
  invitation_id = "i-321cba"
}

resource "foxcon_read_confluent_user" "email" {
  user_email = "user@company.com"
}

resource "foxcon_read_confluent_user" "uid" {
  user_id = "u-abc123"
}