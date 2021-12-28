package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Scan struct {
	FullName *FullName              `json:"fullName"`
	Meta     *MetaData              `json:"meta"`
	Spec     map[string]interface{} `json:"spec"`
}

type ScanJsonObject struct {
	Scan Scan `json:"scan"`
}

func (c *Client) CreateScan(clusterName string, managementClusterName string, provisionerName string, scanType string) (*Scan, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/inspection/scans?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, managementClusterName, provisionerName)

	newScan := &Scan{
		FullName: &FullName{
			ClusterName:           clusterName,
			ProvisionerName:       provisionerName,
			ManagementClusterName: managementClusterName,
		},
	}

	switch scanType {
	case "lite":
		newScan.Spec = map[string]interface{}{
			"liteSpec": map[string]interface{}{},
		}
	case "cis":
		newScan.Spec = map[string]interface{}{
			"cisSpec": map[string]interface{}{},
		}
	case "conformance":
		newScan.Spec = map[string]interface{}{
			"conformanceSpec": map[string]interface{}{},
		}
	}

	newScanObject := &ScanJsonObject{
		Scan: *newScan,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(newScanObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := ScanJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Scan, nil
}

func (c *Client) GetScan(name string, clusterName string, managementClusterName string, provisionerName string) (*Scan, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/inspection/scans/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := ScanJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Scan, nil
}

func (c *Client) DeleteScan(name string, clusterName string, managementClusterName string, provisionerName string) error {

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/inspection/scans/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := ScanJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}
