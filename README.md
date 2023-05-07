# Rentals Service

Rentals Service provides rental listings with various filtering, sorting, and pagination options. 
The API is built using Go, with the Gin framework and another popular packages.

## Features

- Get a single rental by ID
- List rentals with filtering options
  - Filter by price range
  - Filter by rental IDs
  - Filter by proximity to a given location (latitude, longitude)
- Sort rentals by price and year (sort=price|price_desc|year|year_desc)
- Paginate rental listings (limit=n, offset=n)
- Input validation for query parameters
- Logging and instrumentation decorators
- Graceful shutdown

## Prerequisites

- [Go 1.20](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Setting Up and Running the Application

1. Clone the repository:

```bash
git clone https://github.com/plar/rentals-api.git
cd rentals-api
```

2. Execute the following command in the terminal:

```bash
$ make doc-build doc-up doc-logs-follow
docker-compose build
...
db_1   | PostgreSQL init process complete; ready for start up.
app_1  | [GIN-debug] GET    /rentals/:id              --> github.com/plar/rentals-api/handler.RentalHandler.GetRentalByID-fm (3 handlers)
app_1  | [GIN-debug] GET    /rentals                  --> github.com/plar/rentals-api/handler.RentalHandler.GetRentals-fm (3 handlers)
...
```

* `doc-build`: Builds the rentals-api application
* `doc-up`: Starts a PostgreSQL instance with sample rental data (see [db/init.sql](db/init.sql)) and also launches the rentals-api HTTP service
* `doc-logs-follow`: Continuously shows logs for the development environment

The API server will start at `http://localhost:8080`.

3. To terminate the Docker environment, execute the following command in the terminal:

```bash
$ make doc-down 
docker-compose down -v --remove-orphans	
Stopping rentals-api_app_1 ... done
Stopping rentals-api_db_1  ... done
Removing rentals-api_app_1 ... done
Removing rentals-api_db_1  ... done
Removing network rentals-api_default
Removing volume rentals-api_database-data
```

## Usage

I used [HTTPie](https://httpie.io/) to test the API but you can use any other tool(curl, postman), etc that can make HTTP requests. 

### Get a single rental by ID

```bash
$ http :8080/rentals/1
HTTP/1.1 200 OK
Content-Length: 658
Content-Type: application/json; charset=utf-8
Date: Sun, 07 May 2023 06:13:39 GMT

{
    "description": "ultrices consectetur torquent posuere phasellus urna faucibus convallis fusce sem felis malesuada luctus diam hendrerit fermentum ante nisl potenti nam laoreet netus est erat mi",
    "id": 1,
    "length": 15,
    "location": {
        "city": "Costa Mesa",
        "country": "US",
        "lat": 33.64,
        "lng": -117.93,
        "state": "CA",
        "zip": "92627"
    },
    "make": "Volkswagen",
    "model": "Bay Window",
    "name": "'Abaco' VW Bay Window: Westfalia Pop-top",
    "price": {
        "day": 16900
    },
    "primary_image_url": "https://res.cloudinary.com/outdoorsy/image/upload/v1528586451/p/rentals/4447/images/yd7txtw4hnkjvklg8edg.jpg",
    "sleeps": 4,
    "type": "camper-van",
    "user": {
        "first_name": "John",
        "id": 1,
        "last_name": "Smith"
    },
    "year": 1978
}
```

### List all rentals

```bash 
$ http :8080/rentals
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 07 May 2023 06:14:05 GMT
Transfer-Encoding: chunked

{
    "Items": [
        {
          <...10 rentals...>
        }
    ],
    "Paginator": {
        "Limit": 10,
        "Offset": 0,
        "TotalItems": 30
    }
}
```

### List rentals with pagination options

```bash
$ http ':8080/rentals?offset=5&limit=3'
HTTP/1.1 200 OK
Content-Length: 2043
Content-Type: application/json; charset=utf-8
Date: Sun, 07 May 2023 06:54:45 GMT

{
    "Items": [
        {
          <...3 rentals...>
        }
    ],
    "Paginator": {
        "Limit": 3,
        "Offset": 5,
        "TotalItems": 30
    }
}
```

### List rentals with filtering, sorting, and pagination options

```bash
$ http ':8080/rentals?price_min=9000&price_max=75000&sort=price&near=33.64,-117.93'
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 07 May 2023 06:28:28 GMT
Transfer-Encoding: chunked

{
    "Items": [
        {
            "id": 15,
            "location": {
                "lat": 34.02,
                "lng": -118.21,
                ...
            },
            "price": {
                "day": 9900
            },
            ...
        },
        {
            "id": 23,
            "location": {
                "lat": 32.73,
                "lng": -117.24,
                ...
            },
            "price": {
                "day": 9900
            },
            ...
        },
        {
            "id": 7,
            "location": {
                "lat": 33.53,
                "lng": -117.63,
                ...
            },
            "price": {
                "day": 15000
            },
            ...
        },
        {
            "id": 1,
            "location": {
                "lat": 33.64,
                "lng": -117.93,
                ...
            },
            "price": {
                "day": 16900
            },
            ...
        },
        {
            "id": 3,
            "location": {
                "lat": 32.83,
                "lng": -117.28,
                ...
            },
            "price": {
                "day": 18000
            },
            ...
        }
    ],
    "Paginator": {
        "Limit": 10,
        "Offset": 0,
        "TotalItems": 5
    }
}

```

## Running Tests

To run tests, navigate to the project root directory and execute:

```bash
$ make test
go test ./...
?   	github.com/plar/rentals-api	[no test files]
?   	github.com/plar/rentals-api/config	[no test files]
ok  	github.com/plar/rentals-api/domain	0.004s
?   	github.com/plar/rentals-api/logs	[no test files]
ok  	github.com/plar/rentals-api/handler	0.014s
?   	github.com/plar/rentals-api/repository/mocks	[no test files]
?   	github.com/plar/rentals-api/service/mocks	[no test files]
ok  	github.com/plar/rentals-api/repository	0.008s
ok  	github.com/plar/rentals-api/service	0.004s
```

This command will run all test suites in the project.

## TODO

To make the Rentals API production-ready, consider implementing the following enhancements:

1. **Configuration management**: Use environment variables, command-line flags, or a configuration file to manage settings such as database credentials, API secrets, and other sensitive information.

1. **HTTPS**: Enable HTTPS by generating or obtaining SSL/TLS certificates and configuring the server to use them.

1. **CORS**: Add CORS configuration to allow or restrict cross-origin requests from specific domains.

1. **Rate limiting**: Implement rate limiting to protect the API from excessive requests and potential denial-of-service attacks.

1. **Logging and monitoring**: Set up logging and monitoring tools to collect and analyze performance metrics, errors, and other important events.

1. **CI/CD**: Set up a continuous integration and deployment pipeline to automate testing and deployment of code changes.

1. **API versioning**: Implement API versioning to maintain backward compatibility and manage changes to the API over time.

1. **Health check and readiness endpoints**: Add health check and readiness endpoints to monitor the health of the API and ensure that it is ready to handle requests.

By addressing these points, the Rentals API will be better prepared for a production environment, ensuring security, stability, and scalability.
