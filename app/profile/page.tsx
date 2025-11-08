"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { ProfileForm } from "@/components/profile/profile-form";
import { ChangePasswordForm } from "@/components/profile/change-password-form";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

export default function ProfilePage() {
  usePageTitle("Profile Settings | WorkZen");
  const { user, isLoading } = useRequireAuth();

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <Skeleton className="h-96 w-full max-w-2xl" />
      </div>
    );
  }

  return (
    <DashboardLayout>
      <div className="mx-auto max-w-4xl space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Profile Settings</h1>
          <p className="text-muted-foreground">
            Manage your account settings and preferences
          </p>
        </div>

        <Tabs defaultValue="profile" className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="profile">Profile</TabsTrigger>
            <TabsTrigger value="security">Security</TabsTrigger>
          </TabsList>
          <TabsContent value="profile" className="space-y-4">
            <ProfileForm user={user!} />
          </TabsContent>
          <TabsContent value="security" className="space-y-4">
            <ChangePasswordForm />
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  );
}
