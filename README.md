# Gim

Modular framework based on dependency injection. The modularization design is borrowed from NestJS.

# Usage

see [examples](./examples)

# Design

It follows the simple inversion of control(IoC) idea of using singletons as providers and inject the required singletons where needed across the whole application.

Every Gim module corresponds to an IoC container where providers inside can take advantage of each other. Above this, a module can export some(or all) of its providers for consumption by parent modules.
