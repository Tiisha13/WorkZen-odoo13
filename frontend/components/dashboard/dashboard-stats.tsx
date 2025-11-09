"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  IconUsers,
  IconBuilding,
  IconClock,
  IconCalendar,
  IconUserCheck,
  IconUserX,
  IconUserPause,
  IconPercentage,
} from "@tabler/icons-react";
import { useEffect, useState } from "react";
import { apiService } from "@/lib/api-service";
import { API_ENDPOINTS } from "@/lib/config";
import { Skeleton } from "@/components/ui/skeleton";
import type { DashboardData } from "@/lib/types";
import {
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

const COLORS = {
  primary: "#3b82f6",
  success: "#22c55e",
  warning: "#f59e0b",
  danger: "#ef4444",
  purple: "#a855f7",
  cyan: "#06b6d4",
};

export function DashboardStats() {
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [userRole, setUserRole] = useState<string>("");

  useEffect(() => {
    const fetchData = async () => {
      try {
        const user = apiService.getUser();
        setUserRole(user?.role || "");

        // Select endpoint based on role
        let endpoint = API_ENDPOINTS.DASHBOARD; // Default for all users
        if (user?.role === "superadmin") {
          endpoint = API_ENDPOINTS.DASHBOARD_SUPERADMIN;
        } else if (user?.role === "admin") {
          endpoint = API_ENDPOINTS.DASHBOARD_ADMIN;
        }

        const result = await apiService.get<{
          success: boolean;
          data: DashboardData;
        }>(endpoint);
        if (result.success && result.data) {
          setData(result.data);
        }
      } catch (error) {
        console.error("Failed to fetch dashboard data:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  // SuperAdmin-specific stats
  const superAdminStats = [
    {
      title: "Total Companies",
      value: data?.total_companies || 0,
      icon: IconBuilding,
      color: "text-blue-600",
      bgColor: "bg-blue-100 dark:bg-blue-900/20",
    },
    {
      title: "Active Companies",
      value: data?.active_companies || 0,
      icon: IconUserCheck,
      color: "text-green-600",
      bgColor: "bg-green-100 dark:bg-green-900/20",
    },
    {
      title: "Pending Approvals",
      value: data?.pending_approvals || 0,
      icon: IconClock,
      color: "text-yellow-600",
      bgColor: "bg-yellow-100 dark:bg-yellow-900/20",
    },
    {
      title: "Total Employees",
      value: data?.total_employees || 0,
      icon: IconUsers,
      color: "text-purple-600",
      bgColor: "bg-purple-100 dark:bg-purple-900/20",
    },
    {
      title: "Present Today",
      value: data?.present_today || 0,
      icon: IconUserCheck,
      color: "text-cyan-600",
      bgColor: "bg-cyan-100 dark:bg-cyan-900/20",
    },
    {
      title: "Attendance Rate",
      value: `${data?.attendance_rate?.toFixed(1) || 0}%`,
      icon: IconPercentage,
      color: "text-emerald-600",
      bgColor: "bg-emerald-100 dark:bg-emerald-900/20",
    },
    {
      title: "Total Departments",
      value: data?.total_departments || 0,
      icon: IconBuilding,
      color: "text-indigo-600",
      bgColor: "bg-indigo-100 dark:bg-indigo-900/20",
    },
    {
      title: "Pending Leaves",
      value: data?.pending_leaves || 0,
      icon: IconCalendar,
      color: "text-orange-600",
      bgColor: "bg-orange-100 dark:bg-orange-900/20",
    },
  ];

  // Company-level stats (for other roles)
  const companyStats = [
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
      icon: IconUserCheck,
      color: "text-purple-600",
      bgColor: "bg-purple-100 dark:bg-purple-900/20",
    },
    {
      title: "Attendance Rate",
      value: `${data?.attendance_rate?.toFixed(1) || 0}%`,
      icon: IconPercentage,
      color: "text-cyan-600",
      bgColor: "bg-cyan-100 dark:bg-cyan-900/20",
    },
    {
      title: "Absent Today",
      value: data?.absent_today || 0,
      icon: IconUserX,
      color: "text-red-600",
      bgColor: "bg-red-100 dark:bg-red-900/20",
    },
    {
      title: "On Leave Today",
      value: data?.on_leave_today || 0,
      icon: IconUserPause,
      color: "text-orange-600",
      bgColor: "bg-orange-100 dark:bg-orange-900/20",
    },
    {
      title: "Pending Leaves",
      value: data?.pending_leaves || 0,
      icon: IconCalendar,
      color: "text-yellow-600",
      bgColor: "bg-yellow-100 dark:bg-yellow-900/20",
    },
    {
      title: "Approved Leaves",
      value: data?.approved_leaves || 0,
      icon: IconClock,
      color: "text-emerald-600",
      bgColor: "bg-emerald-100 dark:bg-emerald-900/20",
    },
  ];

  const stats = userRole === "superadmin" ? superAdminStats : companyStats;

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[...Array(8)].map((_, i) => (
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
        <div className="grid gap-4 md:grid-cols-2">
          {[...Array(4)].map((_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-6 w-48" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-[300px] w-full" />
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">No dashboard data available</p>
      </div>
    );
  }

  // Prepare department stats for bar chart
  const departmentChartData =
    data.department_stats?.map((dept) => ({
      name: dept.name,
      Total: dept.count,
      Present: dept.present,
      Absent: dept.absent,
      "On Leave": dept.on_leave,
    })) || [];

  // Prepare monthly attendance for line/area chart
  const monthlyChartData = data.monthly_attendance || [];

  // Prepare leave type stats for pie chart
  const leaveTypeChartData =
    data.leave_type_stats?.map((leave) => ({
      name: leave.type.charAt(0).toUpperCase() + leave.type.slice(1),
      Pending: leave.pending,
      Approved: leave.approved,
      Rejected: leave.rejected,
      Total: leave.pending + leave.approved + leave.rejected,
    })) || [];

  // Prepare today's attendance for pie chart
  const todayAttendanceData = [
    { name: "Present", value: data.present_today, color: COLORS.success },
    { name: "Absent", value: data.absent_today, color: COLORS.danger },
    { name: "On Leave", value: data.on_leave_today, color: COLORS.warning },
  ];

  return (
    <div className="space-y-6">
      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {stat.title}
              </CardTitle>
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

      {/* Charts Section */}
      <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        {/* Department-wise Statistics / Company-wise for SuperAdmin */}
        {departmentChartData.length > 0 && (
          <Card className="col-span-1 md:col-span-2 lg:col-span-3">
            <CardHeader>
              <CardTitle>
                {userRole === "superadmin"
                  ? "Top Companies by Employee Count"
                  : "Department-wise Attendance"}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={350}>
                <BarChart data={departmentChartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Bar dataKey="Present" fill={COLORS.success} />
                  <Bar dataKey="Absent" fill={COLORS.danger} />
                  <Bar dataKey="On Leave" fill={COLORS.warning} />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        )}

        {/* Monthly Attendance Trend */}
        {monthlyChartData.length > 0 && (
          <Card className="col-span-1 md:col-span-2 lg:col-span-3">
            <CardHeader>
              <CardTitle>
                {userRole === "superadmin"
                  ? "Platform-wide Attendance Trend (Last 6 Months)"
                  : "Monthly Attendance Trend (Last 6 Months)"}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={350}>
                <AreaChart data={monthlyChartData}>
                  <defs>
                    <linearGradient
                      id="colorPresent"
                      x1="0"
                      y1="0"
                      x2="0"
                      y2="1"
                    >
                      <stop
                        offset="5%"
                        stopColor={COLORS.success}
                        stopOpacity={0.8}
                      />
                      <stop
                        offset="95%"
                        stopColor={COLORS.success}
                        stopOpacity={0}
                      />
                    </linearGradient>
                    <linearGradient
                      id="colorAbsent"
                      x1="0"
                      y1="0"
                      x2="0"
                      y2="1"
                    >
                      <stop
                        offset="5%"
                        stopColor={COLORS.danger}
                        stopOpacity={0.8}
                      />
                      <stop
                        offset="95%"
                        stopColor={COLORS.danger}
                        stopOpacity={0}
                      />
                    </linearGradient>
                    <linearGradient id="colorLeave" x1="0" y1="0" x2="0" y2="1">
                      <stop
                        offset="5%"
                        stopColor={COLORS.warning}
                        stopOpacity={0.8}
                      />
                      <stop
                        offset="95%"
                        stopColor={COLORS.warning}
                        stopOpacity={0}
                      />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="month" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Area
                    type="monotone"
                    dataKey="present"
                    stroke={COLORS.success}
                    fillOpacity={1}
                    fill="url(#colorPresent)"
                    name="Present"
                  />
                  <Area
                    type="monotone"
                    dataKey="absent"
                    stroke={COLORS.danger}
                    fillOpacity={1}
                    fill="url(#colorAbsent)"
                    name="Absent"
                  />
                  <Area
                    type="monotone"
                    dataKey="on_leave"
                    stroke={COLORS.warning}
                    fillOpacity={1}
                    fill="url(#colorLeave)"
                    name="On Leave"
                  />
                </AreaChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        )}

        {/* Today's Attendance Distribution */}
        <Card className="col-span-1">
          <CardHeader>
            <CardTitle>Today&apos;s Attendance Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={todayAttendanceData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, value }) => `${name}: ${value}`}
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {todayAttendanceData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Leave Type Statistics */}
        {leaveTypeChartData.length > 0 && (
          <Card className="col-span-1 md:col-span-2 lg:col-span-2">
            <CardHeader>
              <CardTitle>Leave Types Overview</CardTitle>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={leaveTypeChartData} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis type="number" />
                  <YAxis dataKey="name" type="category" width={80} />
                  <Tooltip />
                  <Legend />
                  <Bar dataKey="Pending" stackId="a" fill={COLORS.warning} />
                  <Bar dataKey="Approved" stackId="a" fill={COLORS.success} />
                  <Bar dataKey="Rejected" stackId="a" fill={COLORS.danger} />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
