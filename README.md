# Checkout Payment Gateway

This is an example of a basic payment gateway. A payment gateway will accept authorize requests, capture funds on client card, refund or void the transaction.

## How to run?

The gateway uses Docker containers to build and run the application so no local dependencies are needed.

Steps to run locally:

1. Simply clone the repo `git clone https://github.com/MihaiBlebea/go-checkout.git`

2. Navigate to the folder where the repo was cloned

3. Run `make setup-test` to run the unit tests in a Docker container. Two new fils will appear in your root folder:
    - `cover.out` file containing the test coverage (hard to read for us humans)
    - `cover.html` which will show in much more detail (and easy to read for humans) what % of the code is covered by tests. This file will also be opened in a nice browser window.

4. Run `make setup` to compile and run the project in a Docker container. Tip: Make sure that no other application is running on your host post `8087` (random port number). You can easily change the port by editing the `Makefile` in the root folder (`up` command, line 18). If you consider changing also the container port, then please run `make env-file` and also update the `.env` file that is created in your root folder.

5. Enjoy the payment gateway!

## Technical decisions and architecture

The project contains 2 major packages witch handle different parts of the application:
- server
- gateway

### API endpoints

- GET `/api/v1/health-check` - simple health check endpoint
    - output:
    ```json
    {
        "ok": true
    }
    ```
- POST `/api/v1/authorize` - authorizes a transaction on the user card. `expire_year` and `expire_month` are supplied as integers, while the `cvv` field is a string as it does not represent a quantity. Amount is calculated in pennies for precision, so `200` is £2.
    - input:
    ```json
    {
        "name_on_card": "John Doe",
        "card_number": "4111 1111 1111 1111",
        "expire_year": 2022,
        "expire_month": 4,
        "cvv": "755",
        "amount": 200,
        "currency": "GBP"
    }
    ```
    - success output:
    ```json
    {
        "id": "292ab873-4c07-4f5e-ba59-9902d563c3be",
        "success": true,
        "amount": 200,
        "currency": "GBP"
    }
    ```
    - fail output:
    ```json
    {
        "success": false,
        "message": "Invalid name on card"
    }
    ```
- POST `/api/v1/capture` - captures funds authorized by the previous endpoint. Using the id from the previous endpoint and an amount in pennies, you can start capturing funds on the user credit card. Remaining amount in response represent the funds remaining to be captured on the user card.
    - input:
    ```json
    {
        "id": "1117bbb5-4772-4065-a8f5-5f3134afe299",
        "amount": 20
    }
    ```
    - success output:
    ```json
    {
        "success": true,
        "remaining": 180,
        "currency": "GBP"
    }
    ```
    - fail output:
    ```json
    {
        "success": false,
        "message": "Unavailable amount"
    }
    ```
- POST `/api/v1/refund` - refunds the funds captured. Using the id of the transaction and an ammount in pennies, you can refund the money captured previously on the user card. The refund amount cannot be greater then the one captured. Remaining field on the response specifies the funds remaining to be refunded. Once a transaction is put in refund state (by creating a first refund) no capture operation is allowed.
    - input:
    ```json
    {
        "id": "9e15a26e-a139-4f16-914c-ab7783ba1495",
        "amount": 5
    }
    ```
    - success output:
    ```json
    {
        "success": true,
        "remaining": 15,
        "currency": "GBP"
    }
    ```
    - fail output:
    ```json
    {
        "success": false,
        "message": "Unavailable amount"
    }
    ```
- POST `/api/v1/void` - voids a transaction. Using the id of the transaction you can void a payment. Once th payment is voided, no other operations are permitted. The `balance` field on the response represents the amount remaining to be refunded to the user.
    - input:
    ```json
    {
        "id": "9e15a26e-a139-4f16-914c-ab7783ba1495",
    }
    ```
    - success output:
    ```json
    {
        "success": true,
        "balance": 15,
        "currency": "GBP"
    }
    ```
    - fail output:
    ```json
    {
        "success": false,
        "message": "Unavailable amount"
    }
    ```
- GET `/api/v1/list` - lists all transactions in the payment gateway. List all transactions made using the payment gateway. Notice that the id of the transactions is returned with the response, in case of a production gateway this should be truncated for security reasons or put behind a merchant auth system.
    - output:
    ```json
    {
        "success": true,
        "transactions": [
            {
                "id": "348233c3-9aff-47d8-9fba-483b2b2e27cc",
                "state": 0,
                "amount": 200,
                "captured": 0,
                "refunded": 0,
                "currency": "GBP"
            },
            {
                "id": "292ab873-4c07-4f5e-ba59-9902d563c3be",
                "state": 0,
                "amount": 200,
                "captured": 0,
                "refunded": 0,
                "currency": "GBP"
            }
        ]
    }
    ```