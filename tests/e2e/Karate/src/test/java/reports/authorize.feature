Feature: Check authorize

Background:
    Given url baseUrl

Scenario: Without authorize header
    Given path 'api', 'reports'
    When method get
    Then status 401

Scenario: User without itlab: user claim
    Given header Authorization = 'Bearer ' + noITLabUserAuthToken
    Given path 'api', 'reports'
    When method get
    Then status 401

Scenario: Success get reports
    Given header Authorization = 'Bearer ' + userAuthToken
    Given path 'api', 'reports'
    When method get
    Then status 200
