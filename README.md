# Deck API

This is a RESTful API for managing decks of cards. It includes three endpoints for creating a deck, opening a deck, and drawing cards from a deck.

## Getting Started

To start the application, copy .env.example to .env then, run the following command:

```sh
docker-compose up
```

To start the testing environment, run the following command:

```sh
docker-compose -f docker-compose-test.yml up
```

after the build is complete, you can check the container id by running the following command:

```sh
docker ps
```

Then to run the migrations, run the following command:

```sh
docker exec -it <container-id> /bin/sh -c "cd /app/internal/database/sql/schema && goose postgres 'postgres://postgres:password@db:5432/deck?sslmode=disable' up"
```

for testing environment:

```sh
docker exec -it <container-id> /bin/sh -c "cd /app/internal/database/sql/schema && goose postgres 'postgres://postgres:password@db:5432/deck_test?sslmode=disable' up"
```

for running the tests, run the following command:

```sh
docker exec -it <container-id> /bin/sh -c "go test ./..."
```

## API Endpoints

### Create Deck

- **URL:** `/v1/decks`
- **Method:** `POST`
- **URL Params:** `shuffled=[boolean]` `cards=[string]`
- **Data Params:** None
- **Success Response:**
  - **Code:** 201
  - **Content:** `{ "deck_id": "3p40eAr6nT7Hk0CmpSSZlJ", "shuffled": false, "remaining": 52 }`
- **Description:** This endpoint creates a new deck of cards. You can optionally shuffle the deck and specify a subset of cards to include in the deck. Adding Card Codes that are not valid will be ignored or duplicated will be ignored.

### Open Deck

- **URL:** `/v1/decks/{deck_id}`
- **Method:** `GET`
- **URL Params:** `deck_id=[uuid]`
- **Data Params:** None
- **Success Response:**
  - **Code:** 200
  - **Content:** `{ "deck_id": "3p40eAr6nT7Hk0CmpSSZlJ", "shuffled": false, "remaining": 52, "cards": [...] }`
- **Description:** This endpoint opens a deck of cards. It returns the deck ID, whether the deck is shuffled, the number of remaining cards, and the remaining cards in the deck. The cards are sorted by their order of addition to the deck.

### Draw Cards

- **URL:** `/v1/decks-draw`
- **Method:** `POST`
- **URL Params:** None
- **Data Params:** `{ "deck_id": "3p40eAr6nT7Hk0CmpSSZlJ", "count": 5 }`
- **Success Response:**
  - **Code:** 200
  - **Content:** `{ "cards": [{ "code": "AS", "value": "ACE", "suit": "SPADES" }, ...] }`
- **Description:** This endpoint draws cards from a deck. You need to specify the number of cards to draw in the request body.

Please replace the placeholders with actual values when making requests to the API.
