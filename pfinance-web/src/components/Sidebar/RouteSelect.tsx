"use client";

import React, { useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { IconType } from "react-icons";
import { FiList, FiHome } from "react-icons/fi";

const routes = [
  { title: "Dashboard", Icon: FiHome, path: "/" },
  { title: "Form", Icon: FiList, path: "/add" },
];

export const RouteSelect = () => {
  const router = useRouter();
  const pathname = usePathname();

  return (
    <div className="space-y-1">
      {routes.map(({ title, Icon, path }) => (
        <Route
          key={title}
          Icon={Icon}
          title={title}
          selected={pathname === path}
          onClick={() => router.push(path)}
        />
      ))}
    </div>
  );
};

const Route = ({
  selected,
  Icon,
  title,
  onClick,
}: {
  selected: boolean;
  Icon: IconType;
  title: string;
  onClick: () => void;
}) => {
  return (
    <button
      onClick={onClick}
      className={`flex items-center justify-start gap-2 w-full rounded px-2 py-1.5 text-sm transition-all ${
        selected
          ? "bg-white text-stone-950 shadow"
          : "hover:bg-stone-200 bg-transparent text-stone-500"
      }`}
    >
      <Icon className={selected ? "text-violet-500" : ""} />
      <span>{title}</span>
    </button>
  );
};
