"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { apiService } from "@/lib/api-service";
import { API_ENDPOINTS } from "@/lib/config";
import { toast } from "sonner";
import {
  IconBuilding,
  IconMail,
  IconPhone,
  IconWorld,
  IconCheck,
  IconClock,
  IconX,
} from "@tabler/icons-react";

interface Company {
  id: string;
  name: string;
  email: string;
  phone?: string;
  industry?: string;
  website?: string;
  logo_url?: string;
  is_approved: boolean;
  is_active: boolean;
  created_at: number;
  updated_at?: number;
}

interface CompaniesResponse {
  success: boolean;
  message: string;
  data: Company[];
  meta: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

const formatTimeAgo = (timestamp: number): string => {
  const date = new Date(timestamp);
  const now = new Date();
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (seconds < 60) return "just now";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes} minute${minutes > 1 ? "s" : ""} ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} hour${hours > 1 ? "s" : ""} ago`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days} day${days > 1 ? "s" : ""} ago`;
  const months = Math.floor(days / 30);
  if (months < 12) return `${months} month${months > 1 ? "s" : ""} ago`;
  const years = Math.floor(days / 365);
  return `${years} year${years > 1 ? "s" : ""} ago`;
};

export function CompanyList() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    total_pages: 0,
  });

  const fetchCompanies = async () => {
    try {
      setLoading(true);
      const response = await apiService.get<CompaniesResponse>(
        `${API_ENDPOINTS.COMPANIES}?page=${pagination.page}&limit=${pagination.limit}`
      );

      if (response.success && response.data) {
        setCompanies(response.data);
        if (response.meta) {
          setPagination(response.meta);
        }
      }
    } catch (error) {
      console.error("Failed to fetch companies:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (companyId: string) => {
    try {
      setActionLoading(companyId);
      const response = await apiService.patch<{
        success: boolean;
        message: string;
      }>(`${API_ENDPOINTS.COMPANIES}/${companyId}/approve`, {});

      if (response.success) {
        toast.success(response.message || "Company approved successfully");
        // Update local state
        setCompanies((prev) =>
          prev.map((company) =>
            company.id === companyId
              ? { ...company, is_approved: true }
              : company
          )
        );
      } else {
        toast.error(response.message || "Failed to approve company");
      }
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to approve company";
      toast.error(errorMessage);
    } finally {
      setActionLoading(null);
    }
  };

  const handleDeactivate = async (companyId: string) => {
    try {
      setActionLoading(companyId);
      const response = await apiService.patch<{
        success: boolean;
        message: string;
      }>(`${API_ENDPOINTS.COMPANIES}/${companyId}/deactivate`, {});

      if (response.success) {
        toast.success(response.message || "Company deactivated successfully");
        // Update local state
        setCompanies((prev) =>
          prev.map((company) =>
            company.id === companyId
              ? { ...company, is_active: false }
              : company
          )
        );
      } else {
        toast.error(response.message || "Failed to deactivate company");
      }
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to deactivate company";
      toast.error(errorMessage);
    } finally {
      setActionLoading(null);
    }
  };

  useEffect(() => {
    fetchCompanies();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pagination.page]);

  if (loading) {
    return (
      <div className="space-y-4">
        {[...Array(3)].map((_, i) => (
          <Card key={i}>
            <CardHeader>
              <Skeleton className="h-6 w-48" />
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <Skeleton className="h-4 w-64" />
                <Skeleton className="h-4 w-48" />
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (companies.length === 0) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12">
          <IconBuilding className="h-12 w-12 text-muted-foreground mb-4" />
          <p className="text-lg font-medium text-muted-foreground">
            No companies found
          </p>
          <p className="text-sm text-muted-foreground mt-1">
            Companies will appear here once they are registered
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      {/* Summary Stats */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total Companies
            </CardTitle>
            <IconBuilding className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{pagination.total}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active</CardTitle>
            <IconCheck className="h-4 w-4 text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {companies.filter((c) => c.is_active && c.is_approved).length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Pending Approval
            </CardTitle>
            <IconClock className="h-4 w-4 text-orange-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {companies.filter((c) => !c.is_approved).length}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Company List */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {companies.map((company) => (
          <Card
            key={company.id}
            className="flex flex-col overflow-hidden transition-shadow hover:shadow-md"
          >
            {/* Header with Company Info and Status */}
            <CardHeader className="pb-3">
              <div className="flex flex-col gap-2">
                <div className="flex items-center justify-between">
                  <div className="flex size-12 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                    <IconBuilding className="size-6 text-primary" />
                  </div>
                  <div className="flex gap-1.5">
                    <Badge
                      variant={company.is_approved ? "default" : "secondary"}
                      className="h-5 text-xs"
                    >
                      {company.is_approved ? (
                        <>
                          <IconCheck className="mr-1 size-2.5" />
                          Approved
                        </>
                      ) : (
                        <>
                          <IconClock className="mr-1 size-2.5" />
                          Pending
                        </>
                      )}
                    </Badge>
                    <Badge
                      variant={company.is_active ? "default" : "destructive"}
                      className="h-5 text-xs"
                    >
                      {company.is_active ? "Active" : "Inactive"}
                    </Badge>
                  </div>
                </div>
                <div className="min-w-0">
                  <CardTitle className="text-base truncate leading-tight">
                    {company.name}
                  </CardTitle>
                  <p className="text-xs text-muted-foreground mt-0.5 truncate">
                    {company.industry || "No industry specified"}
                  </p>
                </div>
              </div>
            </CardHeader>

            {/* Company Details */}
            <CardContent className="flex-1 space-y-2 pb-3">
              <div className="space-y-2">
                {/* Email */}
                <div className="flex items-center gap-2 min-w-0">
                  <IconMail className="size-4 shrink-0 text-muted-foreground" />
                  <span className="text-xs text-muted-foreground shrink-0">
                    Email:
                  </span>
                  <span className="text-sm font-medium truncate">
                    {company.email}
                  </span>
                </div>

                {/* Phone */}
                {company.phone && (
                  <div className="flex items-center gap-2 min-w-0">
                    <IconPhone className="size-4 shrink-0 text-muted-foreground" />
                    <span className="text-xs text-muted-foreground shrink-0">
                      Phone:
                    </span>
                    <span className="text-sm font-medium truncate">
                      {company.phone}
                    </span>
                  </div>
                )}

                {/* Website */}
                {company.website && (
                  <div className="flex items-center gap-2 min-w-0">
                    <IconWorld className="size-4 shrink-0 text-muted-foreground" />
                    <span className="text-xs text-muted-foreground shrink-0">
                      Website:
                    </span>
                    <a
                      href={company.website}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-sm font-medium text-primary hover:underline truncate"
                    >
                      {company.website}
                    </a>
                  </div>
                )}

                {/* Registered */}
                <div className="flex items-center gap-2">
                  <IconClock className="size-4 shrink-0 text-muted-foreground" />
                  <span className="text-xs text-muted-foreground shrink-0">
                    Registered:
                  </span>
                  <span className="text-sm font-medium">
                    {formatTimeAgo(company.created_at)}
                  </span>
                </div>
              </div>

              {/* Action Buttons */}
              {(!company.is_approved || company.is_active) && (
                <div className="flex flex-col gap-2 pt-2 border-t mt-auto">
                  {!company.is_approved && (
                    <Button
                      onClick={() => handleApprove(company.id)}
                      disabled={actionLoading === company.id}
                      className="w-full"
                      size="sm"
                    >
                      <IconCheck className="mr-1.5 size-4" />
                      {actionLoading === company.id
                        ? "Approving..."
                        : "Approve"}
                    </Button>
                  )}
                  {company.is_active && (
                    <Button
                      onClick={() => handleDeactivate(company.id)}
                      disabled={actionLoading === company.id}
                      variant="destructive"
                      className="w-full"
                      size="sm"
                    >
                      <IconX className="mr-1.5 size-4" />
                      {actionLoading === company.id
                        ? "Deactivating..."
                        : "Deactivate"}
                    </Button>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Pagination Info */}
      {pagination.total_pages > 1 && (
        <Card>
          <CardContent className="py-4">
            <div className="flex items-center justify-between text-sm text-muted-foreground">
              <span>
                Showing {companies.length} of {pagination.total} companies
              </span>
              <span>
                Page {pagination.page} of {pagination.total_pages}
              </span>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
