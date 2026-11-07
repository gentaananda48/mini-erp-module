import type { OrderListItem } from "../../../types/api";

interface OrderSummaryProps {
  orders: OrderListItem[];
  totalOrders: number;
}

export function OrderSummary({ orders, totalOrders }: OrderSummaryProps) {
  const createdOrders = orders.filter((order) => order.status === "CREATED").length;

  const cancelledOrders = orders.filter((order) => order.status === "CANCELLED").length;

  const totalValue = orders.filter((order) => order.status === "CREATED").reduce((sum, order) => sum + order.total, 0);

  const formatCurrency = (amount: number) =>
    new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      maximumFractionDigits: 0,
    }).format(amount);

  return (
    <div className="summary-grid">
      <div className="summary-card">
        <span className="summary-label">Total Orders</span>
        <h2>{totalOrders}</h2>
      </div>

      <div className="summary-card">
        <span className="summary-label">Created Orders</span>
        <h2>{createdOrders}</h2>
      </div>

      <div className="summary-card">
        <span className="summary-label">Cancelled Orders</span>
        <h2>{cancelledOrders}</h2>
      </div>

      <div className="summary-card">
        <span className="summary-label">Created Order Value</span>
        <h2>{formatCurrency(totalValue)}</h2>
      </div>
    </div>
  );
}
