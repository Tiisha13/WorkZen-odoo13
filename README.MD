## 🧭 **System Overview**

**WorkZen – Smart Human Resource Management System (HRMS)**
An all-in-one HR management platform designed to modernize HR operations for startups, institutions, and SMEs. It integrates **attendance**, **leave**, **payroll**, and **analytics** in one place — offering a scalable, role-based system with dashboards for data-driven workforce management.

---

## 🎯 **Vision & Mission**

To simplify how organizations handle people, processes, and payroll by:

- Reducing manual HR dependencies
- Improving transparency and accuracy
- Supporting data-driven HR decisions
- Delivering a clean, user-friendly digital HR experience

---

## ⚙️ **Core Modules**

### 1. **User & Role Management**

- **User Registration & Login** (secure authentication)
- **Roles:**

  - Employee
  - HR Officer
  - Payroll Officer
  - Admin

- **Profile Management:** Each user can view/edit their personal details (based on access rights)
- **Admin Controls:** Create, edit, delete, or assign user roles.

---

### 2. **Attendance & Leave Management**

- Employees can **mark daily attendance**
- **View monthly attendance logs**
- Apply for **leave/time-off** (vacation, sick leave, etc.)
- HR or Payroll officer can **approve/reject leave**
- Data syncs with Payroll for automated salary calculation

---

### 3. **Payroll Management**

- Displays **salary breakdown**, **deductions**, and **payout summary**
- Payroll officer or Admin can **generate and edit monthly reports**
- **Payslip generation** after each Payrun
- Includes statutory deductions:

  - **Provident Fund (PF)** – 12% employer/employee contribution
  - **Professional Tax** – monthly state tax on income

---

### 4. **Dashboard & Analytics**

- **Employee Dashboard:** Attendance, Leave, and Payslip overview
- **Admin Dashboard:** Company-wide HR metrics and summaries
- **Analytics Visuals:** Attendance trends, payroll statistics, and leave distribution via charts

---

## 👥 **Roles & Permissions**

| Role                | Capabilities                                              | Restrictions                |
| ------------------- | --------------------------------------------------------- | --------------------------- |
| **Admin**           | Full CRUD across modules, manage users & roles            | None                        |
| **Employee**        | Apply/view leave, view attendance/performance             | No payroll or admin access  |
| **HR Officer**      | Manage employee profiles & attendance                     | No payroll/system settings  |
| **Payroll Officer** | Approve leaves, generate payslips, manage payroll reports | No employee profile changes |

---

## 🧮 **Key HR Terminologies**

| Term                    | Meaning                                                 |
| ----------------------- | ------------------------------------------------------- |
| **Payroll**             | Processing of employee salaries, bonuses, deductions    |
| **Payrun**              | Specific payroll cycle (monthly/weekly)                 |
| **Payslip**             | Official document showing earnings & deductions         |
| **Time-Off**            | Approved leave duration                                 |
| **Wage**                | Compensation based on hours/attendance                  |
| **Provident Fund (PF)** | 12% of basic salary as employee + employer contribution |
| **Professional Tax**    | Monthly tax deducted by employer                        |

---

## 🧩 **Workflow Summary (from Excalidraw Diagram)**

### 🔸 **1. Login & Role Access**

- Users login → Redirected to their respective dashboards based on role
- **Admin** sees full control panel
- **Employee** sees attendance + leave tab
- **HR Officer** sees employee data
- **Payroll Officer** sees payrun & reports

### 🔸 **2. Attendance Flow**

- Employee → Marks attendance daily
- Data stored in attendance logs
- HR Officer reviews attendance reports

### 🔸 **3. Leave Flow**

- Employee applies → HR Officer approves/rejects
- Approved leaves reflected in **Attendance & Payroll**

### 🔸 **4. Payroll Flow**

- Payroll Officer initiates **Payrun** → System calculates salaries
  (Base pay – PF – Professional Tax – Deductions + Allowances)
- Generates **Payslips** for all employees

### 🔸 **5. Reports & Dashboard**

- Admin & Payroll Officer can generate analytics reports
- Visual dashboards display attendance %, salary trends, leave stats

---

## 🧑‍💼 **Hackathon Objective**

To build a **working prototype** implementing the above modules with:

- Secure authentication
- Role-based dashboards
- CRUD operations for HR data
- Real-time payroll calculations
- Visual analytics integration

**Deliverable:** Source code hosted on GitHub with proper commits.

---

## 📊 **Reference Link (from the document)**

Mockup Design:
🔗 [https://link.excalidraw.com/l/65VNwvy7c4X/7gxoB8JymIS](https://link.excalidraw.com/l/65VNwvy7c4X/7gxoB8JymIS)
