"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  IconPlus,
  IconUsers,
  IconCalendar,
  IconClock,
  IconFileText,
} from "@tabler/icons-react";
import { useAuth } from "@/lib/auth-context";
import Link from "next/link";

export function QuickActions() {
  const { hasRole } = useAuth();

  const actions = [
    {
      title: "Add Employee",
      description: "Create new employee account",
      icon: IconUsers,
      href: "/dashboard/users/new",
      roles: ["superadmin", "admin", "hr"],
      color: "bg-blue-500",
    },
    {
      title: "Mark Attendance",
      description: "Check in/out for today",
      icon: IconClock,
      href: "/dashboard/attendance",
      roles: ["superadmin", "admin", "hr", "employee"],
      color: "bg-green-500",
    },
    {
      title: "Request Leave",
      description: "Submit leave application",
      icon: IconCalendar,
      href: "/dashboard/leaves/new",
      roles: ["superadmin", "admin", "hr", "employee"],
      color: "bg-purple-500",
    },
    {
      title: "Upload Document",
      description: "Add new document",
      icon: IconFileText,
      href: "/dashboard/documents/new",
      roles: ["superadmin", "admin", "hr", "employee"],
      color: "bg-orange-500",
    },
  ];

  const visibleActions = actions.filter((action) => hasRole(action.roles));

  return (
    <Card>
      <CardHeader>
        <CardTitle>Quick Actions</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-4">
        {visibleActions.map((action) => (
          <Link key={action.title} href={action.href}>
            <Button
              variant="outline"
              className="h-auto w-full justify-start gap-4 p-4 hover:bg-accent"
            >
              <div
                className={`flex h-10 w-10 items-center justify-center rounded-lg ${action.color} text-white`}
              >
                <action.icon className="h-5 w-5" />
              </div>
              <div className="flex-1 text-left">
                <p className="font-medium">{action.title}</p>
                <p className="text-sm text-muted-foreground">
                  {action.description}
                </p>
              </div>
              <IconPlus className="h-4 w-4" />
            </Button>
          </Link>
        ))}
      </CardContent>
    </Card>
  );
}
