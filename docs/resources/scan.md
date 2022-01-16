---
page_title: "TMC: tmc_cluster_scan"
layout: "tmc"
subcategory: "Tanzu Cluster Inspection"
description: |-
  Creates a Cluster Scan in the TMC platform
---

# Resource: tmc_cluster_scan

The TMC Cluster Scan resource allows to run a scan for a given cluster in Tanzu Mission Control (TMC).There are three types of scans that are currently supported `lite`, `cis` and `conformance`

```terraform
resource "tmc_cluster_scan" "example" {
  cluster_name       = "example-cluster"
  management_cluster = "example-aws-hosted"
  provisioner_name   = "example-aws-provisioner"
  type               = "lite"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required) (Forces Replacement) The name of the Tanzu Cluster for which the nodepool is to be created.
* `management_cluster` - (Required) (Forces Replacement) Name of the management cluster used to provision the cluster.
* `provisioner_name` - (Required) (Forces Replacement) Name of the provisioner to be used.
* `type` - (Required) (Forces Replacement) Type of the scan to be performed.

## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the cluster scan.
* `name` - The name of the cluster scan performed.