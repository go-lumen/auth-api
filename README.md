# Goauth account management microservice

Goauth is an account management microservice written in go. It uses kafka to send events when an action is performed and is meant to be used in an eventsourcing architecture.

## Kafka Topics

co.restmark.goauth.user-creation : Where each user creation is sent.
