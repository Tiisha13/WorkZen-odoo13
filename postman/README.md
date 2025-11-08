# WorkZen HRMS - Postman Collections

This directory contains role-specific Postman collections for testing the WorkZen HRMS Backend API.

## 📦 Collections Overview

### 1. **SuperAdmin.postman_collection.json**

**Role:** Platform Administrator  
**Access Level:** Full platform access across all companies

**Features:**

- Company management (create, approve, deactivate)
- View all companies in the system
- Platform-wide dashboard statistics
- System health monitoring

**Test Credentials:**

- Email: `superadmin@workzen.com`
- Password: `SuperAdmin@123`

---

### 2. **Admin.postman_collection.json**

**Role:** Company Administrator  
**Access Level:** Full company-level access

**Features:**

- User management (create HR, Payroll, Employees)
- Attendance monitoring
- Leave approval/rejection
- Salary structure management
- Payroll configuration and processing
- Document management
- Company dashboard

**Test Credentials:**

- Email: `admin@workzen.com`
- Password: `Admin@123`

---

### 3. **HR.postman_collection.json**

**Role:** Human Resources Officer  
**Access Level:** Employee records, attendance, and leave management

**Features:**

- View all employee records
- Monitor attendance across company
- Attendance statistics and summaries
- Leave approval/rejection
- View pending leave applications
- Document management
- Employee-specific filtering

**Test Credentials:**

- Email: `hr@workzen.com`
- Password: `TempPassword123` (created by Admin)

---

### 4. **Payroll.postman_collection.json**

**Role:** Payroll Officer  
**Access Level:** Salary and payroll management

**Features:**

- Create and update salary structures
- Payroll configuration
- Monthly payrun generation
- Mark payroll as paid
- Employee bank details management
- Attendance viewing (for payroll calculations)

**Test Credentials:**

- Email: `payroll@workzen.com`
- Password: `TempPassword123` (created by Admin)

---

### 5. **Employee.postman_collection.json**

**Role:** Regular Employee  
**Access Level:** Self-service features only

**Features:**

- Check-in/Check-out attendance
- View my attendance history
- Apply for leaves (casual, sick, etc.)
- View my leave applications
- View my salary structure
- View my payroll
- Upload personal documents
- Profile management
- Change password

**Test Credentials:**

- Email: `employee@workzen.com`
- Password: `TempPassword123` (created by Admin)

---

## 🚀 Quick Start

### 1. Import Collections

1. Open Postman
2. Click **Import** button
3. Select **Folder** tab
4. Choose the `postman` directory
5. All 5 collections will be imported

### 2. Set Base URL (Optional)

Each collection has a `base_url` variable set to `http://localhost:3000/api/v1`

To change:

1. Open collection
2. Go to **Variables** tab
3. Update `base_url` value

### 3. Test Workflow

#### **Start with SuperAdmin:**

```
1. Login as SuperAdmin
   → Token auto-saves to collection variable
2. View all companies
3. Check platform dashboard
```

#### **Switch to Admin:**

```
1. Login as Admin
   → Token auto-saves
2. Create HR user
3. Create Payroll user
4. Create Employee user
5. View admin dashboard
```

#### **Test HR Role:**

```
1. Login as HR
2. View all employees
3. Monitor attendance
4. Approve/reject leaves
```

#### **Test Payroll Role:**

```
1. Login as Payroll
2. Create salary structures
3. Configure payroll
4. Generate monthly payrun
5. Mark payroll as paid
```

#### **Test Employee Role:**

```
1. Login as Employee
2. Check-in attendance
3. Apply for leave
4. View my payroll
5. Check-out attendance
```

---

## 🔑 Authentication

All collections use **Bearer Token** authentication:

- Each collection has its own token variable:

  - `superadmin_token`
  - `admin_token`
  - `hr_token`
  - `payroll_token`
  - `employee_token`

- Tokens are **automatically saved** after successful login
- All subsequent requests use the saved token
- Token expires after 24 days (configurable in backend)

---

## 📊 Collection Variables

Each collection maintains its own variables:

