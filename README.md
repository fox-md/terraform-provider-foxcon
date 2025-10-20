## foxcon Provider

`foxcon` provider extends Confluent official [confluentinc/confluent](https://registry.terraform.io/providers/confluentinc/confluent/latest/docs) provider.

`foxcon` includes below resources:
- Normalization configuration for subject.
- Normalization configuration for schema registry.
- Confluent invitation resource that acts as original, however also deletes user from Confluent on resource deletion.
- `foxcon_confluent_read_user` that reads user details from Confluent on resources creation and deletes user from Confluent on resource deletion.
- Cleanup of schema versions. Can be performed for soft-deleted or all non-latest versions.

## Badges

[![License](https://img.shields.io/github/license/fox-md/terraform-provider-foxcon)](/LICENSE)
[![Release](https://img.shields.io/github/v/release/fox-md/terraform-provider-foxcon.svg)](https://github.com/fox-md/terraform-provider-foxcon/releases/latest)
