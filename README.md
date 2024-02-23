<!-- TOC -->
  * [About](#about)
  * [Requirements](#requirements)
      * [External services, system packages and clients:](#external-services-system-packages-and-clients)
  * [Architecture](#architecture)
  * [File structure](#file-structure)
  * [Configuration file:](#configuration-file)
  * [Compiling the program](#compiling-the-program)
  * [Running the application](#running-the-application)
  * [Database migrations](#database-migrations)
  * [Fixtures](#fixtures)
    * [Generate](#generate)
    * [Apply](#apply)
<!-- TOC -->


## About

Web application for storing and managing midi files created for Final Fantasy XVI bards. 


## Requirements

golang 1.22

#### External services, system packages and clients:
- Discord - used for user registration and authentication
- Spotify - used for gathering song related information
- (soon) MySQL 8.3.0. Currently using sqlite3 for fast development

### Deploying
- This application will be deployed on a Kubernetes cluster. (k8s configs wip)
- It can also be deployed on any of the compatible platforms

## Architecture

FFXIV Bard is a web application that implements the Hexagonal architecture design principles.

The [hexagonal architecture](docs/hexagonal_architecture.png) is based on three principles and techniques:

- Explicitly separate User-Side, Business Logic, and Server-Side
- Dependencies are going from User-Side and Server-Side to the Business Logic
- We isolate the boundaries by using Ports and Adapters

More info on the hexagonal
architecture [here](https://blog.octo.com/hexagonal-architecture-three-principles-and-an-implementation-example/).


## File structure

```
├── cmd/ - Contains entry points from commands to workers to rest. Can access both domain and infrastructure components.
├── config/ - Application's parameters configuration.
├── container/ - Service container holding registered services.
├── docs/ - Documentation resources.
├── domain/ - Contains domain components with the strict rule of never using external dependecies.
├── infrastructure/ - Contains logic for external clients like APIs or storage solutions.
└── port/ - Contains Ports(Interfaces)/DTOs.
```


## Configuration file:
See the [configuration template file](config/config-sample.toml).

## Compiling the program
```make compile```
The program will compile binaries for Linux/Darwin/Windows platforms, each supporting both AMD64 and ARM64 architecture. 

## Running the application
 - Make sure you set-up a config/config.toml file before you start compiling
 - All resources (html, css, js, toml) files are embeded in the binary.
 - Executing the program with the following arguments will start the server ```<path/to/binary> server --start --port 80 --pool 10 ```

## Database migrations
- Executing the program with the following arguments will run the database migrations ```<path/to/binary> migrate --up```
- Executing the program with the following arguments will undo the database migrations ```<path/to/binary> migrate --down```

## Fixtures
### Generate
- Executing the program with the following arguments will generate the fixture files ```<path/to/binary> fixtures --generate --count 1000```
### Apply
- Executing the program with the following arguments will apply the generated fixture files ```<path/to/binary> fixtures --execute```
