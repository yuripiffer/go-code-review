# Coupon Service
A Go-based microservice for managing discount coupons with a RESTful API.
It handles coupon creation, retrieval, and application to baskets according to Role-Based Access Control (RBAC) principles.

## Prerequisites
- Docker Engine or
- Golang installed in your machine (check go.mod for go version)


## Quick start with Docker
1. Create the .env file (see details below)
2. Start Docker Engine
3. Run the service:
```bash
make docker/up
```
The service will be available at `http://localhost:8080/api/`

## API Endpoints
Below is the description of the three endpoints implemented in this service up to now. 
Be aware that `Authorization header and RBAC is only required when running the API in production` mode.
In other modes (i.e. development or test), the API will allow requests without any authorization header.

### 1. Create Coupon
- **POST** `/coupon`
- Creates a new discount coupon
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer <token>` (must have "admin" role)
- Request body:
```json
{
    "code": "SAVE10",
    "discount": 10,
    "min_basket_value": 100
}
```
- Response Status: `201 Created`, no content
- curl example (with "admin" role): 
```shell
curl --location 'http://localhost:8080/api/coupon' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1c2VyMTIzIiwicm9sZXMiOlsiYWRtaW4iXSwiZXhwIjoxNzY3MjYxNjAwfQ.AELcDhcB227f5FuLaXqv0ZmR0WDHjtunM065nQmzN7Y' \
--data '{
    "code": "COUPON100",
    "discount": 10,
    "minimum_basket_value": 50
}'
```

### 2. Get Coupons
- **GET** `/coupons`
- Retrieves coupon information
- Headers:
- `Content-Type: application/json`
- `Authorization: Bearer <token>` (must have "admin" or "user" role)
- Request body:
```json
{
    "codes": ["COUPON123", "PROMO456"]
}
```
- Response body:
```json
{
    "coupons": [
        {
            "id": "uuid-123",
            "code": "COUPON123",
            "discount": 10,
            "min_basket_value": 100
        },
        {
            "id": "uuid-456",
            "code": "PROMO456",
            "discount": 20,
            "min_basket_value": 200
        }
    ]
}
```
- curl example (with "user" role):
```shell
curl --location --request GET 'http://localhost:8080/api/coupons' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1c2VyMTIzIiwicm9sZXMiOlsidXNlciJdLCJleHAiOjE3NjcyNjE2MDB9.NL-Fn4mGtR5nBXojQRPma3k6qsz1uvukK3Me4DGc97M' \
--data '{
    "codes": ["COUPON100"]
}'
```

### 3. Apply Coupon
- **POST** `/coupons/validation`
- Applies a coupon to a basket
- Headers:
- `Content-Type: application/json`
- `Authorization: Bearer <token>` (must have "admin" or "user" role)
- Request body:
```json
{
    "value": 150,
    "code": "SAVE10TODAY"
}
```
- Response body:
```json
{
  "value": 150,
  "applied_discount": 10,
  "application_successful": true
}
```
- curl example (with "user" role):
```shell
curl --location 'http://localhost:8080/api/coupon/validation' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1c2VyMTIzIiwicm9sZXMiOlsidXNlciJdLCJleHAiOjE3NjcyNjE2MDB9.NL-Fn4mGtR5nBXojQRPma3k6qsz1uvukK3Me4DGc97M' \
--data '{
    "code": "COUPON100",
    "value": 1000
}'
```

## Data persistence
This is a experimental project, so the data is stored in memory.
The project structure enables the implementation of different data persistence layers in the future (i,e, Redis, Amazon DynamoDB, etc.).

## Development

### API Documentation
Swagger documentation is available at folder `docs/`.

### Mock Generation
The project uses [github.com/matryer/moq](https://github.com/matryer/moq) to generate mocks for interfaces.

To regenerate mocks, ensure you have moq installed:
```bash
go install github.com/matryer/moq@latest
```

### Using the Makefile commands:
#### Running Tests
```bash
make test
```

#### Generate the Swagger Documentation and tidy up
```bash
make generate
```
## Environment Variables
```shell
API_PORT=8080
API_ENV=production
JWT_SECRET="kvAJWS5rbnVxnkzVE6xOOIiBrMpytZOauEX8yOJPl20="
```




