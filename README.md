# http-broadcaster

## Description
Un démon simple écrit en Go qui prend une requête PURGE en entrée et la transmet à plusieurs serveurs varnish.

## Déploiement
Le projet se déploie via les pipelines Gitlab en stg et en prod, via un déclenchement manuel.
Le sidecar vault va déposer un fichier (/vault/secrets/.env) contenant les variables d'environnements.

## Configuration
Arguments de lancement :
* -l (--log) : emplacement du fichier de log. (Default : /app/http-broadcaster.log)
* -e (--envfile) : emplacement du fichier de variables d'environnement. (Default : /vault/secrets/.env)
* --metrics : active l'exposition des métriques prometheus sur /metrics. (Default : false)

### Variables d'environement

La liste de serveurs Varnish peut être fournie directement dans un fichier d'env :
* VARNISH_SERVERS: list of varnish backend servers. Ex "http://10.13.32.1:6081,http://10.13.32.2:6081"
* CLIENT_LIST : list of IPs in CIDR format (0.0.0.0/32) of authorized client to do purge/ban requests.

## Fonctionnalites

* Génère la liste des serveurs Varnish à partir des variables d'environnement.
* Ecoute sur le port 6081.
* Healthcheck disponible sur l'uri /healthcheck pour vérifier son bon fonctionnement. Renvoie un code HTTP 200 et le message "OK".
* Metriques Prometheus disponible (désactivée par défaut).
* Traite les requêtes entrantes en récupérant 3 éléments et en les intégrant à la requête transmise aux serveurs Varnish :
- La méthode (PURGE ou BAN par exemple)
- L'url : / pour BAN, /codes/api/greffes/0101 par exemple pour PURGE.
- Le header X-Cache-Tags : dans le cas d'un BAN ce header contient une valeur.

## Usage
Les interactions se font via le protocol HTTP. Les applications où les utilisateurs envoient une requête de méthode PURGE vers le démon.
Une fois le traitement d'une requête effectuée, le démon renvoie 200 si tout est ok, 405 dans le cas contraire.
