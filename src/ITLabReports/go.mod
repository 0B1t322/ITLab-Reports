module github.com/RTUITLab/ITLab-Reports

go 1.18

require (
	github.com/MicahParks/keyfunc v1.1.0
	github.com/clarketm/json v1.17.1
	github.com/go-kit/kit v0.12.0
	github.com/golang-jwt/jwt/v4 v4.4.1
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.4.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.8.0
	github.com/swaggo/http-swagger v1.2.6
	github.com/swaggo/swag v1.8.1
	github.com/xakep666/mongo-migrate v0.2.1
	go.mongodb.org/mongo-driver v1.9.1
	golang.org/x/exp v0.0.0-20220414153411-bcd21879b8fd
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.8.1 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.5 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samber/do v1.4.1 // indirect
	github.com/swaggo/files v0.0.0-20220728132757-551d4a08d97a // indirect
	github.com/swaggo/gin-swagger v1.5.3 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.12 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/genproto v0.0.0-20210917145530-b395a37504d4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require github.com/0B1t322/QueryParser v0.1.0

// Replace for proto package

require (
	github.com/0B1t322/MongoBuilder v0.1.8
	github.com/0B1t322/RepoGen v0.0.4
	github.com/RTUITLab/ITLab v1.0.0
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/samber/lo v1.33.0
	github.com/samber/mo v1.5.1
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
)

replace github.com/RTUITLab/ITLab v1.0.0 => ./pkg/ITLab
