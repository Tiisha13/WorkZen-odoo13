"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  IconUsers,
  IconBuilding,
  IconClock,
  IconCalendar,
} from "@tabler/icons-react";
import { useEffect, useState } from "react";
import { apiService } from "@/lib/api-service";
import { API_ENDPOINTS } from "@/lib/config";
import { Skeleton } from "@/components/ui/skeleton";

interface DashboardData {
  total_employees: number;
  total_departments: number;
  present_today: number;
  pending_leaves: number;
}

export function DashboardStats() {
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const user = apiService.getUser();
        const endpoint =
          user?.role === "superadmin"
            ? API_ENDPOINTS.DASHBOARD_SUPERADMIN
            : API_ENDPOINTS.DASHBOARD;

        const result = await apiService.get<{
          success: boolean;
          data: DashboardData;
        }>(endpoint);
        if (result.success && result.data) {
          setData(result.data);
        }
      } catch (error) {
        console.error("Failed to fetch dashboard data:", error);
        // Set default data on error
        setData({
          total_employees: 0,
          total_departments: 0,
          present_today: 0,
          pending_leaves: 0,
        });
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const stats = [
    {
      title: "Total Employees",
      value: data?.total_employees || 0,
      icon: IconUsers,
      color: "text-blue-600",
      bgColor: "bg-blue-100 dark:bg-blue-900/20",
    },
    {
      title: "Departments",
      value: data?.total_departments || 0,
      icon: IconBuilding,
      color: "text-green-600",
      bgColor: "bg-green-100 dark:bg-green-900/20",
    },
    {
      title: "Present Today",
      value: data?.present_today || 0,
      icon: IconClock,
      color: "text-purple-600",
      bgColor: "bg-purple-100 dark:bg-purple-900/20",
    },
    {
      title: "Pending Leaves",
      value: data?.pending_leaves || 0,
      icon: IconCalendar,
      color: "text-orange-600",
      bgColor: "bg-orange-100 dark:bg-orange-900/20",
    },
  ];

  if (loading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {[...Array(4)].map((_, i) => (
          <Card key={i}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <Skeleton className="h-4 w-24" />
              <Skeleton className="h-10 w-10 rounded-lg" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-16" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat) => (
        <Card key={stat.title}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
            <div
              className={`flex h-10 w-10 items-center justify-center rounded-lg ${stat.bgColor}`}
            >
              <stat.icon className={`h-5 w-5 ${stat.color}`} />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stat.value}</div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
