# Coupon Service
A Go-based microservice for managing discount coupons with a RESTful API.

## Prerequisites
- Docker Engine


## Quick Start
1. Start Docker Engine
2. Run the service:
```bash
make docker/up
```
The service will be available at `http://localhost:8080/api`

## API Endpoints

### 1. Create Coupon
- **POST** `/coupon`
- Creates a new discount coupon
- Request body:
```json
{
    "code": "SAVE10",
    "discount": 10,
    "min_basket_value": 100
}
```
- Response Status: `201 Created`, no content

### 2. Get Coupons
- **GET** `/coupons`
- Retrieves coupon information
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

### 3. Apply Coupon
- **POST** `/coupons/validation`
- Applies a coupon to a basket
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

## Development

### Running Tests
```bash
make test
```

### API Documentation
Swagger documentation is available at folder `docs/`.

### Mock Generation
The project uses [github.com/matryer/moq](https://github.com/matryer/moq) to generate mocks for interfaces.

To regenerate mocks, ensure you have moq installed:
```bash
go install github.com/matryer/moq@latest
```


