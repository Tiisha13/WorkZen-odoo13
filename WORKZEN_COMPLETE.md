# WorkZen HRMS - Full Stack Application

## 🚀 Overview

WorkZen is a comprehensive Human Resource Management System (HRMS) built with modern technologies. It features complete RBAC (Role-Based Access Control) authentication and full CRUD functionality across all modules.

## 📚 Tech Stack

### Backend

- **Framework**: Go Fiber v2.52.9
- **Database**: MongoDB
- **Authentication**: JWT with custom password hashing (SHA-256/512)
- **API**: RESTful JSON (126 handlers registered)
- **Server**: http://127.0.0.1:5000

### Frontend

- **Framework**: Next.js 16.0.1 with App Router & Turbopack
- **Language**: TypeScript 5.9.3
- **UI Library**: Shadcn UI (Radix + Tailwind CSS 4.1.17)
- **Icons**: Tabler Icons, Lucide React
- **State Management**: React Context API
- **Notifications**: Sonner
- **Charts**: Recharts 2.15.4

## 🎯 Features

### ✅ Complete Modules

1. **Authentication System**

   - Login with username/email + password
   - JWT token management
   - Role-based access control (5 roles)
   - Protected routes with auto-redirect
   - Change password functionality

2. **Dashboard**

   - Real-time statistics (Total Employees, Departments, Present Today, Pending Leaves)
   - Role-based data display
   - Quick overview of system status

3. **User Management** (`/dashboard/users`)

   - List all users with search
   - Create new users
   - Edit user details
   - Delete users
   - Role assignment
   - Status management (Active/Inactive)
   - Designation tracking

4. **Department Management** (`/dashboard/departments`)

   - List all departments
   - Create new departments
   - Edit department information
   - Delete departments
   - Employee count tracking

5. **Attendance System** (`/dashboard/attendance`)

   - Check-in/Check-out functionality
   - Today's attendance status
   - Attendance history table
   - Working hours calculation
   - Status badges (Present, Absent, Leave, Half Day)

6. **Leave Management** (`/dashboard/leaves`)

   - Request leave with type selection
   - Leave balance display
   - HR/Admin approval interface
   - Leave history with status
   - Multiple leave types (Sick, Casual, Annual, Unpaid)

7. **Payroll Management** (`/dashboard/payroll`)

   - Generate payroll for selected month
   - View payroll records
   - Download payslips
   - Salary breakdown (Basic, Allowances, Deductions)
   - Payment status tracking

8. **Document Management** (`/dashboard/documents`)

   - Upload documents with categories
   - View document list
   - Download documents
   - Document viewer
   - Privacy controls (Private/Public)
   - Categories: General, Contract, Policy, Payslip, Certificate

9. **Profile Page** (`/dashboard/profile`)

   - View personal information
   - Employment details
   - Role and status display

10. **Settings Page** (`/dashboard/settings`)
    - Update profile information
    - Change password
    - Notification preferences

## 👥 User Roles & Permissions

### 5-Level RBAC System

1. **SuperAdmin**

   - Full system access
   - Company management
   - All CRUD operations

2. **Admin**

   - User management
   - Department management
   - View all modules

3. **HR**

   - User management
   - Department management
   - Leave approval
   - Attendance management

4. **Payroll**

   - Payroll management
   - Salary management
   - Generate payslips

5. **Employee**
   - View own data
   - Check-in/Check-out
   - Request leaves
   - View documents
   - Update profile

## 🔑 Test Credentials

### SuperAdmin

- Username: `superadmin`
- Password: `SuperAdmin@123`

### Demo Admin

- Username: `demoadmin`
- Password: `Admin@123`

## 🛠️ Setup & Installation

### Backend Setup

```bash
cd backend

# Install dependencies (if not using vendor)
go mod download

# Run the server
go run main.go
```

Backend will start on: **http://127.0.0.1:5000**

### Frontend Setup

```bash
cd frontend

# Install dependencies
pnpm install

# Run development server
pnpm dev
```

Frontend will start on: **http://localhost:3000**

## 📁 Project Structure

```
WorkZen-odoo13/
├── backend/
│   ├── config/          # Configuration management
│   ├── constants/       # Application constants
│   ├── controllers/     # HTTP request handlers
│   ├── databases/       # MongoDB connection & models
│   ├── encryptions/     # Password hashing & encryption
│   ├── helpers/         # Helper functions (including is_deleted fix)
│   ├── http/           # HTTP response utilities
│   ├── middlewares/    # Authentication & RBAC middleware
│   ├── routers/        # Route definitions
│   ├── services/       # Business logic layer
│   └── main.go         # Application entry point
│
└── frontend/
    ├── app/
    │   ├── dashboard/
    │   │   ├── page.tsx              # Dashboard home
    │   │   ├── users/page.tsx        # User management
    │   │   ├── departments/page.tsx  # Department management
    │   │   ├── attendance/page.tsx   # Attendance system
    │   │   ├── leaves/page.tsx       # Leave management
    │   │   ├── payroll/page.tsx      # Payroll management
    │   │   ├── documents/page.tsx    # Document management
    │   │   ├── profile/page.tsx      # User profile
    │   │   └── settings/page.tsx     # Settings
    │   ├── login/page.tsx            # Login page
    │   └── layout.tsx                # Root layout
    ├── components/
    │   ├── dashboard/    # Dashboard components
    │   ├── layout/       # Layout components (sidebar, header)
    │   └── ui/          # Shadcn UI components
    └── lib/
        ├── api-service.ts     # HTTP client
        ├── auth-context.tsx   # Authentication state
        ├── config.ts          # API endpoints & constants
        ├── hooks.ts           # Custom hooks
        └── types.ts           # TypeScript interfaces
```

