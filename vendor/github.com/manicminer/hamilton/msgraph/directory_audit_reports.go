package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/manicminer/hamilton/odata"
)

// DirectoryAuditReportsClient performs operations on directory Audit reports.
type DirectoryAuditReportsClient struct {
	BaseClient Client
}

// NewDirectoryAuditReportsClient returns a new DirectoryAuditReportsClient.
func NewDirectoryAuditReportsClient(tenantId string) *DirectoryAuditReportsClient {
	return &DirectoryAuditReportsClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// List returns a list of Directory audit report logs, optionally queried using OData.
func (c *DirectoryAuditReportsClient) List(ctx context.Context, query odata.Query) (*[]DirectoryAudit, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/auditLogs/directoryAudits",
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("DirectoryAuditReportsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var data struct {
		DirectoryAuditReports []DirectoryAudit `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &data.DirectoryAuditReports, status, nil
}

// Get retrieves a Directory audit report.
func (c *DirectoryAuditReportsClient) Get(ctx context.Context, id string) (*DirectoryAudit, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/auditLogs/directoryAudits/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("DirectoryAuditReportsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var directoryAuditReport DirectoryAudit
	if err := json.Unmarshal(respBody, &directoryAuditReport); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &directoryAuditReport, status, nil
}
