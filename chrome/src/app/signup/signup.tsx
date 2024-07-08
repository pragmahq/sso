"use client";
import { useState, useEffect } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import Image from "next/legacy/image";
import Link from "next/link";
import axios from "axios";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { FaSpinner } from "react-icons/fa6";
import { TbLogin2 } from "react-icons/tb";
import Loader from "@/components/loader";
import ThemeSwitcher from "@/components/theme/theme-switcher";
import { useToast } from "@/components/ui/use-toast";
interface FormData {
  email: string;
  password: string;
  confirmPassword: string;
  [key: string]: string | undefined;
}

export default function SignupPage({ token }: { token: string }) {
  const sp = useSearchParams();
  const router = useRouter();
  const inviteCode = sp.get("invite");
  const { toast } = useToast();
  const [loaded, setLoaded] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState<FormData>({
    email: "",
    password: "",
    confirmPassword: "",
  });
  const [errors, setErrors] = useState<FormData>({
    email: "",
    password: "",
    confirmPassword: "",
  });

  useEffect(() => {
    if (!inviteCode) {
      router.replace("/404");
      return;
    }

    const validateTokenAndInvite = async () => {
      try {
        if (token) {
          await axios.get(
            `${process.env.NEXT_PUBLIC_BACKEND_URL!}/api/auth/validate`,
            { withCredentials: true }
          );
          router.replace("/");
        }
        await axios.get(
          `${process.env
            .NEXT_PUBLIC_BACKEND_URL!}/api/auth/validate-invite/${inviteCode}`
        );
        setLoaded(true);
      } catch (error) {
        if (token) {
          await axios.get("/api/logout", { withCredentials: true });
        }
        router.push("/404");
      }
    };

    validateTokenAndInvite();
  }, [inviteCode, token, router]);

  useEffect(() => {
    validatePasswords();
  }, [formData.password, formData.confirmPassword]);

  const validateEmail = (email: string) =>
    /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(
      email.toLowerCase()
    );

  const validatePassword = (password: string) => password.length >= 8; // Add more conditions as needed

  const validatePasswords = () => {
    const { password, confirmPassword } = formData;
    let passwordError = "";
    let confirmPasswordError = "";

    if (password && !validatePassword(password)) {
      passwordError = "Password must be at least 8 characters long";
    }

    if (confirmPassword && password !== confirmPassword) {
      confirmPasswordError = "Passwords do not match";
    }

    setErrors((prev) => ({
      ...prev,
      password: passwordError,
      confirmPassword: confirmPasswordError,
    }));
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));

    if (name === "email") {
      setErrors((prev) => ({
        ...prev,
        email:
          value && !validateEmail(value)
            ? "Please enter a valid email address"
            : "",
      }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (Object.values(errors).some((error) => error) || !validateForm()) {
      return;
    }
    setIsLoading(true);

    try {
      await axios.post(
        `${process.env.NEXT_PUBLIC_BACKEND_URL!}/api/auth/register`,
        {
          email: formData.email,
          password: formData.password,
          inviteCode,
        }
      );
      router.push("/");
    } catch (error: any) {
      if (error.response?.status === 409) {
        setErrors((prev) => ({
          ...prev,
          email: "Email already in use",
        }));
      } else {
        toast({
          title: "An error occurred",
          description: "Please try again",
          variant: "destructive",
        });
      }
    }
  };

  const validateForm = () => {
    const newErrors = {
      email: !formData.email ? "Email is required" : "",
      password: !formData.password ? "Password is required" : "",
      confirmPassword: !formData.confirmPassword
        ? "Please confirm your password"
        : "",
    };
    setErrors((prev) => ({ ...prev, ...newErrors }));
    return !Object.values(newErrors).some((error) => error);
  };

  if (!loaded) return <Loader />;

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
            <div className="flex flex-col space-y-2 text-center border-[0.5px] rounded-lg py-[3rem] px-3">
              <div className="mb-3">
                <h1 className="text-2xl font-semibold tracking-tight">
                  Authenticate
                </h1>
                <p className="text-sm">
                  You're one step away from having a Pragma ID.
                </p>
              </div>
              <form
                onSubmit={handleSubmit}
                className="flex flex-col text-left mt-10 gap-3 items-center"
              >
                {["email", "password", "confirmPassword"].map((field) => (
                  <div key={field} className="w-[90%]">
                    <Input
                      type={field === "email" ? "email" : "password"}
                      name={field}
                      placeholder={
                        field.charAt(0).toUpperCase() +
                        field.slice(1).replace(/([A-Z])/g, " $1")
                      }
                      className={`border-[0.6px] w-full ${
                        errors[field] ? "border-red-500" : ""
                      }`}
                      onChange={handleInputChange}
                      value={formData[field]}
                      required
                    />
                    {errors[field] && (
                      <p className="text-red-500 text-sm mt-1">
                        {errors[field]}
                      </p>
                    )}
                  </div>
                ))}
                <Input
                  className="w-[90%] border-[0.6px]"
                  value={inviteCode!}
                  disabled
                />
                <Button
                  type="submit"
                  className="dark:bg-white bg-black text-white dark:text-black w-[50%]"
                  disabled={
                    isLoading || Object.values(errors).some((error) => error)
                  }
                >
                  {isLoading ? (
                    <FaSpinner className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <TbLogin2 className="mr-2 h-4 w-4" />
                  )}
                  Sign Up
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
