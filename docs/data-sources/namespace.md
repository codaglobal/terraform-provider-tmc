---
page_title: "TMC: tmc_namespace"
layout: "tmc"
subcategory: "Tanzu Namespace"
description: |-
  Get information on a specific namespace of a cluster in Tanzu Mission Control (TMC)
---

# Data Source: tmc_namespace

The TMC Namespace data resource can be used to get the information of a namespace for a cluster in Tanzu Mission Control (TMC). 

```terraform
data "tmc_namespace" "example" {
  cluster_name       = "example-cluster"
  management_cluster = "example-hosted"
  provisioner_name   = "example-provisioner"
  name               = "example-ns"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the namespace. Changing the name forces recreation of this resource.
* `cluster_name` - (Required) The name of the Tanzu Cluster for which the namespace is to be created.
* `management_cluster` - (Required) Name of the management cluster used to provision the cluster.
* `provisioner_name` - (Required) Name of the provisioner to be used.

## Attributes Reference

* `id` - The UID of the Tanzu Cluster.
* `description` - The description of the nodepool.
* `workspace_name` - Name of the workspace for the created namespace.
* `labels` - A map of labels to assign to the resource.
