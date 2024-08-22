curl -X POST \
  https://api.github.com/user/hooks \
  -H "Authorization: token YOUR_PERSONAL_ACCESS_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  -d '{
  "name": "web",
  "active": true,
  "events": [
    "repository"
  ],
  "config": {
    "url": "https://your-gitea-instance:port/migrate-webhook",
    "content_type": "json"
  }
}'
