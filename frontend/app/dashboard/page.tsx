"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { DashboardStats } from "@/components/dashboard/dashboard-stats";
import { CompanyList } from "@/components/dashboard/company-list";

export default function DashboardPage() {
  usePageTitle("Dashboard | WorkZen");
  const { user } = useRequireAuth();

  // SuperAdmin sees platform-wide stats and company list
  if (user?.role === "superadmin") {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Platform Dashboard</h1>
          <p className="text-muted-foreground">
            Platform-wide statistics and company management
          </p>
        </div>

        <DashboardStats />

        <div className="mt-8">
          <h2 className="text-2xl font-bold mb-4">Companies Management</h2>
          <CompanyList />
        </div>
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
