# DriveLite Backend

Go Gin + Supabase

## Run local

```bash
cp .env.example .env
go run server/main.go
```

## Environment variables

```
SUPABASE_URL=
SUPABASE_SECRET_KEY=
PORT=8080
APPLICATION_PASSWORD=
HOST_EMAIL=
SMTP_HOST=
CLOUDINARY_NAME=
CLOUDINARY_API_KEY=
CLOUDINARY_SECRET_KEY=
```

## Project Tree

```
backend/
├── server/
│   └── main.go                      # Entrypoint
│
├── config/
│   ├── config.go                    # Load env config
│   └── supabase-connection.go       # Supabase client connection
│
├── internal/
│   ├── handler/                     # HTTP Handlers (controller)
│   │   ├── handle_account.go
│   │   ├── handle_document.go
│   │   ├── handle_find.go
│   │   ├── handle_path.go
│   │   └── handle_usage_bucket.go
│   │
│   ├── middleware/
│   │   ├── auth_jwt.go              # JWT verify (Supabase JWT)
│   │   └── set_jwt.go
│   │
│   ├── repository/                  # Tương tác DB (Supabase/PostgreSQL)
│   │   ├── bucket_usage.go
│   │   ├── file_repo.go
│   │   ├── folder_repo.go
│   │   ├── get_path.go
│   │   ├── star_repo.go
│   │   └── trash_repo.go
│   │
│   ├── router/
│   │   └── router.go                # Đăng ký routes + middleware
│   │
│   └── services/                    # Business logic
│       ├── account_service.go
│       ├── create_url.go
│       ├── file_service.go
│       ├── find_service.go
│       ├── folder_service.go
│       ├── login.go
│       ├── password_service.go
│       ├── register.go
│       └── share_service.go
│
├── .env
├── .env.example
├── .gitignore
├── go.mod
└── go.sum
```