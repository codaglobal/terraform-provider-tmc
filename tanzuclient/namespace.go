package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NamespaceSpec struct {
	WorkspaceName string `json:"workspaceName"`
}

type Namespace struct {
	FullName *FullName      `json:"fullName"`
	Meta     *MetaData      `json:"meta"`
	Spec     *NamespaceSpec `json:"spec"`
}

type NamespaceJsonObject struct {
	Namespace Namespace `json:"namespace"`
}

type NamespaceOpts struct {
	Description       string
	Labels            map[string]interface{}
	ManagementCluster string
	ProvisionerName   string
	ClusterName       string
	WorkspaceName     string
}

func (c *Client) CreateNamespace(name string, opts NamespaceOpts) (*Namespace, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/namespaces?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, opts.ClusterName, opts.ManagementCluster, opts.ProvisionerName)

	newNamespace := &Namespace{
		FullName: &FullName{
			Name:                  name,
			ProvisionerName:       opts.ProvisionerName,
			ManagementClusterName: opts.ManagementCluster,
			ClusterName:           opts.ClusterName,
		},
		Meta: &MetaData{
			Description: opts.Description,
			Labels:      opts.Labels,
		},
		Spec: &NamespaceSpec{
			WorkspaceName: opts.WorkspaceName,
		},
	}

	newNamespaceObject := &NamespaceJsonObject{
		Namespace: *newNamespace,
	}

	json_data, err := json.Marshal(newNamespaceObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := NamespaceJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Namespace, nil
}

func (c *Client) GetNamespace(name string, clusterName string, managementClusterName string, provisionerName string) (*Namespace, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/namespaces/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := NamespaceJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Namespace, nil
}

func (c *Client) DeleteNamespace(name string, clusterName string, managementClusterName string, provisionerName string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/namespaces/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := NamespaceJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}
