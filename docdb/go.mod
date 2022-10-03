module platform/services/docdb

go 1.15

require (
	github.com/AppsFlyer/go-sundheit v0.2.0
	github.com/btubbs/datetime v0.1.1
	github.com/go-openapi/errors v0.19.9
	github.com/go-openapi/loads v0.20.0
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/spec v0.20.1
	github.com/go-openapi/strfmt v0.20.0
	github.com/go-openapi/swag v0.19.13
	github.com/go-openapi/validate v0.20.1
	github.com/gocql/gocql v0.0.0-20201215165327-e49edf966d90
	github.com/jessevdk/go-flags v1.4.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/scylladb/gocqlx/v2 v2.3.0
	gitlab.com/youwol/platform/libs/go-libs v0.0.0-20210120092013-386bda2eb61a
	go.uber.org/zap v1.16.0
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	golang.org/x/tools v0.1.0 // indirect
	gopkg.in/inf.v0 v0.9.1
)

replace gitlab.com/youwol/platform/libs/go-libs => ../libs

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.4.3
