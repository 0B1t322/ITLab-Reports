Feature: Creating reports

Background:
    Given url baseUrl
    Given path 'api', 'reports'
    Given header Authorization = 'Bearer ' + users.plain.accessToken
    * def uniqText = function() { return 'uniq text ' + new Date().getTime() }

Scenario: Create self report
    * def reportText = uniqText()
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
    * def reportText = uniqText()
    * def employee = 'another-employee'
    Given request { text: '#(reportText)' }
    Given param implementer = employee
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
            implementer: '#(employee)'
        },
        text: '#(reportText)'
    }
    """

Scenario: Can not create empty string report
    * def reportText = ''
    Given request { text: '#(reportText)' }
    When method post
    Then status 400

Scenario: Can not create report with empty body
    Given request { }
    When method post
    Then status 400
