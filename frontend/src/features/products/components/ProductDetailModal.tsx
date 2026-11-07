import { useEffect, useState } from "react";
import type { Product, StockMovement } from "../../../types/api";
import { updateProduct, deleteProduct, getStockMovements } from "../services/productService";

interface ProductDetailModalProps {
  product: Product;
  onClose: () => void;
  onSuccess: () => void;
}

interface MovementPagination {
  page: number;
  limit: number;
  total: number;
}

const formatCurrency = (value: number) =>
  new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);

const formatDate = (value: string) =>
  new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));

export function ProductDetailModal({ product, onClose, onSuccess }: ProductDetailModalProps) {
  const [isEditing, setIsEditing] = useState(false);

  const [sku, setSku] = useState(product.sku);
  const [name, setName] = useState(product.name);
  const [price, setPrice] = useState(String(product.price));

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Stock movement history
  const [movements, setMovements] = useState<StockMovement[]>([]);
  const [movementPagination, setMovementPagination] = useState<MovementPagination | null>(null);
  const [movementPage, setMovementPage] = useState(1);
  const [movementLoading, setMovementLoading] = useState(true);
  const [movementError, setMovementError] = useState<string | null>(null);
  const [movementRetryCount, setMovementRetryCount] = useState(0);

  // Reset form when product changes
  useEffect(() => {
    setSku(product.sku);
    setName(product.name);
    setPrice(String(product.price));
    setIsEditing(false);
    setError(null);
  }, [product]);

  // Fetch stock movement history
  useEffect(() => {
    let cancelled = false;

    async function fetchMovements() {
      setMovementLoading(true);
      setMovementError(null);

      try {
        const response = await getStockMovements(product.id, {
          page: movementPage,
          limit: 5,
        });

        if (!cancelled) {
          setMovements(response.data);
          setMovementPagination({
            page: response.page,
            limit: response.limit,
            total: response.total,
          });
        }
      } catch (err) {
        if (!cancelled) {
          setMovementError(err instanceof Error ? err.message : "Failed to fetch stock movements.");
        }
      } finally {
        if (!cancelled) {
          setMovementLoading(false);
        }
      }
    }

    fetchMovements();

    return () => {
      cancelled = true;
    };
  }, [product.id, movementPage, movementRetryCount]);

  async function handleUpdate() {
    const parsedPrice = Number(price);

    if (!sku.trim() || !name.trim()) {
      setError("SKU and product name are required.");
      return;
    }

    if (!Number.isSafeInteger(parsedPrice) || parsedPrice <= 0) {
      setError("Price must be a positive whole number.");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      await updateProduct(product.id, {
        sku: sku.trim(),
        name: name.trim(),
        price: parsedPrice,
      });

      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update product.");
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete() {
    const confirmed = window.confirm(`Are you sure you want to delete ${product.name}?`);

    if (!confirmed) return;

    setLoading(true);
    setError(null);

    try {
      await deleteProduct(product.id);
      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete product.");
    } finally {
      setLoading(false);
    }
  }

  function handleCancelEdit() {
    setSku(product.sku);
    setName(product.name);
    setPrice(String(product.price));
    setError(null);
    setIsEditing(false);
  }

  const totalMovementPages = movementPagination ? Math.ceil(movementPagination.total / movementPagination.limit) : 0;

  return (
    <div className="modal-overlay">
      <div className="modal-container">
        <div className="modal-header">
          <h2>{isEditing ? "Edit Product" : "Product Detail"}</h2>

          <button type="button" className="modal-close" onClick={onClose} disabled={loading}>
            ×
          </button>
        </div>

        <div className="modal-body">
          {error && <div className="form-error">{error}</div>}

          {isEditing ? (
            <>
              <div className="form-group">
                <label>SKU</label>
                <input value={sku} onChange={(e) => setSku(e.target.value)} disabled={loading} />
              </div>

              <div className="form-group">
                <label>Product Name</label>
                <input value={name} onChange={(e) => setName(e.target.value)} disabled={loading} />
              </div>

              <div className="form-group">
                <label>Price (IDR)</label>
                <input
                  type="text"
                  inputMode="numeric"
                  value={price}
                  onChange={(e) => {
                    const value = e.target.value;
                    if (/^\d*$/.test(value)) setPrice(value);
                  }}
                  disabled={loading}
                />
              </div>
            </>
          ) : (
            <>
              <div className="product-detail-list">
                <div>
                  <span>SKU</span>
                  <strong>{product.sku}</strong>
                </div>

                <div>
                  <span>Product Name</span>
                  <strong>{product.name}</strong>
                </div>

                <div>
                  <span>Price</span>
                  <strong>{formatCurrency(product.price)}</strong>
                </div>

                <div>
                  <span>Stock</span>
                  <strong>{product.stock}</strong>
                </div>

                <div>
                  <span>Status</span>
                  <strong>{product.is_active ? "Active" : "Inactive"}</strong>
                </div>
              </div>

              <div className="movement-history">
                <div className="movement-history-header">
                  <h3>Stock Movement History</h3>
                  <span>{movementPagination?.total ?? 0} records</span>
                </div>

                {movementLoading ? (
                  <p>Loading stock movements...</p>
                ) : movementError ? (
                  <div className="form-error">
                    <p>{movementError}</p>
                    <button type="button" className="secondary-button" onClick={() => setMovementRetryCount((prev) => prev + 1)}>
                      Retry
                    </button>
                  </div>
                ) : movements.length === 0 ? (
                  <p>No stock movement history yet.</p>
                ) : (
                  <>
                    <div className="movement-history-list">
                      {movements.map((movement) => (
                        <div className="movement-history-item" key={movement.id}>
                          <div className="movement-history-info">
                            <strong>{movement.type === "in" ? "Stock In" : "Stock Out"}</strong>
                            <span>{formatDate(movement.createdAt)}</span>
                            <span>{movement.note || "-"}</span>
                          </div>

                          <div className="movement-history-quantity">
                            <strong className={movement.type === "in" ? "movement-in" : "movement-out"}>
                              {movement.type === "in" ? "+" : "-"}
                              {movement.qty}
                            </strong>
                          </div>
                        </div>
                      ))}
                    </div>

                    <div className="movement-history-pagination">
                      <button type="button" className="secondary-button" disabled={movementPage <= 1} onClick={() => setMovementPage((prev) => prev - 1)}>
                        Previous
                      </button>

                      <span>
                        Page {movementPagination?.page ?? movementPage} of {totalMovementPages}
                      </span>

                      <button type="button" className="secondary-button" disabled={movementPage >= totalMovementPages} onClick={() => setMovementPage((prev) => prev + 1)}>
                        Next
                      </button>
                    </div>
                  </>
                )}
              </div>
            </>
          )}
        </div>

        <div className="modal-footer">
          {isEditing ? (
            <>
              <button type="button" className="secondary-button" onClick={handleCancelEdit} disabled={loading}>
                Cancel
              </button>

              <button type="button" className="primary-button" onClick={handleUpdate} disabled={loading}>
                {loading ? "Saving..." : "Save Changes"}
              </button>
            </>
          ) : (
            <>
              <button type="button" className="danger-button" onClick={handleDelete} disabled={loading}>
                Delete
              </button>

              <button type="button" className="primary-button" onClick={() => setIsEditing(true)} disabled={loading}>
                Edit Product
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
