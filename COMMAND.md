### Application

```sh {"id":"01J4ZXDSM6K2ANMY450KGJGHK6","interactive":"false"}
go mod tidy && go mod vendor
```

```sh {"id":"01J4ZXP1TTGW63J4WMJK259N16"}
go run cmd/main.go
```

### Api Documentation

```sh {"id":"01J4ZXFVKMVZ9PBAMCTWM0BE20","interactive":"false"}
# swagger
swag init -o ../internal/swagger --parseDependency --dir ./,../internal/controller
```

### Migration

- `export GOOSE_DRIVER=postgres`
- `export GOOSE_DBSTRING="user=postgres dbname=bamis sslmode=disable password=localpassword host=localhost"`

```sh {"id":"01J4ZY95T42KCEZJXHX8NJY6AQ","interactive":"false"}
cd pkg/schema
goose up
```

```sh {"id":"01J4ZY9NHCAN0AGX7B952C7CNM","interactive":"false"}
cd pkg/schema
goose down
```

```sh {"id":"01J4ZYVR5X9HR4QFTKVN15DDCH"}
cd pkg/schema
goose status
```

### Git

```sh {"id":"01J4ZZXFJAQNH16SXBFTCTK6N2","interactive":"false"}
# Remove deleted branches
git fetch origin --prune
git branch
```