# DriveLite Backend

Go Gin + Supabase

## Run local
cp .env.example .env
go run cmd/server/main.go

## Environment variables
SUPABASE_URL=
SUPABASE_SECRET_KEY=
PORT=8080

## Project Tree 
backend/
├── server/
│   └── main.go                  # Entrypoint
│             
│
├── config/
│   └── config.go                # Load env, Supabase, Cloudinary config
│
├── internal/
│   ├── domain/                  # Entities thuần (không phụ thuộc framework)
│   │   ├── file.go
│   │   ├── folder.go
│   │   └── user.go
│   │
│   ├── dto/                     # Request/Response structs
│   │   ├── file_dto.go
│   │   ├── folder_dto.go
│   │   └── auth_dto.go
│   │
│   ├── repository/              # Tương tác DB (Supabase/PostgreSQL)
│   │   ├── interface.go         # Định nghĩa interface
│   │   ├── file_repo.go
│   │   ├── folder_repo.go
│   │   └── user_repo.go
│   │
│   ├── service/                 # Business logic
│   │   ├── file_service.go
│   │   ├── folder_service.go
│   │   ├── auth_service.go
│   │   └── share_service.go
│   │
│   ├── handler/                 # HTTP Handlers (controller)
│   │   ├── file_handler.go
│   │   ├── folder_handler.go
│   │   ├── auth_handler.go
│   │   └── share_handler.go
│   │
│   ├── middleware/
│   │   ├── auth.go              # JWT verify (Supabase JWT)
│   │   ├── cors.go
│   │   ├── ratelimit.go
│   │   └── logger.go
│   │
│   ├── pkg/
│   │   ├── cloudinary/          # Cloudinary wrapper
│   │   │   └── cloudinary.go
│   │   ├── supabase/            # Supabase client
│   │   │   └── client.go
│   │   └── response/            # Chuẩn hóa JSON response
│   │       └── response.go
│   │
│   └── router/
│       └── router.go            # Đăng ký routes + middleware
│
├── migrations/                  # SQL migration files
│   ├── 001_create_users.sql
│   ├── 002_create_folders.sql
│   └── 003_create_files.sql
│
├── .env
├── .env.example
├── go.mod
└── Dockerfile