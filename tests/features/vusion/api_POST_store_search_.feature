Feature: vusion API_POST /stores/search

  Scenario: Get stores by text search
    Given I set headers:
      | Apikey                    | $API_KEY                     |
      | Ocp-Apim-Subscription-Key | $VUSION_PRO_SUBSCRIPTION_KEY |
    When I send POST /stores/search with JSON:
      """
      {
        "page": 1,
        "pageSize": 50,
        "search": "000206 - Montijo"
      }
      """
    Then the HTTP status should be 200
    * the "search_000206" response should match the snapshot
