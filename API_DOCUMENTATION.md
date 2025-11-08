# WorkZen HRMS - API Documentation

## Base URL

```
http://localhost:3000/api/v1
```

## Authentication

All endpoints except `/auth/signup` and `/auth/login` require JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

---

## 📋 Health Check

### Check API Status

```http
GET /health
```

**Response:**

```json
{
  "status": "healthy",
  "service": "WorkZen HRMS"
}
```

---

## 🔐 Authentication Endpoints

### 1. Signup (Create Company & Admin User)

```http
POST /auth/signup
```

**Request Body:**

```json
{
  "company_name": "Acme Corp",
  "email": "admin@acme.com",
  "password": "SecurePass123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+919876543210"
}
```

**Response:** `201 Created`

```json
{
  "success": true,
  "message": "Signup successful, awaiting approval"
}
```

### 2. Login

```http
POST /auth/login
```

**Request Body:**

```json
{
  "username": "admin@acme.com",
  "password": "SecurePass123"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "507f1f77bcf86cd799439011",
      "login_id": "EMP001",
      "email": "admin@acme.com",
      "role": "admin",
      "company": "507f191e810c19729de860ea"
    }
  }
}
```

### 3. Get Current User Profile

```http
GET /auth/me
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "User profile retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "login_id": "EMP001",
    "email": "admin@acme.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "admin",
    "status": "active",
    "company": "507f191e810c19729de860ea"
  }
}
```

### 4. Change Password

```http
POST /auth/change-password
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "old_password": "OldPass123",
  "new_password": "NewPass456"
}
```

**Response:** `200 OK`

---

## 🏢 Company Management (SuperAdmin Only)

### 1. Create Company

```http
POST /companies
Authorization: Bearer <token>
Role: SuperAdmin
```

**Request Body:**

```json
{
  "name": "New Company Ltd",
  "email": "contact@newcompany.com",
  "phone": "+919876543210",
  "address": "123 Business Street"
}
```

### 2. List All Companies

```http
GET /companies?page=1&limit=10
Authorization: Bearer <token>
Role: SuperAdmin
```

**Response:**

```json
{
  "success": true,
  "data": {
    "companies": [...],
    "total": 50,
    "page": 1,
    "limit": 10
  }
}
```

### 3. Get Company by ID

```http
GET /companies/:id
Authorization: Bearer <token>
```

### 4. Approve Company

```http
PATCH /companies/:id/approve
Authorization: Bearer <token>
Role: SuperAdmin
```

### 5. Deactivate Company

```http
PATCH /companies/:id/deactivate
Authorization: Bearer <token>
Role: SuperAdmin
```

---

## 👥 User Management

### 1. Create User (Employee)

```http
POST /users
Authorization: Bearer <token>
Role: Admin, HR
```

**Request Body:**

```json
{
  "email": "john.doe@company.com",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+919876543210",
  "role": "employee",
  "department_id": "507f1f77bcf86cd799439011",
  "manager_id": "507f191e810c19729de860ea"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "user": {...},
    "password": "TempPass123"
  }
}
```

### 2. List Users

```http
GET /users?page=1&limit=10
Authorization: Bearer <token>
Role: Admin, HR
```

### 3. Get User by ID

```http
GET /users/:id
Authorization: Bearer <token>
```

### 4. Update User Status

```http
PATCH /users/:id/status
Authorization: Bearer <token>
Role: Admin
```

**Request Body:**

```json
{
  "status": "inactive"
}
```

### 5. Update Bank Details

```http
PATCH /users/:id/bank
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "account_number": "1234567890",
  "ifsc_code": "HDFC0001234",
  "bank_name": "HDFC Bank",
  "pan_no": "ABCDE1234F",
  "uan_no": "123456789012"
}
```

### 6. Delete User (Soft Delete)

```http
DELETE /users/:id
Authorization: Bearer <token>
Role: Admin
```

---

## ⏰ Attendance Management

### 1. Check In

```http
POST /attendance/check-in
Authorization: Bearer <token>
```

**Response:**

```json
{
  "success": true,
  "message": "Check-in successful",
  "data": {
    "id": "...",
    "date": "2025-11-08",
    "check_in": "09:00:00",
    "status": "present"
  }
}
```

### 2. Check Out

```http
POST /attendance/check-out
Authorization: Bearer <token>
```

### 3. Get My Attendance

