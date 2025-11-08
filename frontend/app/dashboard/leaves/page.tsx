"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { useState, useEffect } from "react";
import { apiService } from "@/lib/api-service";
import { API_ENDPOINTS } from "@/lib/config";
import { useAuth } from "@/lib/auth-context";
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { IconPlus, IconCheck, IconX } from "@tabler/icons-react";
import { toast } from "sonner";
import type { Leave } from "@/lib/types";

export default function LeavesPage() {
  usePageTitle("Leave Management | WorkZen");
  useRequireAuth();

  const { user } = useAuth();
  const [leaves, setLeaves] = useState<Leave[]>([]);
  const [loading, setLoading] = useState(true);
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const [formData, setFormData] = useState({
    leave_type: "sick",
    start_date: "",
    end_date: "",
    reason: "",
  });

  useEffect(() => {
    fetchLeaves();
  }, []);

  const fetchLeaves = async () => {
    try {
      setLoading(true);
      const response = await apiService.get<{
        success: boolean;
        data: Leave[];
      }>(API_ENDPOINTS.LEAVES);
      setLeaves(response.data || []);
    } catch (error) {
      toast.error("Failed to fetch leaves");
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await apiService.post(API_ENDPOINTS.LEAVES, {
        ...formData,
        user_id: user?.id,
      });
      toast.success("Leave request submitted successfully");
      setIsDialogOpen(false);
      resetForm();
      fetchLeaves();
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to submit leave"
      );
    }
  };

  const handleApprove = async (id: string) => {
    try {
      await apiService.patch(`${API_ENDPOINTS.LEAVES}/${id}/approve`, {
        status: "approved",
      });
      toast.success("Leave approved");
      fetchLeaves();
    } catch (error) {
      toast.error("Failed to approve leave");
      console.error(error);
    }
  };

  const handleReject = async (id: string) => {
    try {
      await apiService.patch(`${API_ENDPOINTS.LEAVES}/${id}/reject`, {
        status: "rejected",
      });
      toast.success("Leave rejected");
      fetchLeaves();
    } catch (error) {
      toast.error("Failed to reject leave");
      console.error(error);
    }
  };

  const resetForm = () => {
    setFormData({
      leave_type: "sick",
      start_date: "",
      end_date: "",
      reason: "",
    });
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const getStatusBadge = (status: string) => {
    const config: Record<
      string,
      {
        variant: "default" | "secondary" | "destructive" | "outline";
        label: string;
        className?: string;
      }
    > = {
      pending: {
        variant: "outline",
        label: "Pending",
        className: "border-yellow-200 bg-yellow-50 text-yellow-700",
      },
      approved: {
        variant: "outline",
        label: "Approved",
        className: "border-green-200 bg-green-50 text-green-700",
      },
      rejected: {
        variant: "outline",
        label: "Rejected",
        className: "border-red-200 bg-red-50 text-red-700",
      },
    };
    const conf = config[status] || {
      variant: "secondary" as const,
      label: status,
      className: "",
    };
    return (
      <Badge variant={conf.variant} className={conf.className}>
        {conf.label}
      </Badge>
    );
  };

  const getTypeBadge = (type: string) => {
    const config: Record<
      string,
      {
        variant: "default" | "secondary" | "destructive" | "outline";
        label: string;
        className?: string;
      }
    > = {
      sick: {
        variant: "outline",
        label: "Sick",
        className: "border-blue-200 bg-blue-50 text-blue-700",
      },
      casual: {
        variant: "outline",
        label: "Casual",
        className: "border-purple-200 bg-purple-50 text-purple-700",
      },
      annual: {
        variant: "outline",
        label: "Annual",
        className: "border-green-200 bg-green-50 text-green-700",
      },
      unpaid: {
        variant: "outline",
        label: "Unpaid",
        className: "border-gray-200 bg-gray-50 text-gray-700",
      },
    };
    const conf = config[type] || {
      variant: "secondary" as const,
      label: type,
      className: "",
    };
    return (
      <Badge variant={conf.variant} className={conf.className}>
        {conf.label}
      </Badge>
    );
  };

  const canApprove =
    user?.role === "hr" ||
    user?.role === "admin" ||
    user?.role === "superadmin";

  return (
    <div className="flex flex-col gap-6 p-6">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold tracking-tight">
            Leave Management
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            Request and manage leaves
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={resetForm}>
              <IconPlus className="w-4 h-4 mr-2" />
              Request Leave
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Request Leave</DialogTitle>
              <DialogDescription>Submit a new leave request</DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="leave_type">Leave Type *</Label>
                <Select
                  value={formData.leave_type}
                  onValueChange={(value) =>
                    setFormData({ ...formData, leave_type: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="sick">Sick Leave</SelectItem>
                    <SelectItem value="casual">Casual Leave</SelectItem>
                    <SelectItem value="annual">Annual Leave</SelectItem>
                    <SelectItem value="unpaid">Unpaid Leave</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="start_date">Start Date *</Label>
                <Input
                  id="start_date"
                  type="date"
                  value={formData.start_date}
                  onChange={(e) =>
                    setFormData({ ...formData, start_date: e.target.value })
                  }
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="end_date">End Date *</Label>
                <Input
                  id="end_date"
                  type="date"
                  value={formData.end_date}
                  onChange={(e) =>
                    setFormData({ ...formData, end_date: e.target.value })
                  }
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="reason">Reason *</Label>
                <Textarea
                  id="reason"
                  value={formData.reason}
                  onChange={(e) =>
                    setFormData({ ...formData, reason: e.target.value })
                  }
                  required
                  rows={3}
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
                <Button type="submit">Submit Request</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <div className="bg-card rounded-lg border shadow-sm">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Employee</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Start Date</TableHead>
              <TableHead>End Date</TableHead>
              <TableHead>Days</TableHead>
              <TableHead>Reason</TableHead>
              <TableHead>Status</TableHead>
              {canApprove && (
                <TableHead className="text-right">Actions</TableHead>
              )}
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell
                  colSpan={canApprove ? 8 : 7}
                  className="h-32 text-center"
                >
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                  </div>
                </TableCell>
              </TableRow>
            ) : leaves.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={canApprove ? 8 : 7}
                  className="h-32 text-center"
                >
                  <div className="flex flex-col items-center justify-center text-muted-foreground">
                    <p className="text-sm">No leave requests found</p>
                    <p className="text-xs mt-1">
                      Submit a new leave request to get started
                    </p>
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              leaves.map((leave) => (
                <TableRow key={leave.id} className="hover:bg-muted/50">
                  <TableCell className="font-medium">
                    {leave.user?.first_name} {leave.user?.last_name}
                  </TableCell>
                  <TableCell>{getTypeBadge(leave.leave_type)}</TableCell>
                  <TableCell className="text-sm">
                    {formatDate(leave.start_date)}
                  </TableCell>
                  <TableCell className="text-sm">
                    {formatDate(leave.end_date)}
                  </TableCell>
                  <TableCell className="text-sm">{leave.days || 0}</TableCell>
                  <TableCell className="max-w-xs truncate text-sm">
                    {leave.reason}
                  </TableCell>
                  <TableCell>{getStatusBadge(leave.status)}</TableCell>
                  {canApprove && (
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-2">
                        {leave.status === "pending" && (
                          <>
                            <Button
                              size="icon"
                              variant="ghost"
                              onClick={() => handleApprove(leave.id)}
                              className="h-8 w-8 text-green-600 hover:text-green-700"
                            >
                              <IconCheck className="w-4 h-4" />
                            </Button>
                            <Button
                              size="icon"
                              variant="ghost"
                              onClick={() => handleReject(leave.id)}
                              className="h-8 w-8 text-red-600 hover:text-red-700"
                            >
                              <IconX className="w-4 h-4" />
                            </Button>
                          </>
                        )}
                      </div>
                    </TableCell>
                  )}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
