---
page_title: "TMC: tmc_observability_credential"
layout: "tmc"
subcategory: "Tanzu Observability Credential"
description: |-
  Creates and manages a Tanzu Mission Control (TMC) Observability Account Credential.
---

# Resource: tmc_observability_credential

The TMC Observability Credential resource allows requesting the creation of a Observability credential in Tanzu Mission Control (TMC). It also deals with managing the attributes and lifecycle of the resource.

## Example Usage
# Create a Tanzu Observability Credential in the Tanzu platform.
```terraform
resource "tmc_observability_credential" "example" {
  name = "example"
  api_token = "mock-token"
  observability_url = "mock-url"
}
```

## Argument Reference

* `name` - (Required) The name of the Tanzu Observability Credential. Please note that the credential name must be unique across all Credential Types in Tanzu Mission Control.
* `api_token` - (Sensitive) API Token of the Tanzu Observability Credential.
* `observability_url` - (Required) URL of the Tanzu Observability Instance.

## Attributes Reference

* `id` - Unique Identifier (UID) of the Tanzu Observability Credential in the TMC platform.
* `capability` - Capability of the Tanzu Observability Credential.
* `status` - Status of the Tanzu Observability Credential.