```http
GET /attendance/me?month=2025-11
Authorization: Bearer <token>
```

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "date": "2025-11-08",
      "check_in": "09:00:00",
      "check_out": "18:00:00",
      "work_hours": 9.0,
      "status": "present"
    }
  ]
}
```

### 4. List All Attendance (HR/Admin)

```http
GET /attendance?page=1&limit=10&employee_id=xxx&date=2025-11-08
Authorization: Bearer <token>
Role: Admin, HR
```

### 5. Get Attendance Summary

```http
GET /attendance/summary?date=2025-11-08
Authorization: Bearer <token>
Role: Admin, HR
```

**Response:**

```json
{
  "success": true,
  "data": {
    "present": 45,
    "absent": 3,
    "on_leave": 2
  }
}
```

---

## 🏖️ Leave Management

### 1. Apply Leave

```http
POST /leaves
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "leave_type": "casual_leave",
  "from_date": "2025-11-10",
  "to_date": "2025-11-12",
  "reason": "Personal work"
}
```

### 2. List Leaves

```http
GET /leaves?page=1&limit=10&employee_id=xxx&status=pending
Authorization: Bearer <token>
Role: Admin, HR
```

### 3. Approve Leave

```http
PATCH /leaves/:id/approve
Authorization: Bearer <token>
Role: Admin, HR
```

### 4. Reject Leave

```http
PATCH /leaves/:id/reject
Authorization: Bearer <token>
Role: Admin, HR
```

---

## 💰 Salary Structure Management

### 1. Create Salary Structure

```http
POST /salary-structure
Authorization: Bearer <token>
Role: Admin, Payroll (with CanModifySalaryInfo)
```

**Request Body:**

```json
{
  "employee_id": "507f1f77bcf86cd799439011",
  "monthly_wage": 50000,
  "effective_from": "2025-11-01"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "employee_id": "...",
    "monthly_wage": 50000,
    "basic_salary": 25000,
    "house_rent_allowance": 12500,
    "standard_allowance": 1800,
    "performance_bonus": 2000,
    "leave_travel_allowance": 1500,
    "fixed_allowance": 7200,
    "total_earnings": 50000,
    "pf_employee": 3000,
    "pf_employer": 3000,
    "professional_tax": 200,
    "net_pay": 46800
  }
}
```

### 2. Get Salary Structure

```http
GET /salary-structure/:employee_id
Authorization: Bearer <token>
```

### 3. Update Salary Structure

```http
PATCH /salary-structure/:employee_id
Authorization: Bearer <token>
Role: Admin, Payroll (with CanModifySalaryInfo)
```

**Request Body:**

```json
{
  "monthly_wage": 55000,
  "effective_from": "2025-12-01"
}
```

---

## 💵 Payroll Management

### 1. Create/Update Payroll Configuration

```http
POST /payroll/configuration
Authorization: Bearer <token>
Role: Admin
```

**Request Body:**

```json
{
  "pf_employee_percent": 12.0,
  "pf_employer_percent": 12.0,
  "professional_tax": 200.0,
  "default_basic_percent": 50.0,
  "default_hra_percent": 50.0,
  "default_standard_allowance": 1800.0,
  "default_performance_bonus": 2000.0,
  "default_lta": 1500.0
}
```

### 2. Get Payroll Configuration

```http
GET /payroll/configuration
Authorization: Bearer <token>
Role: Admin, Payroll
```

### 3. Create Payrun (Generate Monthly Payroll)

```http
POST /payruns
Authorization: Bearer <token>
Role: Admin, Payroll
```

**Request Body:**

```json
{
  "month": "2025-11"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "...",
    "month": "2025-11",
    "total_employees": 50,
    "processed_count": 48,
    "total_payroll": 2340000,
    "missing_bank_count": 2,
    "missing_manager_count": 1,
    "status": "generated"
  }
}
```

### 4. List Payruns

```http
GET /payruns?page=1&limit=10
Authorization: Bearer <token>
Role: Admin, Payroll
```

### 5. Get Employee Payroll

```http
GET /payrolls/:employee_id?month=2025-11
Authorization: Bearer <token>
```

**Response:**

```json
{
  "success": true,
  "data": {
    "employee_id": "...",
    "month": "2025-11",
    "basic_salary": 25000,
    "house_rent_allowance": 12500,
    "gross_salary": 50000,
    "pf_employee": 3000,
    "professional_tax": 200,
    "total_deductions": 3200,
    "net_pay": 46800,
    "status": "processed",
    "has_bank_account": true,
    "has_manager": true
  }
}
```

### 6. Mark Payroll as Paid

```http
PATCH /payrolls/:id/mark-paid
Authorization: Bearer <token>
Role: Admin, Payroll
```

---

## 📄 Document Management

### 1. Upload Document

```http
POST /documents/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Data:**

- `file`: (file upload)
- `category`: resume | id_proof | payslip | policy | report | other
- `description`: Optional description
- `employee_id`: Optional employee ID

**Response:**

```json
{
  "success": true,
  "message": "Document uploaded successfully",
  "data": {
    "id": "...",
    "file_name": "contract.pdf",
    "file_path": "assets/uploads/507f.../contract/2025/11/uuid.pdf",
    "category": "contract",
    "size": 102400
  }
}
```

### 2. List Documents

```http
GET /documents?page=1&limit=10&category=payslip&employee_id=xxx
Authorization: Bearer <token>
Role: Admin, HR
```

### 3. Delete Document

