package tanzuclient

import (
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
