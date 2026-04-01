Feature: System Health
  In order to ensure the service is reliable
  As an operator
  I want to verify the system health status

  Scenario: The system is healthy
    Given the application is running
    When I request the health status
    Then the response should be "OK"
