"use client";

import Image from "next/image";
import { useState, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { TbLogin2 } from "react-icons/tb";
import { FaSpinner } from "react-icons/fa6";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import ThemeSwitcher from "@/components/theme/theme-switcher";
import { useToast } from "@/components/ui/use-toast";
import Link from "next/link";
import axios from "axios";
import { deleteCookie } from "cookies-next";

export default function AuthenticationPage({
  token,
}: {
  token: string | undefined;
}) {
  const [isLoading, setIsLoading] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailError, setEmailError] = useState("");
  const [redirectUrl, setRedirectUrl] = useState<string | null>(null);

  const { toast } = useToast();
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const redirectu = searchParams.get("redirect") || "/";
    setRedirectUrl(redirectu);

    if (token) {
      axios
        .get(`${process.env.NEXT_PUBLIC_BACKEND_URL!}/api/auth/validate`, {
          withCredentials: true,
        })
        .then((response) => {
          if (response.status === 200) redirect(redirectu, token);
        })
        .catch(() => deleteCookie("Token"));
    }
  }, [searchParams, router, token]);

  const redirect = (url: string, token: string) => {
    if (url.startsWith("/")) return router.push(url);
    window.location.href = `${url}?token=${token}`;
  };

  const validateEmail = (email: string) =>
    /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(
      email.toLowerCase()
    );

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newEmail = e.target.value;
    setEmail(newEmail);
    setEmailError(
      newEmail && !validateEmail(newEmail)
        ? "Please enter a valid email address"
        : ""
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateEmail(email)) {
      setEmailError("Please enter a valid email address");
      return;
    }
    setIsLoading(true);

    try {
      const response = await axios.post(
        `${process.env.NEXT_PUBLIC_BACKEND_URL!}/api/auth/login`,
        {
          headers: { "Content-Type": "application/json" },
          email,
          password,
          credentials: "include",
        }
      );

      if (response.status === 200) {
        redirect(redirectUrl || "/", response.data.token);
      } else if (response.status === 401) {
        toast({
          title: "Bad Credentials.",
          description: "Authentication failed. Please try again.",
          variant: "destructive",
        });
        setPassword("");
      } else {
        throw new Error("Login failed");
      }
    } catch (error) {
      console.error("Login error:", error);
      toast({
        title: "Login Error",
        description: "An unexpected error occurred. Please try again later.",
        variant: "destructive",
      });
      setPassword("");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-screen h-screen flex items-center justify-center flex-col">
      <div className="absolute top-4 right-4">
        <ThemeSwitcher />
      </div>
      <div className="container relative flex-col items-center justify-center md:grid lg:max-w-none lg:px-0">
        <div className="lg:p-8">
          <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
            <div className="flex justify-center items-center gap-3">
              <Image src="/logo.svg" width={80} height={80} alt="pragma logo" />
              <h1 className="text-3xl font-bold">PRAGMA</h1>
            </div>
            <div className="flex flex-col space-y-2 text-center border-[0.5px] rounded-lg py-[5rem] px-3">
              <h1 className="text-2xl font-semibold tracking-tight">
                Authenticate
              </h1>
              <p>Sign into all your pragma apps!</p>
              <form
                onSubmit={handleSubmit}
                className="flex flex-col text-left mt-10 gap-3 items-center"
              >
                <div className="w-[90%]">
                  <Input
                    type="email"
                    id="email"
                    placeholder="Email"
                    className={`border-[0.6px] w-full ${
                      emailError ? "border-red-500" : ""
                    }`}
                    onChange={handleEmailChange}
                    value={email}
                    required
                  />
                  {emailError && (
                    <p className="text-red-500 text-sm mt-1">{emailError}</p>
                  )}
                </div>
                <Input
                  type="password"
                  id="password"
                  placeholder="Password"
                  className="border-[0.6px] w-[90%]"
                  onChange={(e) => setPassword(e.target.value)}
                  value={password}
                  required
                />
                <Button
                  type="submit"
                  className="dark:bg-white bg-black text-white dark:text-black w-[50%]"
                  disabled={isLoading || !!emailError}
                >
                  {isLoading ? (
                    <FaSpinner className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <TbLogin2 className="mr-2 h-4 w-4" />
                  )}
                  Login
                </Button>
                <hr className="w-[70%] mt-3" />
                <p className="text-zinc-400 dark:text-gray-400 text-sm">
                  By continuing, you agree to the{" "}
                  <Link
                    href="https://pragmahq.com/privacy"
                    className="underline font-semibold"
                    target="_blank"
                  >
                    Privacy Policy
                  </Link>
                  .
                </p>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
