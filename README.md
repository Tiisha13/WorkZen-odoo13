# 🚀 WorkZen HRMS

A modern, full-stack **Human Resource Management System** built with **Go (Fiber)**, **MongoDB**, **Next.js 16**, and **TypeScript**. WorkZen provides comprehensive HR operations including attendance tracking, leave management, payroll processing, and employee management with a beautiful, responsive UI.

## ✨ Features

### 🔐 Authentication & Authorization

- JWT-based secure authentication
- Role-based access control (RBAC) with 5 role levels
- Email verification system
- Password encryption with bcrypt
- Multi-tenancy support (SaaS architecture)

### 👥 User & Employee Management

- Complete CRUD operations for employees
- Auto-generated employee codes and passwords
- Manager hierarchy support
- Department assignment
- Bank details management for payroll
- User profile with avatar support
- Role-based data visibility (lower roles cannot see higher roles)

### ⏰ Attendance Tracking

- Real-time check-in/check-out system
- Automatic working hours calculation
- Daily attendance records
- Attendance reports and analytics
- Admin/SuperAdmin excluded from attendance (they are managers, not employees)
- Status tracking (Present, Absent, On Leave)

### 🏖️ Leave Management

- Leave application workflow
- Multiple leave types (Sick, Casual, Vacation)
- Leave approval system for HR/Admin
- Leave balance tracking
- Leave status (Pending, Approved, Rejected)
- Calendar integration for leave dates

### 💰 Payroll & Salary

- Automated salary structure calculation
- Monthly payrun generation
- Salary components (Basic, HRA, Allowances, PF, Tax)
- Payroll processing with deductions
- Salary slips generation
- Bank account validation

### 📄 Document Management

- Secure file upload system
- Multiple document categories (Resume, ID Proof, Payslip, Policy, Report)
- File preview (Images, PDFs, Videos, Audio)
- Download functionality with authentication
- Organized storage: `/assets/uploads/{company}/{category}/{year}/{month}/`
- Auto-create upload directories on server start

### 🏢 Department Management

- Department CRUD operations
- Employee count per department
- Department-wise statistics

### 📊 Interactive Dashboard

- Real-time statistics and metrics
- Department-wise attendance charts (Bar charts)
- Monthly attendance trends (Area charts with gradients)
- Today's attendance distribution (Pie charts)
- Leave type statistics (Stacked bar charts)
- Attendance rate calculation
- Role-based filtering (shows only relevant employee data)
- Responsive 3-column layout on desktop

### 🎨 Modern UI/UX

- Built with Shadcn UI components
- Recharts for interactive data visualization
- Responsive design (Mobile, Tablet, Desktop)
- Dark mode support
- Consistent table styling across all pages
- Loading states and skeletons
- Toast notifications for user feedback
- Error boundaries for graceful error handling

## 🏗️ Architecture

### Backend (Go + Fiber)

```
backend/
├── main.go                 # Application entry point
├── config/                 # Configuration management
├── constants/              # App constants (HTTP, errors, etc.)
├── controllers/            # HTTP request handlers
├── databases/
│   ├── init.go            # MongoDB connection
│   ├── models/            # Data models
│   ├── collections/       # Collection names
│   └── seed/              # Database seeding
├── encryptions/            # AES, password hashing, ID encryption
├── helpers/                # Utility functions (JWT, time, email)
├── http/                   # HTTP response utilities
├── middlewares/            # Auth, RBAC middlewares
├── routers/                # Route definitions
├── services/               # Business logic layer
└── assets/uploads/         # File storage
```

### Frontend (Next.js + TypeScript)

```
frontend/
├── app/                    # Next.js App Router
│   ├── dashboard/         # Protected dashboard pages
│   │   ├── users/
│   │   ├── attendance/
│   │   ├── leaves/
│   │   ├── payroll/
│   │   ├── departments/
│   │   ├── documents/
│   │   └── profile/
│   ├── login/
│   ├── signup/
│   └── verify-email/
├── components/             # Reusable components
│   ├── ui/                # Shadcn UI components
│   └── dashboard/         # Dashboard-specific components
├── lib/                    # Utilities
│   ├── api-service.ts     # API client
│   ├── auth-context.tsx   # Auth state management
│   ├── config.ts          # API endpoints
│   └── types.ts           # TypeScript interfaces
└── hooks/                  # Custom React hooks
```

## 🚀 Getting Started

### Prerequisites

- **Go** 1.21 or higher
- **Node.js** 18.x or higher
- **MongoDB** 4.4 or higher
- **pnpm** (or npm/yarn)

### 1. Clone Repository

```bash
git clone https://github.com/Tiisha13/WorkZen-odoo13.git
cd WorkZen-odoo13
```

### 2. Backend Setup

```bash
cd backend

# Install Go dependencies
go mod download
go mod vendor

# Configure MongoDB (edit config.yml)
cp config.example.yml config.yml

# Run the server
go run main.go
# or build and run
go build -o workzen
./workzen
```

