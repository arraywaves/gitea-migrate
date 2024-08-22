# Gitea Migrate

`Work in progress.`

Gitea-Migrate is a Go application that creates an endpoint for GitHub webhooks to automatically mirror repositories to a Gitea instance. This tool simplifies the process of keeping your Gitea repositories in sync with their GitHub counterparts.

**Note: This project is currently a work in progress. While it should function as described, currently anyone could access your endpoint.**

## Features

- Listens for GitHub webhook events
- Automatically creates mirror repositories in Gitea
- Supports private repositories
- Mirrors wiki, labels, issues, pull requests, and releases

## Prerequisites

- Go 1.22 or higher.
- Access to a Gitea instance.
- GitHub account with repositories you want to mirror.

## Setup

1. Clone this repository:
   ```bash
   git clone https://github.com/arraywaves/gitea-migrate.git
   cd gitea-migrate
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the root directory with the following content:
   ```env
   GITEA_API_URL=https://your-gitea-instance.com/api/v1
   GITEA_USER=your-gitea-username
   GITEA_TOKEN=your-gitea-access-token
   GITHUB_USER=your-github-username
   GITHUB_TOKEN=your-github-personal-access-token
   ```

   Replace the placeholder values with your actual credentials.

4. Build the application:
   ```bash
   go build -o gitea-migrate
   ```

5. Run the application:
   ```bash
   ./gitea-migrate
   ```

   The server will start on port **8080** by default, you can change the port with the PORT environment variable.

   *If you're using https locally set the PORT variable to 443 and change the method ListenAndServe() to ListenAndServeTLS("cert.pem", "key.pem") in **main.go**.*

   *To generate a local SSL certificate interactively you can use **openssl** in your project directory:*
   `openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem`

## Usage

1. Host **gitea-migrate** on your server (preferably on the same domain as Gitea).

2. When the endpoint is triggered, Gitea Migrate will automatically create a mirror repository in your Gitea instance.

3. You can test the webhook locally using the provided script:
   ```bash
   curl -X POST http://localhost:8080/migrate-webhook \
      -H "Content-Type: application/json" \
      -d '{
        "repository": {
          "name": "repository-name",
          "clone_url": "https://github.com/your-username/repository-name.git"
        }
      }'
   ```

   Make sure to modify the script with your repository details before running.

## Project Structure

```
gitea-migrate/
├── api/
│   ├── handlers.go
│   └── router.go
├── logic/
│   └── mirror.go
├── main.go
├── go.mod
├── go.sum
└── .env
```

## Security Note

**Important:** The current version of this project does not yet implement origin restriction for incoming webhook requests. This means that, in its current state, the endpoint could potentially process requests from unknown sources. A planned improvement is to add security measures to restrict and validate the origin of incoming webhook requests.

## Future Improvements

- Implement origin restriction for enhanced security
- Add more comprehensive error handling and logging
- Expand test coverage
- Implement rate limiting to prevent abuse

## License

This project is licensed under the MIT License.
