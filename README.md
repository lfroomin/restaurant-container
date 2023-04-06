restaurant-container is a microservice using a REST API
deployed in a Docker container.

The API is documented using the OAS3 (Swagger) specification,
and the model is generated from the specification. The
specification is in the internal/model/restaurant-api.yaml file.

The basic CRUD endpoints exist for the restaurant entity.
- Create - create a restaurant
- Read - get a restaurant
- Update - update a restaurant
- Delete - delete a restaurant

When a restaurant is created or updated, if it contains
an address, the address is used to look up the geocode
coordinates of the address (lat, lon).

The frameworks/packages/services used:
- gin
- viper
- Dynamo DB
- Location (used for geocoding)

The Dynamo DB database is the same that is created in the
restaurant-serverless project SAM template.

To update the generated model when the OAS3 specification is
changed, do the following:
- From the internal/model folder execute `go generate`

**Unit Tests**
- From the project root folder execute `go test ./...`
