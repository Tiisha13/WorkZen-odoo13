// API Configuration
export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://127.0.0.1:5000";
export const API_VERSION = "/api/v1";

export const API_ENDPOINTS = {
  // Auth
  LOGIN: `${API_BASE_URL}${API_VERSION}/auth/login`,
  SIGNUP: `${API_BASE_URL}${API_VERSION}/auth/signup`,
  LOGOUT: `${API_BASE_URL}${API_VERSION}/auth/logout`,
  ME: `${API_BASE_URL}${API_VERSION}/auth/me`,
  CHANGE_PASSWORD: `${API_BASE_URL}${API_VERSION}/auth/change-password`,
  VERIFY_EMAIL: `${API_BASE_URL}${API_VERSION}/auth/verify-email`,
  RESEND_VERIFICATION: `${API_BASE_URL}${API_VERSION}/auth/resend-verification`,

  // Users
  USERS: `${API_BASE_URL}${API_VERSION}/users`,

  // Departments
  DEPARTMENTS: `${API_BASE_URL}${API_VERSION}/departments`,

  // Attendance
  ATTENDANCES: `${API_BASE_URL}${API_VERSION}/attendance/me`,
  ATTENDANCE_CHECKIN: `${API_BASE_URL}${API_VERSION}/attendance/check-in`,
  ATTENDANCE_CHECKOUT: `${API_BASE_URL}${API_VERSION}/attendance/check-out`,
  ATTENDANCE_RESET: `${API_BASE_URL}${API_VERSION}/attendance/reset`,
  ATTENDANCE_LIST: `${API_BASE_URL}${API_VERSION}/attendance`,

  // Leaves
  LEAVES: `${API_BASE_URL}${API_VERSION}/leaves`,

  // Payroll
  PAYROLL: `${API_BASE_URL}${API_VERSION}/payruns`,
  PAYROLLS: `${API_BASE_URL}${API_VERSION}/payrolls`,
  SALARY: `${API_BASE_URL}${API_VERSION}/salary-structure`,

  // Documents
  DOCUMENTS: `${API_BASE_URL}${API_VERSION}/documents`,
  DOCUMENTS_UPLOAD: `${API_BASE_URL}${API_VERSION}/documents/upload`,

  // Companies
  COMPANIES: `${API_BASE_URL}${API_VERSION}/companies`,

  // Dashboard
  DASHBOARD: `${API_BASE_URL}${API_VERSION}/dashboard`,
  DASHBOARD_ADMIN: `${API_BASE_URL}${API_VERSION}/dashboard/admin`,
  DASHBOARD_SUPERADMIN: `${API_BASE_URL}${API_VERSION}/dashboard/superadmin`,
};

export enum UserRole {
  SUPERADMIN = "superadmin",
  ADMIN = "admin",
  HR = "hr",
  PAYROLL = "payroll",
  EMPLOYEE = "employee",
}

export enum UserStatus {
  ACTIVE = "active",
  INACTIVE = "inactive",
}

export const ROLE_LABELS: Record<UserRole, string> = {
  [UserRole.SUPERADMIN]: "Super Admin",
  [UserRole.ADMIN]: "Admin",
  [UserRole.HR]: "HR",
  [UserRole.PAYROLL]: "Payroll",
  [UserRole.EMPLOYEE]: "Employee",
};
