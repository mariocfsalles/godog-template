Feature: vlink API_GET /stores/:store_id/productLabelling/labels

  Background: Set variables
    Given I set vars:
      | store_id | bomdia_pt.009648 |

  Scenario: Get label information filtered
    Given I set headers:
      | Apikey                    | $API_KEY                    |
      | Ocp-Apim-Subscription-Key | $VLINK_PRO_SUBSCRIPTION_KEY |
    * I set query params:
      | search   | labelId:9C1110EB AND storeId:bomdia_pt.009648 |
      | sort     | -modificationDate                              |
      | page     | 1                                              |
      | pageSize | 1                                              |
      | includes | *                                              |
      | excludes | whatever                                       |
    When I send GET /stores/:store_id/productLabelling/labels
    Then the HTTP status should be 200
    * the "label_filtered" response should match the snapshot
