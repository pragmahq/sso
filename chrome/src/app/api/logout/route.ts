import { cookies } from "next/headers";

export async function GET() {
  cookies().delete("Token");

  return Response.json({ message: "Logged Out." });
}
