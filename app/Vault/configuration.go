// Package vault provides functions to retrieve info from Hashicorp Vault
package vault

import (
	"context"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

// InitVaultConnection builds a client connection to vault and return it.
// It uses the AppRole authentication method.
// APPROLE_ROLEID and APPROLE_SECRETID are fed from environment variables.
func InitVaultConnection() *vault.Client {
	config := vault.DefaultConfig()
	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatal("unable to initialize Vault client: %w", err)
		return client
	}
	roleID := os.Getenv("APPROLE_ROLEID")
	if roleID == "" {
		log.Fatal("no role ID was provided in APPROLE_ROLEID env var")
		return client
	}
	secretID := &auth.SecretID{FromEnv: "APPROLE_SECRETID"}

	appRoleAuth, err := auth.NewAppRoleAuth(
		roleID,
		secretID,
	)
	if err != nil {
		log.Fatal("unable to initialize AppRole auth method: %w", err)
		return client
	}
	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		log.Fatal("unable to login to AppRole auth method: %w", err)
		return client
	}
	if authInfo == nil {
		log.Fatal("no auth info returned after login")
		return client
	}
	return client
}
