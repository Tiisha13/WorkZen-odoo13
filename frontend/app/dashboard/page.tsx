"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { DashboardStats } from "@/components/dashboard/dashboard-stats";

export default function DashboardPage() {
  usePageTitle("Dashboard | WorkZen");
  const { user } = useRequireAuth();

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
