"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { useState, useEffect } from "react";
import { apiService } from "@/lib/api-service";
import { API_ENDPOINTS } from "@/lib/config";
import { Button } from "@/components/ui/button";
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
import { IconPlus, IconDownload } from "@tabler/icons-react";
import { toast } from "sonner";
import type { Salary } from "@/lib/types";

export default function PayrollPage() {
  usePageTitle("Payroll Management | WorkZen");
  useRequireAuth(["superadmin", "admin", "payroll"]);

  const [payrolls, setPayrolls] = useState<Salary[]>([]);
  const [loading, setLoading] = useState(true);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedMonth, setSelectedMonth] = useState<string>(
    new Date().toISOString().slice(0, 7)
  );

  useEffect(() => {
    fetchPayrolls();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedMonth]);

  const fetchPayrolls = async () => {
    try {
      setLoading(true);
      const response = await apiService.get<{
        success: boolean;
        data: Salary[];
      }>(`${API_ENDPOINTS.PAYROLL}?month=${selectedMonth}`);
      setPayrolls(response.data || []);
    } catch (error) {
      toast.error("Failed to fetch payroll data");
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const handleGeneratePayroll = async () => {
    try {
      await apiService.post(`${API_ENDPOINTS.PAYROLL}/generate`, {
        month: selectedMonth,
      });
      toast.success("Payroll generated successfully");
      setIsDialogOpen(false);
      fetchPayrolls();
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to generate payroll"
      );
    }
  };

  const handleDownloadPayslip = async (id: string) => {
    try {
      toast.info("Downloading payslip...");
      // In a real app, this would download a PDF
      await apiService.get(`${API_ENDPOINTS.PAYROLL}/${id}/payslip`);
      toast.success("Payslip downloaded");
    } catch (error) {
      toast.error("Failed to download payslip");
      console.error(error);
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
    });
  };

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      paid: "bg-green-100 text-green-800",
      pending: "bg-yellow-100 text-yellow-800",
      processing: "bg-blue-100 text-blue-800",
    };
    return (
      <Badge className={colors[status] || "bg-gray-100 text-gray-800"}>
        {status.toUpperCase()}
      </Badge>
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Payroll Management</h1>
          <p className="text-muted-foreground">
            Manage employee salaries and payroll
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <IconPlus className="w-4 h-4 mr-2" />
              Generate Payroll
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Generate Payroll</DialogTitle>
              <DialogDescription>
                Generate payroll for selected month
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="month">Select Month</Label>
                <input
                  id="month"
                  type="month"
                  value={selectedMonth}
                  onChange={(e) => setSelectedMonth(e.target.value)}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setIsDialogOpen(false)}
              >
                Cancel
              </Button>
              <Button onClick={handleGeneratePayroll}>Generate</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <div className="flex items-center space-x-4">
        <Label htmlFor="filter-month">Filter by Month:</Label>
        <input
          id="filter-month"
          type="month"
          value={selectedMonth}
          onChange={(e) => setSelectedMonth(e.target.value)}
          className="flex h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
        />
      </div>

      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Employee</TableHead>
              <TableHead>Employee ID</TableHead>
              <TableHead>Basic Salary</TableHead>
              <TableHead>Allowances</TableHead>
              <TableHead>Deductions</TableHead>
              <TableHead>Net Salary</TableHead>
              <TableHead>Month</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={9} className="text-center py-8">
                  Loading...
                </TableCell>
              </TableRow>
            ) : payrolls.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={9}
                  className="text-center py-8 text-muted-foreground"
                >
                  No payroll records found for selected month
                </TableCell>
              </TableRow>
            ) : (
              payrolls.map((payroll) => (
                <TableRow key={payroll.id}>
                  <TableCell className="font-medium">
                    {payroll.user?.first_name} {payroll.user?.last_name}
                  </TableCell>
                  <TableCell>{payroll.user_id}</TableCell>
                  <TableCell>{formatCurrency(payroll.basic_salary)}</TableCell>
                  <TableCell>
                    {formatCurrency(payroll.allowances || 0)}
                  </TableCell>
                  <TableCell>
                    {formatCurrency(payroll.deductions || 0)}
                  </TableCell>
                  <TableCell className="font-semibold">
                    {formatCurrency(payroll.net_salary)}
                  </TableCell>
                  <TableCell>{formatDate(payroll.month)}</TableCell>
                  <TableCell>{getStatusBadge(payroll.status)}</TableCell>
                  <TableCell className="text-right">
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={() => handleDownloadPayslip(payroll.id)}
                    >
                      <IconDownload className="w-4 h-4" />
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