```http
DELETE /documents/:id
Authorization: Bearer <token>
Role: Admin
```

---

## 📊 Dashboard

### 1. Admin Dashboard (Company-Level Stats)

```http
GET /dashboard/admin
Authorization: Bearer <token>
Role: Admin
```

**Response:**

```json
{
  "success": true,
  "data": {
    "total_employees": 50,
    "active_employees": 48,
    "inactive_employees": 2,
    "present_today": 45,
    "absent_today": 3,
    "on_leave_today": 2,
    "pending_leaves": 5,
    "approved_leaves": 120,
    "rejected_leaves": 8,
    "missing_bank_accounts": 2,
    "missing_managers": 1,
    "total_payroll_this_year": 28000000
  }
}
```

### 2. SuperAdmin Dashboard (Platform-Wide Stats)

```http
GET /dashboard/superadmin
Authorization: Bearer <token>
Role: SuperAdmin
```

**Response:**

```json
{
  "success": true,
  "data": {
    "total_companies": 25,
    "active_companies": 23,
    "pending_approvals": 2,
    "total_employees": 1250,
    "total_payroll_processed": 15000,
    "total_payruns_generated": 300,
    "platform_revenue": 0
  }
}
```

---

## 🔑 Role-Based Access Control

### Roles:

1. **SuperAdmin** - Platform administrator
2. **Admin** - Company administrator
3. **HR** - HR personnel
4. **Payroll** - Payroll officer
5. **Employee** - Regular employee

### Access Matrix:

| Endpoint             | SuperAdmin | Admin | HR  | Payroll | Employee |
| -------------------- | ---------- | ----- | --- | ------- | -------- |
| Company Management   | ✅         | ❌    | ❌  | ❌      | ❌       |
| Create/Delete Users  | ❌         | ✅    | ✅  | ❌      | ❌       |
| View Users           | ❌         | ✅    | ✅  | ❌      | ❌       |
| Attendance (Self)    | ❌         | ✅    | ✅  | ✅      | ✅       |
| Attendance (All)     | ❌         | ✅    | ✅  | ❌      | ❌       |
| Leave Apply          | ❌         | ✅    | ✅  | ✅      | ✅       |
| Leave Approve/Reject | ❌         | ✅    | ✅  | ❌      | ❌       |
| Salary Structure     | ❌         | ✅    | ❌  | ✅\*    | ❌       |
| Payroll Config       | ❌         | ✅    | ❌  | ❌      | ❌       |
| Generate Payrun      | ❌         | ✅    | ❌  | ✅      | ❌       |
| Documents            | ❌         | ✅    | ✅  | ❌      | Own Only |
| Dashboard            | ✅         | ✅    | ❌  | ❌      | ❌       |

\*With `CanModifySalaryInfo` permission

---

## 📝 Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "message": "Invalid request body"
}
```

### 401 Unauthorized

```json
{
  "success": false,
  "message": "Invalid credentials"
}
```

### 403 Forbidden

```json
{
  "success": false,
  "message": "Insufficient permissions"
}
```

### 404 Not Found

```json
{
  "success": false,
  "message": "Resource not found"
}
```

### 500 Internal Server Error

```json
{
  "success": false,
  "message": "Internal server error"
}
```

---

## 🚀 Quick Start

1. **Start the server:**

   ```bash
   go run main.go
   ```

2. **Create SuperAdmin** (Manual database entry required)

3. **Signup as Company:**

   ```bash
   curl -X POST http://localhost:3000/api/v1/auth/signup \
     -H "Content-Type: application/json" \
     -d '{
       "company_name": "Acme Corp",
       "email": "admin@acme.com",
       "password": "SecurePass123",
       "first_name": "John",
       "last_name": "Doe",
       "phone": "+919876543210"
     }'
   ```

4. **SuperAdmin approves company**

5. **Login:**

   ```bash
   curl -X POST http://localhost:3000/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{
       "username": "admin@acme.com",
       "password": "SecurePass123"
     }'
   ```

6. **Use the token in subsequent requests**

---

## 📦 Database Collections

- **companies** - Company information
- **users** - Employee and admin users
- **departments** - Department structure
- **attendances** - Daily attendance records
- **leaves** - Leave applications
- **salary_structures** - Employee salary configurations
- **payroll_configurations** - Company payroll settings
- **payruns** - Monthly payroll generations
- **payrolls** - Individual payroll records
- **documents** - Uploaded documents
- **activity_logs** - Audit trail

---

## 🛠️ Technical Details

- **Framework:** Go Fiber v2
- **Database:** MongoDB
- **Authentication:** JWT (HS512, 24-day expiry)
- **Password Hashing:** bcrypt equivalent
- **File Storage:** Local filesystem (/assets/uploads)
- **Time Format:** YYYY-MM-DD (dates), HH:MM:SS (times)
- **Currency:** INR

---

## 📞 Support

For issues or questions, please contact the development team.
