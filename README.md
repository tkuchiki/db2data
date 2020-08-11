# db2data

## Usage

```console
$ ./db2data --help
usage: db2data --dbname=DBNAME --query=QUERY [<flags>]

Database dump to json / yaml

Flags:
      --help                Show context-sensitive help (also try --help-long and --help-man).
      --dbuser="root"       Database user
      --dbpass=DBPASS       Database password
      --dbhost="localhost"  Database host
      --dbport=3306         Database port
      --dbsock=DBSOCK       Database socket
      --dbname=DBNAME       Database name
      --query=QUERY         SQL
      --pkey=PKEY           Primary key
      --pkey-type=string    Primary key type [int, float, string]
      --composite-key=COMPOSITE-KEY
                            Composite key(comma separated)
  -d, --delimiter="-"       Delimiter
      --output=json         Output file format [json, yaml]
      --types=TYPES         Set types (column:type,...) type=[number, string, bool]
      --version             Show application version.
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
  }
}
```

```console
$ ./db2data --dbpass=password --dbport=13306 --dbhost=127.0.0.1 --dbname=test --composite-key="id,name" --query 'SELECT * FROM users' | jq .
{
  "1-alice": {
    "id": 1,
    "name": "alice"
  },
  "2-bob": {
    "id": 2,
    "name": "bob"
  }
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

$ ./db2data --dbname isubata --query "SELECT id, mime, name FROM image WHERE id = 1 OR id =2" --pkey id --pkey-type int --output yaml
1:
  id: 1
  mime: image/jpeg
  name: default.png
2:
  id: 2
  mime: image/jpeg
  name: 1ce0c4ff504f19f267e877a9e244d60ac0bf1a41.png
```
