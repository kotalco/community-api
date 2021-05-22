# :fire: Kotal API

Kotal api server listens on port `3000` and responds to versioned api calls.

Running the api server against real k8s cluster requires:

- [kotal operator](https://github.com/kotalco/kotal) to deployed in the cluster
- api server to be deployed with correct role and role bindings

API server is listening to calls that has the form of **HTTP METHOD** `/api/{version}/{protocol}/{resource}/?{name}` for example:

- POST `/api/v1/ethereum/nodes` to create node
- GET `/api/v1/ethereum/nodes` to list all nodes
- GET `/api/v1/ethereum/nodes/my-node` to get node by name
- PUT `/api/v1/ethereum/nodes/my-node` to update node by name
- DELETE `/api/v1/ethereum/nodes/my-node` to delete node by name

## :rocket: Running the API server

### :floppy_disk: From Source Code

API server can run simple by cloning this repository, then:

```
go run main.go
```

To run the mocking API server, use the `MOCK=true` environment variable:

```
MOCK=true go run main.go
```

### :framed_picture: From Docker Image

To run the API server from the docker image:

```
docker run -p 3000:3000 -e MOCK=true kotalco/api:0.1
```

## :telephone_receiver: Sample cURL Calls

Create a new node:

```
curl -X POST -d '{"name": "my-node", "network": "mainnet", "client": "parity"}' -H 'content-type: application/json' localhost:3000/api/v1/ethereum/nodes
```

Get node by name:

```
curl localhost:3000/api/v1/ethereum/nodes/my-node
```

List all nodes:

```
curl localhost:3000/api/v1/ethereum/nodes
```

Delete node by name:

```
curl -X DELETE localhost:3000/api/v1/ethereum/nodes/my-node
```

Update node by name:

```
curl -X PUT -d '{"client": "geth"}' -H 'content-type: application/json' localhost:3000/api/v1/ethereum/nodes/my-node
```
