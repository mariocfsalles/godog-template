Feature: vlink API_GET /stores/:store_id/productLabelling/labels/:id

  Background: Set variables
    Given I set vars:
      | store_id | bomdia_pt.009648 |

  Scenario: Get labels information by id
    Given I set headers:
      | Apikey                    | $API_KEY                    |
      | Ocp-Apim-Subscription-Key | $VLINK_PRO_SUBSCRIPTION_KEY |
    * I set query params:
      | search | labelId:9C1110EB AND storeId:bomdia_pt.009648 |
    When I send GET /stores/:store_id/productLabelling/labels
    * the HTTP status should be 200
    * I store the "label.labelId" from the response body into "labelId"
    When I send GET /stores/:store_id/productLabelling/labels/:labelId
    Then the HTTP status should be 200
    * the "label_9C1110EB" response should match the snapshot