## 🔧 API Endpoints

### Authentication

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/signup` - Company signup
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/me` - Get current user
- `POST /api/v1/auth/change-password` - Change password

### Users

- `GET /api/v1/users` - List all users
- `POST /api/v1/users` - Create new user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Departments

- `GET /api/v1/departments` - List all departments
- `POST /api/v1/departments` - Create department
- `GET /api/v1/departments/:id` - Get department
- `PUT /api/v1/departments/:id` - Update department
- `DELETE /api/v1/departments/:id` - Delete department

### Attendance

- `GET /api/v1/attendances` - List attendances
- `GET /api/v1/attendances/today` - Get today's attendance
- `POST /api/v1/attendance/checkin` - Check-in
- `POST /api/v1/attendance/checkout` - Check-out

### Leaves

- `GET /api/v1/leaves` - List leaves
- `POST /api/v1/leaves` - Request leave
- `PUT /api/v1/leaves/:id/approve` - Approve leave
- `PUT /api/v1/leaves/:id/reject` - Reject leave

### Payroll

- `GET /api/v1/payroll` - List payroll records
- `POST /api/v1/payroll/generate` - Generate payroll
- `GET /api/v1/payroll/:id/payslip` - Download payslip

### Documents

- `GET /api/v1/documents` - List documents
- `POST /api/v1/documents` - Upload document
- `GET /api/v1/documents/:id/view` - View document
- `GET /api/v1/documents/:id/download` - Download document

### Dashboard

- `GET /api/v1/dashboard` - Get dashboard statistics

## 🐛 Recent Fixes

### Backend Issues Fixed

1. **is_deleted Query Problem**

   - Created `helpers.NotDeletedFilter()` for backward compatibility
   - Updated 14+ queries across all services
   - Fixed seed data to include `IsDeleted: false`

2. **Authentication**
   - Fixed login query with complex `$and` / `$or` logic
   - Improved password verification

### Frontend Features

- Complete authentication system with JWT
- Role-based route protection
- All CRUD modules implemented
- Responsive UI with Shadcn components
- Loading states and error handling
- Toast notifications

## 🎨 UI Components

### Shadcn UI Components Used

- Button, Input, Label, Textarea
- Table, Dialog, Card, Badge, Tabs
- Select, Skeleton
- Custom Sidebar and Navigation

### Icons

- Tabler Icons for main UI elements
- Lucide React for additional icons

## 📝 Development Notes

### State Management

- **Authentication**: React Context API (`auth-context.tsx`)
- **API Calls**: Custom service class (`api-service.ts`)
- **Local Storage**: JWT token, user data, company data

### Route Protection

- Custom hooks: `useRequireAuth()`, `useRedirectIfAuthenticated()`
- Automatic redirection based on authentication state
- Role-based menu filtering

### TypeScript

- Complete type safety with interfaces
- API response types
- Component props types

## 🚦 Running Status

- ✅ Backend: Running on http://127.0.0.1:5000 (PID 9838)
- ✅ Frontend: Running on http://localhost:3000
- ✅ Database: MongoDB connected
- ✅ All TypeScript errors: Resolved
- ✅ All modules: Functional

## 📸 Features Highlight

1. **Comprehensive CRUD**: All modules support Create, Read, Update, Delete
2. **Search & Filter**: Search functionality in users, departments, documents
3. **Real-time Updates**: Immediate feedback with toast notifications
4. **Responsive Design**: Works on desktop and mobile
5. **Role-Based UI**: Menu items and features shown based on user role
6. **Professional UI**: Clean, modern interface with Shadcn components

## 🔮 Future Enhancements

- [ ] Employee onboarding workflow
- [ ] Performance review module
- [ ] Training & development tracking
- [ ] Expense management
- [ ] Time tracking & project management
- [ ] Advanced reporting & analytics
- [ ] Mobile app (React Native)
- [ ] Email notifications
- [ ] File preview for documents
- [ ] Advanced search with filters

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

This project is proprietary software.

## 👨‍💻 Developer

**WorkZen Development Team**

- Full Stack HRMS Solution
- Built with ❤️ using Go, Next.js, and MongoDB

---

**Last Updated**: December 2024
**Version**: 1.0.0
