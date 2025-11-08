import { Metadata } from "next";
import { Fragment, PropsWithChildren } from "react";

export const metadata: Metadata = {
  title: "Resend Verification | WorkZen",
  description: "Resend your email verification link.",
};

const Layout = ({ children }: PropsWithChildren) => {
  return <Fragment>{children}</Fragment>;
};

export default Layout;
