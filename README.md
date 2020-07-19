# Go Application Template

## Directory structure
- transport: All external communications (CLI, http services and server, GRPC, ...)
- db: database related functionality (connection, migrations, seeders, ...)
- monitoring: Logger and prometheus metrics
- models: Application entities (sqlboiler generated models or handwritten structs)
- config: application configuration

## Philosophy:
This template is a bundle of the libraries and structure that I use in almost all projects, so it's tailored to my needs, 
it's not a general purpose template or framework for everyone so don't expect to find everything you need.

## Used libraries:
- SQLBoiler: for generating an ORM using a database-first approach. (recommended only for applications with a lot of tables)
- go-migrate: handling database migrations.
- Cobra: command line functionality
- Echo: Http web server
- logrus: logging tool (planning to switch to zap)
- Viper: configuration

## CLI Built in commands
### Migration
- new
- up
- down
### sqlboiler
- generate: generate models using sqlboiler based on configured sql database.
### swag
- generate: generate  models for http handlers based on swagger API spec
### GRPC:
- generate: generate GRPC go files based on given proto file
### Dependencies
- golang-migrate command line utility: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate


