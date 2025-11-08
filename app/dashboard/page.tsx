"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { DashboardStats } from "@/components/dashboard/dashboard-stats";
import { CompanyList } from "@/components/dashboard/company-list";

export default function DashboardPage() {
  usePageTitle("Dashboard | WorkZen");
  const { user } = useRequireAuth();

  // SuperAdmin sees only company list
  if (user?.role === "superadmin") {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Companies Management</h1>
          <p className="text-muted-foreground">
            View and manage all companies on the platform
          </p>
        </div>

        <CompanyList />
      </div>
    );
  }

  // Other roles see dashboard stats
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">
          Welcome back, {user?.first_name}!
        </h1>
        <p className="text-muted-foreground">
          Here&apos;s what&apos;s happening with your organization today.
        </p>
      </div>

      <DashboardStats />
    </div>
  );
}
