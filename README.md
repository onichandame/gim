# Gim(Go-In-Modules)

Modular framework based on dependency injection. The modularization design is borrowed from NestJS.

# Usage

see [examples](./examples)

# Architecture

The internal system of Gim is (partly) described here. This section is written for the users who want to know how it works without reading the code.

## Pre-requisite

By reading the following sections, it is assumed that the user has got the following ideas in mind:

- Dependency Injection(DI)
- Inversion of Control(IoC)
- Golang's reflect system

## Components

In Gim, there are 2 basic entities: the container and the concrets inside the container. The concrets can be declared as the instance or a constructor returning an instance. When declared as a constructor, its parameters can be other concrets in the container which will be injected on bootstrap.

- Module: corresponds to a container
- Provider: corresponds to a concret

## Bootstrap

Before a module is bootstrapped, its children modules will be bootstrapped first. Then the exported providers of all children modules will be loaded into the current module. If a child exports a submodule, all of the submodule's exported providers will be loaded to the current module.

During bootstrap, all the providers declared as concrets will be loaded first. Then all the providers declared as constructors will be fed with the needed dependencies and executed to get the provided concret. The dependencies must be the sibling providers in the module that have been loaded.
