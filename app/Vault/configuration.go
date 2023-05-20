package Vault

import (
    "fmt"
    "os"
    vault "github.com/hashicorp/vault/api"
    auth "github.com/hashicorp/vault/api/auth/approle"
)

// getVarnishList retrieve the list of varnish servers to send PURGE to.
// It uses the AppRole authentication method.
func getVarnishList() (string, error) {
	config := vault.DefaultConfig() // modify for more granular configuration

	client, err := vault.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("unable to initialize Vault client: %w", err)
	}

        // Get roleID and secretID from ENV vars
	roleID := os.Getenv("APPROLE_ROLE_ID")
	if roleID == "" {
		return "", fmt.Errorf("no role ID was provided in APPROLE_ROLE_ID env var")
	}
	secretID := os.Getenv("APPROLE_SECRET_ID")
	if secretID == "" {
		return "", fmt.Errorf("no secret ID was provided in APPROLE_SECRET_ID env var")
	}

	appRoleAuth, err := auth.NewAppRoleAuth(
		roleID,
		secretID,
		auth.WithWrappingToken(), // Only required if the secret ID is response-wrapped.
	)
	if err != nil {
		return "", fmt.Errorf("unable to initialize AppRole auth method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return "", fmt.Errorf("unable to login to AppRole auth method: %w", err)
	}
	if authInfo == nil {
		return "", fmt.Errorf("no auth info was returned after login")
	}

	// get secret from the default mount path for KV v2 in dev mode, "secret"
	secret, err := client.KVv2("app").Get(context.Background(), "http-broadcaster/stg/varnish_list")
	if err != nil {
		return "", fmt.Errorf("unable to read secret: %w", err)
	}

	// data map can contain more than one key-value pair,
	// in this case we're just grabbing one of them
	value, ok := secret.Data["list"].(string)
	if !ok {
		return "", fmt.Errorf("value type assertion failed: %T %#v", secret.Data["list"], secret.Data["list"])
	}

	return value, nil
}
