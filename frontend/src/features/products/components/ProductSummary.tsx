import type { Product } from "../../../types/api";

interface ProductSummaryProps {
  products: Product[];
  totalProducts: number;
}

const currencyFormatter = new Intl.NumberFormat("id-ID", {
  style: "currency",
  currency: "IDR",
  maximumFractionDigits: 0,
});

export function ProductSummary({ products, totalProducts }: ProductSummaryProps) {
  const activeProducts = products.filter((product) => product.is_active).length;

  const totalStockValue = products.reduce((total, product) => total + product.price * product.stock, 0);

  const lowStockCount = products.filter((product) => product.stock <= 10).length;

  const summaries = [
    {
      label: "Total Products",
      value: totalProducts.toLocaleString("id-ID"),
    },
    {
      label: "Active Products",
      value: activeProducts.toLocaleString("id-ID"),
    },
    {
      label: "Total Stock Value",
      value: currencyFormatter.format(totalStockValue),
    },
    {
      label: "Low Stock",
      value: lowStockCount.toLocaleString("id-ID"),
    },
  ];

  return (
    <section className="summary-grid">
      {summaries.map((summary) => (
        <article key={summary.label} className="summary-card">
          <p>{summary.label}</p>
          <h2>{summary.value}</h2>
        </article>
      ))}
    </section>
  );
}
