/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/huaweicloud/external-sfs/pkg/logger"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack"

	"github.com/gophercloud/gophercloud"
	nativeopenstack "github.com/gophercloud/gophercloud/openstack"
)

// CloudCredentials define
type CloudCredentials struct {
	Global struct {
		AuthURL        string `gcfg:"auth-url"`
		Username       string
		UserID         string `gcfg:"user-id"`
		Password       string
		TenantID       string `gcfg:"tenant-id"`
		TenantName     string `gcfg:"tenant-name"`
		DomainID       string `gcfg:"domain-id"`
		DomainName     string `gcfg:"domain-name"`
		Region         string
		AccessKey      string `gcfg:"access-key"`
		SecretKey      string `gcfg:"secret-key"`
		CACertFile     string `gcfg:"cacert-file"`
		ClientCertFile string `gcfg:"cert"`
		ClientKeyFile  string `gcfg:"key"`
		EndpointType   string `gcfg:"endpoint-type"`
		Insecure       bool
	}

	CloudClient     *golangsdk.ProviderClient
	OpenStackClient *gophercloud.ProviderClient
}

// Validate CloudCredentials
func (c *CloudCredentials) Validate() error {
	validEndpoint := false
	validEndpoints := []string{
		"internal", "internalURL",
		"admin", "adminURL",
		"public", "publicURL",
		"",
	}

	for _, endpoint := range validEndpoints {
		if c.Global.EndpointType == endpoint {
			validEndpoint = true
		}
	}

	if !validEndpoint {
		return fmt.Errorf("Invalid endpoint type provided")
	}

	err := c.newCloudClient()
	if err != nil {
		return err
	}

	return c.newOpenStackClient()
}

// newCloudClient returns new cloud client
func (c *CloudCredentials) newCloudClient() error {
	ao := golangsdk.AuthOptions{
		DomainID:         c.Global.DomainID,
		DomainName:       c.Global.DomainName,
		IdentityEndpoint: c.Global.AuthURL,
		Password:         c.Global.Password,
		TenantID:         c.Global.TenantID,
		TenantName:       c.Global.TenantName,
		Username:         c.Global.Username,
		UserID:           c.Global.UserID,
		// allow to renew tokens
		AllowReauth: true,
	}

	client, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}

	config := &tls.Config{}
	if c.Global.CACertFile != "" {
		caCert, _, err := ReadContents(c.Global.CACertFile)
		if err != nil {
			return fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.Global.Insecure {
		config.InsecureSkipVerify = true
	}

	if c.Global.ClientCertFile != "" && c.Global.ClientKeyFile != "" {
		clientCert, _, err := ReadContents(c.Global.ClientCertFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Cert: %s", err)
		}
		clientKey, _, err := ReadContents(c.Global.ClientKeyFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return err
		}

		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()
	}

	// if OS_DEBUG is set, log the requests and responses
	var osDebug bool
	if os.Getenv("OS_DEBUG") != "" {
		osDebug = true
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	client.HTTPClient = http.Client{
		Transport: &logger.LogRoundTripper{
			Rt:      transport,
			OsDebug: osDebug,
		},
	}

	err = openstack.Authenticate(client, ao)
	if err != nil {
		return err
	}

	c.CloudClient = client

	return nil
}

// newOpenStackClient returns new native openstack client
func (c *CloudCredentials) newOpenStackClient() error {
	ao := gophercloud.AuthOptions{
		DomainID:         c.Global.DomainID,
		DomainName:       c.Global.DomainName,
		IdentityEndpoint: c.Global.AuthURL,
		Password:         c.Global.Password,
		TenantID:         c.Global.TenantID,
		TenantName:       c.Global.TenantName,
		Username:         c.Global.Username,
		UserID:           c.Global.UserID,
		// allow to renew tokens
		AllowReauth: true,
	}

	client, err := nativeopenstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}

	config := &tls.Config{}
	if c.Global.CACertFile != "" {
		caCert, _, err := ReadContents(c.Global.CACertFile)
		if err != nil {
			return fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.Global.Insecure {
		config.InsecureSkipVerify = true
	}

	if c.Global.ClientCertFile != "" && c.Global.ClientKeyFile != "" {
		clientCert, _, err := ReadContents(c.Global.ClientCertFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Cert: %s", err)
		}
		clientKey, _, err := ReadContents(c.Global.ClientKeyFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return err
		}

		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()
	}

	// if OS_DEBUG is set, log the requests and responses
	var osDebug bool
	if os.Getenv("OS_DEBUG") != "" {
		osDebug = true
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	client.HTTPClient = http.Client{
		Transport: &logger.LogRoundTripper{
			Rt:      transport,
			OsDebug: osDebug,
		},
	}

	err = nativeopenstack.Authenticate(client, ao)
	if err != nil {
		return err
	}

	c.OpenStackClient = client

	return nil
}

// getEndpointType returns cloud endpoint type
func (c *CloudCredentials) getEndpointType() golangsdk.Availability {
	if c.Global.EndpointType == "internal" || c.Global.EndpointType == "internalURL" {
		return golangsdk.AvailabilityInternal
	}
	if c.Global.EndpointType == "admin" || c.Global.EndpointType == "adminURL" {
		return golangsdk.AvailabilityAdmin
	}
	return golangsdk.AvailabilityPublic
}

// getNativeEndpointType returns native openstack endpoint type
func (c *CloudCredentials) getNativeEndpointType() gophercloud.Availability {
	if c.Global.EndpointType == "internal" || c.Global.EndpointType == "internalURL" {
		return gophercloud.AvailabilityInternal
	}
	if c.Global.EndpointType == "admin" || c.Global.EndpointType == "adminURL" {
		return gophercloud.AvailabilityAdmin
	}
	return gophercloud.AvailabilityPublic
}

// SFSV2Client return sfs v2 client
func (c *CloudCredentials) SFSV2Client() (*golangsdk.ServiceClient, error) {
	return openstack.NewSharedFileSystemV2(c.CloudClient, golangsdk.EndpointOpts{
		Region:       c.Global.Region,
		Availability: c.getEndpointType(),
	})
}

// NetworkingV1Client return native networking v1 client
func (c *CloudCredentials) NetworkingV1Client() (*golangsdk.ServiceClient, error) {
	return openstack.NewNetworkV1(c.CloudClient, golangsdk.EndpointOpts{
		Region:       c.Global.Region,
		Availability: c.getEndpointType(),
	})
}

// ComputeV2Client return native compute v2 client
func (c *CloudCredentials) ComputeV2Client() (*gophercloud.ServiceClient, error) {
	return nativeopenstack.NewComputeV2(c.OpenStackClient, gophercloud.EndpointOpts{
		Region:       c.Global.Region,
		Availability: c.getNativeEndpointType(),
	})
}
