package wrappers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	commonParams "github.com/checkmarx/ast-cli/internal/params"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	failedToGetGroups = "Failed to parse list response"
)

type GroupsHTTPWrapper struct {
	path string
}

type Group struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func NewHTTPGroupsWrapper(path string) GroupsWrapper {
	tenant := viper.GetString(commonParams.TenantKey)
	tenantPath := strings.Replace(path, "organization", tenant, 1)

	return &GroupsHTTPWrapper{path: tenantPath}
}

func (g *GroupsHTTPWrapper) Get(groupName string) ([]Group, error) {
	clientTimeout := viper.GetUint(commonParams.ClientTimeoutKey)
	reportPath := fmt.Sprintf("%s?groupName=%s", g.path, groupName)
	resp, err := SendHTTPRequest(http.MethodGet, reportPath, nil, true, clientTimeout)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	switch resp.StatusCode {
	case http.StatusBadRequest, http.StatusInternalServerError:
		errorMsg := ErrorMsg{}
		err = decoder.Decode(&errorMsg)
		if err != nil {
			return nil, errors.Wrapf(err, err.Error())
		}
		return nil, errors.Errorf("%s: CODE: %d, %s", failedToGetGroups, errorMsg.Code, errorMsg.Message)
	case http.StatusOK:
		var groups []Group
		err = decoder.Decode(&groups)
		if err != nil {
			return nil, errors.Wrapf(err, failedToGetGroups)
		}
		return groups, nil
	default:
		return nil, errors.Errorf("response status code %d", resp.StatusCode)
	}
}