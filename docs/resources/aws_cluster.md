---
page_title: "TMC: tmc_aws_cluster"
layout: "tmc"
subcategory: "TKG Cluster"
description: |-
  Creates and manages a Cluster using AWS in the TMC platform
---

# Resource: tmc_cluster

The TMC Cluster resource allows requesting the creation of a AWS cluster in Tanzu Mission Control (TMC). It also deals with managing the attributes and lifecycle of the cluster.

!> **Note**: The AWS cluster will be in ready state after a successful apply. But the health of the cluster will be **unknown** until a [nodepool](aws_nodepool.md) resource is created for this cluster

```terraform
resource "tmc_aws_cluster" "example" {
  name               = "example-cluster"
  cluster_group      = "example"
  management_cluster = "example-aws-hosted"
  provisioner_name   = "example-aws-provisioner"
  labels = {
    env = "example"
    createdby = "terraform"
  }
  version            = "1.20.8-1-amazon2"
  credential_name    = "example-aws-cred"
  ssh_key            = "example-key"
  region             = "us-east-1"
  control_plane_spec {
    instance_type      = "m5.large"
    availability_zones = ["us-east-1a"]
    vpc_cidrblock      = "10.0.0.0/16"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) (Forces Replacement) The name of the Tanzu Cluster. Changing the name forces recreation of this resource.
* `description` - (Optional) (Forces Replacement) The description of the Tanzu Cluster.
* `labels` - (Optional) A map of labels to assign to the resource.
* `cluster_group` - (Required) A map of labels to assign to the resource.
* `management_cluster` - (Required) (Forces Replacement) Name of the management cluster used to provision the cluster.
* `provisioner_name` - (Required) (Forces Replacement) Name of the provisioner to be used.
* `version` - (Required) (Forces Replacement) Version of Kubernetes to be used in the cluster.
* `credential_name` - (Required) (Forces Replacement) Name of the AWS Credentials, that is already available in Tanzu Mission Control (TMC). This will be used to provision the required resources in AWS.
* `ssh_key` - (Required) (Forces Replacement) Name of the SSH key pair to be used for the EC2 instances. This key pair can then be used to access the EC2 instances.
* `region` - (Required) (Forces Replacement) AWS region of the cluster.
* `pod_cidrblock` - (Optional) (Forces Replacement) Pod CIDR for Kubernetes pods. Defaults to 192.168.0.0/16.
* `service_cidrblock` - (Optional) (Forces Replacement) Service CIDR for Kubernetes services. Defaults to 10.96.0.0/12.
* [`control_plane_spec`](#control_plane_spec) - (Required) (Forces Replacement) Contains information related to the Control Plane of the cluster

## Nested Blocks

#### `control_plane_spec`

#### Arguments

* `availability_zones` - (Required) A list of availability zones for the cluster's control plane

!> **Note**: The availability zones specified for the cluster must have atleast one private and one public
subnet

* `instance_type` - (Required) Instance type of the EC2 nodes to be used as part of the control plane.
* `vpc_cidrblock` - (Optional) CIDR block of the AWS VPC to be used for the control plane.

!> **Note**: Only one of `vpc_cidrblock` or `vpc_id` can be specified

* `vpc_id` - (Optional) VPC Id in which the cluster is created

!> **Note**: Only one of `vpc_cidrblock` or `vpc_id` can be specified

* `public_subnets` - (Optional) List of public subnets in the VPC. One and only one subnet must be specified for each AZs given in the `availability_zones` field. If `vpc_id` is specified, then this is a **required** field
* `private_subnets` - (Optional) List of private subnets in the VPC. One and only one subnet must be specified for each AZs given in the `availability_zones` field. If `vpc_id` is specified, then this is a **required** field

## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the Tanzu Cluster.
* `resource_version` - An identifier used to track changes to the resource
