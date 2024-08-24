module github.com/onflow/flow-playground-api

go 1.20

require (
	github.com/99designs/gqlgen v0.17.5
	github.com/Masterminds/semver v1.5.0
	github.com/TV4/logrus-stackdriver-formatter v0.1.0
	github.com/alecthomas/chroma v0.8.1
	github.com/getsentry/sentry-go v0.18.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/httplog v0.2.5
	github.com/go-chi/render v1.0.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/sessions v1.2.0
	github.com/gorilla/websocket v1.5.0
	github.com/hashicorp/golang-lru v1.0.2
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/onflow/cadence v1.0.0-preview.32
	github.com/onflow/flow-emulator v1.0.0-preview.31
	github.com/onflow/flow-go v0.35.10-crescendo-preview.25.0.20240604172940-c504b454e576
	github.com/onflow/flow-go-sdk v1.0.0-preview.34
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.18.0
	github.com/rs/cors v1.8.0
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.9.0
	github.com/vektah/gqlparser/v2 v2.4.2
	gorm.io/driver/postgres v1.3.10
	gorm.io/driver/sqlite v1.3.6
	gorm.io/gorm v1.23.9
	gotest.tools v2.2.0+incompatible
)


replace github.com/SaveTheRbtz/mph => github.com/SaveTheRbtz/mph v0.1.1-0.20240117162131-4166ec7869bc
