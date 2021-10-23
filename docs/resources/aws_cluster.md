---
page_title: "TMC: tmc_aws_cluster"
layout: "tmc"
subcategory: "TKG Cluster"
description: |-
  Creates and manages a Cluster using AWS in the TMC platform
---

# Resource: tmc_cluster

The TMC Cluster resource allows requesting the creation of a AWS cluster in Tanzu Mission Control (TMC). It also deals with managing the attributes and lifecycle of the cluster.

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
  availability_zones = ["us-east-1a"]
  instance_type      = "m5.large"
  vpc_cidrblock      = "10.0.0.0/16"
  version            = "1.20.8-1-amazon2"
  credential_name    = "example-aws-cred"
  ssh_key            = "example-key"
  region             = "us-east-1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Tanzu Cluster. Changing the name forces recreation of this resource.
* `description` - (Optional) The description of the Tanzu Cluster.
* `labels` - (Optional) A map of labels to assign to the resource.
* `cluster_group` - (Required) A map of labels to assign to the resource.
* `management_cluster` - (Required) Name of the management cluster used to provision the cluster.
* `provisioner_name` - (Required) Name of the provisioner to be used.
* `availability_zones` - (Required) A list of availability zones for the cluster's control plane

!> **Note**: The availability zones specified for the cluster must have atleast one private and one public
subnet

* `instance_type` - (Required) Instance type of the EC2 nodes to be used as part of the control plane.
* `vpc_cidrblock` - (Required) CIDR block of the AWS VPC to be used for the control plane.
* `version` - (Required) Version of Kubernetes to be used in the cluster.
* `credential_name` - (Required) Name of the AWS Credentials, that is already available in Tanzu Mission Control (TMC). This will be used to provision the required resources in AWS.
* `ssh_key` - (Required) Name of the SSH key pair to be used for the EC2 instances. This key pair can then be used to access the EC2 instances.
* `region` - (Required) AWS region of the cluster.
* `pod_cidrblock` - (Optional) - Pod CIDR for Kubernetes pods. Defaults to 192.168.0.0/16.
* `service_cidrblock` - (Optional) - Service CIDR for Kubernetes services. Defaults to 10.96.0.0/12.


## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the Tanzu Cluster.