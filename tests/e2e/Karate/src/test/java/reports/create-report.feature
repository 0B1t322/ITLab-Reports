@ignore
Feature: Creating one report

Background:
    Given url baseUrl
    Given path 'api', 'reports'
    Given header Authorization = 'Bearer ' + users.plain.accessToken

Scenario:
    Given request { text: '#(text)' }
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
        text: '#(text)'
    }
    """
