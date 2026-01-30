# Digital Wallet Backend Service

Backend service untuk aplikasi dompet digital. User bisa register, login, cek saldo, dan tarik dana.

## Tech Stack

- Go 1.21+
- Gin Framework
- MySQL 8.0+
- JWT Authentication
- Raw SQL (no ORM)

## Instalasi

### 1. Clone & Install Dependencies

```bash
git clone <repository-url>
cd wallet-test
go mod download
```

### 2. Setup Database

```bash
# Login ke MySQL
mysql -u root -p

# Buat database
CREATE DATABASE wallet_db;
exit
```

### 3. Konfigurasi Environment

```bash
cp .env.example .env
```

Edit `.env`:
```env
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_NAME=wallet_db
DATABASE_USER=root
DATABASE_PASSWORD=your_password

JWT_SECRET=your-secret-key
SERVER_PORT=8080
```

### 4. Run Migrations

```bash
cat migrations/*.sql | mysql -u root -p wallet_db
```

### 5. [Optional] Seed Test User

```bash
mysql -u root -p wallet_db < migrations/005_seed_test_user.sql
```

Test user:
- Email: `test@example.com`
- Password: `password123`
- PIN: `123456`

### 6. Run Server

```bash
go run cmd/main.go
```

Server jalan di `http://localhost:8080`

## Troubleshooting

**Database connection error**
- Cek MySQL sudah running
- Cek credentials di `.env`

**Table not found**
- Run migrations: `cat migrations/*.sql | mysql -u root -p wallet_db`

**Invalid token**
- Token expire 24 jam, login lagi

**Insufficient balance**
- Cek saldo dulu: `GET /api/balance`
