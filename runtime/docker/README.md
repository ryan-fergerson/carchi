## Overview

The `Dockerfile` and `docker-compose.yml` files are used to define and orchestrate the creation and running of your Carchi application and its PostgreSQL database in Docker containers.

### Dockerfile

The `Dockerfile` is used to create a Docker image for the Carchi application. It begins by pulling the `golang:1.20.5` image, which provides a Go environment where the application can be built. It sets the working directory to `/carchi`, copies the `go.mod` and `go.sum` files from your local machine to the Docker image, and then downloads the Go modules. Afterward, it copies the rest of your application files into the Docker image. It then builds your Go application, exposing port `8080` for the application to be accessed, and finally sets the entrypoint command to start the application.

### docker-compose.yml

The `docker-compose.yml` file is used to define the services that make up your Carchi application. It specifies two services: `app` and `db`.

The `app` service builds the Docker image for the Carchi application using the `Dockerfile`, and it maps the container's port `8080` to your local machine's port `5555`. This means you can access the application's web interface at `localhost:5555`. This service depends on the `db` service and sets environment variables for the database host, name, user, password, and port.

The `db` service uses the `postgres:15.3` image to create a PostgreSQL database. It restarts always to ensure it's running when needed, maps the container's port `5432` to your local machine's port `5556`, and sets environment variables for the database. It also sets up a volume for data persistence and initializes the database using a SQL script.

### Running the services

To build and start the Carchi application and PostgreSQL database, navigate to the directory containing the `docker-compose.yml` file and run the following command:

```bash
docker-compose up -d
```

This command will build the Docker images (if not already built), create the Docker containers, and start the services in detached mode, meaning they will run in the background.

After running this command, you should be able to access the Carchi application's web interface at localhost:5555.