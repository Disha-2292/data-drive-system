# data-drive-system


A simplified version of Google Drive built using **Go (Gin Framework)** for the backend and **Refine (React + Ant Design)** for the frontend. This system supports user authentication, file and folder CRUD operations, and allows users to create nested folders of unlimited depth.

---

## 🚀 Features

### 🔐 User Authentication
- Register and login functionality.
- Authenticated users can manage their own files and folders.

### 🗂️ CRUD Operations
- Create, Read, Update, Delete for:
  - Files
  - Folders
- Each file/folder is associated with a user.

### 🧩 Nested Folders
- Supports folder structures of unlimited depth (e.g., folders inside folders like Google Drive).

### 🧠 Role-Based Access Control (RBAC)
- Admin users have access to all files and folders.
- Regular users can only manage their own files.

### 💬 File Sharing & Permissions
- Share files with other users.
- Manage access permissions (e.g., view-only, edit).

---

## 🛠️ Technologies Used

| Layer       | Technology                           |
|------------|---------------------------------------|
| Backend     | [Go](https://golang.org/), [Gin](https://gin-gonic.com/), [GORM](https://gorm.io/), MySQL |
| Frontend    | [Refine](https://refine.dev/) (React + Ant Design) |
| Auth        | JWT-based authentication |
| DB Access   | GORM + raw SQL for SELECTs |
| File Upload | Local file system storage (basic implementation) |

---

## ⚙️ Setup Instructions

### ✅ Prerequisites
- Go 1.18+
- MySQL Server
- Node.js (for frontend)

---

### 🧩 Backend Setup

```bash
cd backend
go mod tidy

Created a .env file in backend folder with:

DB_USER=disha
DB_PASSWORD=pasted in .env file
DB_NAME=data_drive
DB_PORT=3306
SECRET_KEY=pasted in .env file

To run the Server-> go run main.go or go run .

🌐 Frontend Setup

cd frontend
npm install
npm run dev

Frontend runs at: http://localhost:5173

📫 API Endpoints 

POST | /auth/register | Register new user | auth not required
POST | /auth/login | Login existing user | auth not required
Auth required to all rest of apis
GET | /api/me | Get logged-in user info | ✅
POST | /api/files/create | Create file/folder | ✅ (Admin)
GET | /api/files/all | Get all files | ✅ (Admin)
GET | /api/files | Get logged-in user's files | ✅ (User)
GET | /api/files/:id | Get file by ID | ✅ (User)
PUT | /api/files/:id | Update file | ✅ (Admin)
DELETE | /api/files/:id | Delete file | ✅ (Admin)
POST | /api/files/upload | Upload file | ✅ (User)
GET | /api/files/download/:id | Download file | ✅ (User)
POST | /api/files/:id/share | Share file with another user | ✅ (User)
GET | /api/files/:id/permissions | Get file permissions | ✅ (User)
PUT | /api/files/:id/permissions | Update file permissions | ✅ (User)
GET | /api/files/search | Search user's files | ✅ (User)


👩‍💻 Author

Disha Gohil

