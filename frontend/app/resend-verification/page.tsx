"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { IconMail, IconLoader2, IconCheck } from "@tabler/icons-react";
import { apiService } from "@/lib/api-service";
import { toast } from "sonner";
import Link from "next/link";

export default function ResendVerificationPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [sent, setSent] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email.trim()) {
      toast.error("Please enter your email address");
      return;
    }

    try {
      setLoading(true);
      const response = await apiService.resendVerification(email.trim());

      if (response.success) {
        setSent(true);
        toast.success(
          response.message || "Verification email sent successfully!"
        );
      } else {
        toast.error(response.message || "Failed to send verification email");
      }
    } catch (err) {
      const errorMessage =
        err instanceof Error
          ? err.message
          : "Failed to send verification email";
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-md">
        <Card>
          <CardHeader className="text-center">
            <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
              <IconMail className="h-8 w-8 text-green-600" />
            </div>
            <CardTitle>Resend Verification Email</CardTitle>
            <CardDescription>
              Enter your email address to receive a new verification link
            </CardDescription>
          </CardHeader>
          <CardContent>
            {!sent ? (
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="email">Email Address</Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="you@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    disabled={loading}
                  />
                </div>

                <Button type="submit" className="w-full" disabled={loading}>
                  {loading ? (
                    <>
                      <IconLoader2 className="mr-2 h-4 w-4 animate-spin" />
                      Sending...
                    </>
                  ) : (
                    "Send Verification Email"
                  )}
                </Button>

                <div className="text-center text-sm text-muted-foreground">
                  <p>
                    Already verified?{" "}
                    <button
                      type="button"
                      onClick={() => router.push("/login")}
                      className="text-green-600 hover:underline"
                    >
                      Go to Login
                    </button>
                  </p>
                </div>
              </form>
            ) : (
              <div className="space-y-4">
                <div className="flex flex-col items-center justify-center space-y-4 py-4">
                  <div className="flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
                    <IconCheck className="h-8 w-8 text-green-600" />
                  </div>
                  <div className="text-center">
                    <h3 className="text-lg font-semibold text-green-600">
                      Email Sent!
                    </h3>
                    <p className="text-sm text-muted-foreground">
                      We&apos;ve sent a verification link to{" "}
                      <strong>{email}</strong>
                    </p>
                  </div>
                </div>

                <div className="space-y-3 text-center text-sm text-muted-foreground">
                  <p>
                    Please check your email and click the verification link.
                  </p>
                  <p className="text-xs">
                    Didn&apos;t receive the email? Check your spam folder or try
                    again in a few minutes.
                  </p>
                </div>

                <div className="flex flex-col space-y-2">
                  <Button
                    onClick={() => {
                      setSent(false);
                      setEmail("");
                    }}
                    variant="outline"
                    className="w-full"
                  >
                    Send to Different Email
                  </Button>
                  <Button
                    onClick={() => router.push("/login")}
                    variant="ghost"
                    className="w-full"
                  >
                    Go to Login
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        <div className="mt-4 text-center">
          <p className="text-sm text-muted-foreground">
            Need help?{" "}
            <Link
              href="mailto:support@workzen.com"
              className="text-green-600 hover:underline"
            >
              Contact Support
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
