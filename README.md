# GolangContainers

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
