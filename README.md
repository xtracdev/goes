Go Event Sourcing

This project defines the base types and interfaces for working with event sourcing. The approach for this and related projects is loosely defined - packages using event sourcing are responsible for proving certain methods and observing conventions, as opposed to having a general framework that provides the hooks for packages to tap into for event sourcing.

The samples project will outline the responsibilities of packages who want to use event sourcing, most likely in conjunction with one of the projects that implement the event store.
