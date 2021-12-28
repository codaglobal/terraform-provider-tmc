package tanzuclient

import (
	"fmt"
	"net/http"
)

type AdminKubeConfig struct {
	KubeConfig string `json:"kubeconfig"`
}

func (c *Client) GetAdminKubeConfig(clusterName string, managementClusterName string, provisionerName string) (*AdminKubeConfig, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/adminkubeconfig?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, managementClusterName, provisionerName)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := AdminKubeConfig{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
