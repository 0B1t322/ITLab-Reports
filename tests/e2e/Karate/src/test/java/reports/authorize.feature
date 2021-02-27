Feature: Check authorize

Background:
    Given url baseUrl

Scenario: Without authorize header
    Given path 'api', 'reports'
    When method get
    Then status 401

Scenario: Success get reports
    Given header Authorization = 'Bearer ' + authToken
    Given path 'api', 'reports'
    When method get
    Then status 200
