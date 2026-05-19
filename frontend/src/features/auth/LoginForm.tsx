"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useSIAKADLogin } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { AlertCircle, LogIn } from "lucide-react";

const loginSchema = z.object({
  nim: z.string().min(1, "NIM is required"),
  password: z.string().min(1, "Password is required"),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export function LoginForm() {
  const login = useSIAKADLogin();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = (data: LoginFormValues) => {
    login.mutate(data);
  };

  return (
    <Card className="w-full max-w-md shadow-lg border-gray-100">
      <CardHeader className="text-center pb-2">
        <div className="lg:hidden flex justify-center mb-3">
          <div className="h-12 w-12 rounded-xl gradient-primary flex items-center justify-center">
            <LogIn className="h-6 w-6 text-white" />
          </div>
        </div>
        <CardTitle className="text-xl">Sign In</CardTitle>
        <CardDescription className="text-sm">
          Use your SIAKAD credentials to authenticate, then link your GitHub
          account.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <label
              htmlFor="nim"
              className="text-sm font-medium text-gray-700"
            >
              NIM (Student ID)
            </label>
            <Input
              id="nim"
              placeholder="Enter your NIM"
              {...register("nim")}
              className="h-11"
            />
            {errors.nim && (
              <p className="text-sm text-red-600 flex items-center gap-1">
                <AlertCircle className="h-3.5 w-3.5" />
                {errors.nim.message}
              </p>
            )}
          </div>
          <div className="space-y-2">
            <label
              htmlFor="password"
              className="text-sm font-medium text-gray-700"
            >
              Password
            </label>
            <Input
              id="password"
              type="password"
              placeholder="Enter your password"
              {...register("password")}
              className="h-11"
            />
            {errors.password && (
              <p className="text-sm text-red-600 flex items-center gap-1">
                <AlertCircle className="h-3.5 w-3.5" />
                {errors.password.message}
              </p>
            )}
          </div>
          {login.isError && (
            <div className="p-3 rounded-lg bg-red-50 border border-red-200">
              <p className="text-sm text-red-700 flex items-center gap-1.5">
                <AlertCircle className="h-4 w-4 shrink-0" />
                Authentication failed. Please check your credentials.
              </p>
            </div>
          )}
          <Button
            type="submit"
            className="w-full h-11 gradient-primary text-white border-0 shadow-sm hover:shadow-md transition-shadow"
            disabled={login.isPending}
          >
            {login.isPending ? "Authenticating..." : "Continue with SIAKAD"}
          </Button>
          <p className="text-xs text-center text-gray-400 pt-2">
            By signing in, you agree to use this platform responsibly.
          </p>
        </form>
      </CardContent>
    </Card>
  );
}
