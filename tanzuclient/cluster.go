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
			Id        string `json:"id,omitempty"`
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
			HighAvailability  bool     `json:"highAvailability,omitempty"`
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
	Status   *Status      `json:"status"`
}

type ClusterJSONObject struct {
	Cluster Cluster `json:"cluster"`
}

// Options interface for passing arguments to the
// functions neccessary to perform on the Cluster
type ClusterOpts struct {
	Region            string
	Version           string
	CredentialName    string
	AvailabilityZones []interface{}
	InstanceType      string
	VpcCidrBlock      string
	PodCidrBlock      string
	ServiceCidrBlock  string
	SshKey            string
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

func (c *Client) CreateCluster(name string, managementClusterName string, provisionerName string, cluster_group string, description string, labels map[string]interface{}, opts *ClusterOpts) (*Cluster, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters", c.baseURL)

	awsSpec := buildAwsJsonObject(opts)

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

func (c *Client) UpdateCluster(name string, managementClusterName string, provisionerName string, cluster_group string, description string, resourceVersion string, labels map[string]interface{}, opts *ClusterOpts) (*Cluster, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, name, managementClusterName, provisionerName)

	awsSpec := buildAwsJsonObject(opts)

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

func buildAwsJsonObject(opts *ClusterOpts) AWSCluster {

	var newAwsSpec AWSCluster

	newAwsSpec.Distribution.ProvisionerCredentialName = opts.CredentialName
	newAwsSpec.Distribution.Region = opts.Region
	newAwsSpec.Distribution.Version = opts.Version

	newAwsSpec.Settings.Network.ClusterNetwork.Pods = make([]struct {
		CidrBlocks string "json:\"cidrBlocks\""
	}, 1)
	newAwsSpec.Settings.Network.ClusterNetwork.Pods[0].CidrBlocks = opts.PodCidrBlock
	newAwsSpec.Settings.Network.ClusterNetwork.Services = make([]struct {
		CidrBlocks string "json:\"cidrBlocks\""
	}, 1)
	newAwsSpec.Settings.Network.ClusterNetwork.Services[0].CidrBlocks = opts.ServiceCidrBlock

	newAwsSpec.Settings.Network.Provider.Vpc.CidrBlock = opts.VpcCidrBlock

	newAwsSpec.Settings.Security.SshKey = opts.SshKey

	newAwsSpec.Topology.ControlPlane.InstanceType = opts.InstanceType

	var azList []string
	for i := 0; i < len(opts.AvailabilityZones); i++ {
		azList = append(azList, opts.AvailabilityZones[i].(string))
	}

	newAwsSpec.Topology.ControlPlane.AvailabilityZones = azList
	if len(azList) > 1 {
		newAwsSpec.Topology.ControlPlane.HighAvailability = true
	}

	return newAwsSpec
}
