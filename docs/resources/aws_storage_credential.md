---
page_title: "TMC: tmc_aws_storage_credential"
layout: "tmc"
subcategory: "AWS Storage Credential"
description: |-
  Creates and manages a Tanzu Mission Control (TMC) AWS Storage Account Credential.
---

# Resource: tmc_aws_storage_credential

The TMC AWS Storage Credential resource allows requesting the creation of a AWS Storage credential in Tanzu Mission Control (TMC). It also deals with managing the attributes and lifecycle of the resource.

## Example Usage
# Create a AWS account Storage Credential in the Tanzu platform.
```terraform
resource "tmc_aws_storage_credential" "example" {
  name = "example"
  access_key_id = "xxxx"
  secret_access_key = "yyyy"
}
```

## Argument Reference

* `name` - (Required) The name of the AWS Storage Account Credential. Please note that the credential name must be unique across all Credential Types in Tanzu Mission Control.
* `access_key_id` - (Sensitive) AWS Access Key ID of the AWS Storage Account Credential.
* `secret_access_key` - (Sensitive) AWS Access Key ID of the AWS Storage Account Credential.

## Attributes Reference

* `id` - Unique Identifier (UID) of the AWS Storage Account Credential in the TMC platform.
* `capability` - Capability of the AWS Storage Account Credential.
* `credential_provider` - Provider of the AWS Storage Account Credential.
* `status` - Status of the AWS Storage Account Credential.