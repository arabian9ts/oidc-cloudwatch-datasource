package plugin

import (
	"encoding/json"
	"errors"

	"github.com/arabian9ts/oidc-cloudwatch-datasource/pkg/google"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

const OIDCTokenIssuerGoogle = "GOOGLE"

type DatasourceConfig struct {
	OIDCTokenIssuer  string `json:"issuer"`
	CredentialsPath  string `json:"credentialsPath"`
	AssumeRole       string `json:"assumeRole"`
	STSRegion        string `json:"stsRegion"`
	MonitoringRegion string `json:"region"`
}

func (sc *DatasourceConfig) Validate() error {
	if sc.OIDCTokenIssuer == "" {
		return errors.New("oidc token issuer is required")
	}
	if sc.AssumeRole == "" {
		return errors.New("assume role is required")
	}
	if sc.STSRegion == "" {
		return errors.New("sts region is required")
	}
	if sc.MonitoringRegion == "" {
		return errors.New("monitoring region is required")
	}
	return nil
}

func GetSettings(s backend.DataSourceInstanceSettings) (*DatasourceConfig, error) {
	config := &DatasourceConfig{}
	if err := json.Unmarshal(s.JSONData, config); err != nil {
		return nil, err
	}
	return config, config.Validate()
}

func newIssuer(cfg *DatasourceConfig) (TokenIssuer, error) {
	switch cfg.OIDCTokenIssuer {
	case OIDCTokenIssuerGoogle:
		return google.NewTokenIssuer(&google.Config{
			CredsFilePath: cfg.CredentialsPath,
			Audience:      PluginID,
		}), nil
	}
	return nil, errors.New("invalid issuer")
}
