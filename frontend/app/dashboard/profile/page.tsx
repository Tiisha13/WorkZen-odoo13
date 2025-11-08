"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { useAuth } from "@/lib/auth-context";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  IconMail,
  IconPhone,
  IconBriefcase,
  IconBuilding,
  IconUser,
} from "@tabler/icons-react";

export default function ProfilePage() {
  usePageTitle("User Profile | WorkZen");
  useRequireAuth();
  const { user, company } = useAuth();

  const getRoleBadge = (role: string) => {
    const colors: Record<string, string> = {
      superadmin: "bg-purple-100 text-purple-800",
      admin: "bg-blue-100 text-blue-800",
      hr: "bg-green-100 text-green-800",
      payroll: "bg-yellow-100 text-yellow-800",
      employee: "bg-gray-100 text-gray-800",
    };
    return (
      <Badge className={colors[role] || colors.employee}>
        {role.charAt(0).toUpperCase() + role.slice(1)}
      </Badge>
    );
  };

  const getStatusBadge = (status: string) => {
    return (
      <Badge variant={status === "active" ? "default" : "secondary"}>
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </Badge>
    );
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Profile</h1>
        <p className="text-muted-foreground">
          View and manage your profile information
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Personal Information</CardTitle>
            <CardDescription>Your basic profile details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center space-x-3">
              <IconUser className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Full Name</p>
                <p className="font-medium">
                  {user?.first_name} {user?.last_name}
                </p>
              </div>
            </div>
            <div className="flex items-center space-x-3">
              <IconMail className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Email</p>
                <p className="font-medium">{user?.email}</p>
              </div>
            </div>
            <div className="flex items-center space-x-3">
              <IconPhone className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Phone</p>
                <p className="font-medium">{user?.phone || "Not provided"}</p>
              </div>
            </div>
            <div className="flex items-center space-x-3">
              <IconUser className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Username</p>
                <p className="font-medium">{user?.username}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Employment Details</CardTitle>
            <CardDescription>Your work-related information</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center space-x-3">
              <IconBriefcase className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Designation</p>
                <p className="font-medium">
                  {user?.designation || "Not assigned"}
                </p>
              </div>
            </div>
            <div className="flex items-center space-x-3">
              <IconBuilding className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Company</p>
                <p className="font-medium">{company?.name || "Not assigned"}</p>
              </div>
            </div>
            <div className="flex items-center space-x-3">
              <div className="w-5 h-5" />
              <div>
                <p className="text-sm text-muted-foreground">Role</p>
                {getRoleBadge(user?.role || "employee")}
              </div>
            </div>
            <div className="flex items-center space-x-3">
              <div className="w-5 h-5" />
              <div>
                <p className="text-sm text-muted-foreground">Status</p>
                {getStatusBadge(user?.status || "active")}
              </div>
            </div>
            <div className="flex items-center space-x-3">
              <div className="w-5 h-5" />
              <div>
                <p className="text-sm text-muted-foreground">Employee Code</p>
                <p className="font-medium">
                  {user?.employee_code || "Not assigned"}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
