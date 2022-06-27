package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Volume struct {
	Name         string `json:"name"`
	MountPath    string `json:"mountPath"`
	Capacity     int    `json:"capacity"`
	StorageClass string `json:"storageClass"`
}

type VsphereControlPlane struct {
	Volumes      []Volume `json:"volumes,omitempty"`
	Class        string   `json:"class"`
	StorageClass string   `json:"storageClass"`
}

type Vsphere struct {
	Settings struct {
		Network struct {
			Pods struct {
				CidrBlocks []string `json:"cidrBlocks"`
			} `json:"pods"`
			Services struct {
				CidrBlocks []string `json:"cidrBlocks"`
			} `json:"services"`
		} `json:"network"`
	} `json:"settings"`
	Distribution struct {
		Version string `json:"version"`
	} `json:"distribution"`
	Topology struct {
		ControlPlane VsphereControlPlane `json:"controlPlane"`
		NodePools    []VsphereNodepool   `json:"nodePools"`
	} `json:"topology"`
}

type VsphereNodepool struct {
	Spec VsphereNodeSpec `json:"spec"`
	Info struct {
		Name string `json:"name"`
	} `json:"info"`
}

type VsphereNodeSpec struct {
	NodeCount string              `json:"workerNodeCount"`
	NodeSpec  VsphereControlPlane `json:"tkgServiceVsphere"`
}

type VsphereSpec struct {
	ClusterGroupName  string  `json:"clusterGroupName"`
	TkgVsphereService Vsphere `json:"tkgServiceVsphere"`
}

type VsphereCluster struct {
	FullName *FullName    `json:"fullName"`
	Meta     *MetaData    `json:"meta"`
	Spec     *VsphereSpec `json:"spec"`
}

type VsphereJsonObject struct {
	Cluster VsphereCluster `json:"cluster"`
}

type VpshereNodepoolOpts struct {
	Name            string
	Class           string
	StorageClass    string
	WorkerNodeCount int
}

type VsphereOpts struct {
	Version          string
	Class            string
	StorageClass     string
	PodCidrBlock     string
	ServiceCidrBlock string
	NodepoolOpts     []VpshereNodepoolOpts
}

func (c *Client) GetVsphereCluster(fullName string, managementClusterName string, provisionerName string) (*VsphereCluster, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, fullName, managementClusterName, provisionerName)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := VsphereJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Cluster, nil
}

func (c *Client) CreateVsphereCluster(name string, managementClusterName string, provisionerName string, cluster_group string, description string, labels map[string]interface{}, opts *VsphereOpts) (*VsphereCluster, error) {

	nodePoolSpec := makeNodePoolSpec(opts.NodepoolOpts)

	newCluster := &VsphereCluster{
		FullName: &FullName{
			Name:                  name,
			ManagementClusterName: managementClusterName,
			ProvisionerName:       provisionerName,
		},
		Meta: &MetaData{
			Description: description,
			Labels:      labels,
		},
		Spec: &VsphereSpec{
			ClusterGroupName: cluster_group,
			TkgVsphereService: Vsphere{
				Settings: struct {
					Network struct {
						Pods struct {
							CidrBlocks []string "json:\"cidrBlocks\""
						} "json:\"pods\""
						Services struct {
							CidrBlocks []string "json:\"cidrBlocks\""
						} "json:\"services\""
					} "json:\"network\""
				}{
					Network: struct {
						Pods struct {
							CidrBlocks []string "json:\"cidrBlocks\""
						} "json:\"pods\""
						Services struct {
							CidrBlocks []string "json:\"cidrBlocks\""
						} "json:\"services\""
					}{
						Pods: struct {
							CidrBlocks []string "json:\"cidrBlocks\""
						}{
							CidrBlocks: []string{opts.PodCidrBlock},
						},
						Services: struct {
							CidrBlocks []string "json:\"cidrBlocks\""
						}{
							CidrBlocks: []string{opts.ServiceCidrBlock},
						},
					},
				},
				Distribution: struct {
					Version string "json:\"version\""
				}{
					Version: opts.Version,
				},
				Topology: struct {
					ControlPlane VsphereControlPlane "json:\"controlPlane\""
					NodePools    []VsphereNodepool   "json:\"nodePools\""
				}{
					ControlPlane: VsphereControlPlane{
						Class:        opts.Class,
						StorageClass: opts.StorageClass,
					},
					NodePools: nodePoolSpec,
				},
			},
		},
	}

	newClusterObject := &VsphereJsonObject{
		Cluster: *newCluster,
	}

	json_data, err := json.Marshal(newClusterObject) // returns []byte
	if err != nil {
		return nil, err
	}

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters", c.baseURL)

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := VsphereJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Cluster, nil
}

func (c *Client) DeleteVsphereCluster(name string, managementClusterName string, provisionerName string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := VsphereJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}

func makeNodePoolSpec(vpshereNodepoolOpts []VpshereNodepoolOpts) []VsphereNodepool {
	npSpec := make([]VsphereNodepool, 0)

	for i := 0; i < len(vpshereNodepoolOpts); i++ {
		toAppend := &VsphereNodepool{
			Spec: VsphereNodeSpec{
				NodeCount: strconv.Itoa(vpshereNodepoolOpts[i].WorkerNodeCount),
				NodeSpec: VsphereControlPlane{
					Class:        vpshereNodepoolOpts[i].Class,
					StorageClass: vpshereNodepoolOpts[i].StorageClass,
				},
			},
			Info: struct {
				Name string "json:\"name\""
			}{
				Name: vpshereNodepoolOpts[i].Name,
			},
		}

		npSpec = append(npSpec, *toAppend)
	}

	return npSpec
}
