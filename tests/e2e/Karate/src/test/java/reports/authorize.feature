Feature: Check authorize

Background:
    Given url baseUrl
    # plain secured path
    Given path 'api', 'reports'

Scenario: Without authorize header
    When method get
    Then status 401
@ignore # Drops app, commented to use another tests
Scenario: User without itlab: user claim
    Given header Authorization = 'Bearer ' + users.incorrect.accessToken
    When method get
    Then status 401

Scenario: Success get reports
    Given header Authorization = 'Bearer ' + users.plain.accessToken
    When method get
    Then status 200

Scenario: Success get reports by admin user
    Given header Authorization = 'Bearer ' + users.admin.accessToken
    When method get
    Then status 200
