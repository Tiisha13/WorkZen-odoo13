export interface User {
  id: string;
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  is_super_admin?: boolean;
  designation?: string;
  department_id?: string;
  employee_code?: string;
  status: string;
  phone?: string;
  company?: string;
  email_verified: boolean;
}

export interface Company {
  id: string;
  name: string;
  email: string;
  phone?: string;
  industry?: string;
  is_approved: boolean;
  is_active: boolean;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
  company?: Company;
}

export interface SignupRequest {
  company_name: string;
  email: string;
  phone: string;
  industry?: string;
  first_name: string;
  last_name: string;
  password: string;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export interface ApiResponse<T = null | undefined | unknown | object> {
  success: boolean;
  message: string;
  data?: T;
  meta?: { [key: string]: number }; // Optional metadata
}

export interface Department {
  id: string;
  name: string;
  description?: string;
  head_id?: string;
  company_id?: string;
  employee_count?: number;
  is_deleted?: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface Attendance {
  id: string;
  user_id: string;
  user?: User;
  date: string;
  check_in: string;
  check_out?: string;
  working_hours?: number;
  status: string; // present, absent, leave, half_day
  notes?: string;
  created_at?: string;
  updated_at?: string;
}

export interface Leave {
  id: string;
  user_id: string;
  user?: User;
  leave_type: string; // sick, casual, annual, unpaid
  start_date: string;
  end_date: string;
  days?: number;
  reason: string;
  status: string; // pending, approved, rejected
  approved_by?: string;
  approved_at?: string;
  created_at?: string;
  updated_at?: string;
}

export interface Salary {
  id: string;
  user_id: string;
  user?: User;
  basic_salary: number;
  allowances?: number;
  deductions?: number;
  net_salary: number;
  month: string;
  status: string; // pending, processing, paid
  payment_date?: string;
  created_at?: string;
  updated_at?: string;
}

export interface Document {
  id: string;
  user_id: string;
  user?: User;
  title: string;
  description?: string;
  file_name: string;
  file_path: string;
  file_type: string;
  file_size: number;
  category: string;
  uploaded_by: string;
  is_private: boolean;
  created_at?: string;
  updated_at?: string;
}
