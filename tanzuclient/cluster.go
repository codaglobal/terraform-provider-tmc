package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Network struct {
	ClusterNetwork struct {
		Pods []struct {
			CidrBlocks string `json:"cidrBlocks"`
		} `json:"pods"`
		Services []struct {
			CidrBlocks string `json:"cidrBlocks"`
		} `json:"services"`
	} `json:"cluster"`
	Provider struct {
		Vpc struct {
			CidrBlock string `json:"cidrBlock"`
		} `json:"vpc"`
		Subnets []struct {
			Id       string `json:"id"`
			IsPublic bool   `json:"isPublic"`
		} `json:"subnets,omitempty"`
	} `json:"provider"`
}

type AWSCluster struct {
	Distribution struct {
		ProvisionerCredentialName string `json:"provisionerCredentialName"`
		Region                    string `json:"region"`
		Version                   string `json:"version"`
	} `json:"distribution"`
	Settings struct {
		Network  Network `json:"network"`
		Security struct {
			SshKey string `json:"sshKey"`
		} `json:"security"`
	} `json:"settings"`
	Topology struct {
		ControlPlane struct {
			AvailabilityZones []string `json:"availabilityZones"`
			InstanceType      string   `json:"instanceType"`
		} `json:"controlPlane"`
		NodePools []struct {
			Spec struct {
				WorkerNodeCount string `json:"workerNodeCount"`
				NodeTkgAws      struct {
					InstanceType     string                   `json:"instanceType"`
					AvailabilityZone string                   `json:"availabilityZone"`
					SubnetId         string                   `json:"subnetId"`
					Version          string                   `json:"version"`
					NodePlacement    []map[string]interface{} `json:"nodePlacement,omitempty"`
				} `json:"tkgAws"`
			} `json:"spec"`
			Info struct {
				Name string `json:"name"`
			} `json:"info"`
		} `json:"nodePools,omitempty"`
	} `json:"topology"`
}

type ClusterSpec struct {
	ClusterGroupName string     `json:"clusterGroupName"`
	TkgAws           AWSCluster `json:"tkgAws,omitempty"`
}

type Cluster struct {
	FullName *FullName    `json:"fullName"`
	Meta     *MetaData    `json:"meta"`
	Spec     *ClusterSpec `json:"spec"`
}

type ClusterJSONObject struct {
	Cluster Cluster `json:"cluster"`
}

func (c *Client) GetCluster(fullName string, managementClusterName string, provisionerName string) (*Cluster, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, fullName, managementClusterName, provisionerName)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := ClusterJSONObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Cluster, nil
}

func (c *Client) CreateCluster(name string, managementClusterName string, provisionerName string, cluster_group string, description string, labels map[string]interface{}, spec map[string]interface{}) (*Cluster, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters", c.baseURL)

	awsSpec := buildAwsJsonObject(spec)

	newCluster := &Cluster{
		FullName: &FullName{
			Name:                  name,
			ManagementClusterName: managementClusterName,
			ProvisionerName:       provisionerName,
		},
		Meta: &MetaData{
			Description: description,
			Labels:      labels,
		},
		Spec: &ClusterSpec{
			ClusterGroupName: cluster_group,
			TkgAws:           awsSpec,
		},
	}

	newClusterObject := &ClusterJSONObject{
		Cluster: *newCluster,
	}

	json_data, err := json.Marshal(newClusterObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := ClusterJSONObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Cluster, nil
}

func buildAwsJsonObject(spec map[string]interface{}) AWSCluster {

	// subnets := spec["subnets"].([]interface{})

	var newAwsSpec AWSCluster

	newAwsSpec.Distribution.ProvisionerCredentialName = spec["credential_name"].(string)
	newAwsSpec.Distribution.Region = spec["region"].(string)
	newAwsSpec.Distribution.Version = spec["version"].(string)

	newAwsSpec.Settings.Network.ClusterNetwork.Pods = make([]struct {
		CidrBlocks string "json:\"cidrBlocks\""
	}, 1)
	newAwsSpec.Settings.Network.ClusterNetwork.Pods[0].CidrBlocks = "192.168.0.0/16"
	newAwsSpec.Settings.Network.ClusterNetwork.Services = make([]struct {
		CidrBlocks string "json:\"cidrBlocks\""
	}, 1)
	newAwsSpec.Settings.Network.ClusterNetwork.Services[0].CidrBlocks = "10.96.0.0/12"

	newAwsSpec.Settings.Network.Provider.Vpc.CidrBlock = spec["vpc_cidrblock"].(string)
	// newAwsSpec.Settings.Network.Provider.Subnets = make([]struct {
	// 	Id       string "json:\"id\""
	// 	IsPublic bool   "json:\"isPublic\""
	// }, len(subnets))

	// for i := 0; i < len(subnets); i++ {
	// 	newAwsSpec.Settings.Network.Provider.Subnets[i].Id = subnets[i].(map[string]interface{})["id"].(string)
	// 	newAwsSpec.Settings.Network.Provider.Subnets[i].IsPublic = subnets[i].(map[string]interface{})["is_public"].(bool)
	// }

	newAwsSpec.Settings.Security.SshKey = spec["ssh_key"].(string)

	newAwsSpec.Topology.ControlPlane.InstanceType = spec["instance_type"].(string)
	newAwsSpec.Topology.ControlPlane.AvailabilityZones = []string{spec["availability_zones"].([]interface{})[0].(string)}

	return newAwsSpec

}

func (c *Client) DeleteCluster(name string, managementClusterName string, provisionerName string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := ClusterJSONObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateCluster(name string, managementClusterName string, provisionerName string, cluster_group string, description string, resourceVersion string, labels map[string]interface{}, spec map[string]interface{}) (*Cluster, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, name, managementClusterName, provisionerName)

	awsSpec := buildAwsJsonObject(spec)

	newCluster := &Cluster{
		FullName: &FullName{
			Name:                  name,
			ManagementClusterName: managementClusterName,
			ProvisionerName:       provisionerName,
		},
		Meta: &MetaData{
			ResourceVersion: resourceVersion,
			Description:     description,
			Labels:          labels,
		},
		Spec: &ClusterSpec{
			ClusterGroupName: cluster_group,
			TkgAws:           awsSpec,
		},
	}

	newClusterObject := &ClusterJSONObject{
		Cluster: *newCluster,
	}

	json_data, err := json.Marshal(newClusterObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := ClusterJSONObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Cluster, nil
}
