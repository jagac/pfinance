import React from "react";
import { FiCalendar } from "react-icons/fi";

export const TopBar = () => {
    const today  = new Date();
  return (
    <div className="border-b px-4 mb-4 mt-2 pb-4 border-stone-200">
      <div className="flex items-center justify-between p-0.5">
        <div>
          <span className="text-sm font-bold block">🚀 Good morning, Jagos!</span>
          <span className="text-xs block text-stone-500">
            {today.toLocaleDateString("en-US")}
          </span>
        </div>
      </div>
    </div>
  );
};
