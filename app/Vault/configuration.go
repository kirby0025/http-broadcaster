// Package vault provides functions to retrieve info from Hashicorp Vault
package vault

import (
	"context"
	"log"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

// GetVarnishListFromVault retrieve the list of varnish servers to send PURGE to.
// It uses the AppRole authentication method.
// APPROLE_ROLE_ID and APPROLE_SECRET_ID are fed from environment variables.
func GetVarnishListFromVault() []string {
	config := vault.DefaultConfig()
	var value []string

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatal("unable to initialize Vault client: %w", err)
		return value
	}

	roleID := os.Getenv("APPROLE_ROLE_ID")
	if roleID == "" {
		log.Fatal("no role ID was provided in APPROLE_ROLE_ID env var")
		return value
	}
	secretID := &auth.SecretID{FromEnv: "APPROLE_SECRET_ID"}

	appRoleAuth, err := auth.NewAppRoleAuth(
		roleID,
		secretID,
	)
	if err != nil {
		log.Fatal("unable to initialize AppRole auth method: %w", err)
		return value
	}

	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		log.Fatal("unable to login to AppRole auth method: %w", err)
		return value
	}
	if authInfo == nil {
		log.Fatal("no auth info returned after login")
		return value
	}

	secret, err := client.KVv2("app").Get(context.Background(), "http-broadcaster/stg/varnish_list")
	if err != nil {
		log.Fatal("unable to read secret: %w", err)
		return value
	}

	// selecting list key from retrieved secret
	list, ok := secret.Data["list"].(string)
	if !ok {
		log.Fatal("value type assertion failed: %T %#v", secret.Data["list"], secret.Data["list"])
		return value
	}
	value = strings.Split(string(list), ",")
	return value
}
