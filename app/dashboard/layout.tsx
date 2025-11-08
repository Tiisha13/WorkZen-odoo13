import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { PropsWithChildren } from "react";

const Layout = ({ children }: PropsWithChildren) => {
  return <DashboardLayout>{children}</DashboardLayout>;
};

export default Layout;
