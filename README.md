# dns-prototype

---

### setup

1. Install Go

   - via Installer https://go.dev/doc/install
   - via Scoop `scoop bucket add main & scoop install main/go`

2. Install Go plugin on VS Code or any other compatible Editor (Jetbrains GoLand)
3. Clone Project https://github.com/lucabarletta/dns-prototype.git
4. Open project in terminal or cmd and run `go run main.go`
5. Go documentation https://go.dev/doc/effective_go

---

### examples

#### GET example

```curl
curl --request GET \
  --url http://localhost:9090/domain
```

```json
{
    "created": "2023-12-20T17:03:57.4278276+01:00",
    "domain": "domain",
    "subName": "subName",
    "name": "name",
    "type": "type",
    "record": "192.168.1.1",
    "ttl": 3600
}
```

This dns entry is hardcoded and available after startup

---

#### POST example

```curl
curl --request POST \
  --url http://localhost:9090/test123 \
  --header 'Content-Type: application/json' \
  --data '{
  "domain": "domain",
  "subName": "subName",
  "name": "name",
  "type": "type",
  "record": "192.168.1.1",
  "ttl": 3600
}'
```

```json
{
    "created": "2023-12-20T17:04:35.202372+01:00",
    "domain": "domain",
    "subName": "subName",
    "name": "name",
    "type": "type",
    "record": "192.168.1.1",
    "ttl": 3600
}
```

## TODO

- Docker compose
- documentation
