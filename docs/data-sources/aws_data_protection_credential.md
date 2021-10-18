---
page_title: "TMC: tmc_aws_data_protection_credential"
layout: "tmc"
subcategory: "AWS Data Protection Credential"
description: |-
  Get information on a Tanzu Mission Control (TMC) AWS Data Protection Account Credential.
---

# Data Source: tmc_aws_data_protection_credential

Use this data source to get the details about a AWS Data Protection Credential in TMC platform.

## Example Usage
# Get details of a aws data protection credential account in the Tanzu platform.
```terraform
data "tmc_aws_data_protection_credential" "example" {
  name = "example"
}
```

## Argument Reference

* `name` - (Required) The name of the AWS Data Protection Account Credential to lookup in the TMC platform. If no credential is found with this name, an error will be returned.

## Attributes Reference

* `id` - Unique Identifiers (UID) of the found AWS Data Protection Account Credential in the TMC platform.
* `iam_role_arn` - AWS IAM Role Arn of the found AWS Data Protection Account Credential
* `capability` - Capability of the found AWS Data Protection Account Credential
* `credential_provider` - Provider of the found AWS Data Protection Account Credential
* `status` - Status of the found AWS Data Protection Account Credential