# WorkZen HRMS - Backend API

A comprehensive Human Resource Management System built with Go (Fiber v2) and MongoDB.

## 🚀 Features

### Core Modules

- **Authentication & Authorization** - JWT-based auth with role-based access control
- **Company Management** - Multi-tenancy support with company-level isolation
- **User Management** - Employee CRUD with bank details, manager hierarchy
- **Attendance Tracking** - Check-in/out with automatic work hours calculation
- **Leave Management** - Leave application workflow with HR approval
- **Salary Structure** - Automated salary component calculation
- **Payroll Processing** - Monthly payrun generation with deductions
- **Document Management** - File upload with organized storage
- **Analytics Dashboard** - Company and platform-wide statistics

### Key Features

- ✅ Role-based access control (5 roles)
- ✅ Multi-tenancy architecture
- ✅ Automatic salary calculations (Basic, HRA, Allowances, PF, Tax)
- ✅ Attendance integration with leave approvals
- ✅ Warning system for missing bank accounts/managers
- ✅ Soft delete for users
- ✅ Comprehensive audit logging
- ✅ File upload with organized directory structure

## 📋 Prerequisites

- **Go** 1.21 or higher
- **MongoDB** 4.4 or higher
- **Git**

## 🔧 Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Tiisha13/WorkZen-odoo13.git
   cd WorkZen-odoo13
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   go mod vendor
   ```

3. **Configure MongoDB:**

   Update `config.yml` with your MongoDB connection string:

   ```yaml
   database:
     uri: "mongodb://localhost:27017"
     name: "workzen_hrms"

   server:
     port: 3000

   jwt:
     secret: "your-secret-key-here"
   ```

4. **Create uploads directory:**
   ```bash
   mkdir -p assets/uploads
   ```

## 🏃 Running the Application

### Development Mode

```bash
go run main.go
```

### Production Build

```bash
go build -o workzen-api main.go
./workzen-api
```

The server will start on `http://localhost:3000`

## 📚 API Documentation

