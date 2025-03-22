"use client";

import Link from "next/link";
import { useRouter } from "next/router";
import React, { useEffect, useState } from "react";
import { FiArrowUpRight, FiDollarSign, FiMoreHorizontal } from "react-icons/fi";

interface Asset {
  id: number;
  name: string;
  type: string;
  ticker?: string;
  price?: number;
  amount: number;
  currency: string;
  interest_rate?: number;
  compounding_frequency?: string;
  interest_start?: string;
  created_at: string;
  returns?: number;
}

export const AssetsTable = () => {
  const [assets, setAssets] = useState<Asset[]>([]);

  useEffect(() => {
    const fetchAssets = async () => {
      try {
        const res = await fetch("https://pfinance.jagactechlab.com/api/assets/all");
        const data: Asset[] = await res.json();
        const retRes = await fetch(`https://pfinance.jagactechlab.com/api/returns`);
        const retData = await retRes.json();
        const assetsWithReturns = data.map(asset => ({
          ...asset,
          returns: retData.returns?.[asset.id] || 0,
        }));

        setAssets(assetsWithReturns);
      } catch (error) {
        console.error("Error fetching assets:", error);
      }
    };

    fetchAssets();
  }, []);


  return (
    <div className="col-span-12 p-4 rounded border border-stone-300">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="flex items-center gap-1.5 font-medium">
          <FiDollarSign /> Your Investments
        </h3>
        <Link href={"/add"}>
        <button className="text-sm text-violet-500 hover:underline">
          Add Asset
        </button>
        </Link>
      </div>
      <table className="w-full table-auto">
        <TableHead />
        <tbody>
          {assets.map((asset, index) => (
            <TableRow key={asset.id} asset={asset} order={index + 1} />
          ))}
        </tbody>
      </table>
    </div>
  );
};

const TableHead = () => (
  <thead>
    <tr className="text-sm font-normal text-stone-500">
      <th className="text-start p-1.5">Name</th>
      <th className="text-start p-1.5">Type</th>
      <th className="text-start p-1.5">Price</th>
      <th className="text-start p-1.5">Amount</th>
      <th className="text-start p-1.5">Currency</th>
      <th className="text-start p-1.5">Returns</th>
      <th className="w-8"></th>
    </tr>
  </thead>
);

const TableRow = ({ asset, order }: { asset: Asset; order: number }) => {
  return (
    <tr className={order % 2 ? "bg-stone-100 text-sm" : "text-sm"}>
      <td className="p-1.5">{asset.name}</td>
      <td className="p-1.5">{asset.type}</td>
      <td className="p-1.5">${asset.price?.toFixed(2) || "N/A"}</td>
      <td className="p-1.5">{asset.amount}</td>
      <td className="p-1.5">{asset.currency}</td>
      <td className={`p-1.5 font-medium ${asset.returns! >= 0 ? "text-green-600" : "text-red-600"}`}>
        {asset.returns ? `$${asset.returns.toFixed(2)}` : "N/A"}
      </td>
      <td className="w-8">
        <button className="hover:bg-stone-200 transition-colors grid place-content-center rounded text-sm size-8">
          <FiMoreHorizontal />
        </button>
      </td>
    </tr>
  );
};
export default AssetsTable;
