---
page_title: "TMC: tmc_cluster_backup"
layout: "tmc"
subcategory: "Cluster Backups"
description: |-
  Creates and manages a Cluster Backup in the TMC platform
---

# Resource: tmc_cluster_group

The TMC Cluster Backup resource allows requesting the creation of a cluster backup in Tanzu Mission Control (TMC). It also deals with managing the attributes and lifecycle of the cluster backup.

```terraform
resource "tmc_cluster_backup" "example" {
  name        = "example"
  labels = {
    Environment = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Tanzu Cluster Backup.
* `management_cluster_name` - (Required) Name of the Tanzu Management Cluster.
* `provisioner_name` - (Required) Name of the Tanzu provisioner associated with the backup.
* `cluster_name` - (Required) Name of the Tanzu cluster to backup.
* `labels` - (Optional) A map of labels to assign to the resource.
* `included_namespaces` - (Optional) The namespaces to be included for backup from. If empty, all namespaces are included.
* `excluded_namespaces` - (Optional) The namespaces to be excluded in the backup.
* `included_resources` - (Optional) The name list for the resources included into backup. If empty, all resources are included.
* `excluded_resources` - (Optional) The name list for the resources excluded in backup.
* `label_selector` - (Optional) Label query over a set of resources to be included in backup.
* `retention_period` - (Required) The backup retention period in seconds. (e.g., `3600s`, `32700s`)
* `storage_location` - (Required) The name of a BackupStorageLocation where the backup will be stored.
* `snapshot_volumes` - (Required) A boolean flag which specifies whether cloud snapshots of any PV's referenced in the set of objects included in the Backup are taken.
* `volume_snapshot_locations` - (Required) A list containing names of VolumeSnapshotLocations associated with this backup. Should not be left empty if `snapshot_volumes` is set to true.
* `include_cluster_scoped_resources` - (Optional) A boolean flag which specifies whether cluster-scoped resources were included for consideration in the backup. Defaults to true.

(Changing any arugment associated with the resource forces recreation as there is no update method associated with the backup resource as of now)

## Attributes Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The UID of the Tanzu Cluster Backup.
* `status` - Status of the found Cluster Backup.