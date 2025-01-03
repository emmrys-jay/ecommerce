# Ecommerce

## Description

A simple RESTful Ecommerce API written in Go programming language. The structure of this project was inspired by [Bagas Hizbullah](https://github.com/bagashiz) through this article:

https://dev.to/bagashiz/building-restful-api-with-hexagonal-architecture-in-go-1mij

It uses the Go standard library with [Go-chi](https://gin-gonic.com/) as the lightweight router and [PostgreSQL](https://www.postgresql.org/) as the database with [pgx](https://github.com/jackc/pgx/) as the driver and [Squirrel](https://github.com/Masterminds/squirrel/) as the query builder. 

## Getting Started

The development process relies heavily on Make. Make is used to manage dependencies, run migrations, build the application, and run tests.

1. Ensure you have [Go](https://go.dev/) 1.23 or higher and [Make](https://www.gnu.org/software/make/) installed on your machine:

    ```bash
    go version && make --version
    ```
    If they are not installed, install them using the Links above for your specific architecture.
2. Create a copy of the `config-sample.yml` file and rename it to `config.yml`:

    ```bash
    cp config-sample.yml config.yml
    ```
    Update configuration values as needed.

3. Update the `Makefile` variable `CONFIG_FILE` to point to your `config.yml` file.

## Database Setup

To set up the database, you have three options:


*Option 1: Using Docker Compose*

You can use Docker Compose to start a Postgres database container. To do this:


1. Install Docker and Docker Compose on your machine.
2. Run the following command to start the database container:
    ```
    make service-up
    ```
You can use `make service-down` to stop the postgres container.

3. The database will be available at `localhost:5433`. The port `5433` was used to avoid conflict with any locally installed postgres instance, but you can modify it in your `config.yml` file.

*Option 2: Using a Local Postgres Database*

If you already have a Postgres database installed locally, you can use it instead:


1. Ensure your local Postgres instance is running and accessible.
2. Ensure the username, password and database in your `config.yml` file is created for your local Postgres instance. You can use `psql` to connect to your local instance using the default user:
    ```bash
    psql -d postgres
    ```
3. Update the `PORT` environment variable in the `config.yml` file to point to the port of your local Postgres instance database.


## Starting the Application

1. Run database migrations using (Optional: It also runs on application startup):

    ```bash
    make migrate-up
    ```
    You can drop all tables in the database using `make migrate-down`. You can ignore the migrate commands because the application automatically runs your database migrations on startup.

2. Install all dependencies necessary for running the commands to start up the application on main:

    ```bash
    make install
    ```

    Some of the commands are written for MacOS ARM 64 architecture. Update them accordingly for any other operating system.

3. Run the project in development mode using:

    ```bash
    make dev
    ```
4. Run the project in production mode using:
    ```bash
    make start
    ```

## Documentation

API documentation can be found in `docs/` directory. To view the documentation, open the browser and go to `http://127.0.0.1:8080/swagger/index.html`. The documentation is generated using [swaggo](https://github.com/swaggo/swag/) with [http-swagger](https://github.com/swaggo/http-swagger/) middleware.

Also, you can view the [Postman](https://www.postman.com/) documentation [here](https://documenter.getpostman.com/view/27735481/2sAYJ6CKfj)

## Learning References

1. [Hexagonal architecture in Go](https://dev.to/bagashiz/building-restful-api-with-hexagonal-architecture-in-go-1mij)
