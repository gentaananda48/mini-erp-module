import type { Product } from "../../../types/api";

interface ProductTableProps {
  products: Product[];
  onDetail: (product: Product) => void;
  onAdjustment: (product: Product) => void;
}

const currencyFormatter = new Intl.NumberFormat("id-ID", {
  style: "currency",
  currency: "IDR",
  maximumFractionDigits: 0,
});

function formatLastMovement(value: string | null): string {
  if (!value) return "Never";

  const timestamp = new Date(value).getTime();

  if (Number.isNaN(timestamp)) return "—";

  const diffSeconds = Math.floor((Date.now() - timestamp) / 1000);

  if (diffSeconds < 60) return "Just now";

  const minutes = Math.floor(diffSeconds / 60);

  if (minutes < 60) return `${minutes}m ago`;

  const hours = Math.floor(minutes / 60);

  if (hours < 24) return `${hours}h ago`;

  const days = Math.floor(hours / 24);

  return `${days}d ago`;
}

export function ProductTable({ products, onDetail, onAdjustment }: ProductTableProps) {
  return (
    <div className="table-wrapper">
      <table className="product-table">
        <thead>
          <tr>
            <th>SKU</th>
            <th>Product Name</th>
            <th>Price</th>
            <th>Stock</th>
            <th>Status</th>
            <th>Last Movement</th>
            <th>Actions</th>
          </tr>
        </thead>

        <tbody>
          {products.map((product) => (
            <tr key={product.id}>
              <td className="sku-cell">{product.sku}</td>
              <td className="product-name">{product.name}</td>
              <td>{currencyFormatter.format(product.price)}</td>

              <td>
                <span className={product.stock <= 10 ? "stock-low" : "stock-normal"}>{product.stock}</span>
              </td>

              <td>
                <span className={`status-badge ${product.is_active ? "active" : "inactive"}`}>{product.is_active ? "Active" : "Inactive"}</span>
              </td>

              <td>{formatLastMovement(product.last_movement_at)}</td>

              <td>
                <div style={{ display: "flex", gap: "8px" }}>
                  <button type="button" className="action-button" onClick={() => onDetail(product)}>
                    Detail
                  </button>

                  <button type="button" className="secondary-button" onClick={() => onAdjustment(product)}>
                    Adjustment
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
