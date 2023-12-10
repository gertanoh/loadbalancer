# Load Balancer Project

## Introduction

This project is a custom implementation of a load balancer in Go. It is designed to distribute incoming network traffic across multiple servers to ensure no single server bears too much demand. By spreading the requests across multiple servers, it increases responsiveness and availability of applications.
It is my extended solution of [challenge-load-balancer](https://codingchallenges.fyi/challenges/challenge-load-balancer/).

[Serf](https://www.serf.io/) is used for decentralised server discovery. Heatlh checks is performed with Serf membership builtin gossip protocol

## Features

- **Traffic Distribution**: Evenly distributes client requests or network load efficiently across multiple servers.
- **Fault Tolerance**: Automatically reroutes traffic in case of server failure.
- **Scalability**: Easily scales out by adding more servers to the pool.
- **Health Checks**: Regularly checks the health of the backend servers to ensure traffic is sent only to the healthy ones.

## To-Do

- [ ] Add support for sticky sessions.
- [ ] Provide Docker container support.
- [ ] Clean up serf printing too many debug logs.
