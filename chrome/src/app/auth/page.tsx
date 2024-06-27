import { cookies } from "next/headers";
import AuthenticationPage from "./auth";

export default function AuthPage() {
  const token = cookies().get("Token")?.value;
  return <AuthenticationPage token={token} />;
}
