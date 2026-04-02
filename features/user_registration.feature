Feature: User Registration
  As a new customer
  I want to register an account
  In order to access the system

  Scenario: Successful registration results in persistent state
    Given the registration service is available
    When I register a user with email "alice@example.com" and name "Alice"
    Then the response status should be 201
    And the user "alice@example.com" should exist in the database
