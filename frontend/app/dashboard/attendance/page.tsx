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
        const today = new Date().toISOString().split("T")[0];
        const todayRecord = response.data.find((att) =>
          att.date.startsWith(today)
        );
        if (todayRecord) setTodayAttendance(todayRecord);
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

  const formatTime = (dateString: string | undefined) => {
    if (!dateString) return "-";
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return "Invalid Date";
    return date.toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return "-";
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return "Invalid Date";
    return date.toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      present: "bg-green-100 text-green-800",
      absent: "bg-red-100 text-red-800",
      leave: "bg-yellow-100 text-yellow-800",
      half_day: "bg-orange-100 text-orange-800",
    };
    return (
      <Badge className={colors[status] || "bg-gray-100 text-gray-800"}>
        {status.replace("_", " ").toUpperCase()}
      </Badge>
    );
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Attendance</h1>
        <p className="text-muted-foreground">Track your daily attendance</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Today&apos;s Attendance</CardTitle>
          <CardDescription>Check in and check out for today</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div className="space-y-2">
              {todayAttendance ? (
                <>
                  <div className="flex items-center gap-4">
                    <div>
                      <p className="text-sm text-muted-foreground">Check-in</p>
                      <p className="text-lg font-semibold">
                        {formatTime(todayAttendance.check_in)}
                      </p>
                    </div>
                    {todayAttendance.check_out && (
                      <>
                        <div className="text-muted-foreground">→</div>
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
                            <div className="text-muted-foreground">•</div>
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
            <div className="flex gap-2">
              {!todayAttendance ? (
                <Button onClick={handleCheckIn} disabled={checkingIn}>
                  <IconClock className="w-4 h-4 mr-2" />
                  {checkingIn ? "Checking In..." : "Check In"}
                </Button>
              ) : !todayAttendance.check_out ? (
                <Button
                  onClick={handleCheckOut}
                  disabled={checkingIn}
                  variant="destructive"
                >
                  <IconClockStop className="w-4 h-4 mr-2" />
                  {checkingIn ? "Checking Out..." : "Check Out"}
                </Button>
              ) : (
                <Badge variant="outline" className="text-green-600 px-4 py-2">
                  <IconCalendar className="w-4 h-4 mr-2" />
                  Completed
                </Badge>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="border rounded-lg">
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
                <TableCell colSpan={5} className="text-center py-8">
                  Loading...
                </TableCell>
              </TableRow>
            ) : attendances.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className="text-center py-8 text-muted-foreground"
                >
                  No attendance records found
                </TableCell>
              </TableRow>
            ) : (
              attendances.map((attendance) => (
                <TableRow key={attendance.id}>
                  <TableCell className="font-medium">
                    {formatDate(attendance.date)}
                  </TableCell>
                  <TableCell>{formatTime(attendance.check_in)}</TableCell>
                  <TableCell>
                    {attendance.check_out
                      ? formatTime(attendance.check_out)
                      : "-"}
                  </TableCell>
                  <TableCell>
                    {attendance.working_hours
                      ? `${attendance.working_hours.toFixed(2)}h`
                      : "-"}
                  </TableCell>
                  <TableCell>{getStatusBadge(attendance.status)}</TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
