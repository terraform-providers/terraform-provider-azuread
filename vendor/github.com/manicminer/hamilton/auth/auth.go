package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/pkcs12"
	"golang.org/x/oauth2"

	"github.com/manicminer/hamilton/environments"
)

// Authorizer is anything that can return an access token for authorizing API connections
type Authorizer interface {
	Token() (*oauth2.Token, error)
}

type Api int

const (
	MsGraph Api = iota
	AadGraph
)

// NewAuthorizer returns a suitable Authorizer depending on what is defined in the Config
// Authorizers are selected for authentication methods in the following preferential order:
// - Client certificate authentication
// - Client secret authentication
// - Azure CLI authentication
//
// Whether one of these is returned depends on whether it is enabled in the Config, and whether sufficient
// configuration fields are set to enable that authentication method.
//
// For client certificate authentication, specify TenantID, ClientID and ClientCertData / ClientCertPath.
// For client secret authentication, specify TenantID, ClientID and ClientSecret.
// MSI authentication (if enabled) using the Azure Metadata Service is then attempted
// Azure CLI authentication (if enabled) is attempted last
//
// It's recommended to only enable the mechanisms you have configured and are known to work in the execution
// environment. If any authentication mechanism fails due to misconfiguration or some other error, the function
// will return (nil, error) and later mechanisms will not be attempted.
func (c *Config) NewAuthorizer(ctx context.Context, api Api) (Authorizer, error) {
	if c.EnableClientCertAuth && strings.TrimSpace(c.TenantID) != "" && strings.TrimSpace(c.ClientID) != "" && (len(c.ClientCertData) > 0 || strings.TrimSpace(c.ClientCertPath) != "") {
		a, err := NewClientCertificateAuthorizer(ctx, c.Environment, api, c.Version, c.TenantID, c.ClientID, c.ClientCertData, c.ClientCertPath, c.ClientCertPassword)
		if err != nil {
			return nil, fmt.Errorf("could not configure ClientCertificate Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	if c.EnableClientSecretAuth && strings.TrimSpace(c.TenantID) != "" && strings.TrimSpace(c.ClientID) != "" && strings.TrimSpace(c.ClientSecret) != "" {
		a, err := NewClientSecretAuthorizer(ctx, c.Environment, api, c.Version, c.TenantID, c.ClientID, c.ClientSecret)
		if err != nil {
			return nil, fmt.Errorf("could not configure ClientCertificate Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	if c.EnableMsiAuth {
		a, err := NewMsiAuthorizer(ctx, c.Environment, api, c.MsiEndpoint)
		if err != nil {
			return nil, fmt.Errorf("could not configure MSI Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	if c.EnableAzureCliToken {
		a, err := NewAzureCliAuthorizer(ctx, api, c.TenantID)
		if err != nil {
			return nil, fmt.Errorf("could not configure AzureCli Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	return nil, fmt.Errorf("no Authorizer could be configured, please check your configuration")
}

// NewAzureCliAuthorizer returns an Authorizer which authenticates using the Azure CLI.
func NewAzureCliAuthorizer(ctx context.Context, api Api, tenantId string) (Authorizer, error) {
	conf, err := NewAzureCliConfig(api, tenantId)
	if err != nil {
		return nil, err
	}
	return conf.TokenSource(ctx), nil
}

// NewMsiAuthorizer returns an authorizer which uses managed service identity to for authentication.
func NewMsiAuthorizer(ctx context.Context, environment environments.Environment, api Api, msiEndpoint string) (Authorizer, error) {
	conf, err := NewMsiConfig(ctx, resource(environment, api), msiEndpoint)
	if err != nil {
		return nil, err
	}
	return conf.TokenSource(ctx), nil
}

// NewClientCertificateAuthorizer returns an authorizer which uses client certificate authentication.
func NewClientCertificateAuthorizer(ctx context.Context, environment environments.Environment, api Api, tokenVersion TokenVersion, tenantId, clientId string, pfxData []byte, pfxPath, pfxPass string) (Authorizer, error) {
	if len(pfxData) == 0 {
		var err error
		pfxData, err = ioutil.ReadFile(pfxPath)
		if err != nil {
			return nil, fmt.Errorf("could not read pkcs12 store at %q: %s", pfxPath, err)
		}
	}

	key, cert, err := pkcs12.Decode(pfxData, pfxPass)
	if err != nil {
		return nil, fmt.Errorf("could not decode pkcs12 credential store: %s", err)
	}

	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unsupported non-rsa key was found in pkcs12 store %q", pfxPath)
	}

	conf := ClientCredentialsConfig{
		ClientID:    clientId,
		PrivateKey:  x509.MarshalPKCS1PrivateKey(priv),
		Certificate: cert.Raw,
		Scopes:      scopes(environment, api),
		TokenURL:    TokenEndpoint(environment.AzureADEndpoint, tenantId, tokenVersion),
	}
	if tokenVersion == TokenVersion1 {
		conf.Resource = resource(environment, api)
	}
	return conf.TokenSource(ctx, ClientCredentialsAssertionType), nil
}

// NewClientSecretAuthorizer returns an authorizer which uses client secret authentication.
func NewClientSecretAuthorizer(ctx context.Context, environment environments.Environment, api Api, tokenVersion TokenVersion, tenantId, clientId, clientSecret string) (Authorizer, error) {
	conf := ClientCredentialsConfig{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       scopes(environment, api),
		TokenURL:     TokenEndpoint(environment.AzureADEndpoint, tenantId, tokenVersion),
	}
	if tokenVersion == TokenVersion1 {
		conf.Resource = resource(environment, api)
	}
	return conf.TokenSource(ctx, ClientCredentialsSecretType), nil
}

func TokenEndpoint(endpoint environments.AzureADEndpoint, tenant string, version TokenVersion) (e string) {
	if tenant == "" {
		tenant = "common"
	}
	e = fmt.Sprintf("%s/%s/oauth2", endpoint, tenant)
	if version == TokenVersion2 {
		e = fmt.Sprintf("%s/%s", e, "v2.0")
	}
	e = fmt.Sprintf("%s/token", e)
	return
}

func scopes(env environments.Environment, api Api) (s []string) {
	switch api {
	case MsGraph:
		s = []string{fmt.Sprintf("%s/.default", env.MsGraph.Endpoint)}
	case AadGraph:
		s = []string{fmt.Sprintf("%s/.default", env.AadGraph.Endpoint)}
	}
	return
}

func resource(env environments.Environment, api Api) (r string) {
	switch api {
	case MsGraph:
		r = fmt.Sprintf("%s/", env.MsGraph.Endpoint)
	case AadGraph:
		r = fmt.Sprintf("%s/", env.AadGraph.Endpoint)
	}
	return
}
