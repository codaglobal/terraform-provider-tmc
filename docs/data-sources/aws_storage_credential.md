---
page_title: "TMC: tmc_aws_storage_credential"
layout: "tmc"
subcategory: "AWS Storage Credential"
description: |-
  Get information on a Tanzu Mission Control (TMC) AWS Storage Account Credential.
---

# Data Source: tmc_aws_storage_credential

Use this data source to get the details about a AWS Storage Credential in TMC platform.

## Example Usage
# Get details of a aws storage credential account in the Tanzu platform.
```terraform
data "tmc_aws_storage_credential" "example" {
  name = "example"
}
```

## Argument Reference

* `name` - (Required) The name of the AWS Storage Account Credential to lookup in the TMC platform. If no credential is found with this name, an error will be returned.

## Attributes Reference

* `id` - Unique Identifiers (UID) of the found AWS Storage Account Credential in the TMC platform.
* `capability` - Capability of the found AWS Storage Account Credential
* `credential_provider` - Provider of the found AWS Storage Account Credential
* `status` - Status of the found AWS Storage Account Credential