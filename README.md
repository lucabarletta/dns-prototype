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
curl --location 'http://localhost:9090/domain.i2p'
```

```json
{
  // todo: add response
}
```

This dns entry is hardcoded and available after startup

---

#### PUT example

```curl
curl --location \
  --request PUT 'http://localhost:9090/domain.i2p/m6elnkiizogiz5wq4perd7aslir5rdu7jmtwxlxua5aofa43zyva'
```

```json
{
    "ident": "543E2jGQ"
}
```

## TODO

- Docker compose
- documentation
