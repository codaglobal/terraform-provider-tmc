---
page_title: "TMC: tmc_observability_credential"
layout: "tmc"
subcategory: "Tanzu Observability Credential"
description: |-
  Get information on a Tanzu Mission Control (TMC) Observability Credential.
---

# Data Source: tmc_observability_credential

Use this data source to get the details about an Observability Credential in TMC platform.

## Example Usage
# Get details of an observability credential account in the Tanzu platform.
```terraform
data "tmc_observability_credential" "example" {
  name = "example"
}
```

## Argument Reference

* `name` - (Required) The name of the Observability Credential to lookup in the TMC platform. If no credential is found with this name, an error will be returned.

## Attributes Reference

* `id` - Unique Identifiers (UID) of the found Observability Credential in the TMC platform.
* `observability_url` - URL of the found Observability Credential
* `capability` - Capability of the found Observability Credential
* `status` - Status of the found Tanzu Observability Instance

