PostAPI
=======

PostAPI is a RESTful API that manages posts using Couchbase as the database.

Prerequisites
-------------

-   Docker
-   Docker Compose
-   Golang 1.20
-   gosec (not mandatory)

Configuration
-------------

The application uses configuration files in YAML format located in the `.config` directory. There are three different configuration files for different environments:

-   `dev.yaml`: development environment
-   `staging.yaml`: staging environment
-   `prod.yaml`: production environment

You should copy the relevant configuration file to `config.yaml` in the root directory of the project and modify it to suit your needs.

Running the Application
-----------------------

You can start the application and its dependencies (Couchbase server) using Docker Compose:


```
make up
```

This will start the Couchbase server and the PostAPI service.

Running Unit Tests
------------------

You can run unit tests using the following command:

```
make unit-test
```

Running Integration Tests
-------------------------

You can run integration tests using the following command:

```
make integration-test
```

Running Database Tests
----------------------

You can run database tests using the following command:


```
make db-test
```

Code Coverage
-------------

You can generate code coverage reports using the following command:

```
make code-coverage
```

This will generate a coverage report in HTML format and print the total coverage percentage to the console.

Linting
-------

You can run linting checks using the following command:

```
make lint
```

Security
--------

To run security checks with `gosec`:

```
make security-check
```

License
-------

This project is licensed under the MIT License - see the `LICENSE` file for details.