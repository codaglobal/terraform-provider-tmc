---
page_title: "TMC: tmc_management_cluster"
layout: "tmc"
subcategory: "Tanzu Management Clusters"
description: |-
  Creates and Registers a Management Cluster in the TMC platform
---

# Resource: tmc_management_cluster

The TMC Management Cluster resource allows requesting the creation of a management cluster in Tanzu Mission Control (TMC). 

!> **Note**: This resource does not support `update` operation and hence will be destroyed and recreated for every change.

```terraform
resource "tmc_management_cluster" "example" {
  name                     = "tf-mgmt-cluster"
  description              = "terraform created mgmt cluster"
  kubernetes_provider_type = "tkg"
  default_cluster_group    = "default"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Tanzu Management Cluster. Changing the name forces recreation of this resource.
* `description` - (Optional) The description of the Tanzu Management Cluster.
* `labels` - (Optional) A map of labels to assign to the resource.
* `kubernetes_provider_type` - (Required) Type of cluster to be registered into TMC. Can be one of `tkg`, `tkgservice`, `tkghosted` or `other`
* `default_cluster_group` - (Required) Default cluster group for the workload clusters.

## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the Tanzu Management Cluster.
* `registration_url` - An URL to fetch the Tanzu Agent installation YAML which is necessary to establish connection to the registered cluster.
