package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NodeName struct {
	OrgID                 string `json:"orgId"`
	ClusterName           string `json:"clusterName"`
	ManagementClusterName string `json:"managementClusterName"`
	ProvisionerName       string `json:"provisionerName"`
	Name                  string `json:"name"`
}

type AwsNodeSpec struct {
	InstanceType     string `json:"instanceType"`
	AvailabilityZone string `json:"availabilityZone"`
	Version          string `json:"version"`
}

type AwsNodePool struct {
	NodeLabels      map[string]interface{} `json:"cloudLabels,omitempty"`
	CloudLabels     map[string]interface{} `json:"nodeLabels,omitempty"`
	WorkerNodeCount string                 `json:"workerNodeCount"`
	NodeTkgAws      AwsNodeSpec            `json:"tkgAws"`
}

type NodePool struct {
	FullName *NodeName    `json:"fullName"`
	Meta     *MetaData    `json:"meta"`
	Spec     *AwsNodePool `json:"spec"`
}

type NodePoolJsonObject struct {
	NodePool NodePool `json:"nodepool"`
}

func (c *Client) CreateNodePool(name string, managementClusterName string, provisionerName string, clusterName string, description string, cloudLabels map[string]interface{}, nodeLabels map[string]interface{}, nodeCount int, spec map[string]interface{}) (*NodePool, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/nodepools", c.baseURL, clusterName)

	awsNodeSpec := buildAwsNodeSpec(spec)

	newNodePool := &NodePool{
		FullName: &NodeName{
			ClusterName:           clusterName,
			Name:                  name,
			ManagementClusterName: managementClusterName,
			ProvisionerName:       provisionerName,
		},
		Meta: &MetaData{
			Description: description,
		},
		Spec: &AwsNodePool{
			NodeLabels:      nodeLabels,
			CloudLabels:     cloudLabels,
			WorkerNodeCount: fmt.Sprint(nodeCount),
			NodeTkgAws:      awsNodeSpec,
		},
	}

	newNodePoolObject := &NodePoolJsonObject{
		NodePool: *newNodePool,
	}

	json_data, err := json.Marshal(newNodePoolObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := NodePoolJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.NodePool, nil

}

func (c *Client) GetNodePool(name string, clusterName string, managementClusterName string, provisionerName string) (*NodePool, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/nodepools/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := NodePoolJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.NodePool, nil
}

func (c *Client) UpdateNodePool(name string, managementClusterName string, provisionerName string, clusterName string, description string, cloudLabels map[string]interface{}, nodeLabels map[string]interface{}, nodeCount int, spec map[string]interface{}) (*NodePool, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/nodepools/%s", c.baseURL, clusterName, name)

	awsNodeSpec := buildAwsNodeSpec(spec)

	newNodePool := &NodePool{
		FullName: &NodeName{
			ClusterName:           clusterName,
			Name:                  name,
			ManagementClusterName: managementClusterName,
			ProvisionerName:       provisionerName,
		},
		Meta: &MetaData{
			Description: description,
		},
		Spec: &AwsNodePool{
			NodeLabels:      nodeLabels,
			CloudLabels:     cloudLabels,
			WorkerNodeCount: fmt.Sprint(nodeCount),
			NodeTkgAws:      awsNodeSpec,
		},
	}

	newNodePoolObject := &NodePoolJsonObject{
		NodePool: *newNodePool,
	}

	json_data, err := json.Marshal(newNodePoolObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := NodePoolJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.NodePool, nil
}

func (c *Client) DeleteNodePool(name string, clusterName string, managementClusterName string, provisionerName string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/nodepools/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := NodePoolJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}

func buildAwsNodeSpec(spec map[string]interface{}) AwsNodeSpec {
	var newAwsNodeSpec AwsNodeSpec

	newAwsNodeSpec.AvailabilityZone = spec["availability_zone"].(string)
	newAwsNodeSpec.InstanceType = spec["instance_type"].(string)
	newAwsNodeSpec.Version = spec["version"].(string)

	return newAwsNodeSpec
}
