Feature: Publishing and consuming service order events

  Background:
    Given the topic "orders.events" is accessible

  Scenario: 1) Creating a service order publishes an event and I can consume it
    And I send POST /orders with JSON:
      """
      {
        "customer": "Acme",
        "items": [
          "x",
          "y"
        ]
      }
      """
    Then the HTTP status should be 201
    Then the response body should be:
      """
      {
        "customer": "Acme",
        "id": "$ANY_ULID",
        "items": [
          "x",
          "y"
        ],
        "status": "OPEN"
      }
      """
    And I store the "id" from the response body into "order_id"
    And there must be an event on topic "orders.events" of type "OrderCreated" for "order_id" within 5s

  Scenario: 2) Updating a service order publishes an event and I can consume it
    And I have an order created via API:
      """
      {
        "customer": "Umbrella",
        "items": [
          "a"
        ]
      }
      """
    When I send PUT /orders/{order_id}/status with JSON:
      """
      {
        "status": "DONE"
      }
      """
    Then the HTTP status should be 200
    And there must be an event on topic "orders.events" of type "OrderStatusUpdated" for "order_id" within 5s

  Scenario: 3) Printing Kafka events for order creation
    Given the topic "orders.events" is accessible
    And I start printing Kafka events matching "OrderCreated"
    And I clear any pending Kafka events
    When I send POST /orders with JSON:
      """
      {
        "customer": "Acme",
        "items": [
          "x",
          "y"
        ]
      }
      """
    Then the HTTP status should be 201
    And I store the "id" from the response body into "order_id"
    And there must be an event on topic "orders.events" of type "OrderCreated" for "order_id" within 5s

  Scenario: 4) Printing Kafka events from the end
    Given the topic "orders.events" is accessible from the end
    And I start printing Kafka events matching "OrderCreated"
    And I clear any pending Kafka events
    When I send POST /orders with JSON:
      """
      {
        "customer": "Acme",
        "items": [
          "x",
          "y"
        ]
      }
      """
    Then the HTTP status should be 201
    And I store the "id" from the response body into "order_id"
    And there must be an event on topic "orders.events" of type "OrderCreated" for "order_id" within 5s

  Scenario: 5) Get order by ID after update
    Given I have an order created via API:
      """
      {
        "customer": "Umbrella",
        "items": [
          "a"
        ]
      }
      """
    When I send PUT /orders/{order_id}/status with JSON:
      """
      {
        "status": "DONE"
      }
      """
    When I send GET /orders/{order_id}
    Then the HTTP status should be 200
    Then the response body should be:
      """
      {
        "customer": "Umbrella",
        "id": "$ANY_ULID",
        "items": [
          "a"
        ],
        "status": "DONE",
        "createdAt": "$ANY_TIMESTAMP",
        "updatedAt": "$ANY_TIMESTAMP"
      }
      """
