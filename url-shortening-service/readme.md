# URL Shortening Service
Sample solution for the [url shortening service](https://roadmap.sh/projects/url-shortening-service) challenge from [roadmap.sh](https://roadmap.sh/golang/projects).


## Overview
This project is a url shortening service built using Go that allows users to shorten long urls.

## Requirements
- Go 1.16+
- PostgreSQL

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/muhammedkucukaslan/roadmap.sh-projects.git
   cd roadmap.sh-projects/url-shortening-service
   ```

2. Set up the database:
   - Ensure PostgreSQL is running
   - Set up environment variables (see Configuration section)

3. Run the application:
   ```bash
   make run
   ```
The application will run on port 3000.

##Example

```bash
curl -X POST http://localhost:3000/shorten -H "Content-Type: application/json" -d '{"url": "https://www.google.com"}'

Response:
{   
    "shortCode": "oaihIUL"
}
```


```bash
curl http://localhost:3000/shorten/1

Response:
{
    "id": "13",
    "url": "https://www.google.com",
    "shortCode": "oaihIUL",
    "accessCount": 4,
    "createdAt": "2025-02-22T14:11:29.102007Z",
    "updatedAt": "2025-02-22T14:13:31.054372Z"
}
```


