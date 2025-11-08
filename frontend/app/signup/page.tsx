import { SignupForm } from "@/components/signup-form";

export const metadata = {
  title: "Sign Up | WorkZen",
  description: "Create a new account to access WorkZen.",
};

export default function Page() {
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <SignupForm />
      </div>
    </div>
  );
}
