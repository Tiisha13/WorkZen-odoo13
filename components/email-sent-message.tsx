"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { IconCheck, IconMail } from "@tabler/icons-react";

interface EmailSentMessageProps {
  email: string;
  onResend?: () => void;
  onBackToLogin?: () => void;
}

export function EmailSentMessage({
  email,
  onResend,
  onBackToLogin,
}: EmailSentMessageProps) {
  return (
    <Card>
      <CardHeader className="text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
          <IconCheck className="h-8 w-8 text-green-600" />
        </div>
        <CardTitle>Check Your Email</CardTitle>
        <CardDescription>
          We&apos;ve sent a verification link to <strong>{email}</strong>
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2 text-center text-sm text-muted-foreground">
          <p>Click the link in the email to verify your account.</p>
          <p className="text-xs">
            The verification link will expire in 24 hours.
          </p>
        </div>

        <div className="rounded-lg border border-blue-200 bg-blue-50 p-4">
          <div className="flex items-start space-x-3">
            <IconMail className="mt-0.5 h-5 w-5 shrink-0 text-blue-600" />
            <div className="text-sm text-blue-900">
              <p className="font-medium">Can&apos;t find the email?</p>
              <ul className="mt-2 space-y-1 text-xs">
                <li>• Check your spam or junk folder</li>
                <li>• Make sure the email address is correct</li>
                <li>• Wait a few minutes and try again</li>
              </ul>
            </div>
          </div>
        </div>

        {onResend && (
          <Button onClick={onResend} variant="outline" className="w-full">
            Resend Verification Email
          </Button>
        )}

        {onBackToLogin && (
          <Button onClick={onBackToLogin} variant="ghost" className="w-full">
            Back to Login
          </Button>
        )}
      </CardContent>
    </Card>
  );
}
