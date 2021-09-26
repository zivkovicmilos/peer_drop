![Banner](.github/banner.png)

## peer_drop.

peer_drop is a safe peer-to-peer decentralized file sharing solution. It leverages enterprise grade security and
verification practices to make sure files are safely distributed from one machine to the other.

Users can create _workspaces_ and share their workspace mnemonic so others can join.

## Installation

### Client server install

The latest version of golang is recommended.

To install the client server dependencies, in the `peer_drop/client` directory run:

```bash
go mod download
````

Finally, to build the executable, run:

```bash
go build main.go
````

### Web UI install

Yarn is used as the package manager, so make sure it is installed on your system. To install the frontend dependencies,
in the `peer_drop/web-ui` directory run:

```bash
yarn install
````

## How to run 📝

There are 2 types of nodes in a peer_drop system:

* **Client node**: handles regular UI actions, connects to other client nodes, does file transfer and handshaking
  protocols
* **Rendezvous node**: serves as a bootstrap node for client node discovery, doesn't take part in file transfer

In order for _client nodes_ to discover each other, a _rendezvous node_ needs to be already up and running.

### Deploying a rendezvous node

A rendezvous node can be deployed using the following command:

```bash
go run main.go --directory=app_data_rendezvous --host=46.101.229.130 --grpc-port=10001 --libp2p-port=10002 --rendezvous 
```

Flag description:

- `directory` flag specifies the working directory for the rendezvous node data
- `grpc-port` flag specifies the available gRPC port for stream protocols
- `libp2p-port` flag specifies the available libp2p port the libp2p host instance can use
- `rendezvous` flag specifies that this node is starting in _rendezvous_ mode
- `host` flag specifies the host that libp2p should use, if running on a cloud server this is the public IP

When starting the rendezvous node, take note of the `multiaddr` the output provides, for example:

```
Rendezvous node started with ID: /ip4/46.101.229.130/tcp/10002/p2p/QmdBhgtDwRVkJJ5DbfVHUm4FuhcmpDQ4Qebx66x48MmYu8
```

This ID needs to be used when configuring / starting the client nodes.

At least one rendezvous node is required for normal operation. To start multiple rendezvous nodes, pass in
the `--rendezvous-node` flag for every new node, with the `multiaddr` of the running rendezvous node.

### Deploying a client node

Before running the client node, make sure it has a set rendezvous node address. In `peer_drop/client/config/config.go`,
set the rendezvous `multiaddr`s for the `DefaultRendezvousNodes` value:

```go
// Default rendezvous nodes that are already up and running
var (
DefaultRendezvousNodes = []string{
"/ip4/46.101.229.130/tcp/10002/p2p/QmdBhgtDwRVkJJ5DbfVHUm4FuhcmpDQ4Qebx66x48MmYu8", // DigitalOcean rendezvous
}
)
```

A client node can be started using the following command:

```bash
go run main.go
```

The command accepts all the previously mentioned flags for setting up the directory and host / port information. The
default host / HTTP port is `localhost:5000`.

### Starting the Web interface

Before starting the web interface, create a file named `.env` inside the `peer_drop/web-ui` directory. Populate the file
so the frontend knows where the client server is running:

```
REACT_APP_CLIENT_API_BASE_URL=<host/port>
```

An example of a correct `.env` file for a client server running on `localhost:5000` would be:

```
REACT_APP_CLIENT_API_BASE_URL=http://localhost:5000
```

### Starting multiple client nodes on the same machine

When starting multiple client nodes, or multiple client nodes with a rendezvous node, be sure to give each starting
instance a **separate directory and port** through the mentioned flags when starting the clients.

The UI can also be started multiple times, just make a copy of the *web-ui* folder to a different location, and in the
new `.env` file add:

```
REACT_APP_CLIENT_API_BASE_URL=<host/port>

PORT=<port number>

```

The default `PORT` value is `3000`. The UI's need to be running on different ports, and reference the correct client
servers in order for the system to function correctly.

Example `.env` of a UI that's referencing a client server running on `localhost:5000`:

```
REACT_APP_CLIENT_API_BASE_URL=http://localhost:5000

PORT=5000

```

Each individual client node / UI needs to be started separately.