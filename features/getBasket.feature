@get
Feature: As a user when I try to get a existing basket, I would like to retrieve desired basket

    Scenario: Basket does not exists
        Given I have an invalid basket
        When I try to retrive the invalid basket
        Then I should receive a error message

    Scenario: Basket exists
        Given I have a valid basket
        When I try to retrive the basket
        Then I should receive desired basket