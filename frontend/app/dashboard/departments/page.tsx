"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { useState, useEffect } from "react";
import { apiService } from "@/lib/api-service";
import { API_ENDPOINTS } from "@/lib/config";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
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
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  IconPlus,
  IconEdit,
  IconTrash,
  IconSearch,
  IconUsers,
} from "@tabler/icons-react";
import { toast } from "sonner";
import type { Department } from "@/lib/types";

export default function DepartmentsPage() {
  usePageTitle("Department Management | WorkZen");
  useRequireAuth(["superadmin", "admin", "hr"]);

  const [departments, setDepartments] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingDept, setEditingDept] = useState<Department | null>(null);

  const [formData, setFormData] = useState({
    name: "",
    description: "",
    head_id: "",
  });

  useEffect(() => {
    fetchDepartments();
  }, []);

  const fetchDepartments = async () => {
    try {
      setLoading(true);
      const response = await apiService.get<{
        success: boolean;
        data: Department[];
      }>(API_ENDPOINTS.DEPARTMENTS);
      setDepartments(response.data || []);
    } catch (error) {
      toast.error("Failed to fetch departments");
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingDept) {
        await apiService.patch(
          `${API_ENDPOINTS.DEPARTMENTS}/${editingDept.id}`,
          formData
        );
        toast.success("Department updated successfully");
      } else {
        await apiService.post(API_ENDPOINTS.DEPARTMENTS, formData);
        toast.success("Department created successfully");
      }
      setIsDialogOpen(false);
      resetForm();
      fetchDepartments();
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to save department"
      );
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this department?")) return;
    try {
      await apiService.delete(`${API_ENDPOINTS.DEPARTMENTS}/${id}`);
      toast.success("Department deleted successfully");
      fetchDepartments();
    } catch (error) {
      toast.error("Failed to delete department");
      console.error(error);
    }
  };

  const resetForm = () => {
    setFormData({ name: "", description: "", head_id: "" });
    setEditingDept(null);
  };

  const handleEdit = (dept: Department) => {
    setEditingDept(dept);
    setFormData({
      name: dept.name,
      description: dept.description || "",
      head_id: dept.head_id || "",
    });
    setIsDialogOpen(true);
  };

  const filteredDepartments = departments.filter(
    (dept) =>
      dept.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (dept.description &&
        dept.description.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Department Management</h1>
          <p className="text-muted-foreground">Manage company departments</p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={resetForm}>
              <IconPlus className="w-4 h-4 mr-2" />
              Add Department
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {editingDept ? "Edit Department" : "Add New Department"}
              </DialogTitle>
              <DialogDescription>
                {editingDept
                  ? "Update department information"
                  : "Create a new department"}
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">Department Name *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <Input
                  id="description"
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                />
              </div>
              <DialogFooter>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setIsDialogOpen(false)}
                >
                  Cancel
                </Button>
                <Button type="submit">
                  {editingDept ? "Update Department" : "Create Department"}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Departments</CardTitle>
              <CardDescription>
                {filteredDepartments.length} department
                {filteredDepartments.length !== 1 ? "s" : ""} total
              </CardDescription>
            </div>
            <div className="relative w-72">
              <IconSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
              <Input
                placeholder="Search departments..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border bg-card shadow-sm overflow-hidden">
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Department Name</TableHead>
                    <TableHead>Description</TableHead>
                    <TableHead>
                      <div className="flex items-center gap-2">
                        <IconUsers className="w-4 h-4" />
                        Employees
                      </div>
                    </TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {loading ? (
                    <TableRow>
                      <TableCell colSpan={4} className="text-center py-12">
                        <div className="flex flex-col items-center gap-2">
                          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                          <p className="text-sm text-muted-foreground">
                            Loading departments...
                          </p>
                        </div>
                      </TableCell>
                    </TableRow>
                  ) : filteredDepartments.length === 0 ? (
                    <TableRow>
                      <TableCell
                        colSpan={4}
                        className="text-center py-12 text-muted-foreground"
                      >
                        <div className="flex flex-col items-center gap-2">
                          <IconSearch className="w-12 h-12 opacity-20" />
                          <p className="font-medium">No departments found</p>
                          <p className="text-sm">
                            {searchTerm
                              ? "Try adjusting your search"
                              : "Click 'Add Department' to create one"}
                          </p>
                        </div>
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredDepartments.map((dept) => (
                      <TableRow
                        key={dept.id}
                        className="hover:bg-muted/50 transition-colors"
                      >
                        <TableCell className="font-medium">
                          {dept.name}
                        </TableCell>
                        <TableCell className="text-muted-foreground">
                          {dept.description || (
                            <span className="italic">No description</span>
                          )}
                        </TableCell>
                        <TableCell>
                          <Badge variant="secondary" className="font-mono">
                            {dept.employee_count || 0}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex items-center justify-end gap-2">
                            <Button
                              size="sm"
                              variant="ghost"
                              onClick={() => handleEdit(dept)}
                              className="hover:bg-blue-50 hover:text-blue-600"
                            >
                              <IconEdit className="w-4 h-4" />
                            </Button>
                            <Button
                              size="sm"
                              variant="ghost"
                              onClick={() => handleDelete(dept.id)}
                              className="hover:bg-red-50 hover:text-red-600"
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
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
