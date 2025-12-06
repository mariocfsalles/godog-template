@products
Feature: vlink API_GET /stores/:store_id/productLabelling/products/:id

  Background: Set variables
    Given I set vars:
      | store_id | bomdia_pt.009648 |

  Scenario: Get products labelling information by id
    Given I set headers:
      | Apikey                    | $API_KEY                    |
      | Ocp-Apim-Subscription-Key | $VLINK_PRO_SUBSCRIPTION_KEY |
    * I set query params:
      | search | coca |
    * I send GET /stores/:store_id/productLabelling/products
    * the HTTP status should be 200
    * I store the "product.itemId" from the response body into "itemId"
    When I send GET /stores/:store_id/productLabelling/products/:itemId
    Then the HTTP status should be 200
    * the "product_sku_2391674" response should match the snapshot



