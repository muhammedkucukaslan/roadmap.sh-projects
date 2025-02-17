## Caching Proxy Server
This is a sample solution for [caching proxy](https://roadmap.sh/projects/caching-server) . It serves as a basic solution for caching HTTP requests and responses, using an in-memory map to store the cache. The challenge is from [roadmap.sh](https://roadmap.sh/golang/projects).

## Features
- In-Memory Caching:
  - The server uses an in-memory map to store cached responses. This approach allows for fast lookups and reduces the time needed to retrieve cached data.

- Proxy Functionality:
  - The server forwards client requests to the target server and caches the responses. If the same request is made again, the server returns the cached response, saving the time and resources of making a new request to the target server.

- Simple Design:
  - This is a minimalistic implementation aimed at demonstrating the core concepts of a caching proxy server. It is not suitable for production use without further enhancements.

## Usage

### Start the server:

```bash
go run . --port <PORT> --origin <ORIGIN_URL>
```
### To clear the cache::

```bash
go run . --clear-cache
```
