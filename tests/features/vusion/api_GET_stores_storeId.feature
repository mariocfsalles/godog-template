Feature: vusion API_GET /stores/:store_id

  Background: Set variables
    Given I set vars:
      | store_id | bomdia_pt.009648 |

  Scenario: Get store by ID
    Given I set headers:
      | Apikey                    | $API_KEY                     |
      | Ocp-Apim-Subscription-Key | $VUSION_PRO_SUBSCRIPTION_KEY |
    When I send GET /stores/:store_id
    Then the HTTP status should be 200
    * the "storeId_bomdia_pt.009648" response should match the snapshot
