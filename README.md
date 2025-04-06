# DuelTerm
Extreme shooting in console between '@' and '&amp;" without any profit.
 ![logo](https://github.com/user-attachments/assets/f49ef782-0272-498a-a742-509ae2f26ece)

 ## How to Run

### 1. Run the Server
```bash
go run cmd/server/main.go
```

### 2. Run Two Clients

#### Option A — Localhost (same PC)
```bash
go run cmd/client/main.go
go run cmd/client/main.go
```

#### Option B — LAN
On second machine:
```bash
go run cmd/client/main.go 192.168.0.X:8080
```

#### Option C — Internet
Use `ngrok`:
```bash
ngrok tcp 8080
go run cmd/client/main.go 0.tcp.ngrok.io:PORT
```

---