| Variable      | Purpose           | Auto-Set           |
| ------------- | ----------------- | ------------------ |
| `base_url`    | API endpoint      | No                 |
| `*_token`     | Auth token        | Yes (on login)     |
| `user_id`     | Current user ID   | Yes (on login)     |
| `employee_id` | Target employee   | Yes (when created) |
| `company_id`  | Company ID        | Yes (SuperAdmin)   |
| `leave_id`    | Leave application | Yes (when created) |
| `payroll_id`  | Payroll record    | Yes (when fetched) |

---

## 🧪 Testing Scenarios

### Scenario 1: Complete Employee Onboarding

```
Admin Collection:
1. Create Employee
2. Create Salary Structure
3. Update Bank Details

Employee Collection:
4. Login as new employee
5. Change password
6. Check-in attendance
7. Upload resume
```

### Scenario 2: Leave Management Flow

```
Employee Collection:
1. Apply for leave

HR Collection:
2. View pending leaves
3. Approve leave

Employee Collection:
4. Check leave status
```

### Scenario 3: Monthly Payroll Processing

```
Payroll Collection:
1. Create payroll configuration
2. Generate monthly payrun
3. Get employee payroll
4. Mark as paid
```

### Scenario 4: Platform Management

```
SuperAdmin Collection:
1. View all companies
2. Approve pending companies
3. Check platform statistics
4. Deactivate inactive companies
```

---

## 📝 Notes

### Token Management

- Tokens are saved per collection
- Login once per role to start testing
- Change password updates don't affect token
- Re-login if token expires

### Variable Chaining

- Variables auto-populate across requests
- Employee IDs captured when creating users
- Leave IDs captured when applying leaves
- Payroll IDs captured when fetching payroll

### Error Handling

- All requests include proper error responses
- Check Postman console for detailed errors
- 401 = Token expired or invalid
- 403 = Insufficient permissions
- 404 = Resource not found

### File Uploads

- Document upload uses `multipart/form-data`
- Select file in Postman's body → form-data
- Set `category` and `description` fields
- Requires valid employee_id

---

## 🔒 Security Notes

1. **Change Default Passwords** after first login
2. **Never commit** credentials to version control
3. **Use environment variables** for sensitive data
4. **Rotate tokens** regularly in production
5. **Test with dummy data** only

---

## 🐛 Troubleshooting

### Issue: 401 Unauthorized

**Solution:** Login again to refresh token

### Issue: 403 Forbidden

**Solution:** Check role permissions for the endpoint

### Issue: Token not saving

**Solution:** Check the "Tests" tab in login request

### Issue: Employee ID not found

**Solution:** Run "List All Users" first to populate variable

### Issue: File upload fails

**Solution:** Ensure file is selected in form-data body

---

## 📚 Additional Resources

- [API Documentation](../API_DOCUMENTATION.md)
- [Quick Reference](../QUICK_REFERENCE.md)
- [Seed Guide](../SEED_GUIDE.md)
- [README](../README.md)

---

## 🎯 Testing Checklist

### SuperAdmin

- [ ] Login successful
- [ ] List companies
- [ ] Create company
- [ ] Approve company
- [ ] View platform dashboard

### Admin

- [ ] Login successful
- [ ] Create HR user
- [ ] Create Payroll user
- [ ] Create Employee
- [ ] Update bank details
- [ ] Create salary structure
- [ ] Configure payroll
- [ ] Generate payrun
- [ ] View admin dashboard

### HR

- [ ] Login successful
- [ ] View all employees
- [ ] Monitor attendance
- [ ] View pending leaves
- [ ] Approve leave
- [ ] Reject leave
- [ ] Upload document

### Payroll

- [ ] Login successful
- [ ] Create salary structure
- [ ] Update salary structure
- [ ] Configure payroll
- [ ] Generate payrun
- [ ] Mark payroll as paid

### Employee

- [ ] Login successful
- [ ] Check-in attendance
- [ ] Check-out attendance
- [ ] Apply leave
- [ ] View my leaves
- [ ] View my payroll
- [ ] Upload document
- [ ] Change password

---

**Happy Testing! 🚀**
