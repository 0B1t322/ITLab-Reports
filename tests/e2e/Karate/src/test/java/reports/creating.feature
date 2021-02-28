Feature: Creating reports

Background:
    Given url baseUrl
    Given path 'api', 'reports'
    Given header Authorization = 'Bearer ' + users.plain.accessToken

Scenario: Create self report
    * def reportText = 'report  body ' + uniqueText()
    Given request { text: '#(reportText)' }
    When method post
    Then status 200
    And match response == 
    """
    {
        id: '#present',
        date: '#present',
        archived: false,
        assignees: {
            reporter: '#(users.plain.id)',
            implementer: '#(users.plain.id)'
        },
        text: '#(reportText)'
    }
    """

Scenario: Create report about another employee
    * def reportText = uniqueText()
    * def employee = 'another-employee'
    * def userId = users.plain.id
    * def accessToken = users.plain.accessToken
    Then call read('create-report.feature') { text: '#(reportText)', employee: '#(employee)' }

Scenario: Can not create empty string report
    * def reportText = ''
    Given request { text: '#(reportText)' }
    When method post
    Then status 400

Scenario: Can not create report with empty body
    Given request { }
    When method post
    Then status 400
