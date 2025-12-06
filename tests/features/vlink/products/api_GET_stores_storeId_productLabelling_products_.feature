Feature: vlink API_GET /stores/:store_id/productLabelling/products

  Background: Set variables
    Given I set vars:
      | store_id | bomdia_pt.009648 |

  Scenario: Get products labelling information
    Given I set headers:
      | Apikey                    | $API_KEY                    |
      | Ocp-Apim-Subscription-Key | $VLINK_PRO_SUBSCRIPTION_KEY |
    * I set query params:
      | search   | coca |
      | pageSize | 1    |
    When I send GET /stores/:store_id/productLabelling/products
    Then the HTTP status should be 200
    * the "product_filtered" response should match the snapshot
