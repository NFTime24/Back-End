module deukyunlee/protocol-camp

go 1.18

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/joho/godotenv v1.4.0
	github.com/swaggo/swag v1.8.4
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/tools v0.1.11 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace deukyunlee/protocol-camp => ./pkg/httpHandlers
