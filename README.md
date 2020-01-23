# Lana Basket

> A simple checkout service.

_Lana basket_ uses [Long Ben](https://github.com/orov-io/lbasket) witch in turn uses [BlackBart](https://github.com/orov-io/BlackBart) witch wraps the [GIN](https://github.com/gin-gonic/gin) framework. Use it as a guide to enforce best practices when you build JSON HTTP servers.

## Considerations

* Please, read entire README.md in order to know how to run this service. If you have any question or trouble, please, contact me at hi@orov.io

* About testing:

  * Done for this challenge
    * Unit testing for [checkout](./packages/checkout/) package.
    * Unit testing for [client](./client) package.
    * Functional testing is provided in a indirect way by client packages.
    * Behavioral testing for new and get endpoints. As the logic is the same as the client tests, we only provide some test as a demonstration.

  * Not implemented
    * Unit testing for [lanabadger](./packages/lanabadger/) is straightforward, so is not added.
  
  * Implemented but not shown in coverage tools
    * [Service](./service) package has not unit test, but client test runs against the endpoint and code in this package, so running client test we can check the service package integrity.

* I usually use my own framework on top on some strong utilities, as induced in the first paragraph. This is done in order to achieve a DRY approach. But I know that this approach can take turns in coupling problems in a microservice mesh if we're not careful.
Furthermore,for this project I have added to the framework the possibility of use an internal database (key/value), which I have not needed so far.

* Given the probability of collision of uuids (see [uuid collisions](https://en.wikipedia.org/wiki/Universally_unique_identifier#Collisions) on wikipedia), it is not checked that it already exists to avoid overwriting.

* I use the JSON format to store baskets. I know that an approach with the encoding/base64 and encoding/gob will be better at a performance level, but JSON performance is enough for this toy example. Also, JSON presents a more familiar interface for many developers and we don't know who will touch this code in the future.

* In order to have integrity, a list of available products is added and exposed on `GET $BASE_PATH/$SERVICE_VERSION/$SERVICE_BASE_PATH/products`. I know that it is not a requirement, but our users deserve an enjoyable experience. As is not a requirement to add new products to the list, only basic functionality is provided.

* About the "Dealing with money as integers": Internally (the server) we deal with money as integer, as I guess it's business logic.
But keep on mind that we expose to the client a human readable representation of the money value. This is done so as not to expand such business logic to the client.

* Discounts logic is in code. We are aware that the discounts could have been created in the database and applied dynamically; but given the scope of the challenge it is left in the code. In case our product have a great impact on our customers and they demand more and more products with more discounts, We should consider that the discounts live in the database and are associated to the products dynamically.To do this, we will need to create a discountManager that deals with this logic.

* As it is possible to add an external DB in the future, project is ready to add a __postgres__ DB. _[docker-compose](./docker-compose.yaml)_ file adds a postgres container to provide local development if necessary.

* About thread-safety, we use the [badger](https://github.com/dgraph-io/badger) package as internal db, which provides thread-safety by default. This could also have been achieved by applying __[mutexes](https://gobyexample.com/mutexes)__ on top of a `map[string]interface{}` map.
However, given the reputation of the above-mentioned package, this option is chosen as it is the most tested.

## Quick start

Provide a _.env_ file variables with all variables founds in _[example.env](./example.env)_. Load this _.env_ file to your environment variables.

Be sure that you have all dependencies by running:

```bash
go mod tidy
```

Run the server with below command:

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

You can shutdown the service with:

```bash
make down
```

## Dependencies

As internal dependencies, this module relies in some internal dependencies:

* [BlackBeard](https://github.com/orov-io/BlackBeard), the client utility.
* [BlackBart](https://github.com/orov-io/BlackBart), the server utility.

Also, the intensive use of go modules force us to need go1.13.

For testing, we use below libraries:

* [godog](https://github.com/DATA-DOG/godog) for cucumber based behavioral testing.
* [Convey](https://github.com/smartystreets/goconvey/convey) for unit and functional testing.

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
* ENV (set it to local)
* DATABASE_HOST
* DATABASE_PASSWORD
* DATABASE_USER
* DATABASE_SSL_MODE
* SERVICE_DATABASE_NAME
* DATABASE_MIGRATIONS_DIR
* SERVICE_NAME
* SERVICE_VERSION
* GOPRIVATE (set it to: github.com/orov-io/*)
* BASE_PATH (set it to: http://localhost )
* SERVICE_DESCRIPTION ( optional )
* SERVICE_VERSION (set it to v1)
* SERVICE_BASE_PATH
* ENABLE_BADGER

Database env variables only need to be provided if you are using a POSTGRES database; and, for now, this is not the circumstance.

As the _[example.env](./example.env)_ is provided, you can simply copy it to __.env__ before run the project and then execute:

```bash
source .env
```

## Running tests

Please, be sure that test file has access to your env variables (i.e.: vscode do not do this for you in lot of cases).

To run cucumber test you must first install the [godog](https://github.com/DATA-DOG/godog) binary file:

```bash
go get github.com/DATA-DOG/godog/cmd/godog
```

## About the API

You can find the API specification on the *[open-api_example](open-api_example.yaml)* file.
