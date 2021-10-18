---
page_title: "TMC: tmc_aws_data_protection_credential"
layout: "tmc"
subcategory: "AWS Data Protection Credential"
description: |-
  Creates and manages a Tanzu Mission Control (TMC) AWS Data Protection Account Credential.
---

# Resource: tmc_aws_data_protection_credential

The TMC AWS Data Protection Credential resource allows requesting the creation of a AWS data protection credential in Tanzu Mission Control (TMC). It also deals with managing the attributes and lifecycle of the resource.

## Example Usage
# Create a AWS account Data Protection Credential in the Tanzu platform.
```terraform
resource "tmc_aws_data_protection_credential" "example" {
  name = "example"
  iam_role_arn = "arn:aws:iam::1234567890:role/example_arn"
}
```

## Argument Reference

* `name` - (Required) The name of the AWS Data Protection Account Credential. Please note that the credential name must be unique across all Credential Types in Tanzu Mission Control.
* `iam_role_arn` - (Required) IAM Role ARN of the AWS Data Protection Account Credential.

## Attributes Reference

* `id` - Unique Identifier (UID) of the AWS Data Protection Account Credential in the TMC platform.
* `capability` - Capability of the AWS Data Protection Account Credential.
* `credential_provider` - Provider of the AWS Data Protection Account Credential.
* `status` - Status of the AWS Data Protection Account Credential.