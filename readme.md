## Финальный проект
В этом задании мы будет делать реальное тестовое задание от реальной компании на позицию Golang разработчика. 

Задание предлагается вам в изначальном виде, каким его отправила компания. 

Делайте его тщательно, это очень полезный для вас опыт. После того, как мы закончим работать с ним, я покажу вам финальный вид решения, которое было принято работодателем.

Обратите внимание, что в этом задании требуются тесты - без них задание не считается решенным. 

## Overview of the Task

### Objective: 
Implement a system with two microservices:
- Gateway Service: A RESTful service that acts as a gateway, receiving user requests and forwarding them to the Hashing Service using gRPC.
- Hashing Service: A service that checks if a payload's hash already exists and creates a hash for new payloads.
### Requirements
- Language and Library: The solution must be implemented in Go, leveraging the go-kit library for building microservices.
- Architecture: Follow a clear architecture pattern (e.g., Clean Architecture, Hexagonal) to ensure separation of concerns and maintainability.
### Endpoints:
- CheckHash: Checks if the payload's hash already exists.
- GetHash: Returns the hash for an existing payload.
- CreateHash: Creates and stores a hash for a new payload.
- Communication: The Gateway Service communicates with the Hashing Service using gRPC.
- Testing: The solution should be thoroughly tested, including unit tests and integration tests to cover the logic and communication between services.
