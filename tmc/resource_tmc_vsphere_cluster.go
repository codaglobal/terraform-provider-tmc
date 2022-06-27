package tmc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVsphereCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVsphereClusterCreate,
		ReadContext:   resourceVsphereClusterRead,
		DeleteContext: resourceVsphereClusterDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Cluster",
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource version of the Cluster",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Cluster",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !IsValidTanzuName(v) {
						errs = append(errs, fmt.Errorf("name should contain only lowercase letters, numbers or hyphens and should begin with either an alphabet or number"))
					}
					return
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the Cluster",
			},
			"management_cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of an existing management cluster to be used",
			},
			"provisioner_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of an existing provisioner to be used",
			},
			"cluster_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the cluster group",
			},
			"labels": labelsSchemaImmutable(),
			"pod_cidrblock": {
				Type:        schema.TypeString,
				Description: "CIDR block used by the Cluster's Pods",
				Optional:    true,
				ForceNew:    true,
				Default:     "192.168.0.0/16",
			},
			"service_cidrblock": {
				Type:        schema.TypeString,
				Description: "CIDR block used by the Cluster's Services",
				Optional:    true,
				ForceNew:    true,
				Default:     "10.96.0.0/12",
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Kubernetes version to be used",
				ForceNew:    true,
				Required:    true,
			},
			"control_plane_spec": {
				Type:        schema.TypeList,
				Description: "Contains information related to the Control Plane of the cluster",
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"class": {
							Type:        schema.TypeString,
							Description: "Indicates the size of the VMs to be provisioned",
							Required:    true,
							ForceNew:    true,
						},
						"storage_class": {
							Type:        schema.TypeString,
							Description: "Storage Class to be used for storage of the disks which store the root filesystems of the nodes",
							Required:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"nodepool": {
				Type:        schema.TypeList,
				Description: "Contains specifications for a nodepool which is part of the cluster",
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nodepool_name": {
							Type:        schema.TypeString,
							Description: "Determines the name of the nodepool",
							Required:    true,
							ForceNew:    true,
						},
						"worker_node_count": {
							Type:        schema.TypeInt,
							Description: "Determines the number of worker nodes provisioned",
							Required:    true,
							ForceNew:    true,
						},
						"node_class": {
							Type:        schema.TypeString,
							Description: "Determines the class of the worker node",
							Required:    true,
							ForceNew:    true,
						},
						"node_storage_class": {
							Type:        schema.TypeString,
							Description: "Determines the storage policy used for the worker node",
							Required:    true,
							ForceNew:    true,
						},
					},
				},
			},
		},
	}
}

func resourceVsphereClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels").(map[string]interface{})
	cluster_group := d.Get("cluster_group").(string)
	controlPlaneSpec := d.Get("control_plane_spec").([]interface{})[0].(map[string]interface{})

	nodePoolOpts := makeNodepoolOpts(d.Get("nodepool").([]interface{}))

	opts := &tanzuclient.VsphereOpts{
		Version:          d.Get("version").(string),
		Class:            controlPlaneSpec["class"].(string),
		StorageClass:     controlPlaneSpec["storage_class"].(string),
		PodCidrBlock:     d.Get("pod_cidrblock").(string),
		ServiceCidrBlock: d.Get("service_cidrblock").(string),
		NodepoolOpts:     nodePoolOpts,
	}

	vSphereCluster, err := client.CreateVsphereCluster(clusterName, managementClusterName, provisionerName, cluster_group, description, labels, opts)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create vSphere Cluster",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId(vSphereCluster.Meta.UID)

	resourceVsphereClusterRead(ctx, d, m)

	return diags

}

func resourceVsphereClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	cluster, err := client.GetVsphereCluster(clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read vSphere Cluster",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
	}

	d.Set("resource_version", cluster.Meta.ResourceVersion)
	d.Set("description", cluster.Meta.Description)
	d.Set("cluster_group", cluster.Spec.ClusterGroupName)
	if err := d.Set("labels", cluster.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read vSphere cluster",
			Detail:   fmt.Sprintf("Error getting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("version", cluster.Spec.TkgVsphereService.Distribution.Version)
	d.Set("pod_cidrblock", cluster.Spec.TkgVsphereService.Settings.Network.Pods.CidrBlocks[0])
	d.Set("service_cidrblock", cluster.Spec.TkgVsphereService.Settings.Network.Services.CidrBlocks[0])

	spec := make([]map[string]interface{}, 0)
	cp_spec := flatten_vsphere_control_plane_spec(&cluster.Spec.TkgVsphereService.Topology.ControlPlane)
	spec = append(spec, cp_spec)
	np_spec := flatten_vsphere_nodepool_spec(&cluster.Spec.TkgVsphereService.Topology.NodePools)

	if err := d.Set("control_plane_spec", spec); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read vSphere cluster",
			Detail:   fmt.Sprintf("Error getting control plane information for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	if err := d.Set("nodepool", np_spec); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read vSphere cluster",
			Detail:   fmt.Sprintf("Error getting control plane information for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	return diags
}

func resourceVsphereClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	if err := client.DeleteVsphereCluster(clusterName, managementClusterName, provisionerName); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete vSphere Cluster",
			Detail:   fmt.Sprintf("Error deleting resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId("")

	return nil
}

func flatten_vsphere_control_plane_spec(vsphereSpec *tanzuclient.VsphereControlPlane) map[string]interface{} {
	cp_spec := make(map[string]interface{})

	cp_spec["class"] = vsphereSpec.Class
	cp_spec["storage_class"] = vsphereSpec.StorageClass

	return cp_spec
}

func flatten_vsphere_nodepool_spec(vsphereNodepool *[]tanzuclient.VsphereNodepool) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	for i := 0; i < len(*vsphereNodepool); i++ {
		toAppend := make(map[string]interface{})

		toAppend["nodepool_name"] = (*vsphereNodepool)[i].Info.Name
		toAppend["worker_node_count"], _ = strconv.Atoi((*vsphereNodepool)[i].Spec.NodeCount)
		toAppend["node_class"] = (*vsphereNodepool)[i].Spec.NodeSpec.Class
		toAppend["node_storage_class"] = (*vsphereNodepool)[i].Spec.NodeSpec.StorageClass

		result = append(result, toAppend)
	}

	return result
}

func makeNodepoolOpts(arrayOfNodePoolSpec []interface{}) []tanzuclient.VpshereNodepoolOpts {

	npSpec := make([]tanzuclient.VpshereNodepoolOpts, 0)

	for i := 0; i < len(arrayOfNodePoolSpec); i++ {
		toAppend := &tanzuclient.VpshereNodepoolOpts{
			Name:            arrayOfNodePoolSpec[0].(map[string]interface{})["nodepool_name"].(string),
			Class:           arrayOfNodePoolSpec[0].(map[string]interface{})["node_class"].(string),
			StorageClass:    arrayOfNodePoolSpec[0].(map[string]interface{})["node_storage_class"].(string),
			WorkerNodeCount: arrayOfNodePoolSpec[0].(map[string]interface{})["worker_node_count"].(int),
		}

		npSpec = append(npSpec, *toAppend)
	}

	return npSpec
}
