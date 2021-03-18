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

### Server package

Main entry point in the `server` package is the `server.go` file that exposes a simple `NewServer` method which accepts a `gateway` service and a `logger` service, both represented by an interface situated in the `contracts.go` file.

The port adaptor pattern works great with the Go inferred interface and offers a much more decoupled solution compared with other languages that require a specified interface and a "middle-man" class to act as the adaptor. The downside of this implementation is that we lose the adaptor layer that could be used to map different parameters between the port and the third party service.

To mitigate this loss, we could bring in a decorator for the third party service which is used to map params between the port and the third party service.

The logging middleware is situated in the `middleware.go` file and it integrates nicely with the gorilla/mux package that handles the http server.

Inside the server package we also have the validate sub-package that provides a nice `tag` based validation to account for the required params in the request body. This can be easily extended to provide more custom tags or rules. For this example I kept this minimal and just implemented a check for required fields. In case of needing a more complex solution this could be broken into individual validators which could be called behind a factory.

The `handler` package inside the main`server` package offers a nice abstraction for the http handlers needed by the http "framework". It also contains a `contracts.go` file which defines the necesary interfaces that need to be implemented by it's dependencies.

### Gateway package

The gateway package acts like the domain package in this appliation defining all the business rules.

The `api.go` file represents the entry point in this package and defines all the exposed methods and structs.

I have chosen a "decorator" pattern to wrap the private methods and expose them outside of the `gateway` package. This would be a good start to add logging in this layer. I considered this approach but decided against it as it would have created a lot of noise and would have made tracking logs a bit more dificult with so many layers loging the same thing. Definetely worth adding for a bigger and more complex application.

I decided to not link all the methods in the package to the main service struct, as that would have made testing a bit more complex. This is a pattern that I grew very fond of and I think it combines the best out of both worlds, making tests easier to implement (similar to functional languages) and also allow for a more "OOP" friendly code (not really sure how to call it).

I chose to implement the gateway application in memory so this test dos not contain a database or a persistence layer. I considered adding a database but decided against it, as I think this makes the test a bit more interesting and brings some nice problems to solve. For example:

- race conditions reading and writing to the `transactions` map. To mitigate the risk I chosen a mutex lock that we can lock and unlock before reading and writing to it. I chosen to lock just the reads in some of the cases.

- persistence issue: one solution to further add persistence (database) would be to use composition to extend the functionality of the gateway package service. Consider this:

```Go
package persistgtway // from persistent gateway :)

import (
    gtway "github.com/MihaiBlebea/go-checkout/gateway"
)
type PersistentGateway struct {
    Gateway
    *gorm.DB // or any other orm
}

func (ps *PersistentGateway) AuthorizePayment(options gtway.AuthorizeOptions) (string, error) {
	id, err := ps.AuthorizePayment(options)
    if err != nil {
        return "", err
    }

    transactions := ps.ListTransactions()

    // persist the transaction into database

    return id, nil
}
```

## Improvements and TO-DOs

-[] Add more tests and better cover the API endpoints. Not very happy of the API test coverage.
-[] Add persistence layer using gorm or any other ORM abstracted behind repositoris.
-[] Add auth for the API endpoints and limit access based on merchant.
-[] Better handle the void case, as I am not vry familliar with how the void state would affect the transaction I left the money in the transaction and not updated it's captured and refunded field. Not sure if refunding all the captured amounts as part of the void process would work.
-[] Better handle the states of the transaction using a state machine and a set of simple rules for what can and cannot happen depending on current state of the transaction.
-[] Improve the sandbox cards logic.