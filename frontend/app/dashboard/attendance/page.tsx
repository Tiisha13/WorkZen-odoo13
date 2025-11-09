"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { useState, useEffect } from "react";
import { apiService } from "@/lib/api-service";
import { API_ENDPOINTS } from "@/lib/config";
import { useAuth } from "@/lib/auth-context";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { IconClock, IconClockStop, IconCalendar } from "@tabler/icons-react";
import { toast } from "sonner";
import type { Attendance } from "@/lib/types";

export default function AttendancePage() {
  usePageTitle("Attendance Management | WorkZen");
  useRequireAuth();

  const { user } = useAuth();
  const [attendances, setAttendances] = useState<Attendance[]>([]);
  const [loading, setLoading] = useState(true);
  const [todayAttendance, setTodayAttendance] = useState<Attendance | null>(
    null
  );
  const [checkingIn, setCheckingIn] = useState(false);

  useEffect(() => {
    fetchAttendances();
    fetchTodayAttendance();
  }, []);

  const fetchAttendances = async () => {
    try {
      setLoading(true);
      const response = await apiService.get<{
        success: boolean;
        data: Attendance[];
      }>(API_ENDPOINTS.ATTENDANCES);
      // Sort by date, newest first
      const sortedData = (response.data || []).sort(
        (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()
      );
      setAttendances(sortedData);
    } catch (error) {
      console.error("Failed to fetch attendances:", error);
      setAttendances([]);
    } finally {
      setLoading(false);
    }
  };

  const fetchTodayAttendance = async () => {
    try {
      const response = await apiService.get<{
        success: boolean;
        data: Attendance[];
      }>(API_ENDPOINTS.ATTENDANCES);
      // Get today's attendance from the list
      if (response.data && response.data.length > 0) {
        const today = new Date();
        const todayStr = today.toISOString().split("T")[0];

        const todayRecord = response.data.find((att) => {
          if (!att.date) return false;
          // Try to match date in various formats
          const dateStr = String(att.date);
          return dateStr.startsWith(todayStr) || dateStr.includes(todayStr);
        });

        if (todayRecord) {
          console.log("Today's attendance:", todayRecord);
          setTodayAttendance(todayRecord);
        }
      }
    } catch (error) {
      console.error("No attendance today", error);
    }
  };

  const handleCheckIn = async () => {
    try {
      setCheckingIn(true);
      const response = await apiService.post<{
        success: boolean;
        data: Attendance;
      }>(API_ENDPOINTS.ATTENDANCE_CHECKIN, {
        user_id: user?.id,
        date: new Date().toISOString(),
      });

      if (response.success && response.data) {
        // Immediately update today's attendance
        setTodayAttendance(response.data);
        // Add to top of attendances list
        setAttendances((prev) => [response.data, ...prev]);
        toast.success("Checked in successfully");
      }

      // Still fetch to ensure we have the latest data
      await fetchAttendances();
      await fetchTodayAttendance();
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to check in"
      );
    } finally {
      setCheckingIn(false);
    }
  };

  const handleCheckOut = async () => {
    if (!todayAttendance) return;
    try {
      setCheckingIn(true);
      const response = await apiService.post<{
        success: boolean;
        data: Attendance;
      }>(API_ENDPOINTS.ATTENDANCE_CHECKOUT, {
        attendance_id: todayAttendance.id,
        date: new Date().toISOString(),
      });

      if (response.success && response.data) {
        // Immediately update today's attendance
        setTodayAttendance(response.data);
        // Update the record in the attendances list
        setAttendances((prev) =>
          prev.map((att) => (att.id === response.data.id ? response.data : att))
        );
        toast.success("Checked out successfully");
      }

      // Still fetch to ensure we have the latest data
      await fetchAttendances();
      await fetchTodayAttendance();
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to check out"
      );
    } finally {
      setCheckingIn(false);
    }
  };

  const handleReset = async () => {
    if (
      !confirm(
        "Are you sure you want to reset today's attendance? This will allow you to check in again."
      )
    ) {
      return;
    }
    try {
      setCheckingIn(true);
      await apiService.delete(API_ENDPOINTS.ATTENDANCE_RESET);
      toast.success("Attendance reset successfully! You can check in again.");
      setTodayAttendance(null);
      fetchAttendances();
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to reset attendance"
      );
    } finally {
      setCheckingIn(false);
    }
  };

  const formatTime = (timeString: string | undefined) => {
    if (!timeString) return "-";
    try {
      // If it's just a time string like "02:36:03" or "14:30:00"
      if (timeString.match(/^\d{2}:\d{2}:\d{2}$/)) {
        const [hours, minutes] = timeString.split(":").map(Number);
        const period = hours >= 12 ? "PM" : "AM";
        const displayHours = hours % 12 || 12;
        return `${displayHours.toString().padStart(2, "0")}:${minutes
          .toString()
          .padStart(2, "0")} ${period}`;
      }

      // If it includes date (ISO format or UTC format)
      if (timeString.includes("T") || timeString.includes("Z")) {
        const date = new Date(timeString);
        if (!isNaN(date.getTime())) {
          return date.toLocaleTimeString("en-US", {
            hour: "2-digit",
            minute: "2-digit",
            hour12: true,
          });
        }
      }

      return timeString; // Return as is if we can't format it
    } catch (error) {
      console.error("Error formatting time:", error, timeString);
      return "-";
    }
  };

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return "-";
    try {
      // Handle both ISO string and other formats
      let date: Date;
      if (dateString.includes("T") || dateString.includes("Z")) {
        // ISO format or UTC format
        date = new Date(dateString);
      } else if (dateString.includes("-")) {
        // YYYY-MM-DD format
        const [year, month, day] = dateString.split("-").map(Number);
        date = new Date(year, month - 1, day);
      } else {
        // Try parsing as is
        date = new Date(dateString);
      }

      if (isNaN(date.getTime())) return "-";

      return date.toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
      });
    } catch (error) {
      console.error("Error formatting date:", error);
      return "-";
    }
  };

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<
      string,
      {
        variant: "default" | "secondary" | "destructive" | "outline";
        label: string;
        className?: string;
      }
    > = {
      present: {
        variant: "outline",
        label: "Present",
        className: "border-green-200 bg-green-50 text-green-700",
      },
      absent: {
        variant: "outline",
        label: "Absent",
        className: "border-red-200 bg-red-50 text-red-700",
      },
      leave: {
        variant: "outline",
        label: "Leave",
        className: "border-yellow-200 bg-yellow-50 text-yellow-700",
      },
      half_day: {
        variant: "outline",
        label: "Half Day",
        className: "border-orange-200 bg-orange-50 text-orange-700",
      },
    };

    const config = statusConfig[status] || {
      variant: "secondary" as const,
      label: status,
      className: "",
    };

    return (
      <Badge variant={config.variant} className={config.className}>
        {config.label}
      </Badge>
    );
  };

  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Header Section */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold tracking-tight">
            Attendance
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            Track your daily attendance and working hours
          </p>
        </div>
      </div>

      {/* Today's Attendance Card - Hidden for Admin/SuperAdmin */}
      {user?.role !== "admin" && user?.role !== "superadmin" && (
        <Card>
          <CardHeader>
            <CardTitle>Today&apos;s Attendance</CardTitle>
            <CardDescription>Check in and check out for today</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
              <div className="space-y-2">
                {todayAttendance ? (
                  <>
                    <div className="flex flex-wrap items-center gap-4 md:gap-6">
                      <div>
                        <p className="text-sm text-muted-foreground">
                          Check-in
                        </p>
                        <p className="text-lg font-semibold">
                          {formatTime(todayAttendance.check_in)}
                        </p>
                      </div>
                      {todayAttendance.check_out && (
                        <>
                          <div className="hidden md:block text-muted-foreground">
                            →
                          </div>
                          <div>
                            <p className="text-sm text-muted-foreground">
                              Check-out
                            </p>
                            <p className="text-lg font-semibold">
                              {formatTime(todayAttendance.check_out)}
                            </p>
                          </div>
                          {todayAttendance.working_hours && (
                            <>
                              <div className="hidden md:block text-muted-foreground">
                                •
                              </div>
                              <div>
                                <p className="text-sm text-muted-foreground">
                                  Working Hours
                                </p>
                                <p className="text-lg font-semibold">
                                  {todayAttendance.working_hours.toFixed(2)}h
                                </p>
                              </div>
                            </>
                          )}
                        </>
                      )}
                    </div>
                  </>
                ) : (
                  <p className="text-sm text-muted-foreground">
                    You haven&apos;t checked in today
                  </p>
                )}
              </div>
              <div className="flex flex-wrap gap-2">
                {!todayAttendance ? (
                  <Button
                    onClick={handleCheckIn}
                    disabled={checkingIn}
                    size="sm"
                  >
                    <IconClock className="w-4 h-4 mr-2" />
                    {checkingIn ? "Checking In..." : "Check In"}
                  </Button>
                ) : !todayAttendance.check_out ? (
                  <>
                    <Button
                      onClick={handleCheckOut}
                      disabled={checkingIn}
                      variant="destructive"
                      size="sm"
                    >
                      <IconClockStop className="w-4 h-4 mr-2" />
                      {checkingIn ? "Checking Out..." : "Check Out"}
                    </Button>
                    <Button
                      onClick={handleReset}
                      disabled={checkingIn}
                      variant="outline"
                      size="sm"
                    >
                      Reset
                    </Button>
                  </>
                ) : (
                  <>
                    <Badge
                      variant="outline"
                      className="text-green-600 px-4 py-2"
                    >
                      <IconCalendar className="w-4 h-4 mr-2" />
                      Completed
                    </Badge>
                    <Button
                      onClick={handleReset}
                      disabled={checkingIn}
                      variant="outline"
                      size="sm"
                    >
                      Reset
                    </Button>
                  </>
                )}
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Attendance Records Table */}
      <div className="rounded-md border bg-card shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Date</TableHead>
                <TableHead>Check In</TableHead>
                <TableHead>Check Out</TableHead>
                <TableHead>Working Hours</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-32 text-center">
                    <div className="flex items-center justify-center">
                      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                    </div>
                  </TableCell>
                </TableRow>
              ) : attendances.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-32 text-center">
                    <div className="flex flex-col items-center justify-center text-muted-foreground">
                      <p className="text-sm">No attendance records found</p>
                      <p className="text-xs mt-1">Start by checking in today</p>
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                attendances.map((attendance) => (
                  <TableRow key={attendance.id} className="hover:bg-muted/50">
                    <TableCell className="font-medium">
                      {formatDate(attendance.date)}
                    </TableCell>
                    <TableCell className="text-sm">
                      {formatTime(attendance.check_in)}
                    </TableCell>
                    <TableCell className="text-sm">
                      {attendance.check_out ? (
                        formatTime(attendance.check_out)
                      ) : (
                        <span className="text-muted-foreground">-</span>
                      )}
                    </TableCell>
                    <TableCell className="text-sm">
                      {attendance.working_hours ? (
                        `${attendance.working_hours.toFixed(2)}h`
                      ) : (
                        <span className="text-muted-foreground">-</span>
                      )}
                    </TableCell>
                    <TableCell>{getStatusBadge(attendance.status)}</TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </div>
    </div>
  );
}
