# **Image Service**
Image service is a service that allows compressing and converting images and after that getting them for saving.


## Installation
To start using image-service, Install Go and run go get:

```
go get -u github.com/alisavch/image-service
```

## Start service (requires docker)

server app listening on http://localhost:8080
```
make build-containers
```

## Implementing a server 
First, a server instance has to be created:
```
    srv := NewServer(rabbit, currentService)
```
The NewServer constructor actually takes a message broker RabbitMQ and Service:
```
package apiserver

func NewServer(mq AMQP, service *Service) *Server {
	s := &Server{
		router:  mux.NewRouter(),
		mq:      mq,
		service: service,
		logger:  NewLogger(),
	}
	s.ConfigureRouter()
	return s
}
```

Service is a structure containing ServiceOperations and ServiceOperations is an implementation
integrations of Authorization, Image and S3Bucket.

The authorization is a set of functions for authorizing a user i.e. user creation, token generation, token analysis.
The image is a set of functions for integration with an image, such as compressing, converting, uploading, getting, etc.
The S3Bucket is a set of functions for interacting with AWS, specifically uploading and downloading image to an S3 bucket.

Service also has message broker. The broker is what dispatches events to clients.
When you publish a message, the broker distributes it to all connections (subscribers).

Logger is a structured logger for API. Function NewLogger() takes log.NewLogger() that automatically takes  
parameter logrus.New (package github.com/sirupsen/logrus) and you can change this if you need.

The last step is router. Endpoints: 
~~~
POST - /api/sign-up - create user
POST - /api/sign-in - user authorization
GET  - /api/history - get user request history
POST - /api/compress?width={value} - compress image
GET  - /api/compress/{compressedID}?original={value} - get/download compressed or original image
POST - /api/convert - convert image
GET  - /api/convert/{convertedID}?original={value} - get/download converted or original image
~~~

## Testing
Running test:
```
make test
```
Code coverage: To generate code coverage report, execute:
```
go test -coverprofile=cover.out
```
This should print the following after running all the tests.
```
coverage: 77.9% of statements
ok      github.com/alisavch/image-service/internal/apiserver    0.081s
```

You can also save the results as HTML for more detailed code view of the coverage.
```
go tool cover -html=cover.out -o coverage.html
```

This will generate a file called coverage.html. The coverage.html is provided in the repo which is pre-generated.