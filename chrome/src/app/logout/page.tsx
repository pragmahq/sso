"use client";
import axios from "axios";
import { useRouter } from "next/navigation";

export default function LogoutPage() {
  const router = useRouter();

  axios.get("/api/logout").then(() => {
    router.push("/auth");
  });

  return;
}
