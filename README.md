# viewer

A website that showcases information about different car models, their specifications, manufacturers, and more.

This project demonstrates **Clean Architecture** principles to decouple transport, business logic, and data access layers.
Leverage Go Interfaces and Dependency Injection to ensure a modular, testable, and maintainable codebase.
An app has a custom **behavioral recommendation engine**, based on user interactions with the app.
A web interface has responsive **CSS Grid** layouts without relying on heavy client-side frameworks (JS is restricted for the Go module).

## Key learnings:
- Architecting a scalable Go application using Clean Architecture principles to maintain a strict, logical separation between the transport layer, business logic, and data access.
- Developing modular, testable codebases by leveraging Go Interfaces and Dependency Injection for high component maintainability.
- Integrating external REST APIs via an abstracted data layer, allowing for seamless transitions between data sources (API/DB) without impacting core logic.
- Improving application performance through in-memory caching (with future Redis swap in mind) and UX personalized experiences via HTTP cookies.