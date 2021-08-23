# **Image Service**
The Image Service is a service that maintains the upload of images to the file system.

Create .env file in root directory and add following values:
~~~~
DB_USER=postgres
DB_PASSWORD=Password1
DB_HOST=localhost
DB_PORT=5432
DB_NAME=conversion_compression_service

TOKEN_TTL=12h

RABBITMQ_URL=amqp://guest:guest@localhost:5672/
~~~~
Use `make build` to build and `make run` to run project, `make lint` to check code with linter.
