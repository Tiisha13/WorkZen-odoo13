"use client";

import { useState, useEffect } from "react";
import { apiService } from "@/lib/api-service";
import { API_ENDPOINTS } from "@/lib/config";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
} from "@/components/ui/drawer";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { IconPlus, IconEdit, IconTrash, IconSearch } from "@tabler/icons-react";
import { toast } from "sonner";
import { useMediaQuery, usePageTitle, useRequireAuth } from "@/lib/hooks";
import type { User, Department } from "@/lib/types";

export default function UsersPage() {
  usePageTitle("User Management | WorkZen");
  useRequireAuth(["superadmin", "admin", "hr"]);

  const [mounted, setMounted] = useState(false);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  // Individual form states
  const [email, setEmail] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("employee");
  const [designation, setDesignation] = useState("");
  const [phone, setPhone] = useState("");
  const [departmentId, setDepartmentId] = useState("");
  const [departments, setDepartments] = useState<Department[]>([]);

  const currentUser = apiService.getUser();
  const isMobile = useMediaQuery("(max-width: 768px)");

  useEffect(() => {
    setMounted(true);
    fetchUsers();
    fetchDepartments();
  }, []);

  const fetchDepartments = async () => {
    try {
      const response = await apiService.get<{
        success: boolean;
        data: Department[];
      }>(API_ENDPOINTS.DEPARTMENTS);
      setDepartments(response.data || []);
    } catch (error) {
      console.error("Failed to fetch departments:", error);
    }
  };

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await apiService.get<{ success: boolean; data: User[] }>(
        API_ENDPOINTS.USERS
      );
      setUsers(response.data || []);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to fetch users";
      if (!message.includes("Session expired")) {
        toast.error(message);
      }
      console.error("Fetch users error:", error);
      setUsers([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validation
    if (!firstName.trim() || !lastName.trim() || !email.trim()) {
      toast.error("First name, last name, and email are required");
      return;
    }

    if (!editingUser && !password) {
      toast.error("Password is required for new users");
      return;
    }

    const userData = {
      email,
      first_name: firstName,
      last_name: lastName,
      password,
      role,
      designation,
      phone,
      department_id: departmentId || undefined,
    };

    try {
      setSubmitting(true);
      if (editingUser) {
        await apiService.put(
          `${API_ENDPOINTS.USERS}/${editingUser.id}`,
          userData
        );
        toast.success("User updated successfully");
      } else {
        await apiService.post(API_ENDPOINTS.USERS, userData);
        toast.success("User created successfully");
      }
      setIsDialogOpen(false);
      resetForm();
      fetchUsers();
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to save user";
      if (!message.includes("Session expired")) {
        toast.error(message);
      }
    } finally {
      setSubmitting(false);
    }
  };

  const [deleteUser, setDeleteUser] = useState<User | null>(null);

  const handleDelete = async (id: string) => {
    try {
      await apiService.delete(`${API_ENDPOINTS.USERS}/${id}`);
      toast.success("User deleted successfully");
      setDeleteUser(null);
      fetchUsers();
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to delete user";
      if (!message.includes("Session expired")) {
        toast.error(message);
      }
      console.error("Delete user error:", error);
    }
  };

  const resetForm = () => {
    setEmail("");
    setFirstName("");
    setLastName("");
    setPassword("");
    setRole("employee");
    setDesignation("");
    setPhone("");
    setDepartmentId("");
    setEditingUser(null);
  };

  const handleEdit = (user: User) => {
    setEditingUser(user);
    setEmail(user.email);
    setFirstName(user.first_name);
    setLastName(user.last_name);
    setPassword("");
    setRole(user.role);
    setDesignation(user.designation || "");
    setPhone(user.phone || "");
    setDepartmentId(user.department_id || "");
    setIsDialogOpen(true);
  };

  const filteredUsers = users.filter(
    (user) =>
      user.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
      user.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
      `${user.first_name} ${user.last_name}`
        .toLowerCase()
        .includes(searchTerm.toLowerCase())
  );

  const getRoleBadge = (role: string) => {
    const config: Record<string, { className: string; label: string }> = {
      superadmin: {
        className:
          "bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400",
        label: "Super Admin",
      },
      admin: {
        className:
          "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400",
        label: "Admin",
      },
      hr: {
        className:
          "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400",
        label: "HR",
      },
      payroll: {
        className:
          "bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400",
        label: "Payroll",
      },
      employee: {
        className:
          "bg-slate-100 text-slate-800 dark:bg-slate-900/30 dark:text-slate-400",
        label: "Employee",
      },
    };
    const { className, label } = config[role] || config.employee;
    return <Badge className={className}>{label}</Badge>;
  };

  const getStatusBadge = (status: string) => {
    const isActive = status === "active";
    return (
      <Badge
        className={
          isActive
            ? "bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400"
            : "bg-rose-100 text-rose-800 dark:bg-rose-900/30 dark:text-rose-400"
        }
      >
        {isActive ? "Active" : "Inactive"}
      </Badge>
    );
  };

  // Simplified and optimized form
  const UserFormFields = () => (
    <div className="grid gap-4 md:grid-cols-2">
      <div className="space-y-1.5">
        <Label htmlFor="first_name">First Name *</Label>
        <Input
          id="first_name"
          value={firstName}
          onChange={(e) => setFirstName(e.target.value)}
          required
          placeholder="John"
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="last_name">Last Name *</Label>
        <Input
          id="last_name"
          value={lastName}
          onChange={(e) => setLastName(e.target.value)}
          required
          placeholder="Doe"
        />
      </div>
      <div className="space-y-1.5 md:col-span-2">
        <Label htmlFor="email">Email *</Label>
        <Input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          placeholder="john@example.com"
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="password">Password {!editingUser && "*"}</Label>
        <Input
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required={!editingUser}
          placeholder="••••••••"
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="role">Role *</Label>
        <Select value={role} onValueChange={setRole}>
          <SelectTrigger>
            <SelectValue placeholder="Select role" />
          </SelectTrigger>
          <SelectContent>
            {currentUser?.role === "superadmin" && (
              <>
                <SelectItem value="employee">Employee</SelectItem>
                <SelectItem value="hr">HR</SelectItem>
                <SelectItem value="payroll">Payroll</SelectItem>
                <SelectItem value="admin">Admin</SelectItem>
              </>
            )}
            {currentUser?.role === "admin" && (
              <>
                <SelectItem value="employee">Employee</SelectItem>
                <SelectItem value="hr">HR</SelectItem>
                <SelectItem value="payroll">Payroll</SelectItem>
              </>
            )}
            {currentUser?.role === "hr" && (
              <SelectItem value="employee">Employee</SelectItem>
            )}
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="designation">Designation</Label>
        <Input
          id="designation"
          value={designation}
          onChange={(e) => setDesignation(e.target.value)}
          placeholder="Software Engineer"
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="department">Department</Label>
        <Select value={departmentId} onValueChange={setDepartmentId}>
          <SelectTrigger>
            <SelectValue placeholder="Select department" />
          </SelectTrigger>
          <SelectContent>
            {departments.map((dept) => (
              <SelectItem key={dept.id} value={dept.id}>
                {dept.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="phone">Phone</Label>
        <Input
          id="phone"
          type="tel"
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          placeholder="+1234567890"
        />
      </div>
    </div>
  );

  if (!mounted) {
    return null;
  }

  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Header Section */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold tracking-tight">
            Users
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage team members and their permissions
          </p>
        </div>
        <Button
          onClick={() => {
            resetForm();
            setIsDialogOpen(true);
          }}
          size="sm"
        >
          <IconPlus className="w-4 h-4 mr-2" />
          Add User
        </Button>
      </div>

      {/* Search Bar */}
      <div className="relative max-w-md">
        <IconSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground w-4 h-4" />
        <Input
          placeholder="Search by name, email, or username..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="pl-9 h-9"
        />
      </div>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={!!deleteUser} onOpenChange={() => setDeleteUser(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete User</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete{" "}
              <strong>
                {deleteUser?.first_name} {deleteUser?.last_name}
              </strong>
              ? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => deleteUser && handleDelete(deleteUser.id)}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Users Table Card */}
      <div className="bg-card rounded-lg border shadow-sm">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>User</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Role</TableHead>
              <TableHead>Designation</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} className="h-32 text-center">
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                  </div>
                </TableCell>
              </TableRow>
            ) : filteredUsers.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="h-32 text-center">
                  <div className="flex flex-col items-center justify-center text-muted-foreground">
                    <p className="text-sm">No users found</p>
                    {searchTerm && (
                      <p className="text-xs mt-1">Try adjusting your search</p>
                    )}
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              filteredUsers.map((user) => (
                <TableRow key={user.id} className="hover:bg-muted/50">
                  <TableCell>
                    <div className="flex flex-col">
                      <span className="font-medium">
                        {user.first_name} {user.last_name}
                      </span>
                      <span className="text-xs text-muted-foreground">
                        @{user.username}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="text-sm">{user.email}</TableCell>
                  <TableCell>{getRoleBadge(user.role)}</TableCell>
                  <TableCell className="text-sm">
                    {user.designation || (
                      <span className="text-muted-foreground">-</span>
                    )}
                  </TableCell>
                  <TableCell>{getStatusBadge(user.status)}</TableCell>
                  <TableCell className="text-right">
                    <div className="flex items-center justify-end gap-2">
                      <Button
                        size="icon"
                        variant="ghost"
                        onClick={() => handleEdit(user)}
                        className="h-8 w-8"
                      >
                        <IconEdit className="w-4 h-4" />
                      </Button>
                      <Button
                        size="icon"
                        variant="ghost"
                        onClick={() => setDeleteUser(user)}
                        className="h-8 w-8 text-destructive hover:text-destructive"
                      >
                        <IconTrash className="w-4 h-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Mobile Drawer */}
      {isMobile && (
        <Drawer open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DrawerContent className="max-h-[95vh]">
            <DrawerHeader>
              <DrawerTitle>
                {editingUser ? "Edit User" : "Add New User"}
              </DrawerTitle>
              <DrawerDescription>Fill in the required fields</DrawerDescription>
            </DrawerHeader>
            <form onSubmit={handleSubmit} className="px-4 pb-4 overflow-y-auto">
              <div className="space-y-4 pb-4">
                <UserFormFields />
              </div>
              <DrawerFooter className="px-0 pt-4">
                <Button type="submit" className="w-full" disabled={submitting}>
                  {submitting ? "Saving..." : editingUser ? "Update" : "Create"}
                </Button>
                <DrawerClose asChild>
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    disabled={submitting}
                  >
                    Cancel
                  </Button>
                </DrawerClose>
              </DrawerFooter>
            </form>
          </DrawerContent>
        </Drawer>
      )}

      {/* Desktop Dialog */}
      {!isMobile && (
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>
                {editingUser ? "Edit User" : "Add New User"}
              </DialogTitle>
              <DialogDescription>
                Fill in the required fields below
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-6">
              <UserFormFields />
              <DialogFooter className="gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setIsDialogOpen(false)}
                  disabled={submitting}
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={submitting}>
                  {submitting ? "Saving..." : editingUser ? "Update" : "Create"}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}
