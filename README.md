# http-broadcaster

## Description
Un démon simple écrit en Go qui prend une requête PURGE ou BAN en entrée et la transmet à plusieurs serveurs varnish.

## Installation
* Compiler le programme
```
cd app/
go get ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ../build/http-broadcaster main.go
```

## Configuration

### Variables d'environement
* APPROLE_ROLEID : RoleID utilisé pour se connecter à Vault.
* APPROLE_SECRETID : SecretID utilisé pour se connecter à Vault.
* LIST_METHOD: "vault" ou "file". Définit la méthode de construction de la liste des serveurs varnish.
* VAULT_ADDR : Adresse du serveur Vault.
* VAULT_DATABASE : database where informations are located in Vault.
* VAULT_PATH: path to the secrets containing the informations.
* VARNISH_SERVERS : list of varnish servers.
```
http://10.13.32.1:6081,http://10.13.32.2:6081
```

## Fonctionnalites

* Génère la liste des serveurs Varnish en lisant le fichier "varnish" présent à côté du binaire ou en allant les chercher dans Vault.
* Ecoute sur le port 6081.
* Healthcheck disponible sur l'uri /healthcheck pour vérifier son bon fonctionnement. Renvoie un code HTTP 200 et le message "OK".
* A l'arrivée d'une requête, récupère la méthode, l'url et le contenu du header X-Cache-Tags (facultatif) et les envoie aux serveurs Varnish.

## Usage
Les interactions se font via le protocol HTTP. Les applications où les utilisateurs envoient une requête de méthode PURGE vers le démon.
Une fois le traitement d'une requête effectuée, le démon renvoie 200 si tout est ok, 405 dans le cas contraire.

## Roadmap
* Ajouter une forme d'authentification.
* Ajouter d'autres possibilités que l'envoi à Varnish.
