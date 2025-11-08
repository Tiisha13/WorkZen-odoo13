import type { Metadata } from "next";
import { ThemeProvider } from "@/components/theme-provider";
import { AuthProvider } from "@/lib/auth-context";
import { Toaster } from "@/components/ui/sonner";
import { Fragment, PropsWithChildren } from "react";
import "./globals.css";

export const metadata: Metadata = {
  title: "WorkZen - HRMS",
  description: "WorkZen - Your HRMS Solution",
};

export default function RootLayout({ children }: PropsWithChildren) {
  return (
    <Fragment>
      <html lang="en" suppressHydrationWarning>
        <head />
        <body>
          <ThemeProvider
            attribute="class"
            defaultTheme="system"
            enableSystem
            disableTransitionOnChange
          >
            <AuthProvider>
              {children}
              <Toaster />
            </AuthProvider>
          </ThemeProvider>
        </body>
      </html>
    </Fragment>
  );
}
