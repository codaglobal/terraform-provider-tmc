package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var k8sProviderTypeMap = map[string]string{
	"tkg":        "VMWARE_TANZU_KUBERNETES_GRID",
	"tkgservice": "VMWARE_TANZU_KUBERNETES_GRID_SERVICE",
	"tkghosted":  "VMWARE_TANZU_KUBERNETES_GRID_HOSTED",
	"other":      "KUBERNETES_PROVIDER_UNSPECIFIED",
}

type MgmtClusterSpec struct {
	KubernetesProviderType string `json:"kubernetesProviderType"`
	DefaultClusterGroup    string `json:"defaultClusterGroup"`
}

type ManagementCluster struct {
	FullName *FullName        `json:"fullName"`
	Meta     *MetaData        `json:"meta"`
	Spec     *MgmtClusterSpec `json:"spec"`
	Status   struct {
		RegistrationURL string `json:"registrationUrl,omitempty"`
	}
}

type MgmtClusterJsonObject struct {
	MgmtCluster ManagementCluster `json:"managementCluster"`
}

func (c *Client) CreateMgmtCluster(name string, defaultCg string, k8sProviderType string, description string, labels map[string]interface{}) (*ManagementCluster, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/managementclusters", c.baseURL)

	newMgmtCluster := &ManagementCluster{
		FullName: &FullName{
			Name: name,
		},
		Meta: &MetaData{
			Description: description,
			Labels:      labels,
		},
		Spec: &MgmtClusterSpec{
			KubernetesProviderType: k8sProviderTypeMap[k8sProviderType],
			DefaultClusterGroup:    defaultCg,
		},
	}

	newMgmtClusterObject := &MgmtClusterJsonObject{
		MgmtCluster: *newMgmtCluster,
	}

	json_data, err := json.Marshal(newMgmtClusterObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := MgmtClusterJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.MgmtCluster, nil
}

func (c *Client) GetMgmtCluster(name string) (*ManagementCluster, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/managementclusters/%s", c.baseURL, name)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := MgmtClusterJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.MgmtCluster, nil
}

func (c *Client) DeleteMgmtCluster(name string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/managementclusters/%s", c.baseURL, name)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := MgmtClusterJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}
