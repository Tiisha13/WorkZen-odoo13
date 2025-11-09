"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { ProfileForm } from "@/components/profile/profile-form";
import { ChangePasswordForm } from "@/components/profile/change-password-form";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useState, useEffect } from "react";
import { apiService } from "@/lib/api-service";
import { API_ENDPOINTS } from "@/lib/config";
import { BankDetails } from "@/lib/types";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  IconCurrencyDollar,
  IconBuildingBank,
  IconEdit,
  IconReport,
} from "@tabler/icons-react";
import { toast } from "sonner";

interface Salary {
  id: string;
  monthly_wage: number;
  yearly_wage: number;
  currency: string;
  effective_from: string;
  is_active: boolean;
}

export default function ProfilePage() {
  usePageTitle("Profile Settings | WorkZen");
  const { user, isLoading } = useRequireAuth();

  const [salary, setSalary] = useState<Salary | null>(null);
  const [loadingSalary, setLoadingSalary] = useState(true);
  const [isBankDialogOpen, setIsBankDialogOpen] = useState(false);
  const [bankDetails, setBankDetails] = useState<BankDetails>({
    account_number: "",
    bank_name: "",
    ifsc_code: "",
    branch_name: "",
    pan_no: "",
    uan_no: "",
  });

  useEffect(() => {
    if (user?.id) {
      fetchSalary();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user?.id]);

  // Load bank details from user object
  useEffect(() => {
    if (user?.bank_details) {
      setBankDetails({
        account_number: user.bank_details.account_number || "",
        bank_name: user.bank_details.bank_name || "",
        ifsc_code: user.bank_details.ifsc_code || "",
        branch_name: user.bank_details.branch_name || "",
        pan_no: user.bank_details.pan_no || "",
        uan_no: user.bank_details.uan_no || "",
      });
    }
  }, [user?.bank_details]);

  const fetchSalary = async () => {
    if (!user?.id) return;

    try {
      setLoadingSalary(true);
      const response = await apiService.get<{ success: boolean; data: Salary }>(
        `${API_ENDPOINTS.SALARY}/${user.id}`
      );
      setSalary(response.data);
    } catch (error) {
      // Silently handle 404 - employee simply doesn't have salary structure yet
      if (error instanceof Error && error.message.includes("not found")) {
        setSalary(null);
      } else {
        console.error("Failed to fetch salary:", error);
      }
    } finally {
      setLoadingSalary(false);
    }
  };

  const handleBankDetailsSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!user?.id) return;

    try {
      await apiService.patch(`${API_ENDPOINTS.USERS}/${user.id}/bank`, {
        account_number: bankDetails.account_number,
        bank_name: bankDetails.bank_name,
        ifsc_code: bankDetails.ifsc_code,
        branch_name: bankDetails.branch_name,
        pan_no: bankDetails.pan_no,
        uan_no: bankDetails.uan_no,
      });
      toast.success("Bank details updated successfully");
      setIsBankDialogOpen(false);

      // Refresh user data to show updated bank details
      window.location.reload();
    } catch (error) {
      const message =
        error instanceof Error
          ? error.message
          : "Failed to update bank details";
      toast.error(message);
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(amount);
  };

  if (isLoading || !user) {
    return (
      <DashboardLayout>
        <div className="flex h-screen items-center justify-center">
          <Skeleton className="h-96 w-full max-w-2xl" />
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="mx-auto max-w-4xl space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Profile Settings</h1>
          <p className="text-muted-foreground">
            Manage your account settings and preferences
          </p>
        </div>

        <Tabs defaultValue="profile" className="w-full">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="profile">Profile</TabsTrigger>
            <TabsTrigger value="security">Security</TabsTrigger>
            <TabsTrigger value="salary">Salary & Bank</TabsTrigger>
          </TabsList>
          <TabsContent value="profile" className="space-y-4">
            <ProfileForm user={user} />
          </TabsContent>
          <TabsContent value="security" className="space-y-4">
            <ChangePasswordForm />
          </TabsContent>
          <TabsContent value="salary" className="space-y-4">
            <div className="grid gap-6 md:grid-cols-2">
              {/* Salary Information Card */}
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle className="flex items-center gap-2">
                        <IconCurrencyDollar className="w-5 h-5" />
                        Salary Information
                      </CardTitle>
                      <CardDescription>
                        Your current salary structure
                      </CardDescription>
                    </div>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => {
                        toast.info("Salary reports feature coming soon!");
                      }}
                    >
                      <IconReport className="w-4 h-4 mr-2" />
                      Reports
                    </Button>
                  </div>
                </CardHeader>
                <CardContent>
                  {loadingSalary ? (
                    <div className="flex items-center justify-center py-8">
                      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                    </div>
                  ) : salary ? (
                    <div className="space-y-4">
                      <div>
                        <p className="text-sm text-muted-foreground">
                          Monthly Wage
                        </p>
                        <p className="text-2xl font-bold">
                          {formatCurrency(salary.monthly_wage)}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground">
                          Yearly Wage
                        </p>
                        <p className="text-xl font-semibold">
                          {formatCurrency(salary.yearly_wage)}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground">
                          Effective From
                        </p>
                        <p className="font-medium">
                          {new Date(salary.effective_from).toLocaleDateString(
                            "en-US",
                            {
                              year: "numeric",
                              month: "long",
                              day: "numeric",
                            }
                          )}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground">Status</p>
                        <Badge
                          variant={salary.is_active ? "default" : "secondary"}
                          className={
                            salary.is_active
                              ? "bg-green-100 text-green-800"
                              : "bg-gray-100 text-gray-800"
                          }
                        >
                          {salary.is_active ? "Active" : "Inactive"}
                        </Badge>
                      </div>
                    </div>
                  ) : (
                    <div className="text-center py-8 text-muted-foreground">
                      <IconCurrencyDollar className="w-12 h-12 mx-auto mb-2 opacity-50" />
                      <p>No salary information available</p>
                      <p className="text-sm">Contact HR for salary details</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Bank Details Card */}
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle className="flex items-center gap-2">
                        <IconBuildingBank className="w-5 h-5" />
                        Bank Details
                      </CardTitle>
                      <CardDescription>
                        Manage your banking information
                      </CardDescription>
                    </div>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setIsBankDialogOpen(true)}
                    >
                      <IconEdit className="w-4 h-4 mr-2" />
                      Update
                    </Button>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div>
                      <p className="text-sm text-muted-foreground">
                        Account Number
                      </p>
                      <p className="font-medium">
                        {bankDetails.account_number
                          ? `****${bankDetails.account_number.slice(-4)}`
                          : "Not provided"}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">Bank Name</p>
                      <p className="font-medium">
                        {bankDetails.bank_name || "Not provided"}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">IFSC Code</p>
                      <p className="font-medium">
                        {bankDetails.ifsc_code || "Not provided"}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">
                        Branch Name
                      </p>
                      <p className="font-medium">
                        {bankDetails.branch_name || "Not provided"}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">
                        PAN Number
                      </p>
                      <p className="font-medium">
                        {bankDetails.pan_no || "Not provided"}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">
                        UAN Number
                      </p>
                      <p className="font-medium">
                        {bankDetails.uan_no || "Not provided"}
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>
        </Tabs>

        {/* Bank Details Dialog */}
        <Dialog open={isBankDialogOpen} onOpenChange={setIsBankDialogOpen}>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle>Update Bank Details</DialogTitle>
              <DialogDescription>
                Enter your banking information for salary payments
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleBankDetailsSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="account_number">
                  Account Number <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="account_number"
                  value={bankDetails.account_number}
                  onChange={(e) =>
                    setBankDetails({
                      ...bankDetails,
                      account_number: e.target.value,
                    })
                  }
                  placeholder="1234567890"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="bank_name">
                  Bank Name <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="bank_name"
                  value={bankDetails.bank_name}
                  onChange={(e) =>
                    setBankDetails({
                      ...bankDetails,
                      bank_name: e.target.value,
                    })
                  }
                  placeholder="ABC Bank"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="ifsc_code">
                  IFSC Code <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="ifsc_code"
                  value={bankDetails.ifsc_code}
                  onChange={(e) =>
                    setBankDetails({
                      ...bankDetails,
                      ifsc_code: e.target.value,
                    })
                  }
                  placeholder="ABCD0123456"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="branch_name">Branch Name</Label>
                <Input
                  id="branch_name"
                  value={bankDetails.branch_name}
                  onChange={(e) =>
                    setBankDetails({
                      ...bankDetails,
                      branch_name: e.target.value,
                    })
                  }
                  placeholder="Main Branch"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="pan_no">PAN Number</Label>
                <Input
                  id="pan_no"
                  value={bankDetails.pan_no}
                  onChange={(e) =>
                    setBankDetails({
                      ...bankDetails,
                      pan_no: e.target.value,
                    })
                  }
                  placeholder="ABCDE1234F"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="uan_no">UAN Number</Label>
                <Input
                  id="uan_no"
                  value={bankDetails.uan_no}
                  onChange={(e) =>
                    setBankDetails({
                      ...bankDetails,
                      uan_no: e.target.value,
                    })
                  }
                  placeholder="123456789012"
                />
              </div>
              <DialogFooter className="gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setIsBankDialogOpen(false)}
                >
                  Cancel
                </Button>
                <Button type="submit">Save Changes</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  );
}
