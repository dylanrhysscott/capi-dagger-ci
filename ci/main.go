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
	"fmt"
)

type CapiDaggerCi struct{}

// Runs terraform init, plan and apply to deploy infrastructure
func (m *CapiDaggerCi) DeployInfra(ctx context.Context, path Directory, token *Secret, spacesAccessKey *Secret, spacesSecretKey *Secret) (string, error) {
	tokenCleartext, err := token.Plaintext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed getting token: %s", err)
	}
	spacesAccessKeyCleartext, err := spacesAccessKey.Plaintext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed getting spaces access key: %s", err)
	}
	spacesSecretKeyCleartext, err := spacesSecretKey.Plaintext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed getting spaces secret key: %s", err)
	}
	out, err := dag.Container().
		From("hashicorp/terraform:latest").
		WithDirectory("/infra", &path).
		WithWorkdir("/infra").
		WithEnvVariable("AWS_ACCESS_KEY_ID", spacesAccessKeyCleartext).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", spacesSecretKeyCleartext).
		WithEnvVariable("DIGITALOCEAN_TOKEN", tokenCleartext).
		WithExec([]string{"init"}).
		WithExec([]string{"plan", "-out", "server.plan"}).
		WithExec([]string{"apply", "server.plan"}).
		Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed deploying infra: %s", err)
	}
	return out, nil
}
