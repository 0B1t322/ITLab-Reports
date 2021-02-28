Feature: Receive reports

Background:
    Given url baseUrl
    Given path 'api', 'reports'
    * def userId = 'user-receive-test-' + uniqueText()
    * def accessToken = createUserJwt(userId)
    Given header Authorization = 'Bearer ' + accessToken
    * def sleep = function(pause){ java.lang.Thread.sleep(pause) }

Scenario: Get all reports created by user
    * table reports
        | text       | employee                      |
        | 'Report 1' | 'implementer-' + uniqueText() |
        | 'Report 2' | 'implementer-' + uniqueText() |
        | 'Report 3' | 'implementer-' + uniqueText() |
    * call read('create-report.feature') reports
    When method get
    Then status 200
    And match each response[*].assignees.reporter == userId
    And match response[*].assignees.implementer contains only $reports[*].employee

Scenario: Get all reports created about user
    * def targetUserId = 'target-user-' + uniqueText()
    * table reports
        | text       | employee     |
        | 'Report 1' | targetUserId |
        | 'Report 2' | targetUserId |
        | 'Report 3' | targetUserId |
    * call read('create-report.feature') reports
    * def targetUserToken = createUserJwt(targetUserId)
    Given header Authorization = 'Bearer ' + targetUserToken
    And method get
    Then status 200
    And match each response[*].assignees.reporter == userId
    And match each response[*].assignees.implementer == targetUserId

Scenario: Get all reports created about user and by user
    * def secondUserId = 'second-user-' + uniqueText()
    
    # about second from first
    * call read('create-report.feature') {text: 'Report about second from first', employee: '#(secondUserId)'}
    
    # about first from second
    * def secondUserToken = createUserJwt(secondUserId)
    * def accessToken = secondUserToken
    * def defaultUserId = userId
    * def userId = secondUserId
    * call read('create-report.feature') {text: 'Report about first from second', employee: '#(defaultUserId)'}
    * def userId = defaultUserId

    Given header Authorization = 'Bearer ' + secondUserToken
    And method get
    Then status 200
    And match response[*].assignees.reporter contains ['#(userId)', '#(secondUserId)']
    And match response[*].assignees.implementer contains ['#(userId)', '#(secondUserId)']

Scenario: Get reports in date range
    * def plusSecondsToDate = 
    """
        function(stringDate, seconds) {
            if (!stringDate.endsWith('Z')) {
                stringDate += 'Z';
            }
            return new Date(Date.parse(stringDate) + 1000 * seconds).toISOString();
        }
    """

    * def firstReport = call read('create-report.feature') { text: 'First report', employee: '#(userId)' }
    
    * def beforeFirstReport = plusSecondsToDate(firstReport.response.date, -5)
    * def rightAfterFirstReport = plusSecondsToDate(firstReport.response.date, 1)

    # sleep 3 seconds
    * call sleep 3000 
    * def secondReport = call read('create-report.feature') { text: 'Second report', employee: '#(userId)' }
    * def rightAfterSecondReport = plusSecondsToDate(secondReport.response.date, 20)

    When path 'employee', userId
    When params {dateBegin: '#(beforeFirstReport)', dateEnd: '#(rightAfterFirstReport)'}
    When method get
    Then status 200
    And assert response.length == 1
    And assert response[0].text == 'First report'

    Given path 'api', 'reports', 'employee', userId
    Given header Authorization = 'Bearer ' + accessToken
    When params {dateBegin: '#(rightAfterFirstReport)', dateEnd: '#(rightAfterSecondReport)'}
    When method get
    Then status 200
    And assert response.length == 1
    And assert response[0].text == 'Second report'


    # When path 'api', 'reports', 'employee', userId
    # When params {dateBegin: '#(rightAfterFirstReport)', dateEnd: '#(longAfterFirstReport)'}
    # Given header Authorization = 'Bearer ' + accessToken

