# Gitea Migrate

`Work in progress.`

Gitea-Migrate is a Go application that creates an endpoint for GitHub webhooks to automatically mirror repositories to a Gitea instance. This tool simplifies the process of keeping your Gitea repositories in sync with their GitHub counterparts.

**Note: This project is currently a work in progress. While it should function as described, currently anyone could access your endpoint.**

## Features

- Checks GitHub account for unmirrored repositories.
- Automatically creates mirror repositories in Gitea from new Github repositories.
- Keeps track of mirrored repositories in a new file `./mirrored_repos.json`.
- Supports private repositories (via a Github Personal Access Token).
- Mirrors wiki, labels, issues, pull requests, and releases.
- *In Development* Option for a Github webhook endpoint.

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
   PORT=port-number // default 8080
   POLLING_INTERVAL_MINUTES=time-in-minutes // default 60
   MIRROR_MODE=mirror-mode // default poll (see options below)
   ```

   Replace the placeholder values with your actual credentials.

   **Mirror Mode** has 3 options:
   `MIRROR_MODE=poll`: Only use polling (default), set the `POLLING_INTERVAL_MINUTES` variable to change how often it checks your repositories.
   `MIRROR_MODE=webhook`: Only use webhook (Github Organisations), exposes endpoint `/migrate-webhook` **Note: currently insecure**.
   `MIRROR_MODE=both or unset`: Use both polling and webhook (future default).

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

2. When the server starts, Gitea Migrate will automatically scan and migrate your repos to Gitea.

*For Github organisations there is the endpoint `/migrate-webhook` which is enabled in webhook mode (see Environment Variables above). When triggered, Gitea Migrate will automatically create a mirror repository in your Gitea instance for any newly created repos.*

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
│   └── github_poller.go
│   └── mirror.go
├── scripts/api/
│   └── test-migrate-endpoint.sh
├── main.go
├── go.mod
├── go.sum
└── .env
```

## Security Note

**Important:** The current version of this project does not yet implement origin restriction for incoming webhook requests. This means that in its current state, the webhook endpoint could potentially process requests from unknown sources. A planned improvement is to add security measures to restrict and validate the origin of incoming webhook requests. For now please ensure `MIRROR_MODE` is set to `poll`.

## Future Improvements

- Implement origin restriction for enhanced security (webhook mode)
- Implement rate limiting to prevent abuse (webhook mode)
- Expand test coverage

## License

This project is licensed under the MIT License.
