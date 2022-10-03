module platform/services/storage

go 1.14

require (
	github.com/AppsFlyer/go-sundheit v0.2.0
	github.com/go-openapi/errors v0.19.8
	github.com/go-openapi/loads v0.19.6
	github.com/go-openapi/runtime v0.19.24
	github.com/go-openapi/spec v0.19.14
	github.com/go-openapi/strfmt v0.19.11
	github.com/go-openapi/swag v0.19.12
	github.com/go-openapi/validate v0.19.14
	github.com/jessevdk/go-flags v1.4.0
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/mitchellh/mapstructure v1.3.3
	github.com/patrickmn/go-cache v2.1.0+incompatible
	gitlab.com/youwol/platform/libs/go-libs v0.0.0-20201119104045-97be0e6b6e9e
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
)

replace gitlab.com/youwol/platform/libs/go-libs => ../libs
