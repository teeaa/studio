## Dance studio

### Run with Docker
Start Docker containers by running `docker-compose up -d --build`.
This will create a MySQL server with the correct (empty) database and user for it.
It will also start the dance studio server which will do the initial migration of creating `classes` and `bookings` tables to the database.

The MySQL instance will be running in port localhost:13306 (3306 inside container network) and has user:password/database dancestudio:dancestudio/dancestudio .
The REST API will be running in port http://localhost:8080/

### Run outside Docker
Connect to your favourite MySQL instance by providing the database connection info to the REST server via environment.
DANCESTUDIO_MYSQLUSER, DANCESTUDIO_MYSQLPASSWORD, DANCESTUDIO_MYSQLPASSWORD, DANCESTUDIO_MYSQLADDRESS, DANCESTUDIO_MYSQLDB, DANCESTUDIO_MYSQLPORT

The MySQL schemas for the two tables are located in `/mysql` in project root. 

Build the REST server by running `go build github.com/teeaa/studio/cmd/server/.` in project root. This will create the executable `./server`. To run that (with env vars) run for example `DANCESTUDIO_MYSQLPORT=13306 ./server`

### JSON Payloads
`POST /classes`
For creating/updating classes:
```
{
	"name": "Class name",
	"start_date": "2019-01-01",
	"end_date": "2020-12-31",
	"capacity": 15
}
```

For creating/updating bookings:
`POST /bookings`
```
{
	"name": "Teea the Ballet dancer",
	"booking_date": "2019-07-15",
	"class_id": 1
}
```

Also available:
GET /<classes/bookings>/
GET /<classes/bookings>/<id>
PUT /<classes/bookings>/<id>
DELETE /<classes/bookings>/<id>

Restrictions: The booking date must fall inside the class start and end dates for the booked class. This is checked on creation and update.

### Tests

Run tests by running `go test ./...` in the project root. This will test the helper functions as well as all the database and route methods in both of classes and bookings modules.
The database and web server are mocked so those don't need to be running for the tests to succeed.

