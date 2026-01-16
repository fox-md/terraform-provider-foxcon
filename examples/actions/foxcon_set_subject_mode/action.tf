variable "subject_name" {
  default = "test"
}

locals {
  json_schema = <<-EOT
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "http://example.com/myURI.schema.json",
  "title": "SampleRecord",
  "description": "Sample schema to help",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "myField1": {
      "type": "integer",
      "description": "The integer type is used for integral numbers."
    }
  }
}
EOT
}

resource "confluent_subject_mode" "ro" {
  subject_name = var.subject_name
  mode         = "READONLY"
}

resource "confluent_schema" "this" {

  depends_on = [confluent_subject_mode.ro]

  subject_name = var.subject_name

  format = "JSON"
  schema = jsonencode(local.json_schema)

  lifecycle {
    action_trigger {
      events  = [after_create, after_update]
      actions = [action.foxcon_set_subject_mode.ro]
    }
    action_trigger {
      events  = [before_create, before_update]
      actions = [action.foxcon_set_subject_mode.rw]
    }
  }

}

action "foxcon_set_subject_mode" "ro" {
  config {
    subject_name = var.subject_name
    mode         = "READONLY"
  }
}

action "foxcon_set_subject_mode" "rw" {
  config {
    subject_name = var.subject_name
    mode         = "READWRITE"
  }
}
