---
page_title: "TMC: tmc_cluster_admin_kubeconfig"
layout: "tmc"
subcategory: "TKG Cluster"
description: |-
  Get a base64 encoded string of the kubeconfig for a given cluster.
---

# Data Source: tmc_cluster_admin_kubeconfig

Use this data source to get the kubeconfig for a cluster in TMC platform.

## Example Usage
# Get the kubeconfig for a cluster in the Tanzu platform.
```terraform
data "tmc_cluster_admin_kubeconfig" "example" {
  cluster_name = "example-cluster"
  management_cluster = "example-hosted-cluster"
  provisioner_name = "example-provisioner"
}
```

## Argument Reference

* `cluster_name` - (Required) The name of the Tanzu Cluster for which the kubeconfig is to be fetched.
* `management_cluster` - (Required) Name of the management cluster used to provision the cluster.
* `provisioner_name` - (Required) Name of the provisioner used.

## Attributes Reference

* `kubeconfig` - A base64 encoded string of the kubeconfig necessary to connect to the cluster.