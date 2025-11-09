"use client";

import { useEffect, useState, Suspense } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { IconCheck, IconX, IconLoader2, IconMail } from "@tabler/icons-react";
import { apiService } from "@/lib/api-service";
import { toast } from "sonner";

function VerifyEmailContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [verifying, setVerifying] = useState(true);
  const [verified, setVerified] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const token = searchParams.get("token");

  const verifyEmail = async (verificationToken: string) => {
    try {
      setVerifying(true);
      setError(null);

      const response = await apiService.verifyEmail(verificationToken);

      if (response.success) {
        setVerified(true);
        toast.success(response.message || "Email verified successfully!");

        // Redirect to login after 3 seconds
        setTimeout(() => {
          router.push("/login");
        }, 3000);
      } else {
        setError(response.message || "Verification failed");
        toast.error(response.message || "Verification failed");
      }
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Verification failed";
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setVerifying(false);
    }
  };

  useEffect(() => {
    if (!token) {
      setError("Invalid verification link. No token provided.");
      setVerifying(false);
      return;
    }

    verifyEmail(token);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token]);
  const handleRetry = () => {
    if (token) {
      verifyEmail(token);
    }
  };

  const handleGoToLogin = () => {
    router.push("/login");
  };

  return (
    <div className="flex min-h-screen w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-md">
        <Card>
          <CardHeader className="text-center">
            <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
              <IconMail className="h-8 w-8 text-green-600" />
            </div>
            <CardTitle>Email Verification</CardTitle>
            <CardDescription>
              {verifying
                ? "Verifying your email address..."
                : verified
                ? "Your email has been verified"
                : "Verification failed"}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {verifying && (
              <div className="flex flex-col items-center justify-center space-y-4 py-8">
                <IconLoader2 className="h-12 w-12 animate-spin text-green-600" />
                <p className="text-center text-sm text-muted-foreground">
                  Please wait while we verify your email address...
                </p>
              </div>
            )}

            {!verifying && verified && (
              <div className="space-y-4">
                <div className="flex flex-col items-center justify-center space-y-4 py-4">
                  <div className="flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
                    <IconCheck className="h-8 w-8 text-green-600" />
                  </div>
                  <div className="text-center">
                    <h3 className="text-lg font-semibold text-green-600">
                      Success!
                    </h3>
                    <p className="text-sm text-muted-foreground">
                      Your email has been verified successfully.
                    </p>
                  </div>
                </div>
                <div className="space-y-3">
                  <p className="text-center text-sm text-muted-foreground">
                    You can now log in to your WorkZen account.
                  </p>
                  <p className="text-center text-xs text-muted-foreground">
                    Redirecting to login page in 3 seconds...
                  </p>
                  <Button
                    onClick={handleGoToLogin}
                    className="w-full"
                    size="lg"
                  >
                    Go to Login
                  </Button>
                </div>
              </div>
            )}

            {!verifying && !verified && error && (
              <div className="space-y-4">
                <div className="flex flex-col items-center justify-center space-y-4 py-4">
                  <div className="flex h-16 w-16 items-center justify-center rounded-full bg-red-100">
                    <IconX className="h-8 w-8 text-red-600" />
                  </div>
                  <div className="text-center">
                    <h3 className="text-lg font-semibold text-red-600">
                      Verification Failed
                    </h3>
                    <p className="text-sm text-muted-foreground">{error}</p>
                  </div>
                </div>
                <div className="space-y-2">
                  <p className="text-center text-sm text-muted-foreground">
                    Possible reasons:
                  </p>
                  <ul className="space-y-1 text-xs text-muted-foreground">
                    <li>• The verification link has expired (24 hours)</li>
                    <li>• The email has already been verified</li>
                    <li>• The verification link is invalid</li>
                  </ul>
                </div>
                <div className="flex flex-col space-y-2">
                  <Button
                    onClick={handleRetry}
                    variant="outline"
                    className="w-full"
                  >
                    Try Again
                  </Button>
                  <Button
                    onClick={handleGoToLogin}
                    variant="outline"
                    className="w-full"
                  >
                    Go to Login
                  </Button>
                  <Button
                    onClick={() => router.push("/signup")}
                    variant="ghost"
                    className="w-full"
                  >
                    Create New Account
                  </Button>
                </div>
              </div>
            )}

            {!token && !verifying && (
              <div className="space-y-4">
                <div className="flex flex-col items-center justify-center space-y-4 py-4">
                  <div className="flex h-16 w-16 items-center justify-center rounded-full bg-orange-100">
                    <IconX className="h-8 w-8 text-orange-600" />
                  </div>
                  <div className="text-center">
                    <h3 className="text-lg font-semibold text-orange-600">
                      Invalid Link
                    </h3>
                    <p className="text-sm text-muted-foreground">
                      This verification link is not valid.
                    </p>
                  </div>
                </div>
                <Button
                  onClick={() => router.push("/signup")}
                  className="w-full"
                >
                  Create New Account
                </Button>
              </div>
            )}
          </CardContent>
        </Card>

        <div className="mt-4 text-center">
          <p className="text-sm text-muted-foreground">
            Need help?{" "}
            <a
              href="mailto:support@workzen.com"
              className="text-green-600 hover:underline"
            >
              Contact Support
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}

export default function VerifyEmailPage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-screen w-full items-center justify-center p-6 md:p-10">
          <div className="w-full max-w-md">
            <Card>
              <CardContent className="flex items-center justify-center py-12">
                <IconLoader2 className="h-8 w-8 animate-spin text-primary" />
              </CardContent>
            </Card>
          </div>
        </div>
      }
    >
      <VerifyEmailContent />
    </Suspense>
  );
}
