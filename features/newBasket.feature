@new
Feature: As a user when I call the new basket endpoint, I would like to receive a new empty basket

  Scenario: Valid call
    Given I have a new basket request
    When I receive the response
    Then I should receive a new empty basket