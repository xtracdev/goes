@eventpub
Feature: Event Publishing

    Scenario: All stored events are published
        Given an event sourced aggregate
        When I create and modify an instance of the aggregate
        Then all the events are published

    Scenario: Events not published when loading an aggregate
        Given an event sourced aggregate's event history
        When I instantiate the aggregate from its history
        Then no events are published

    Scenario:
        Given an event store with a registered subscriber
        When the subscriber unsubscribes
        Then previously subscribed callback is not invoked when events are published