# Backend Microservices with GoGin Gateway and gRPC

## 📌 Overview

โปรเจกต์นี้เป็นระบบ Backend ที่พัฒนาด้วยสถาปัตยกรรม Microservices โดยมี Gateway ที่พัฒนาโดยใช้ **GoGin** เพื่อรับ RESTful API จากฝั่ง Frontend และติดต่อกับแต่ละ Microservice ผ่าน **gRPC**

> ระบบใช้ **Docker Compose** ในการจำลองฐานข้อมูลและบริการต่าง ๆ

---

## 🧱 System Architecture

- **Gateway**: GoGin (RESTful API Gateway) – รันที่ `port 8081`
- **Microservices (gRPC)**:
  - `auth-service` – รันที่ `port 50051`
  - `users-service` – รันที่ `port 50052`
  - `reset-password-service` – รันที่ `port 50053`

### Flow การทำงาน

1. ผู้ใช้เรียก API ผ่าน Gateway (REST)
2. Gateway แปลงเป็น gRPC และส่งต่อไปยัง Microservice ที่เกี่ยวข้อง
3. Microservice ตอบกลับผ่าน gRPC แล้ว Gateway แปลงกลับเป็น REST response

---

## 📦 Technologies

- Go (Gin, gRPC)
- MongoDB (ผ่าน Docker)
- Docker & Docker Compose

---

## 🚀 API Documentation

### 🔐 Authentication APIs

#### `POST /register` - ลงทะเบียนผู้ใช้ใหม่
```json
Request:
{
  "email": "test@example.com",
  "password": "securepassword"
}
Response:
{
  "message": "User registered successfully"
}
