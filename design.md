# Design and building a loadbalancer in Go
The features to be supported are : 
- Redirect traffic to two or more servers
- Health check servers 
- Handle servers going offline and coming back online

Load balancer is a server listenning on port 8080 for incoming connections.

How will the load balancer detect the new servers ? 
How will the servers know that a server is alive or down  .
Servers annouces themselves to the load balancer with services registry, what tool for services discovery ? serf vs etcd vs zookeper 
health check is done actively by the load balancer