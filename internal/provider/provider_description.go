// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

var providerDescription string = "`foxcon` provider extends Confluent official [confluentinc/confluent](https://registry.terraform.io/providers/confluentinc/confluent/latest/docs) provider.\n\n" +
	"`foxcon` includes below:" +
	`
- Normalization configuration for subject.
- Normalization configuration for schema registry.
- Confluent invitation resource that acts as original, however also deletes user from Confluent on resource deletion.
- Cleanup of schema versions. Can be performed for soft-deleted or all non-latest versions.
` + "- `foxcon_confluent_read_user` that reads user details from Confluent on resources creation and deletes user from Confluent on resource deletion.\n" +
	"- `foxcon_set_subject_mode` action that sets subject mode adhoc."
