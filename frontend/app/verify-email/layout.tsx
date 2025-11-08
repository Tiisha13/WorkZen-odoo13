import { Metadata } from "next";
import { PropsWithChildren, Fragment } from "react";

export const metadata: Metadata = {
  title: "Verify Email | WorkZen",
  description:
    "Verify your email address to complete your WorkZen registration.",
};

const Layout = ({ children }: PropsWithChildren) => {
  return <Fragment>{children}</Fragment>;
};

export default Layout;