**Backend Configuration (`config.yml`):**

```yaml
database:
  uri: "mongodb://localhost:27017"
  name: "workzen_hrms"

server:
  port: 5000

jwt:
  secret: "your-super-secret-jwt-key-change-in-production"

smtp:
  host: "smtp.gmail.com"
  port: 587
  username: "your-email@gmail.com"
  password: "your-app-password"
  from: "noreply@workzen.com"
```

Backend runs on: **http://localhost:5000**

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
pnpm install

# Run development server
pnpm dev

# Build for production
pnpm build
pnpm start
```

Frontend runs on: **http://localhost:3000**

### 4. Default Login Credentials

**SuperAdmin:**

- Email: `superadmin@workzen.com`
- Password: `SuperAdmin@123`

**Demo Admin (after signup):**

- Use the credentials you create during signup
- Company-level admin access

## 🔑 User Roles & Permissions

| Role           | Level       | Permissions                                        | Dashboard Access                           |
| -------------- | ----------- | -------------------------------------------------- | ------------------------------------------ |
| **SuperAdmin** | 1 (Highest) | Platform-wide management, all companies            | All companies, platform stats              |
| **Admin**      | 2           | Company-level management, see HR/Payroll/Employees | Company employees (excluding other admins) |
| **HR**         | 3           | Employee management, leave approval                | HR, Payroll, Employees                     |
| **Payroll**    | 4           | Payroll processing, salary management              | Payroll, Employees                         |
| **Employee**   | 5 (Lowest)  | Own profile, attendance, leave application         | Own data only                              |

**Role Hierarchy Rules:**

- Lower priority roles cannot see higher priority roles
- Admins are NOT counted as employees (they are managers)
- Dashboard statistics filtered by role visibility
- Attendance, leave, and payroll data respects role hierarchy

## 📡 API Endpoints

### Authentication

- `POST /api/v1/auth/signup` - Company signup
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/verify-email?token={token}` - Email verification
- `POST /api/v1/auth/resend-verification` - Resend verification email

### Users

- `GET /api/v1/users` - List users (role-filtered)
- `POST /api/v1/users` - Create user (auto-generates password)
- `GET /api/v1/users/:id` - Get user details
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user (soft delete)

### Attendance

- `POST /api/v1/attendance/check-in` - Check in (employees only)
- `POST /api/v1/attendance/check-out` - Check out (employees only)
- `GET /api/v1/attendance/me` - My attendance records
- `GET /api/v1/attendance` - All attendance (admin/HR)

### Leaves

- `GET /api/v1/leaves` - List leaves (filtered by role)
- `POST /api/v1/leaves` - Apply for leave
- `PUT /api/v1/leaves/:id/approve` - Approve leave (HR/Admin)
- `PUT /api/v1/leaves/:id/reject` - Reject leave (HR/Admin)

### Payroll

- `GET /api/v1/payruns` - List payruns
- `POST /api/v1/payruns/generate` - Generate monthly payrun
- `GET /api/v1/payruns/:id` - Get payrun details
- `PUT /api/v1/payruns/:id/process` - Process payroll

### Documents

- `POST /api/v1/documents` - Upload document
- `GET /api/v1/documents` - List documents
- `GET /api/v1/documents/:id/view` - View/preview document
- `GET /api/v1/documents/:id/download` - Download document
- `DELETE /api/v1/documents/:id` - Delete document

### Dashboard

- `GET /api/v1/dashboard` - General dashboard
- `GET /api/v1/dashboard/admin` - Admin dashboard (with charts data)
- `GET /api/v1/dashboard/superadmin` - SuperAdmin dashboard

### Departments

- `GET /api/v1/departments` - List departments
- `POST /api/v1/departments` - Create department
- `PUT /api/v1/departments/:id` - Update department
- `DELETE /api/v1/departments/:id` - Delete department

## 🎨 Tech Stack

### Backend

- **Framework:** Go Fiber v2 (Express-inspired)
- **Database:** MongoDB with official Go driver
- **Authentication:** JWT (JSON Web Tokens)
- **Encryption:** AES-256 for ID encryption, bcrypt for passwords
- **Email:** SMTP integration for verification emails
- **File Storage:** Local filesystem with organized structure

### Frontend

- **Framework:** Next.js 16 (App Router, Turbopack)
- **Language:** TypeScript
- **UI Library:** Shadcn UI (Radix UI primitives)
- **Charts:** Recharts 2.15.4
- **Icons:** Tabler Icons
- **Styling:** Tailwind CSS
- **Forms:** React Hook Form (where applicable)
- **State Management:** React Context API
- **HTTP Client:** Fetch API with custom wrapper
- **Notifications:** Sonner (toast notifications)

## 📦 Key Features Implementation

### 1. Role-Based Dashboard Filtering

