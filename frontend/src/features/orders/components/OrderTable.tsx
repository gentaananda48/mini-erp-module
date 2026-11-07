import type { OrderListItem } from "../../../types/api";

interface OrderTableProps {
  orders: OrderListItem[];
  onDetail: (order: OrderListItem) => void;
}

export function OrderTable({ orders, onDetail }: OrderTableProps) {
  function formatCurrency(amount: number) {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      maximumFractionDigits: 0,
    }).format(amount);
  }

  function formatDate(date: string) {
    return new Date(date).toLocaleString("id-ID", {
      dateStyle: "medium",
      timeStyle: "short",
    });
  }

  if (orders.length === 0) {
    return (
      <div className="state-container">
        <h3>No orders found</h3>
        <p>There are currently no orders to display.</p>
      </div>
    );
  }

  return (
    <div className="table-wrapper">
      <table className="product-table">
        <thead>
          <tr>
            <th>Customer</th>
            <th>Status</th>
            <th>Total</th>
            <th>Created At</th>
            <th>Action</th>
          </tr>
        </thead>

        <tbody>
          {orders.map((order) => (
            <tr key={order.id}>
              <td>{order.customer}</td>
              <td>
                <span className={order.status === "CREATED" ? "status-badge status-created" : "status-badge status-cancelled"}>{order.status}</span>
              </td>
              <td>{formatCurrency(order.total)}</td>
              <td>{formatDate(order.createdAt)}</td>
              <td>
                <button type="button" className="action-button" onClick={() => onDetail(order)}>
                  Detail
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
