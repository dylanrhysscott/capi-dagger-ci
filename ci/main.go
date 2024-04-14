// A module for provisioning DOKS clusters and installing CAPI

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"
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
	planOpts := []string{"plan"}
	applyOpts := []string{"apply", "-auto-approve"}
	if destroy {
		planOpts = append(planOpts, "--destroy")
		applyOpts = append(applyOpts, "--destroy")
	}
	authenticatedTerraform := dag.Container().
		From("hashicorp/terraform:latest").
		WithDirectory("/infra", &path).
		WithWorkdir("/infra").
		WithEnvVariable("AWS_ACCESS_KEY_ID", spacesAccessKeyCleartext).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", spacesSecretKeyCleartext).
		WithEnvVariable("DIGITALOCEAN_TOKEN", tokenCleartext).
		WithEnvVariable("CACHEBUSTER", time.Now().String()). // https://archive.docs.dagger.io/0.9/cookbook/#invalidate-cache - cache invalidation not yet supported :(
		WithExec([]string{"init"}).
		WithExec(planOpts)
	if apply {
		authenticatedTerraform = authenticatedTerraform.WithExec(applyOpts)
	}
	return authenticatedTerraform, nil
}

// Installs CAPI into given DOCluster
func (m *CapiDaggerCi) InstallCAPI(ctx context.Context, clusterName string) (*Container, error) {
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
		WithEnvVariable("CACHEBUSTER", time.Now().String()). // https://archive.docs.dagger.io/0.9/cookbook/#invalidate-cache - cache invalidation not yet supported :(
		WithExec([]string{"clusterctl", "init", "--kubeconfig", kubeconfigPath, "--infrastructure", "digitalocean"}), nil
}
