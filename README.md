# Receipt Processing API

## Overview

This API allows users to submit receipt data and receive points based on specific rules. It provides two main endpoints:

1. **Process Receipts**: Submits a receipt and returns a unique ID for the receipt.
2. **Get Points**: Retrieves the points associated with a specific receipt ID.

## Requirements
- Go 
- `gin` package for the Gin framework (use `go get github.com/gin-gonic/gin`)
- `uuid` package for generating unique IDs (use `go get github.com/google/uuid`)

## Configuration
- The application runs on **port 8080** by default. If you wish to change the port, modify the following line in the `main.go` file:
    ```go
    r.Run(":8080") // Change the port here if needed
    ```

## Endpoints

### 1. **Process Receipts**

- **Path**: `/receipts/process`
- **Method**: `POST`
- **Payload**: 
    ```json
    {
      "retailer": "Retailer Name",
      "purchaseDate": "2024-11-23",
      "purchaseTime": "14:30",
      "items": [
        {
          "shortDescription": "Item 1",
          "price": "19.99"
        },
        {
          "shortDescription": "Item 2",
          "price": "9.99"
        }
      ],
      "total": "29.98"
    }
    ```
- **Response**: 
    ```json
    {
      "id": "7fb1377b-b223-49d9-a31a-5a02701dd310"
    }
    ```
  - The response will contain an ID for the receipt.

### 2. **Get Points**

- **Path**: `/receipts/{id}/points`
- **Method**: `GET`
- **Response**:
    ```json
    {
      "points": 32
    }
    ```
  - The response will return the number of points awarded for the receipt with the given ID.

## Points Calculation Rules

The points awarded for a receipt are calculated based on the following rules:

1. **Retailer Name**: 1 point for each alphanumeric character in the retailer's name.
2. **Round Dollar Amount**: 50 points if the total is a round dollar amount (no cents).
3. **Multiple of 0.25**: 25 points if the total is a multiple of 0.25.
4. **Item Count**: 5 points for every two items in the receipt.
5. **Item Description Length**: 5 points if the description length is divisible by 3 and the price is multiplied by 0.2, rounded up to the nearest integer.
6. **Odd Day of Purchase**: 6 points if the purchase date is on an odd day.
7. **Afternoon Purchase**: 10 points if the purchase time is between 2:00 PM and 4:00 PM.

## Testing

You can test the API using **curl**:

# 1. **Test Submitting a Receipt**
```
curl -X POST http://localhost:8080/receipts/process \
  -H "Content-Type: application/json" \
  -d '{
    "retailer": "Retailer Name",
    "purchaseDate": "2024-11-23",
    "purchaseTime": "14:30",
    "items": [
      {"shortDescription": "Item 1", "price": "19.99"},
      {"shortDescription": "Item 2", "price": "9.99"}
    ],
    "total": "29.98"
  }'
```
# 2. **Test Retrieving Points for a Receipt**

# Replace {id} with the actual receipt ID returned from the previous request.
```curl http://localhost:8080/receipts/{id}/points```

