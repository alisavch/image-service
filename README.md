# **Image Service**
Image Service is a service that supports compressing and converting images and uploading them to the file system or AWS S3 Bucket.

Create .env file in root directory and add following values to run in docker container:
~~~~
DB_USER=postgres
DB_PASSWORD=Password1
DB_HOST=postgresql
DB_PORT=5432
DB_NAME=conversion_compression_service

TOKEN_TTL=12h

RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

AWS_REGION=YOUR_REGION
AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY=YOUR_SECRET_ACCESS_KEY
BUCKET_NAME=YOUR_BUCKET_NAME
~~~~
Use `make build` to build and `make run` to run project, `make lint` to check code with linter.
