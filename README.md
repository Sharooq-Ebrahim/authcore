# AuthCore

AuthCore is a production-grade, scalable authentication and authorization service built in Go. It follows **Clean Architecture** principles to ensure maintainability, testability, and decoupling of business logic from infrastructure.

## 🏗️ Architecture

The project is structured into four main layers:

- **Domain Layer (`internal/domain`)**: Contains core entities, custom app errors, and repository/service interfaces. This is the heart of the system and has no dependencies on other layers.
- **Usecase Layer (`internal/usecase`)**: Orchestrates business logic by coordinating entities and repository interfaces.
- **Infrastructure Layer (`internal/infrastructure`)**: Concrete implementations of repositories (PostgreSQL), security (JWT, Bcrypt), and database configurations.
- **Delivery Layer (`internal/delivery`)**: Entry points to the system (HTTP Handlers, Middleware, DTOs).

## 🚀 Features

### 1. Authentication
- **Email/Password**: Secure registration and login using Bcrypt hashing.
- **OAuth2 (Google)**: Support for Google OAuth login. The system automatically links OAuth accounts to existing users by email or creates new users if they don't exist.
- **Extensible Provider System**: A decoupled OAuth adapter system (`internal/oauth`) allows adding new providers (GitHub, etc.) by implementing a simple interface.

### 2. Token Management
- **JWT (JSON Web Tokens)**: Issuance of short-lived Access Tokens and long-lived Refresh Tokens.
- **Refresh Flow**: Dedicated endpoint to rotate tokens using valid refresh tokens.
- **Security**: HS256 signing with configurable expiration and secrets.

### 3. Authorization & RBAC
- **Role-Based Access Control**: Support for roles like `admin` and `user`.
- **Protected Actions**: Administrative actions (like role assignment) are restricted to users with the `admin` role.

### 4. Client Management (B2B)
- **Multi-Client Support**: Ability to register external clients that consume the Auth service.
- **API Key Authentication**: Middleware that validates client credentials (`x-api-key`) for every authentication request.

## 🗄️ Database Structure

Powered by PostgreSQL and managed via `golang-migrate`.

- **`users`**: Core user data (email, password hash, role).
- **`oauth_accounts`**: Links third-party provider IDs to internal user records.
- **`refresh_tokens`**: Stores token metadata for revocation and session management.
- **`clients`**: Organizations or applications using AuthCore.
- **`client_credentials`**: API keys and secrets associated with clients.

## 🔑 Authentication Flow

1. **Client Identity**: Every request to `/register`, `/login`, or `/oauth` must include an `x-api-key`.
2. **User Login**: User provides credentials or completes OAuth dance.
3. **Internal Linking**: For OAuth, the system checks `oauth_accounts` or matches by email in `users`.
4. **Token Issuance**: System returns a JWT Access Token (short-lived) and a Refresh Token (long-lived).
5. **Authorization**: Subsequent requests use the Bearer token to identify the user and their permissions.

## 🛠️ Tech Stack

- **Language**: Go 1.25+
- **Framework**: Standard Library (`net/http`) for maximum performance and minimal bloat.
- **Database**: PostgreSQL with `pgx` driver.
- **Libraries**:
  - `golang-jwt/jwt`: Token management.
  - `golang.org/x/crypto`: Bcrypt hashing.
  - `golang.org/x/oauth2`: OAuth2 flow.
  - `joho/godotenv`: Configuration management.
  - `golang-migrate`: Schema migrations.

## 🚦 API Endpoints

### Auth
- `POST /register`: Register a new user (requires API key).
- `POST /login`: Login with email/password.
- `POST /refresh`: Exchange refresh token for new access/refresh pair.
- `GET /verify`: Validate an access token.
- `GET /profile`: Get current user profile.
- `POST /assign-role`: (Admin Only) Assign roles to users.

### OAuth
- `GET /oauth/{provider}`: Redirect to provider login (e.g., `/oauth/google`).
- `GET /oauth/{provider}/callback`: Handle provider callback and return JWTs.

### Client Management
- `POST /clients`: Create a new client.
- `GET /clients/{id}`: Get client details.
- `POST /clients/{id}/credentials`: Generate a new API Key for a client.
- `GET /clients/{id}/credentials`: List API keys for a client.
