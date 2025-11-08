"use client";

import { useAuth } from "@/lib/auth-context";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
} from "@/components/ui/sidebar";
import {
  IconHome,
  IconUsers,
  IconBuilding,
  IconCalendar,
  IconFileText,
  IconCash,
  IconSettings,
  IconLogout,
  IconClock,
} from "@tabler/icons-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

const menuItems = [
  {
    href: "/dashboard",
    icon: IconHome,
    label: "Dashboard",
    roles: ["superadmin", "admin", "hr", "payroll", "employee"],
  },
  {
    href: "/dashboard/users",
    icon: IconUsers,
    label: "Users",
    roles: ["admin", "hr"],
  },
  {
    href: "/dashboard/departments",
    icon: IconBuilding,
    label: "Departments",
    roles: ["admin", "hr"],
  },
  {
    href: "/dashboard/attendance",
    icon: IconClock,
    label: "Attendance",
    roles: ["admin", "hr", "employee"],
  },
  {
    href: "/dashboard/leaves",
    icon: IconCalendar,
    label: "Leaves",
    roles: ["admin", "hr", "employee"],
  },
  {
    href: "/dashboard/payroll",
    icon: IconCash,
    label: "Payroll",
    roles: ["admin", "payroll"],
  },
  {
    href: "/dashboard/documents",
    icon: IconFileText,
    label: "Documents",
    roles: ["admin", "hr", "employee"],
  },
];

export function AppSidebar() {
  const { user, logout, hasRole, isSuperAdmin } = useAuth();
  const pathname = usePathname();

  const canViewItem = (roles: string[]) => {
    return hasRole(roles);
  };

  // SuperAdmin only sees Dashboard (Companies List)
  const visibleMenuItems = isSuperAdmin()
    ? menuItems.filter((item) => item.href === "/dashboard")
    : menuItems.filter((item) => canViewItem(item.roles));

  return (
    <Sidebar>
      <SidebarHeader className="border-b p-4">
        <div className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <span className="text-lg font-bold">W</span>
          </div>
          <div>
            <p className="text-sm font-semibold">WorkZen</p>
            <p className="text-xs text-muted-foreground">HRMS</p>
          </div>
        </div>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Menu</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {visibleMenuItems.map((item) => (
                <SidebarMenuItem key={item.href}>
                  <SidebarMenuButton asChild isActive={pathname === item.href}>
                    <Link
                      href={item.href}
                      className={cn(
                        "flex items-center gap-2",
                        pathname === item.href && "bg-accent"
                      )}
                    >
                      <item.icon className="h-4 w-4" />
                      <span>{item.label}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarGroup>
          <SidebarGroupLabel>Account</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild isActive={pathname === "/profile"}>
                  <Link href="/profile" className="flex items-center gap-2">
                    <IconSettings className="h-4 w-4" />
                    <span>Profile</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton
                  onClick={logout}
                  className="flex items-center gap-2 text-red-600"
                >
                  <IconLogout className="h-4 w-4" />
                  <span>Logout</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="border-t p-4">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary text-primary-foreground">
            <span className="text-sm font-semibold">
              {user?.first_name?.[0]}
              {user?.last_name?.[0]}
            </span>
          </div>
          <div className="flex-1 overflow-hidden">
            <p className="truncate text-sm font-medium">
              {user?.first_name} {user?.last_name}
            </p>
            <p className="truncate text-xs text-muted-foreground">
              {user?.role}
            </p>
          </div>
        </div>
      </SidebarFooter>
    </Sidebar>
  );
}
