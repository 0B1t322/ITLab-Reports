@ignore
Feature: Creating one report

Background:
    Given url baseUrl
    Given path 'api', 'reports'
    Given header Authorization = 'Bearer ' + accessToken

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
            reporter: '#(userId)',
            implementer: '#(employee)'
        },
        text: '#(text)'
    }
    """
