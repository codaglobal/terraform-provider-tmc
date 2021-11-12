---
page_title: "TMC: tmc_cluster_backup"
layout: "tmc"
subcategory: "Cluster Backups"
description: |-
  Get information on a Tanzu Mission Control (TMC) Cluster Backup.
---

# Data Source: tmc_cluster_group

Use this data source to get the details about a cluster backup in TMC platform.

## Example Usage
# Get details of a cluster backup in the Tanzu platform.
```terraform
data "tmc_cluster_backup" "example" {
  name = "example-bkp"
  management_cluster_name = "example-mgmt-cluster"
  provisioner_name = "example-provisioner"
  cluster_name = "example-cluster"
}
```

## Argument Reference

* `name` - (Required) The name of the cluster backup to lookup in the TMC platform. If no cluster backup is found with this name, an error will be returned.
* `management_cluster_name` - (Required) Name of the Tanzu Management Cluster
* `provisioner_name` - (Required) Name of the Tanzu provisioner associated with the backup
* `cluster_name` - (Required) Name of the Tanzu cluster which was backed up

## Attributes Reference

* `id` - Unique Identifiers (UID) of the found Tanzu Cluster Backup in the TMC platform.
* `included_namespaces` - The namespaces to be included for backup from. If empty, all namespaces are included.
* `excluded_namespaces` - The namespaces to be excluded in the backup.
* `included_resources` - The name list for the resources included into backup. If empty, all resources are included.
* `excluded_resources` - The name list for the resources excluded in backup.
* `label_selector` - Label query over a set of resources to be included in backup.
* `retention_period` - The backup retention period in seconds. For ex., 3600s, 32700s
* `storage_location` - The name of a BackupStorageLocation where the backup has been stored.
* `snapshot_volumes` - A flag which specifies whether cloud snapshots of any PV's referenced in the set of objects included in the Backup are taken.
* `volume_snapshot_locations` - A list containing names of VolumeSnapshotLocations associated with this backup.
* `include_cluster_scoped_resources` - A flag which specifies whether cluster-scoped resources were included for consideration in the backup.
* `status` - Status of the found Cluster Backup.
* `labels` - A mapping of labels of the resource.