```go
// Backend filters data based on user role
func GetAdminDashboard(companyID, userRole) {
    var excludeRoles []models.Role
    switch userRole {
    case models.RoleAdmin:
        excludeRoles = []models.Role{SuperAdmin, Admin}
    case models.RoleHR:
        excludeRoles = []models.Role{SuperAdmin, Admin}
    // ... more filtering logic
    }
}
```

### 2. Encrypted IDs for Security

```go
// All IDs encrypted in API responses
func ConvertDocumentToResponse(doc *models.Document) *DocumentResponse {
    return &DocumentResponse{
        ID:         encryptions.EncryptID(doc.ID.Hex()),
        Company:    encryptions.EncryptID(doc.Company.Hex()),
        // ... other fields
    }
}
```

### 3. Auto-Generated Passwords

```go
// Users created without passwords - auto-generated by backend
func CreateUser(userData) {
    password := helpers.GenerateRandomPassword(12)
    hashedPassword := encryptions.HashPassword(password)
    // ... save user and return plain password once
}
```

### 4. Standardized Table Theme

```tsx
// Consistent table wrapper across all pages
<div className="rounded-md border bg-card shadow-sm overflow-hidden">
  <div className="overflow-x-auto">
    <Table>
      <TableHeader>{/* bg-muted/50 applied */}</TableHeader>
      <TableBody>{/* Consistent padding */}</TableBody>
    </Table>
  </div>
</div>
```

### 5. Interactive Charts with Recharts

```tsx
// Dashboard with multiple chart types
<ResponsiveContainer width="100%" height={350}>
  <BarChart data={departmentStats}>
    <Bar dataKey="Present" fill="#22c55e" />
    <Bar dataKey="Absent" fill="#ef4444" />
  </BarChart>
</ResponsiveContainer>
```

## 🔒 Security Features

- ✅ JWT authentication with token expiry
- ✅ Password hashing with bcrypt (cost factor 10)
- ✅ AES-256 encryption for sensitive IDs
- ✅ CORS configuration
- ✅ Rate limiting (can be implemented)
- ✅ Input validation
- ✅ SQL injection prevention (using MongoDB)
- ✅ XSS prevention (React automatic escaping)
- ✅ Role-based access control (RBAC)
- ✅ Multi-tenancy data isolation

## 📝 Database Seeding

On first run, the system automatically seeds:

**SuperAdmin User:**

- Email: `superadmin@workzen.com`
- Password: `SuperAdmin@123`
- Role: superadmin

**Demo Company:**

- Name: "Acme Corporation"
- Admin user with admin role
- Basic departments (Engineering, HR, Sales, Finance)

## 🧪 Testing

### Backend Testing

```bash
cd backend
go test ./... -v
```

### Frontend Testing

```bash
cd frontend
pnpm test
```

### API Testing with Postman

Postman collections are available in `backend/postman/`:

- `Admin.postman_collection.json`
- `Employee.postman_collection.json`
- `HR.postman_collection.json`
- `Payroll.postman_collection.json`
- `SuperAdmin.postman_collection.json`

## 🚀 Deployment

### Backend Deployment

```bash
# Build binary
go build -o workzen-api main.go

# Set environment variables
export MONGODB_URI="your-mongo-connection-string"
export JWT_SECRET="your-production-secret"

# Run
./workzen-api
```

### Frontend Deployment (Vercel)

```bash
# Build
pnpm build

# Preview production build locally
pnpm start

# Deploy to Vercel
vercel deploy --prod
```

**Environment Variables (`.env.local`):**

```env
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
```

## 📊 Project Statistics

- **Backend:** ~15,000 lines of Go code
- **Frontend:** ~10,000 lines of TypeScript/React
- **API Endpoints:** 50+
- **Database Collections:** 9
- **UI Components:** 40+ reusable components
- **Charts:** 4 interactive chart types

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Express-inspired Go web framework
- [Next.js](https://nextjs.org/) - The React Framework
- [Shadcn UI](https://ui.shadcn.com/) - Beautifully designed components
- [Recharts](https://recharts.org/) - Composable charting library
- [MongoDB](https://www.mongodb.com/) - NoSQL database
- [Tabler Icons](https://tabler-icons.io/) - Customizable icon set

## 📞 Support

For support, email support@workzen.com or open an issue on GitHub.

## 👥 Team

This project was developed by:

- **Shani Sinojiya** - [@Shani-Sinojiya](https://github.com/Shani-Sinojiya)
- **Tisha Patel** - [@Tiisha13](https://github.com/Tiisha13)

## Demo Video

<div style="padding:75.00% 0 0 0;position:relative;"><iframe src="https://player.vimeo.com/video/1135002138?h=5498b9bf86&amp;badge=0&amp;autopause=0&amp;player_id=0&amp;app_id=58479%2Fembed" allow="autoplay; fullscreen; picture-in-picture" allowfullscreen="" frameborder="0" style="position:absolute;top:0;left:0;width:100%;height:100%;"></iframe></div>
