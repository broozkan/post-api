PostAPI
=======

PostAPI is a RESTful API that manages posts using Couchbase as the database.

Prerequisites
-------------

-   Docker
-   Docker Compose
-   Golang 1.20
-   gosec (not mandatory)

Folder Structure
----------------

The project's folder structure is organized as follows:

-   `cmd`: contains the main application file and its tests.
-   `handlers`: contains the HTTP request handlers and their tests.
-   `internal`: contains the internal packages of the application, including its services, repository, models, and configuration, and their tests.
-   `pkg`: contains the shared packages that can be used by different applications. Currently, it only includes the server package that provides a basic HTTP server implementation.
-   `test`: contains the test data used in the unit and integration tests.
-   `Dockerfile`: the Dockerfile used to build the Docker image of the application.
-   `Makefile`: includes various commands to build, test, and run the application.
-   `README.md`: provides instructions and information about the project.
-   `docker-compose.yml`: the Docker Compose file used to run the application and its dependencies.
-   `go.mod` and `go.sum`: the Go module files that define the project's dependencies.

Configuration
-------------

The application uses configuration files in YAML format located in the `.config` directory. There are three different configuration files for different environments:

-   `dev.yaml`: development environment
-   `staging.yaml`: staging environment
-   `prod.yaml`: production environment

You should copy the relevant configuration file to `./config` in the root directory of the project and modify it to suit your needs.



Running the Application
-----------------------

### Couchbase Image Path (skip this if you are not using apple silicon)

If you are using an Apple Silicon-based machine, you need to update the path of the Couchbase image in the `docker-compose.yml` file. There are two different paths for the `couchbase` and `couchbase-arm` folders located in the `.dev/deployment` directory. Follow the steps below to update the image path:

1.  Open the `docker-compose.yml` file in a text editor.
2.  Find the `couchbase` service and replace the `image` field with the path of the `couchbase` or `couchbase-arm` folder depending on your machine architecture.

    yamlCopy code

    ```
    couchbase:
    image: couchbase/server:community-aarch64 # change this to .dev/deployment/couchbase-arm if you are using an Apple Silicon-based machine
    container_name: couchbase
    ```

3.  Save the changes and close the file.

After making these changes, you can run the application with the updated Couchbase image path.

dev.yaml is configured already for running locally

You can start the application and its dependencies (Couchbase server) using Docker Compose:


```
export APP_ENV=dev
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