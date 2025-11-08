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
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
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
import type { User } from "@/lib/types";

export default function UsersPage() {
  usePageTitle("User Management | WorkZen");
  useRequireAuth(["superadmin", "admin", "hr"]);

  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const [formData, setFormData] = useState({
    username: "",
    email: "",
    first_name: "",
    last_name: "",
    password: "",
    role: "employee",
    designation: "",
    phone: "",
  });

  const isMobile = useMediaQuery("(max-width: 768px)");

  useEffect(() => {
    fetchUsers();
  }, []);

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
    if (!formData.username.trim() || !formData.email.trim()) {
      toast.error("Username and email are required");
      return;
    }

    if (!editingUser && !formData.password) {
      toast.error("Password is required for new users");
      return;
    }

    try {
      if (editingUser) {
        await apiService.put(
          `${API_ENDPOINTS.USERS}/${editingUser.id}`,
          formData
        );
        toast.success("User updated successfully");
      } else {
        await apiService.post(API_ENDPOINTS.USERS, formData);
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
    setFormData({
      username: "",
      email: "",
      first_name: "",
      last_name: "",
      password: "",
      role: "employee",
      designation: "",
      phone: "",
    });
    setEditingUser(null);
  };

  const handleEdit = (user: User) => {
    setEditingUser(user);
    setFormData({
      username: user.username,
      email: user.email,
      first_name: user.first_name,
      last_name: user.last_name,
      password: "",
      role: user.role,
      designation: user.designation || "",
      phone: user.phone || "",
    });
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
    const colors: Record<string, string> = {
      superadmin: "bg-purple-100 text-purple-800",
      admin: "bg-blue-100 text-blue-800",
      hr: "bg-green-100 text-green-800",
      payroll: "bg-yellow-100 text-yellow-800",
      employee: "bg-gray-100 text-gray-800",
    };
    return (
      <Badge className={colors[role] || colors.employee}>
        {role.charAt(0).toUpperCase() + role.slice(1)}
      </Badge>
    );
  };

  const getStatusBadge = (status: string) => {
    return (
      <Badge variant={status === "active" ? "default" : "secondary"}>
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </Badge>
    );
  };

  // Form fields component
  const UserFormFields = () => (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div className="space-y-2">
        <Label htmlFor="first_name">First Name *</Label>
        <Input
          id="first_name"
          value={formData.first_name}
          onChange={(e) =>
            setFormData({ ...formData, first_name: e.target.value })
          }
          required
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="last_name">Last Name *</Label>
        <Input
          id="last_name"
          value={formData.last_name}
          onChange={(e) =>
            setFormData({ ...formData, last_name: e.target.value })
          }
          required
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="username">Username *</Label>
        <Input
          id="username"
          value={formData.username}
          onChange={(e) =>
            setFormData({ ...formData, username: e.target.value })
          }
          required
          disabled={!!editingUser}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="email">Email *</Label>
        <Input
          id="email"
          type="email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          required
        />
      </div>
      {!editingUser && (
        <div className="space-y-2">
          <Label htmlFor="password">Password *</Label>
          <Input
            id="password"
            type="password"
            value={formData.password}
            onChange={(e) =>
              setFormData({ ...formData, password: e.target.value })
            }
            required={!editingUser}
          />
        </div>
      )}
      <div className="space-y-2">
        <Label htmlFor="role">Role *</Label>
        <Select
          value={formData.role}
          onValueChange={(value) => setFormData({ ...formData, role: value })}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="employee">Employee</SelectItem>
            <SelectItem value="hr">HR</SelectItem>
            <SelectItem value="payroll">Payroll</SelectItem>
            <SelectItem value="admin">Admin</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2">
        <Label htmlFor="designation">Designation</Label>
        <Input
          id="designation"
          value={formData.designation}
          onChange={(e) =>
            setFormData({ ...formData, designation: e.target.value })
          }
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="phone">Phone</Label>
        <Input
          id="phone"
          value={formData.phone}
          onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
        />
      </div>
    </div>
  );

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">User Management</h1>
          <p className="text-muted-foreground">
            Manage system users and their roles
          </p>
        </div>
        {isMobile ? (
          <Drawer open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DrawerTrigger asChild>
              <Button onClick={resetForm}>
                <IconPlus className="w-4 h-4 mr-2" />
                Add User
              </Button>
            </DrawerTrigger>
            <DrawerContent>
              <DrawerHeader>
                <DrawerTitle>
                  {editingUser ? "Edit User" : "Add New User"}
                </DrawerTitle>
                <DrawerDescription>
                  {editingUser
                    ? "Update user information"
                    : "Create a new user account"}
                </DrawerDescription>
              </DrawerHeader>
              <form onSubmit={handleSubmit} className="space-y-4 px-4">
                <UserFormFields />
                <DrawerFooter>
                  <Button type="submit">
                    {editingUser ? "Update User" : "Create User"}
                  </Button>
                  <DrawerClose asChild>
                    <Button type="button" variant="outline">
                      Cancel
                    </Button>
                  </DrawerClose>
                </DrawerFooter>
              </form>
            </DrawerContent>
          </Drawer>
        ) : (
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
              <Button onClick={resetForm}>
                <IconPlus className="w-4 h-4 mr-2" />
                Add User
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>
                  {editingUser ? "Edit User" : "Add New User"}
                </DialogTitle>
                <DialogDescription>
                  {editingUser
                    ? "Update user information"
                    : "Create a new user account"}
                </DialogDescription>
              </DialogHeader>
              <form onSubmit={handleSubmit} className="space-y-4">
                <UserFormFields />
                <DialogFooter>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setIsDialogOpen(false)}
                  >
                    Cancel
                  </Button>
                  <Button type="submit">
                    {editingUser ? "Update User" : "Create User"}
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>
        )}
      </div>

      <div className="flex items-center space-x-2">
        <div className="relative flex-1 max-w-sm">
          <IconSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
          <Input
            placeholder="Search users..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-10"
          />
        </div>
      </div>

      {/* Alert Dialog for Delete Confirmation */}
      <AlertDialog open={!!deleteUser} onOpenChange={() => setDeleteUser(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the
              user
              <strong className="block mt-2">
                {deleteUser?.first_name} {deleteUser?.last_name} (
                {deleteUser?.username})
              </strong>
              and remove their data from the system.
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

      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Username</TableHead>
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
                <TableCell colSpan={7} className="text-center py-8">
                  Loading...
                </TableCell>
              </TableRow>
            ) : filteredUsers.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={7}
                  className="text-center py-8 text-muted-foreground"
                >
                  No users found
                </TableCell>
              </TableRow>
            ) : (
              filteredUsers.map((user) => (
                <TableRow key={user.id}>
                  <TableCell className="font-medium">
                    {user.first_name} {user.last_name}
                  </TableCell>
                  <TableCell>{user.username}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>{getRoleBadge(user.role)}</TableCell>
                  <TableCell>{user.designation || "-"}</TableCell>
                  <TableCell>{getStatusBadge(user.status)}</TableCell>
                  <TableCell className="text-right space-x-2">
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={() => handleEdit(user)}
                    >
                      <IconEdit className="w-4 h-4" />
                    </Button>
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={() => setDeleteUser(user)}
                      className="text-red-600 hover:text-red-700"
                    >
                      <IconTrash className="w-4 h-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
