# GolangContainers

1. REST endpoint to post a JSON => {firstname, lastname}

2. REST endpoint to get a list of {firstname, lastname} with query to filter for the last 60s, 5mins or 1hr
   *array must be ordered by timestamp
   
Process

1. When the POST request is done, the JSON must be pushed to the NATS (https://nats.io/).

2. Implement a Consumer to read the Queue and combine the firstname and lastname into fullname and push it to a Redis.

3. Query to the Redis and return the list of users with query to filter for the last 60s, 5mins or 1hr
   *array must be ordered by timestamp, The Redis should have a TTL for the dataset inserted of not more than 1hr

## Setup

### Dependencies

Please clone project to your GOPATH: for exampele go/src/github.com/igorextrawest/GolangContainers

```bash
# install go dependency management tool 
go get -u github.com/golang/dep/cmd/dep
  
# download and install docker-compose
https://docs.docker.com/compose/install

# install dependencies
dep init
```
### Launch
```bash
docker-compose up
```

### Examples 
for testing on localhost: port 8080 

POST http://127.0.0.1:8080/v1/user

Request 
```bash
{
	"firstname":"Igor",
	"lastname":"Shcherbyna",
	"address":"UA",
	"gender":"male"
}
```

GET http://127.0.0.1:8080/v1/user?query=1m

Response 
```bash
{
    "values": [
        "Igor Shcherbyna"
    ]
}
```
