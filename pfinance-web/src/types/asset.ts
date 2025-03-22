export interface Asset {
    id?: number;
    type: "Stock" | "Gold" | "Bond" | "Savings" | "Crypto";
    name: string;
    ticker?: string;
    price?: number;
    amount: number;
    currency: string;
    interestRate?: number;
    compoundingFrequency?: "daily" | "monthly" | "quarterly" | "annually";
    interestStart?: string;
    createdAt: string;
  }
