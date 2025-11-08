"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { IconClock } from "@tabler/icons-react";

const activities = [
  {
    id: 1,
    user: "John Doe",
    action: "marked attendance",
    time: "2 hours ago",
    type: "attendance",
  },
  {
    id: 2,
    user: "Jane Smith",
    action: "requested leave",
    time: "3 hours ago",
    type: "leave",
  },
  {
    id: 3,
    user: "Admin",
    action: "processed payroll",
    time: "5 hours ago",
    type: "payroll",
  },
  {
    id: 4,
    user: "Mike Johnson",
    action: "uploaded document",
    time: "6 hours ago",
    type: "document",
  },
];

const getTypeBadge = (type: string) => {
  const variants: Record<
    string,
    { variant: "default" | "secondary" | "outline"; label: string }
  > = {
    attendance: { variant: "default", label: "Attendance" },
    leave: { variant: "secondary", label: "Leave" },
    payroll: { variant: "outline", label: "Payroll" },
    document: { variant: "outline", label: "Document" },
  };
  return variants[type] || variants.document;
};

export function RecentActivity() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Activity</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {activities.map((activity) => {
            const typeBadge = getTypeBadge(activity.type);
            return (
              <div key={activity.id} className="flex items-start gap-4">
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-accent">
                  <IconClock className="h-4 w-4" />
                </div>
                <div className="flex-1 space-y-1">
                  <p className="text-sm">
                    <span className="font-medium">{activity.user}</span>{" "}
                    {activity.action}
                  </p>
                  <div className="flex items-center gap-2">
                    <Badge variant={typeBadge.variant}>{typeBadge.label}</Badge>
                    <span className="text-xs text-muted-foreground">
                      {activity.time}
                    </span>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
