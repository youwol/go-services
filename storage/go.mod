module platform/services/storage

go 1.14

require (
	github.com/AppsFlyer/go-sundheit v0.2.0
	github.com/go-openapi/errors v0.20.3
	github.com/go-openapi/loads v0.21.2
	github.com/go-openapi/runtime v0.24.1
	github.com/go-openapi/spec v0.20.7
	github.com/go-openapi/strfmt v0.21.3
	github.com/go-openapi/swag v0.22.3
	github.com/go-openapi/validate v0.22.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/minio/minio-go/v7 v7.0.39
	github.com/mitchellh/mapstructure v1.4.3
	github.com/patrickmn/go-cache v2.1.0+incompatible
	gitlab.com/youwol/platform/libs/go-libs v0.0.0-20201119104045-97be0e6b6e9e
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
)

replace gitlab.com/youwol/platform/libs/go-libs => ../libs
