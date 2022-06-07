---
page_title: "TMC: tmc_aws_nodepool"
layout: "tmc"
subcategory: "TKG Cluster"
description: |-
  Get information on a specific nodepool of a AWS cluster in Tanzu Mission Control (TMC)
---

# Data Source: tmc_aws_nodepool

The TMC Nodepool data resource can be used to get the information of a nodepool for a AWS cluster in Tanzu Mission Control (TMC). 

```terraform
data "tmc_aws_nodepool" "example" {
  name               = "example-node-pool"
  cluster_name       = "example-cluster"
  management_cluster = "example-aws-hosted"
  provisioner_name   = "example-aws-provisioner"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Nodepool. Changing the name forces recreation of this resource.
* `cluster_name` - (Required) The name of the Tanzu Cluster for which the nodepool is to be created.
* `management_cluster` - (Required) Name of the management cluster used to provision the cluster.
* `provisioner_name` - (Required) Name of the provisioner to be used.

## Attributes Reference

* `id` - The UID of the Tanzu Cluster.
* `description` - The description of the nodepool.
* `node_labels` - A map of node labels to assign to the resource.
* `cloud_labels` - A map of cloud labels to assign to the resource.
* `availability_zone` - The AWS availability zone for the cluster's worker nodes.
* `instance_type` - Instance type of the EC2 nodes to be used as part of the nodepool.
* `version` - Version of Kubernetes to be used in the cluster.
