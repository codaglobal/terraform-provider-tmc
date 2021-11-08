---
page_title: "TMC: tmc_nodepool"
layout: "tmc"
subcategory: "TKG Cluster"
description: |-
  Creates and manages a nodepool for a AWS cluster in the TMC platform
---

# Resource: tmc_cluster

The TMC Cluster resource allows requesting the creation of a nodepool for a AWS cluster in Tanzu Mission Control (TMC). 

```terraform
resource "tmc_aws_nodepool" "example" {
  name               = "default-node-pool"
  cluster_name       = "example-cluster"
  management_cluster = "example-aws-hosted"
  provisioner_name   = "example-aws-provisioner"
  worker_node_count  = 1
  availability_zone = "us-east-1a"
  instance_type     = "m5.large"
  version           = "1.20.8-1-amazon2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Nodepool. Changing the name forces recreation of this resource.
* `description` - (Optional) The description of the nodepool.
* `node_labels` - (Optional) A map of node labels to assign to the resource.
* `cloud_labels` - (Optional) A map of cloud labels to assign to the resource.
* `cluster_name` - (Required) The name of the Tanzu Cluster for which the nodepool is to be created.
* `management_cluster` - (Required) Name of the management cluster used to provision the cluster.
* `provisioner_name` - (Required) Name of the provisioner to be used.
* `availability_zone` - (Required) The AWS availability zone for the cluster's worker nodes.
* `instance_type` - (Required) Instance type of the EC2 nodes to be used as part of the nodepool.
* `version` - (Required) Version of Kubernetes to be used in the cluster.


## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the Tanzu Cluster Group.