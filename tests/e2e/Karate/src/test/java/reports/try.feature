Feature: Try work with karate

Background:
    Given url baseUrl
    Given header Authorization = 'Bearer ' + authToken

Scenario: Try get reports
    Given path 'api', 'reports'
    When method get
    Then status 200
    And match response == []
