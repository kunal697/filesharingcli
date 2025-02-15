# FileShare

<img src="https://socialify.git.ci/kunal697/filesharingcli/image?custom_description=A+secure+file-sharing+backend+service+that+uses+GitHub+as+storage+and+PostgreSQL+for+metadata+management.&description=1&font=Inter&language=1&name=1&owner=1&pattern=Solid&stargazers=1&theme=Dark" alt="filesharingcli" width="640" height="320" />

## Features

- 🔐 Secure site creation and management
- 📤 File upload to GitHub repository
- 📥 File download with authentication
- 🔑 JWT-based authentication
- 🗄️ PostgreSQL database integration
- 🚀 RESTful API endpoints

## Prerequisites

- Go 1.20 or higher
- PostgreSQL database (Neon DB)
- GitHub Account with Personal Access Token

## Environment Variables

The following environment variables are required in `.env`:

```env
DATABASE_URL="your_neon_db_url"
GITHUB_TOKEN="your_github_token"
```

## API Endpoints

### Sites
- `POST /createsite` - Create a new site
  ```json
  {
    "site_name": "your_site_name",
    "password": "your_password"
  }
  ```

- `GET /site/:site_name?password=site_password` - Get site details and files
- `GET /sites` - List all sites

### Files
- `POST /upload/:site_name` - Upload a file (multipart/form-data)
- `GET /getfile/:id` - Download a file (requires auth token)

## Project Structure

```
Filesharing/
├── internal/
│   ├── db/         # Database connection and models
│   ├── handlers/   # Request handlers
│   ├── models/     # Data models
│   ├── routes/     # API routes
│   └── utils/      # Utility functions
├── .env            # Environment variables
├── .gitignore
├── go.mod
└── main.go
```
