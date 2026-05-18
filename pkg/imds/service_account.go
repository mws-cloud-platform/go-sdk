package imds

import (
	"context"
	"encoding/json"
	"maps"
	"slices"

	"go.mws.cloud/util-toolset/pkg/os/env"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

const (
	ErrServiceAccountNotFound = consterr.Error("service account not found for this Compute virtual machine")

	serviceAccountsKey = "instance/service-accounts/?recursive=true"
)

// GetVMServiceAccount returns a service account attached to the current VM.
//
// It returns an error when running outside a VM
// or when no service account is attached.
func GetVMServiceAccount(ctx context.Context, client HTTPClient, env env.Env) (string, error) {
	serviceAccountsString, err := NewClient(client, env).GetWithContext(ctx, serviceAccountsKey)
	if err != nil {
		return "", err
	}

	var serviceAccounts map[string]any
	if err = json.Unmarshal([]byte(serviceAccountsString), &serviceAccounts); err != nil {
		return "", err
	}
	if len(serviceAccounts) == 0 {
		return "", ErrServiceAccountNotFound
	}
	return slices.Min(slices.Collect(maps.Keys(serviceAccounts))), nil
}
