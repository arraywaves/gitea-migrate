# Gitea Migrate

Gitea Migrate is a Go application that automatically mirrors GitHub repositories to a Gitea instance. It primarily uses a polling mechanism to periodically check for new GitHub repositories and mirror them to Gitea. The tool simplifies the process of keeping your Gitea repositories in sync with their GitHub counterparts.
This tool is designed to work with personal GitHub accounts and provides a flexible solution for maintaining a Gitea mirror of your GitHub repositories.

## Features

- Polling mechanism: Regularly checks for new GitHub repositories and mirrors them to Gitea.
- Configurable polling interval: Adjust the frequency of checks to suit your needs (note: Github has an API rate limit of 5000 requests per hour for Personal Access Tokens).
- Tracks mirrored repositories: Maintains a list of already mirrored repositories to avoid duplication.
- Recognizes pre-existing mirrors: Identifies and tracks repositories that were manually mirrored prior to using Gitea Migrate.
- Supports public and private repositories.
- Mirrors wiki, labels, issues, pull requests, and releases.

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
   PORT=set-port-number // default 8080
   POLLING_INTERVAL_MINUTES=set-time-in-minutes // default 60
   MIRROR_MODE=set-mirror-mode-option // default poll (see options below)
   ```

   Replace the placeholder values with your actual credentials.

   **Mirror Mode** has 3 options:
   `MIRROR_MODE=poll (or unset)`: Only use polling (default), set the `POLLING_INTERVAL_MINUTES` variable to change how often it checks your repositories.
   `MIRROR_MODE=webhook`: *In Development* Only use webhook (Github Organisations), exposes endpoint `/migrate-webhook` **Note: Currently insecure with no origin whitelist**.
   `MIRROR_MODE=both`: Use both polling and webhook *(future default)*.

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

-- Webhooks --
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

   Make sure to set `MIRROR_MODE` to `webhook` or `both` and modify the script with your repository details before running.

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
