import * as React from "react";
import { cn } from "@/lib/utils";

interface DataTableWrapperProps {
  children: React.ReactNode;
  className?: string;
}

/**
 * Standardized wrapper for all tables across the application
 * Ensures consistent styling and theme
 */
export function DataTableWrapper({
  children,
  className,
}: DataTableWrapperProps) {
  return (
    <div
      className={cn(
        "rounded-md border bg-card shadow-sm overflow-hidden",
        className
      )}
    >
      <div className="overflow-x-auto">{children}</div>
    </div>
  );
}
