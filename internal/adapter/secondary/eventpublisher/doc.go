// Package eventpublisher includes types and functions to implement an event
// publisher. The event publisher is responsible for sending the events to a
// message broker in a reliable way and, sometimes, inside the boundaries of
// a business transaction. To do it, this implementation uses the Outbox
// pattern approach.
package eventpublisher
