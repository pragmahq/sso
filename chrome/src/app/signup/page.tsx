import Page from "./signup";
import { cookies } from "next/headers";

export default function AuthPage() {
  const token = cookies().get("Token")?.value;

  return <Page token={token || ""} />;
}
