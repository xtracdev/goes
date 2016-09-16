@eventstore
Feature: Event Store

    Scenario: Events for an Aggregate
        Given an aggregate
        And an event store
        When the aggregate has uncommitted events
        And the events are stored
        Then the events for that aggregate can be retrieved
        And the aggregate state can be recreated using the events

    Scenario: Events for a specific aggregate id can be retrieved
        Given two aggregates
        And an event store containing all events for both aggregates
        When the events for an aggregate are retrieved
        Then only the events associated with the specific aggregate are retrieved

    Scenario: Aggregates are versioned
        Given an aggregate
        When I add an event
        Then the aggregate version is incremented
        And the aggregate version is correct when built from event history
        And all the events in the event history have the aggregate id as their source

    Scenario: Concurrency exceptions can occur if an aggregate is modified concurrently
        Given an aggregate
        When it is modified by two concurrent threads of control
        Then the second thread that stored the aggregate gets a concurrency error