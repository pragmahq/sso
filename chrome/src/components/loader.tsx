"use client";

import Image from "next/image";
import { useTheme } from "next-themes";

export default function Loader() {
  const { theme } = useTheme();

  return (
    <div className="w-screen h-screen flex items-center justify-center bg-white dark:bg-black">
      <div className="relative w-16 h-16 discord-spin">
        <Image
          src="/logo.svg"
          alt="Loading"
          layout="fill"
          className={`${theme === "dark" ? "invert" : ""}`}
        />
      </div>
    </div>
  );
}
