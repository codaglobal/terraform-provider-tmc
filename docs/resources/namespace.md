---
page_title: "TMC: tmc_namespace"
layout: "tmc"
subcategory: "Tanzu Namespace"
description: |-
  Creates and manages a namespace for a cluster in the TMC platform
---

# Resource: tmc_namespace

The TMC Namespace resource allows requesting the creation of a namespace for a cluster in Tanzu Mission Control (TMC). 

```terraform
resource "tmc_namespace" "example" {
  name               = "example-ns"
  description        = "terraform created mgmt cluster"
  cluster_name       = "example-cluster"
  management_cluster = "example-hosted"
  provisioner_name   = "example-provisioner"
  workspace_name     = "default"
  labels = {
    "CreatedBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Namespace. Changing the name forces recreation of this resource.
* `description` - (Optional) The description of the namespace.
* `cluster_name` - (Required) The name of the Tanzu Cluster for which the namespace is to be created.
* `management_cluster` - (Required) Name of the management cluster used to provision the cluster.
* `provisioner_name` - (Required) Name of the provisioner to be used.
* `workspace_name` - (Required) Name of the workspace for the created namespace.
* `labels` - (Optional) A map of labels to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the Tanzu Namespace.