# DNS Resolver

A simple DNS resolver implemented in Go, based on [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035).

## Introduction

DNS (Domain Name System) is the internet's phonebook. It translates human-readable domain names (like www.google.com) into IP addresses (like 142.250.192.174) that computers use to identify each other on the network.

A DNS resolver is a program that queries DNS servers to resolve domain names into IP addresses. This project implements a basic DNS resolver from scratch in Go, demonstrating the fundamentals of DNS communication.

## Technical Overview

This DNS resolver implementation follows the specifications outlined in RFC 1035 (https://datatracker.ietf.org/doc/html/rfc1035), which defines the specification for the DNS protocol.

Key technical aspects of this implementation include:

1. **DNS Message Format**: The resolver constructs and parses DNS messages according to the standard format.

2. **Query Construction**: The resolver builds DNS queries with appropriate flags, such as Recursion Desired (RD) for recursive lookups.

3. **UDP Communication**: DNS typically uses UDP for queries and responses. This resolver establishes UDP connections to communicate with DNS servers.

4. **Response Parsing**: The resolver decodes the DNS server's response, extracting relevant information such as IP addresses from resource records.

## How it works

1. The resolver constructs a DNS query message for a given domain name.
2. It sends this query to a DNS server (default: 8.8.8.8, Google's public DNS) via UDP.
3. The server responds with the IP address(es) for the domain.
4. The resolver parses the response, decoding the various sections of the DNS message.
5. Finally, it displays the results, either in a human-readable format or as raw data.

## Usage

1. Clone the repository:

   ```
   git clone https://github.com/prabin-acharya/dns-resolver.git
   cd dns-resolver
   ```

2. Run the resolver:

   ```
   go run main.go <domain-name>
   ```

   Example:

   ```
   go run main.go www.example.com
   ```

3. To see the raw DNS response, use the `--raw` flag:
   ```
   go run main.go --raw www.example.com
   ```

## References

- https://datatracker.ietf.org/doc/html/rfc1035 (RFC 1035: Domain Names - Implementation and Specification)
- https://www.cloudflare.com/learning/dns/what-is-dns
