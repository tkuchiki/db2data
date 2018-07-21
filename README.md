# db2data

## Usage

```console
$ ./db2data --help
usage: db2data --dbname=DBNAME --query=QUERY [<flags>]

Database dump to json / yaml

Flags:
  --help                   Show context-sensitive help (also try --help-long and --help-man).
  --dbuser="root"          Database user
  --dbpass=DBPASS          Database password
  --dbhost="localhost"     Database host
  --dbport=3306            Database port
  --dbsock=DBSOCK          Database socket
  --dbname=DBNAME          Database name
  --query=QUERY            SQL
  --pkey=PKEY              Primary key
  --pkey-type=string       Primary key type [int, float, string]
  --output=json            Output file format [json, yaml]
  --rows-index=ROWS-INDEX  Rows index [int=0, float=0.0, string=rows]
  --version                Show application version.
```

### JSON

```console
$ ./db2data --dbname isubata --query "SELECT id, mime, name FROM image WHERE id = 1 OR id =2" --pkey id --output json | jq .
{
  "1": {
    "id": 1,
    "mime": "image/jpeg",
    "name": "default.png"
  },
  "2": {
    "id": 2,
    "mime": "image/jpeg",
    "name": "1ce0c4ff504f19f267e877a9e244d60ac0bf1a41.png"
  },
  "count": 2,
  "default_rows": 2
}
```

### YAML

```console
$ ./db2data --dbname isubata --query "SELECT id, mime, name FROM image WHERE id = 1 OR id =2" --pkey id --output yaml
"1":
  id: 1
  mime: image/jpeg
  name: default.png
"2":
  id: 2
  mime: image/jpeg
  name: 1ce0c4ff504f19f267e877a9e244d60ac0bf1a41.png
count: 2
default_rows: 2

$ ./db2data --dbname isubata --query "SELECT id, mime, name FROM image WHERE id = 1 OR id =2" --pkey id --pkey-type int --output yaml
1:
  id: 1
  mime: image/jpeg
  name: default.png
2:
  id: 2
  mime: image/jpeg
  name: 1ce0c4ff504f19f267e877a9e244d60ac0bf1a41.png
count: 2
default_rows: 2
```
