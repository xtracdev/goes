@testagg
Feature: TestAgg Event Store

  Scenario: Events for a TestAgg Aggregate
    Given a TestAgg aggregate
    And an event store for storing TestAgg events
    When the TestAgg aggregate has uncommitted events
    And the TestAgg events are stored
    Then the events for the TestAgg aggregate can be retrieved
    And the TestAgg aggregate state can be recreated using the events

  Scenario: TestAgg instances are versioned
    Given a TestAgg aggregate
    When I update foo
    Then the TestAgg aggregate version is incremented
    And the TestAgg aggregate version is correct when built from event history
    And all the events in the Test event history have the aggregate id as their source

