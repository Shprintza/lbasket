# Lana Basket

> A simple checkout service.

_Lana basket_ uses [Long Ben](https://github.com/orov-io/lbasket) witch in turn uses [BlackBart](https://github.com/orov-io/BlackBart) witch wraps the [GIN](https://github.com/gin-gonic/gin) framework. Use it as a guide to enforce best practices when you build JSON HTTP servers.

## Considerations

* A long time ago, I learned not to repeat myself, so Im using my own framework on top on some strong utilities, as induced by the previous paragraph.
Furthermore, for this project I have added to the framework the possibility of use an internal database (key/value), which I have not needed so far.

* Given the probability of collision of a uuid, it is not checked that it already exists to avoid overwriting. If it is needed (see [uuid collisions](https://en.wikipedia.org/wiki/Universally_unique_identifier#Collisions) on wikipedia), please say it to me.

* I use the JSON format to store baskets. I know that an approach with the encoding/base64 and encoding/gob will be better at a performance level, but JSON performance is enough for this toy example.

## Quick start

Provide a _.env_ file variables with all variables founds in _[example.env](./example.env)_. Load this _.env_ file to your environment variables.

Assert that port $PORT is available and run below command:

```Bash
make up logs
```

This command will build the docker images and run the _Lana Basket_ server.

You can test the service pinging it:

```Bash
curl --request GET \
  --url http://localhost:8080/v1/${SERVICE_BASE_PATH}/ping

> {"status":"OK","message":"pong"}
```

## Making your own service

As _Long Ben_ is near to be an out-the-box service, you will need to change some variables.

Some imports on _[main.go](./main.go)_, _[responses.go](./service/responses.go)_, _[handlers.go](./service/handlers.go)_, _[Dockerfile](./Dockerfile)_,  _[client.go](./client/client.go)_ and _[client_test.go](./client/client_test)_ has hardcoded references to lbasket module. Please, change this import sentences in order to make your own service.

Also, you will need to change the module name in the _[go.mod](./go.mod)_ file.

### Where to place my files

Service related go files should be organizing as below:

* Handlers functions => _[handlers.go](./service/handlers.go)_
* Gin routes and groups => _[routes.go](./service/routes.go)_
* Responses functions => _[responses.go](./service/responses.go)_
* Packages providing functionality: _[packages](./packages/README.md)_
* Database migrations => _[migrations](./migrations/README.md)_
* Request and response structs => _[models](./models)_
* Service go client => _[client.go](./client/client.go)_. This is also a good please to provide end-to-end tests, as you can see at _[client_test.go](./client/client_test.go)_.

## Dependencies

As internal dependencies, this module relies in some internal dependencies:

* [BlackBeard](https://github.com/orov-io/BlackBeard), the client utility.
* [BlackBart](https://github.com/orov-io/BlackBart), the server utility.

Also, the intensive use of go modules force us to need go1.13.

## ENV VARIABLES

This backbone relies in some env variables to enable needed modules and be deployed. We can discriminate between __built time__ and __run time__ variables:

### Built time variables

These variables are used in  built time and are only needed on docker build time or on deployed time. Please, be sure that these variables are available when you tried to deploy/stand up your service:

* PORT (only local): Internal port to serve. Used in docker-compose.
* DATABASE_USER & DATABASE_PASSWORD: Used by docker-compose to set the database container.
* SERVICE_NAME: Used both in docker-compose and GCloud&Pipelines deployment.
* SERVICE_DESCRIPTION, SERVICE_VERSION & SERVICE_BASE_PATH: Used to deploy the google endpoints gateway configuration. Gae path will be /${SERVICE_VERSION}/${SERVICE_BASE_PATH}
* GOOGLE_APPLICATION_CREDENTIALS: Used in docker-compose to gain access to your buckets. This is a path to file contains your IAM json file.
* GCLOUD_API_KEYFILE: Used in pipelines to gain access to your gcloud project. Set it with a base64 encoded version of your IAM json file.
* ref: Set it to __$ref__ to apply correctly _envsub_ to your _openapi-appengine-example.yaml_ file.
* NETRC: a base64 encoded file with your bitbucket access token.
* GCLOUD_PROJECT
* INSTANCE_CONNECTION_NAME: Needed for the app.yaml to provide a socket to your database.
* BUCKET_CREDENTIALS: Used on test to gain access to bucket storage. Use it if your test need this permission. Set it with a base64 encoded version of your IAM json file.
* SONAR_CLOUD: Use on test to send the result to sonarcloud.

### Run time variables

Please, see documentation of the [BlackBart](https://github.com/orov-io/BlackBart) to find variables that enables your server capabilities.

* PORT (only local)
* DATABASE_HOST
* DATABASE_PASSWORD
* DATABASE_USER
* DATABASE_SSL_MODE
* SERVICE_DATABASE_NAME
* DATABASE_MIGRATIONS_DIR
* GOOGLE_APPLICATION_CREDENTIALS
* GCLOUD_STORAGE_BUCKET
* FIREBASE_BUCKET
* FIREBASE_BUCKET_FILE_NAME
* SERVICE_NAME
* SERVICE_VERSION
* REDIS_ADDRESS
* REDIS_PASSWORD
