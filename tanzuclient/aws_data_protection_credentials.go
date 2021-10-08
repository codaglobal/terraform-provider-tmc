package tanzuclient

import (
	"fmt"
	"net/http"
)

type CredentialMetaData struct {
	Provider string `json:"provider"`
}

type CredentialCapability struct {
	Provider string `json:"provider"`
}

type AwsCredential struct {
	AccountID string `json:"accountId,omitempty"`
}

type AwsDataProtectionCredential struct {
	FullName *FullName `json:"fullName"`
	MetaData *MetaData `json:"meta"`
	Spec     struct {
		Meta       *CredentialMetaData `json:"meta"`
		Capability string              `json:"capability"`
		Data       struct {
			AwsCredential struct {
				AccountID string `json:"accountId"`
				IamRole   struct {
					Arn string `json:"arn"`
				} `json:"iamRole"`
			} `json:"awsCredential"`
		} `json:"data"`
	} `json:"spec"`
	Status struct {
		Phase string `json:"phase"`
	} `json:"status"`
}

type AwsDataProtectionCredentialResponse struct {
	AwsDataProtectionCredential AwsDataProtectionCredential `json:"credential"`
}

func (c *Client) GetAwsDataProtectionCredential(name string) (*AwsDataProtectionCredential, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/account/credentials/%s", c.baseURL, name)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := AwsDataProtectionCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.AwsDataProtectionCredential, nil
}
