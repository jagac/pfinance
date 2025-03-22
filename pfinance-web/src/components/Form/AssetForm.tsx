"use client";
import { Asset } from "@/types/asset";
import { useState } from "react";
import { FiDollarSign } from "react-icons/fi";
import SelectField from "./Select";
import InputField from "./Input";

export default function AddAssetForm() {
  const [form, setForm] = useState<Asset>({
    type: "Stock",
    name: "",
    ticker: "",
    price: undefined,
    amount: 0,
    currency: "USD",
    interestRate: undefined,
    compoundingFrequency: "monthly",
    interestStart: "",
    createdAt: new Date().toISOString(),
  });

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const dataToSend = {
      ...form,
      price: form.price !== undefined ? Number(form.price) : null,
      interestRate:
        form.interestRate !== undefined ? Number(form.interestRate) : null,
      amount: form.amount !== undefined ? Number(form.amount) : null,
      interestStart: form.interestStart
        ? new Date(form.interestStart).toISOString()
        : undefined,
    };

    const response = await fetch(
      "https://pfinance.jagactechlab.com/api/assets/new",
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(dataToSend),
      }
    );

    if (response.ok) {
      alert("Asset added successfully!");
      setForm({
        type: "Stock",
        name: "",
        ticker: "",
        price: undefined,
        amount: 0,
        currency: "USD",
        interestRate: undefined,
        compoundingFrequency: "monthly",
        interestStart: "",
        createdAt: new Date().toISOString(),
      });
    } else {
      alert("Error adding asset.");
    }
  };

  return (
    <section className="relative bg-white rounded-lg">
      <div className="container mx-auto px-6 sm:px-8">
        <div className="mt-12  p-8 rounded-2xl max-w-4xl mx-auto rounded border border-stone-300">
          <form onSubmit={handleSubmit}>
            <InputField
              label="Asset Name"
              name="name"
              type="text"
              value={form.name}
              onChange={handleChange}
              placeholder="Enter asset name"
              required
            />
            <SelectField
              label="Asset Type"
              name="type"
              value={form.type}
              onChange={handleChange}
              options={["Stock", "Gold", "Bond", "Savings", "Crypto"]}
            />
            {form.type !== "Savings" && (
              <InputField
                label="Ticker"
                name="ticker"
                type="text"
                value={form.ticker}
                onChange={handleChange}
                placeholder="Enter ticker symbol"
              />
            )}
            {form.type !== "Savings" && (
              <InputField
                label="Price per Unit"
                name="price"
                type="number"
                value={form.price}
                onChange={handleChange}
                placeholder="Enter price"
                required
              />
            )}
            <InputField
              label="Amount"
              name="amount"
              type="number"
              value={form.amount}
              onChange={handleChange}
              placeholder="Enter amount"
              required
            />
            <SelectField
              label="Currency"
              name="currency"
              value={form.currency}
              onChange={handleChange}
              options={["USD", "EUR"]}
            />
            {(form.type === "Savings" || form.type === "Bond") && (
              <InputField
                label="Interest Start Date"
                name="interestStart"
                type="date"
                value={form.interestStart}
                onChange={handleChange}
                required
              />
            )}
            {form.type === "Savings" && (
              <>
                <InputField
                  label="Interest Rate (%)"
                  name="interestRate"
                  type="number"
                  value={form.interestRate}
                  onChange={handleChange}
                  placeholder="Enter interest rate"
                  required
                />
                <SelectField
                  label="Compounding Frequency"
                  name="compoundingFrequency"
                  value={form.compoundingFrequency}
                  onChange={handleChange}
                  options={["daily", "monthly", "quarterly", "annually"]}
                />
              </>
            )}
            <button
              type="submit"
              className="mt-4 w-full flex items-center justify-center gap-2 p-2 bg-violet-600 text-white rounded hover:bg-violet-700 transition"
            >
              <FiDollarSign /> Save Asset
            </button>
          </form>
        </div>
      </div>
    </section>
  );
}
