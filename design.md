# Design and building a loadbalancer in Go
The features to be supported are : 
- Redirect traffic to backend servers
- Health check servers 
- Handle servers going offline and coming back online
- Automatic servers discovery

Load balancer is a server listenning on port 8080 for incoming connections.

## Basic server with HTTP keep alive
The load balancer is waiting for incoming requests. HTTP keep alive is setup for connections between servers and load balancers and between the lb and the clients with gracefull shutdown.

## Servers discovery
How will the load balancer detect the new servers ? How will the servers know that a server is alive or down ?

I see multiple alternatives. I can open a port on the load balancer and expect servers to register on this port to announce their presence. The load balancer periodically do health checks to detect servers presences.
For this alternative, I can use etcd or zookeeper to make it totally distributed.
The other alternative is to use serf for servers discovery and health checking. Serf servers discovery includes already health checks(by triggerring events when an agent leaves the cluster). It is also event based, to handle events. Serf is eventual consistent. We do not need strong consistency for our use case.
the servers and the load balancers will form a cluster and the load balancer will receive notifications when a server joins or leaves the cluster.
