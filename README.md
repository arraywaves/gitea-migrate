# Gitea Migrate

Gitea Migrate is a simple Go server application that automatically mirrors GitHub repositories to a Gitea instance by periodically checking for new GitHub repositories and migrating them to Gitea, you could run this just once from your local environment if you don't need to check periodically for new Github repositories. This tool simplifies the process of keeping your Gitea repositories in sync with their GitHub counterparts and is designed to work with personal GitHub accounts, providing a flexible solution for maintaining a Gitea backup of your GitHub repositories. A webhook endpoint is in development for Github organisations.

## Features

- Regularly checks for new GitHub repositories and mirrors them to Gitea.
- Adjust the frequency of checks to suit your needs (note: Github has an API rate limit of 5,000 requests per hour for Personal Access Tokens - this service makes 1 request check).
- Maintains a list of already mirrored repositories to avoid duplication.
- Non-destructive: if a repository with the same name already exists on Gitea it won't be overwritten.
- Supports public and private repositories.
- Mirrors wiki, labels, issues, pull requests, and releases.
- Enables Gitea mirroring by default (periodically sync from Github to Gitea), this can be turned off.

## Future development

- *In Development* Webhook support: An endpoint for GitHub webhooks is in development, which will allow for immediate mirroring when new repositories are created on GitHub Organisations.

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

3. Create an `.env` file in the root directory with the following content:
   ```env
   GITEA_API_URL=https://your-gitea-instance.com/api/v1
   GITEA_USER=your-gitea-username
   GITEA_TOKEN=your-gitea-access-token
   GITHUB_USER=your-github-username
   GITHUB_TOKEN=your-github-personal-access-token
   GH_RATE_LIMIT=set-rate-limit # default 4990 (this is highly unlikely to be reached and is just a precaution)
   PORT=set-port-number # default 8080
   POLLING_INTERVAL_MINUTES=set-time-in-minutes # default 60
   MIGRATE_MODE=set-mirror-mode-option # default poll (see options below)
   ENABLE_MIRROR=set-mirror-mode # default true (if set to false, your new repo won't sync with its Github counterpart. To enable mirror mode later you'll need to delete your Gitea repo and re-run Gitea Migrate with `ENABLE_MIRROR` set to true)
   ```

   Replace the placeholder values with your actual credentials - restart the server if it's already running to use any updated settings.

   ### **Mirror Mode** has 3 options:

   - `MIGRATE_MODE=poll (or unset)`: Only use polling (default), set the `POLLING_INTERVAL_MINUTES` variable to change how often it checks your repositories.
   - `MIGRATE_MODE=webhook`: *In Development* Only use webhook (Github Organisations), exposes endpoint `/migrate-webhook` **Note: Currently insecure with no origin whitelist**.
   - `MIGRATE_MODE=both`: Use both polling and webhook *(future default)*.

4. Build the application:
   ```bash
   go build -o gitea-migrate ./cmd/server
   ```

5. Run the application:
   ```bash
   ./gitea-migrate
   ```

   The server will start on port **8080** by default, you can change the port with the PORT environment variable.

   *If you're using https locally set the PORT variable to 443 and change the method ListenAndServe() to ListenAndServeTLS("cert.pem", "key.pem") in **main.go**.*

   *To generate a local SSL certificate interactively you can use **openssl** in your project directory:*
   `openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem`

   The included Dockerfile should help you get up and running with Docker.

## Usage

1. Host **gitea-migrate** on your server.

2. When the server starts, Gitea Migrate will automatically scan and migrate your repos to Gitea.

-- **Webhooks** --
*For Github organisations there is the endpoint `/migrate-webhook` which is enabled in webhook mode (see Environment Variables above). When triggered, Gitea Migrate will automatically create a mirror repository in your Gitea instance from the request payload (`name`, `clone_url`).*

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

   Make sure to set `MIGRATE_MODE` to `webhook` or `both` and modify the script with your repository details before running.

   -- **Health Check** --
   Health checks can be made at the `/health` endpoint when the server is runnning.
   `curl -v https://your-domain.com/health`

## Project Structure

```
gitea-migrate/
├── cmd/
│   ├── server/
│   │   └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers.go
│   │   └── router.go
│   ├── config/
│   │   └── config.go
│   ├── core/
│   │   ├── giteaservice.go
│   │   ├── githubservice.go
│   │   ├── interfaces.go
│   │   ├── poller.go
│   │   └── ratelimiter.go
├── pkg/
│   ├── models/
│   │   └── repository.go
├── scripts/api/
│   └── test-migrate-endpoint.sh
├── .env
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
├── main.go
└── README.md
```

## Security Note

**Important:** The current version of this project does not yet implement origin restriction for incoming webhook requests. This means that in its current state, the webhook endpoint could potentially process requests from unknown sources. A planned improvement is to add security measures to restrict and validate the origin of incoming webhook requests. For now please ensure `MIGRATE_MODE` is set to `poll`.

## Future Improvements

- Handle pagination for larger Github accounts (which will increase Github requests by 1 per 100 repositories).
- Implement test coverage.
- Enable more migration options.
- Implement origin restriction (CSRF) for security in webhook mode.

## License

This project is licensed under the [MIT License](LICENSE).
