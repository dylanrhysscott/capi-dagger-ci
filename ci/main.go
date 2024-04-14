// A generated module for CapiDaggerCi functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"encoding/base64"
	"fmt"
)

type CapiDaggerCi struct {
	Token           *Secret
	SpacesAccessKey *Secret
	SpacesSecretKey *Secret
}

// Sets up the DO account credentials for subsquent functions
func (m *CapiDaggerCi) WithDOCreds(ctx context.Context, token *Secret, spacesAccessKey *Secret, spacesSecretKey *Secret) (*CapiDaggerCi, error) {
	m.Token = token
	m.SpacesAccessKey = spacesAccessKey
	m.SpacesSecretKey = spacesSecretKey
	return m, nil
}

func (m CapiDaggerCi) fetchPipelineCreds(ctx context.Context) (string, string, string, error) {
	tokenCleartext, err := m.Token.Plaintext(ctx)
	if err != nil {
		return "", "", "", fmt.Errorf("failed getting token: %s", err)
	}
	spacesAccessKeyCleartext, err := m.SpacesAccessKey.Plaintext(ctx)
	if err != nil {
		return "", "", "", fmt.Errorf("failed getting spaces access key: %s", err)
	}
	spacesSecretKeyCleartext, err := m.SpacesSecretKey.Plaintext(ctx)
	if err != nil {
		return "", "", "", fmt.Errorf("failed getting spaces secret key: %s", err)
	}
	return tokenCleartext, spacesAccessKeyCleartext, spacesSecretKeyCleartext, nil
}

// Runs terraform init, plan and apply to deploy infrastructure
func (m *CapiDaggerCi) DeployInfra(
	ctx context.Context,
	path Directory,
	// +optional
	// +default=true
	apply bool,
	// +optional
	// +default=false
	destroy bool) (*Container, error) {
	tokenCleartext, spacesAccessKeyCleartext, spacesSecretKeyCleartext, err := m.fetchPipelineCreds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting pipeline creds: %s", err)
	}
	planOpts := []string{"plan", "-out", "server.plan"}
	if destroy {
		planOpts = append(planOpts, "--destroy")
	}
	authenticatedTerraform := dag.Container().
		From("hashicorp/terraform:latest").
		WithDirectory("/infra", &path).
		WithWorkdir("/infra").
		WithEnvVariable("AWS_ACCESS_KEY_ID", spacesAccessKeyCleartext).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", spacesSecretKeyCleartext).
		WithEnvVariable("DIGITALOCEAN_TOKEN", tokenCleartext).
		WithExec([]string{"init"}).
		WithExec(planOpts)
	if apply {
		authenticatedTerraform = authenticatedTerraform.WithExec([]string{"apply", "server.plan"})
	}
	return authenticatedTerraform, nil
}

// Installs CAPI into given DOCluster
func (m *CapiDaggerCi) InstallCAPI(
	ctx context.Context,
	clusterName string,
	arch string) (*Container, error) {
	kubeconfigPath := "/root/.kube/config"
	tokenCleartext, _, _, err := m.fetchPipelineCreds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting pipeline creds: %s", err)
	}
	getKubeconfig := dag.Container().
		From("digitalocean/doctl:1.105.0").
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", tokenCleartext).
		WithExec([]string{"kubernetes", "cluster", "kubeconfig", "save", clusterName})
	if err != nil {
		return nil, fmt.Errorf("failed getting cluster config: %s", err)
	}
	capiCreds := base64.StdEncoding.EncodeToString([]byte(tokenCleartext))
	return dag.Container().
		From("alpine:latest").
		WithFile(kubeconfigPath, getKubeconfig.File(kubeconfigPath)).
		WithEnvVariable("DIGITALOCEAN_ACCESS_TOKEN", tokenCleartext).
		WithEnvVariable("DO_B64ENCODED_CREDENTIALS", capiCreds).
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"curl", "-L", "https://github.com/digitalocean/doctl/releases/download/v1.105.0/doctl-1.105.0-linux-amd64.tar.gz", "-o", "doctl-1.105.0-linux-amd64.tar.gz"}).
		WithExec([]string{"tar", "xf", "doctl-1.105.0-linux-amd64.tar.gz"}).
		WithExec([]string{"mv", "doctl", "/root/.kube/doctl"}).
		WithExec([]string{"curl", "-L", "https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.6.3/clusterctl-linux-amd64", "-o", "/usr/bin/clusterctl"}).
		WithExec([]string{"chmod", "+x", "/usr/bin/clusterctl", "/root/.kube/doctl"}).
		WithExec([]string{"clusterctl", "init", "--kubeconfig", kubeconfigPath, "--infrastructure", "digitalocean"}), nil
}
