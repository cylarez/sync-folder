```
 ____  _____    _    ____  __  __ _____                 _ 
|  _ \| ____|  / \  |  _ \|  \/  | ____|  _ __ ___   __| |
| |_) |  _|   / _ \ | | | | |\/| |  _|   | '_ ` _ \ / _` |
|  _ <| |___ / ___ \| |_| | |  | | |___ _| | | | | | (_| |
|_| \_\_____/_/   \_\____/|_|  |_|_____(_)_| |_| |_|\__,_|
```
# Demo Server Golang

## Requirements

Golang version >= 1.18

https://go.dev/doc/install

## How to print help

```
go run main.go -help
```

## How to build

```
go build .
```

## How to run

*For security reason, server will NOT run if you haven't defined an API Key.*

Set the API KEY using the `API_KEY` environment variable:
```
export API_KEY=12345
```

OPTIONAL: Server is listening at 127.0.0.1:8080 by default, but you can change it using the `SERVER_ADDRESS` environment variable:
```
export SERVER_ADDRESS=127.0.0.1:9999
```

Then run your build:
```
./server
```

Or just run the source code:
```
go run main.go
```

## How to build a docker image

```
docker build . -t demo-server
```

## How to run the docker image

```
docker run -p 8080:8080 -e API_KEY=12345 demo-server
```

To change the port you can change it using docker port config:
```
docker run -p 9090:8080 -e API_KEY=12345 demo-server
```

## How to run the unit-tests

```
 go test ./...
```