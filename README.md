# http-broadcaster

## Description
Un démon simple écrit en Go qui prend une requête PURGE en entrée et la transmet à plusieurs serveurs varnish.

## Installation
* Compiler le programme
```
cd app/
go get ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ../build/http-broadcaster main.go
```
* Déposer la liste des serveurs varnish à côté du binaire, avec ce format :
```
http://10.13.32.1:6081,http://10.13.32.2:6081
```

## Fonctionnalites

* Génère la liste des serveurs Varnish en lisant le fichier "varnish" présent à côté du binaire.
* Ecoute sur le port 6081.
* Healthcheck disponible sur l'uri /healthcheck pour vérifier son bon fonctionnement. Renvoie un code HTTP 200 et le message "OK".
* Traite toutes les requêtes entrantes comme l'url à purger dans varnish. Par exemple un appel sur http://10.13.101.11:6081/hello/test entrainera une purge de l'uri "/hello/test" sur les serveurs Varnish.

## Usage
Les interactions se font via le protocol HTTP. Les applications où les utilisateurs envoient une requête de méthode PURGE vers le démon.
Une fois le traitement d'une requête effectuée, le démon renvoie 200 si tout est ok, 405 dans le cas contraire.

## Roadmap
* Aller chercher la liste des varnish dans vault.
* Ajouter une forme d'authentification.
* Ajouter d'autres possibilités que l'envoi à Varnish.
