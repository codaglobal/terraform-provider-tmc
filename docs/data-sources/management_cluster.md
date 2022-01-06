---
page_title: "TMC: tmc_management_cluster"
layout: "tmc"
subcategory: "Tanzu Management Clusters"
description: |-
  Get information on a Management Cluster in the TMC platform
---

# Data Source: tmc_management_cluster

Use this data source to get the details about a management cluster in TMC platform.

## Example Usage
# Get details of a cluster group in the Tanzu platform.
```terraform
data "tmc_management_cluster" "example" {
  name = "example-cluster"
}
```

## Argument Reference

* `name` - (Required) The name of the management cluster to lookup in the TMC platform.

## Attributes Reference

* `id` - Unique Identifiers (UID) of the found management cluster group in the TMC platform.
* `description` - Description of the management cluster.
* `labels` - A mapping of labels of the resource.
* `kubernetes_provider_type` - Type of cluster to be registered into TMC. Can be one of `tkg`, `tkgservice`, `tkghosted` or `other`
* `default_cluster_group` - Default cluster group for the workload clusters.
* `registration_url` - An URL to fetch the Tanzu Agent installation YAML which is necessary to establish connection to the registered cluster (if available)