Full API documentation is available in [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

### Quick Test

**Health Check:**

```bash
curl http://localhost:3000/api/v1/health
```

**Signup:**

```bash
curl -X POST http://localhost:3000/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "Test Company",
    "email": "admin@test.com",
    "password": "TestPass123",
    "first_name": "Admin",
    "last_name": "User",
    "phone": "+919876543210"
  }'
```

## 🏗️ Project Structure

```
WorkZen-odoo13/
├── main.go                 # Application entry point
├── config.yml              # Configuration file
├── config/                 # Config loading
├── constants/              # Application constants
├── databases/              # Database initialization
│   ├── models/            # MongoDB models
│   └── collections/       # Collection names
├── services/              # Business logic layer
│   ├── auth_service.go
│   ├── company_service.go
│   ├── user_service.go
│   ├── attendance_service.go
│   ├── leave_service.go
│   ├── salary_service.go
│   ├── payroll_service.go
│   ├── document_service.go
│   └── dashboard_service.go
├── controllers/           # HTTP request handlers
│   ├── auth_controller.go
│   ├── company_controller.go
│   ├── user_controller.go
│   ├── attendance_controller.go
│   ├── leave_controller.go
│   ├── salary_controller.go
│   ├── payroll_controller.go
│   ├── document_controller.go
│   └── dashboard_controller.go
├── middlewares/           # Auth & RBAC middleware
│   ├── auth.go
│   └── rbac.go
├── routers/              # Route registration
│   ├── main.go
│   └── routes.go
├── helpers/              # Utility functions
│   ├── jwt.go
│   ├── salary.go
│   ├── time.go
│   └── loginid.go
├── encryptions/          # Encryption utilities
│   ├── password.go
│   ├── aes.go
│   └── hash.go
├── http/                 # HTTP response helpers
│   ├── success.go
│   └── errors.go
└── assets/               # Static files
    └── uploads/          # Document uploads
```

## 🔐 Security Features

- JWT authentication with HS512 algorithm
- Password hashing using bcrypt-equivalent
- Role-based access control (RBAC)
- Company-scoped data isolation
- Token expiration (24 days)
- Middleware-based authorization

## 👥 User Roles

1. **SuperAdmin** - Platform administrator (manages companies)
2. **Admin** - Company administrator (full company access)
3. **HR** - HR personnel (user management, leaves, attendance)
4. **Payroll** - Payroll officer (salary & payroll management)
5. **Employee** - Regular employee (limited self-service access)

## 💰 Salary Calculation Logic

The system automatically calculates salary components based on monthly wage:

1. **Basic Salary** = 50% of Monthly Wage
2. **HRA** = 50% of Basic Salary
3. **Standard Allowance** = ₹1,800 (fixed)
4. **Performance Bonus** = ₹2,000 (fixed)
5. **LTA** = ₹1,500 (fixed)
6. **Fixed Allowance** = Remaining amount to match wage
7. **PF Employee** = 12% of Basic Salary
8. **PF Employer** = 12% of Basic Salary
9. **Professional Tax** = ₹200 (fixed)
10. **Net Pay** = Gross Salary - (PF Employee + Professional Tax)

## 📊 Database Collections

- `companies` - Company information
- `users` - Employee and admin users
- `departments` - Department hierarchy
- `attendances` - Daily attendance logs
- `leaves` - Leave applications
- `salary_structures` - Salary configurations
- `payroll_configurations` - Payroll settings
- `payruns` - Monthly payroll batches
- `payrolls` - Individual payroll records
- `documents` - Uploaded documents
- `activity_logs` - Audit trail

## 🧪 Testing

### Manual Testing with cURL

See [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) for detailed endpoint examples.

### Postman Collection

Import the API endpoints into Postman:

1. Base URL: `http://localhost:3000/api/v1`
2. Set Authorization header for protected routes
3. Use Bearer Token authentication

## 🔄 Workflow Examples

### 1. Employee Onboarding

```
1. Admin creates user via POST /users
2. System generates login ID and temporary password
3. Employee logs in and changes password
4. Admin creates salary structure via POST /salary-structure
5. System auto-calculates all components
```

### 2. Monthly Payroll Processing

```
1. Payroll officer creates payrun via POST /payruns
2. System processes all active employees
3. Fetches salary structures
4. Calculates deductions (PF, Tax)
5. Generates payroll records
6. Flags missing bank accounts/managers
7. Officer marks payrolls as paid
```

### 3. Leave Application

```
1. Employee applies leave via POST /leaves
2. System validates dates and calculates days
3. HR approves leave via PATCH /leaves/:id/approve
4. System creates attendance records for leave period
5. Attendance status = "on_leave" for those days
```

## 🐛 Troubleshooting

### MongoDB Connection Issues

```bash
# Check MongoDB is running
sudo systemctl status mongod

# Start MongoDB
sudo systemctl start mongod
```

### Port Already in Use

```bash
# Change port in config.yml
server:
  port: 8080
```

### File Upload Errors

```bash
# Ensure uploads directory exists with write permissions
mkdir -p assets/uploads
chmod 755 assets/uploads
```

## 📝 Development Notes

### Adding New Endpoints

1. Create service method in `services/`
2. Create controller method in `controllers/`
3. Register route in `routers/routes.go`
4. Apply appropriate middleware

### Database Schema Changes

1. Update model in `databases/models/`
2. Update service logic
3. Test thoroughly
4. Document in API docs

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

This project is proprietary software. All rights reserved.

## 👨‍💻 Development Team

- Backend API: WorkZen Development Team
- MongoDB Schema Design
- Authentication & Authorization
- Business Logic Implementation

## 📞 Support

For technical support or questions:

- Email: support@workzen.com
- Documentation: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

---

**Built with ❤️ using Go and MongoDB**
