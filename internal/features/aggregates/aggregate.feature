@aggregate
Feature: Event Sourced Aggregate

    Scenario: Aggregate created
        Given an event sourced aggregate
        When I create new instances
        Then the instances have unique IDs
        And there's an uncommitted event
        And the uncommited event's source ID is the aggregate ID

    Scenario: Aggregate from event history
        Given an event sourced aggregate's event history
        When I instantiate the aggregate from its history
        Then the instance state is correct
        And there are no uncommitted events
        And the aggregate version is correct
