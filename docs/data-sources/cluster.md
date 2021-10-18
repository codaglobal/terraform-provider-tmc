---
page_title: "TMC: tmc_cluster"
layout: "tmc"
subcategory: "Cluster"
description: |-
  Get information on a specific cluster in Tanzu Mission Control (TMC)
---

# Data Source: tmc_cluster

The TMC Cluster data resource can be used to get the information of a cluster in Tanzu Mission Control (TMC). 

```terraform
resource "tmc_cluster" "example" {
  name               = "example-cluster"
  cluster_group      = "example"
  management_cluster = "example-aws-hosted"
  provisioner_name   = "example-aws-provisioner"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Tanzu Cluster. Changing the name forces recreation of this resource.
* `cluster_group` - (Required) A map of labels to assign to the resource.
* `management_cluster` - (Required) Name of the management cluster used to provision the cluster.
* `provisioner_name` - (Required) Name of the provisioner to be used.

## Attributes Reference

* `id` - The UID of the Tanzu Cluster.
* `description` - The description of the Tanzu Cluster.
* `labels` - A map of labels to assign to the resource.
* [`tkg-aws`](#tkg-aws) - Contains information for provisioning a cluster using AWS.

## Nested Blocks

#### `tkg-aws`

#### Attributes

* `availability_zones` - A list of availability zones for the cluster's control plane
* `instance_type` - Instance type of the EC2 nodes to be used as part of the control plane.
* `vpc_cidrblock` - CIDR block of the AWS VPC to be used for the control plane.
* `version` - Version of Kubernetes to be used in the cluster.
* `credential_name` - Name of the AWS Credentials, that is already available in Tanzu Mission Control (TMC). This will be used to provision thresources in AWS.
* `ssh_key` - Name of the SSH key pair to be used for the EC2 instances. This key pair can then be used to access the EC2 instances.
* `region` - AWS region of the cluster.
* `pod_cidrblock` - Pod CIDR for Kubernetes pods. Defaults to 192.168.0.0/16.
* `service_cidrblock` - Service CIDR for Kubernetes services. Defaults to 10.96.0.0/12.
