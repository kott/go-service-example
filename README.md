# go-service-example
Outline of a simple HTTP service in Go.

This project will be using the [gin](https://github.com/gin-gonic/gin) framework for the bones of a REST API. The purpose is to allow this to be used as a template such that the groundwork is layed out when creating a new service.

## Getting Started
Not much for now, but might become more involved as we go along.
* Install dependencies: `go get -u ./...`
* Run the service: `docker-compose up --build`

## Configuration
This service is setup to use `<env_name>.env` files to load its configuration (from both `env_file` in docker 
compose & directly). These files can be passed into docker, depending on which environment you are deploying to.
 
The environment variable `SERVICES_PROFILE` should be set to `docker` (e.g. `SERVICES_PROFILE=docker`) if the service is
running within a container. The reason for this is because this tells the service to grab the environment variables as 
its configuration. This would allow someone to change the config by injecting new values for these variables, which is sometimes
desired over a redeploy with the changed `.env` file.

In a local instance we could set `SERVICES_PROFILE=local`  and it will read the `local.env` file for the configuration. 

_Note: this configuration is somewhat opinionated, but can be easily changed._ 

## Database
The docker-compose file defines a database this service can use for local development. When we need to create a
migration, we can do this using the [golang-migrate CLI tool](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate). An example from their documentation: 
`migrate create -ext sql -dir db/migrations -seq create_article_table` would create the up/down migration files for the
article table and adhere to the naming convention of the files. To run the migrations (up/down) can be done with:
`migrate -source file://db/migrations -database postgres://gouser@localhost:5432/example?sslmode=disable up`.

## TODO
- o11y (i.e. request tracing / monitoring)
- Deployment (likely AWS fargate)
- GitHub actions for linting, testing, deploying
- Authentication
- Asynchronous tasks management
