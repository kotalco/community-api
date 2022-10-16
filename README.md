# :fire: Kotal API

Kotal API is used by Kotal Dashboard to manage [Kotal Operator](https://github.com/kotalco/kotal) custom resources like IPFS peers, Polkadot validator nodes, Chainlink nodes, Ethereum nodes ... etc.

Kotal API server listens on port `5000` and responds to versioned API calls.

API server port can be changed using `KOTAL_API_SERVER_PORT` environment variable.

Running the API server against real k8s cluster requires:

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

**NOTE:** This command will run the API server and expects an actual k8s cluster with kubeconfig available in the default kubeconfig dir.

To run the mocking API server, use the `MOCK=true` environment variable:

```
MOCK=true go run main.go
```

### :framed_picture: From Docker Image

To run the API server from the docker image:

```
docker run -p 3000:3000 -e MOCK=true kotalco/community-api:develop
```

## :telephone_receiver: Sample cURL Calls

Create a new node:

```
curl -X POST -d '{"name": "my-node", "network": "mainnet", "client": "besu"}' -H 'content-type: application/json' localhost:3000/api/v1/ethereum/nodes
```

Get node by name:

```
curl localhost:3000/api/v1/ethereum/nodes/my-node
```

List all nodes:

```
curl localhost:3000/api/v1/ethereum/nodes
```

Update node by name:

```
curl -X PUT -d '{"rpc": true}' -H 'content-type: application/json' localhost:3000/api/v1/ethereum/nodes/my-node
```

Delete node by name:

```
curl -X DELETE localhost:3000/api/v1/ethereum/nodes/my-node
```
