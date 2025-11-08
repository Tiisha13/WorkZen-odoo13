"use client";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { User } from "@/lib/types";
import { ROLE_LABELS } from "@/lib/config";

interface ProfileFormProps {
  user: User;
}

export function ProfileForm({ user }: ProfileFormProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Personal Information</CardTitle>
        <CardDescription>
          Your account details and contact information
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label>First Name</Label>
            <Input value={user.first_name} disabled />
          </div>
          <div className="space-y-2">
            <Label>Last Name</Label>
            <Input value={user.last_name} disabled />
          </div>
        </div>

        <div className="space-y-2">
          <Label>Username</Label>
          <Input value={user.username} disabled />
        </div>

        <div className="space-y-2">
          <Label>Email</Label>
          <Input value={user.email} type="email" disabled />
        </div>

        <div className="space-y-2">
          <Label>Phone Number</Label>
          <Input value={user.phone || "Not provided"} disabled />
        </div>

        <div className="space-y-2">
          <Label>Role</Label>
          <Input
            value={
              ROLE_LABELS[user.role as keyof typeof ROLE_LABELS] || user.role
            }
            disabled
          />
        </div>

        <div className="space-y-2">
          <Label>Employee Code</Label>
          <Input value={user.employee_code || "Not assigned"} disabled />
        </div>

        {user.department_id && (
          <div className="space-y-2">
            <Label>Department ID</Label>
            <Input value={user.department_id} disabled />
          </div>
        )}

        <p className="text-sm text-muted-foreground">
          Contact your administrator to update your profile information.
        </p>
      </CardContent>
    </Card>
  );
}
