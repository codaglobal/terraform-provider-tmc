---
page_title: "TMC: tmc_aws_cluster"
layout: "tmc"
subcategory: "TKG Cluster"
description: |-
  Get information on a specific cluster created using AWS in Tanzu Mission Control (TMC)
---

# Data Source: tmc_cluster

The TMC Cluster data resource can be used to get the information of any AWS cluster in Tanzu Mission Control (TMC). 

```terraform
resource "tmc_aws_cluster" "example" {
  name               = "example-cluster"
  management_cluster = "example-aws-hosted"
  provisioner_name   = "example-aws-provisioner"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Tanzu Cluster. Changing the name forces recreation of this resource.
* `management_cluster` - (Required) Name of the management cluster used to provision the cluster.
* `provisioner_name` - (Required) Name of the provisioner to be used.

## Attributes Reference

* `id` - The UID of the Tanzu Cluster.
* `description` - The description of the Tanzu Cluster.
* `labels` - A map of labels to assign to the resource.
* `cluster_group` - A map of labels to assign to the resource.
* `version` - Version of Kubernetes to be used in the cluster.
* `credential_name` - Name of the AWS Credentials, that is already available in Tanzu Mission Control (TMC). This will be used to provision the resources in AWS.
* `ssh_key` - Name of the SSH key pair to be used for the EC2 instances. This key pair can then be used to access the EC2 instances.
* `region` - AWS region of the cluster.
* `pod_cidrblock` - Pod CIDR for Kubernetes pods. Defaults to 192.168.0.0/16.
* `service_cidrblock` - Service CIDR for Kubernetes services. Defaults to 10.96.0.0/12.
* [`control_plane_spec`](#control_plane_spec) - Contains information related to the Control Plane of the cluster


## Nested Blocks

#### `control_plane_spec`

#### Attributes

* `availability_zones` - A list of availability zones for the cluster's control plane
* `instance_type` - Instance type of the EC2 nodes to be used as part of the control plane.
* `vpc_cidrblock` - CIDR block of the AWS VPC to be used for the control plane.
* `vpc_id` - VPC Id in which the cluster is created
* `public_subnets` - List of public subnets in the VPC. One and only one subnet must be specified for each AZs given in the `availability_zones` field. 
* `private_subnets` - List of private subnets in the VPC. One and only one subnet must be specified for each AZs given in the `availability_zones` field.
