---
page_title: "TMC: tmc_vsphere_cluster"
layout: "tmc"
subcategory: "TKG Cluster"
description: |-
  Creates and manages a Vsphere Cluster in the TMC platform
---

# Resource: tmc_cluster

The TMC Cluster resource allows requesting the creation of a Vsphere cluster in Tanzu Mission Control (TMC). It also deals with managing the attributes and lifecycle of the cluster.


```terraform
resource "tmc_vsphere_cluster" "example" {
  name               = "example-vsphere-cluster"
  management_cluster = "example-vpshere-mgmt-cluster"
  provisioner_name   = "example-provisioner"

  version       = "v1.21.6+vmware.1-tkg.1.b3d708a"
  cluster_group = "default"

  control_plane_spec {
    class         = "best-effort-xsmall"
    storage_class = "vsphere-tanzu-example-storage-policy"
  }

  nodepool {
    nodepool_name = "example-nodepool"
    worker_node_count = 1
    node_class = "best-effort-small"
    node_storage_class = "vsphere-tanzu-example-storage-policy"
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
* `region` - (Required) (Forces Replacement) AWS region of the cluster.
* `pod_cidrblock` - (Optional) (Forces Replacement) Pod CIDR for Kubernetes pods. Defaults to 192.168.0.0/16.
* `service_cidrblock` - (Optional) (Forces Replacement) Service CIDR for Kubernetes services. Defaults to 10.96.0.0/12.
* [`control_plane_spec`](#control_plane_spec) - (Required) (Forces Replacement) Contains information related to the Control Plane of the cluster
* [`nodepool`](#nodepool) - (Required) (Forces Replacement) Contains information related to the Nodepool of the cluster

## Nested Blocks

#### `control_plane_spec`

#### Arguments

* `class` - (Required) Indicates the size of the VMs to be provisioned.
* `storage_class` - (Required) Storage Class to be used for storage of the disks which store the root filesystems of the nodes

#### `nodepool`

#### Arguments

* `nodepool_name` - (Required) Determines the name of the nodepool
* `worker_node_count` - (Required) Determines the number of worker nodes provisioned
* `node_class` - (Required) Determines the class of the worker node
* `node_storage_class` - (Required) Determines the storage policy used for the worker node


## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the Tanzu Cluster.
* `resource_version` - An identifier used to track changes to the resource
