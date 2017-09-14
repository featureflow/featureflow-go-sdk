Feature: UserBuilder
  Scenario: Test the User Builder can build a valid user with a key
    Given there is access to the User Builder module
    When the builder is initialised with the key "user"
    And the user is built using the builder
    Then the result user should have a key "user"
    And the result user should have no values

  Scenario: Test the User Builder can build a valid user with a key
    Given there is access to the User Builder module
    When the builder is initialised with the key "user"
    And the builder is given the following attributes
      | key  | value  |
      | age  | 21     |
      | type | beta   |
    And the user is built using the builder
    Then the result user should have an id "user"
    And the result user should have the key "age" with value "21"
    And the result user should have the key "type" with value "beta"

  Scenario: Test the User Builder throws an error when no id is provided
    Given there is access to the User Builder module
    When the builder is initialised with the id ""
    Then the builder should throw an error
