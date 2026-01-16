---
page_title: "Foxcon Provider 1.2.x: Upgrade Guide for Provider Schema configuration"
---
# Foxcon Provider 1.2.x: Upgrade Guide for Provider Schema configuration

## Summary

This guide provides detailed instructions for migrating from resource-level schema configuration to provider-level in the foxcon Terraform Provider.

-> **Note:** 
Resource-level configuration is working the same way. This is an optional migration that removes credentials from terraform state file.

## Provider Configuration Migration Instructions

Before reading further, ensure that your current configuration using resource credentials blocks
successfully runs [`terraform plan`](https://www.terraform.io/docs/commands/plan.html)
without unexpected changes. Run the following command:
```bash
terraform plan
```

Your output should resemble:
```bash
foxcon_subject_normalization.test: Refreshing state...
...

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

At this point schema registry credentials are stored in the state file.

The next step is to manually update your Terraform configuration to do the following:
- Add the `schema_registry_rest_endpoint`, `schema_registry_api_key` and `schema_registry_api_secret` attributes to the `provider "foxcon"` block.
- Remove `cloud_api_key` and `cloud_api_secret` attributes, unless those are needed.

### Terraform Configuration Before Migration

```hcl
provider "foxcon" {
  cloud_api_key = "test"
  cloud_api_secret = "test"
}

resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "http://localhost:8081"
  subject_name = "test"
  normalization_enabled = true
  credentials {
    key = "admin"
    secret = "admin-secret"
  }
}
```

### Terraform Configuration After Migration

Resource-level schema registry details can be removed.

```hcl
provider "foxcon" {
  schema_registry_rest_endpoint = "http://localhost:8081"
  schema_registry_api_key = "admin"
  schema_registry_api_secret = "admin-secret"
}

resource "foxcon_subject_normalization" "test" {
  subject_name = "test"
  normalization_enabled = true
}
```

### Verify the Terraform Plan Output

After making the above changes to your Terraform configuration, run the `terraform plan` command again:


```bash
foxcon_subject_normalization.test: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # foxcon_subject_normalization.test will be updated in-place
  ~ resource "foxcon_subject_normalization" "test" {
      ~ last_updated          = "2025-11-06T22:17:26+02:00" -> (known after apply)
      - rest_endpoint         = "http://localhost:8081" -> null
        # (2 unchanged attributes hidden)

      - credentials {
          - key    = "admin" -> null
          - secret = (sensitive value) -> null
        }
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Terraform has compared your real infrastructure against your configuration and found no differences, so no
changes are needed.
```

-> **Note:** The plan output should display `0 to add` **and** `0 to destroy`, confirming that no resources will be recreated or deleted. Only removed schema connection details should be set to `null`.

## Sanity Check

Check that the upgrade was successful by ensuring that your environment
successfully runs [`terraform plan`](https://www.terraform.io/docs/commands/plan.html)
without unexpected changes. Run the following command:
```bash
terraform plan
```
Your output should resemble:
```bash
foxcon_subject_normalization.test: Refreshing state...

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

In the Terraform state file, the `credentials` and `rest_endpoint` blocks will appear as `null`. This indicates that resource/data-source is now successfully authenticated via the provider’s configuration, rather than resource-specific credentials.
```
{
  ...
  "resources": [
    {
      "mode": "managed",
      "type": "foxcon_subject_normalization",
      "name": "test",
      "provider": "provider[\"registry.terraform.io/fox-md/foxcon\"]",
      "instances": [
        {
          "schema_version": 0,
          "attributes": {
            "credentials": null,
            "last_updated": "2025-11-06T22:18:32+02:00",
            "normalization_enabled": false,
            "rest_endpoint": null,
            "subject_name": "test"
          },
          "sensitive_attributes": [],
          "identity_schema_version": 0
        }
      ]
    }
  ]
  ...
}
```

If you run into any problems, [report an issue](https://github.com/fox-md/terraform-provider-foxcon/issues).
