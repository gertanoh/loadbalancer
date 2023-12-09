# Load Balancer Project

## Introduction

This project is a custom implementation of a load balancer in Go. It is designed to distribute incoming network traffic across multiple servers to ensure no single server bears too much demand. By spreading the requests across multiple servers, it increases responsiveness and availability of applications.
It is my extended solution of [challenge-load-balancer](https://codingchallenges.fyi/challenges/challenge-load-balancer/)

## Features

- **Traffic Distribution**: Evenly distributes client requests or network load efficiently across multiple servers.


## To-Do

- [ ] Makefile
- [ ] Automatic servers discovery
- [ ] Health Checks: Regularly checks the health of the backend servers to ensure traffic is sent only to the healthy ones (Active).
- [ ] Add support for sticky sessions.
- [ ] Implement SSL termination.
- [ ] Provide Docker container support.
- [ ] Scalability: Easily scales out by adding more servers to the pool.
