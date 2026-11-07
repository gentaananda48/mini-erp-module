import { useState } from "react";
import type { Product, CreateStockMovementPayload } from "../../../types/api";
import { createStockMovement } from "../services/productService";

interface StockAdjustmentModalProps {
  product: Product;
  onClose: () => void;
  onSuccess: () => void;
}

export function StockAdjustmentModal({ product, onClose, onSuccess }: StockAdjustmentModalProps) {
  const [type, setType] = useState<CreateStockMovementPayload["type"]>("in");

  const [qty, setQty] = useState("");
  const [note, setNote] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const parsedQty = Number(qty);

  const isValidQty = qty.trim() !== "" && Number.isSafeInteger(parsedQty) && parsedQty > 0;

  const afterAdjustment = type === "in" ? product.stock + (isValidQty ? parsedQty : 0) : product.stock - (isValidQty ? parsedQty : 0);

  async function handleSubmit() {
    if (!isValidQty) {
      setError("Quantity must be a positive whole number.");
      return;
    }

    if (type === "out" && parsedQty > product.stock) {
      setError(`Insufficient stock. Available stock: ${product.stock}.`);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      await createStockMovement(product.id, {
        type,
        qty: parsedQty,
        note: note.trim(),
      });

      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create stock adjustment.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="modal-overlay">
      <div className="modal-container">
        <div className="modal-header">
          <h2>Stock Adjustment</h2>

          <button type="button" className="modal-close" onClick={onClose} disabled={loading}>
            ×
          </button>
        </div>

        <div className="modal-body">
          {error && <div className="form-error">{error}</div>}

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
              <span>Current Stock</span>
              <strong>{product.stock}</strong>
            </div>
          </div>

          <div className="form-group">
            <label>Movement Type</label>

            <div style={{ display: "flex", gap: "20px", marginTop: "8px" }}>
              <label style={{ display: "flex", alignItems: "center", gap: "6px" }}>
                <input
                  type="radio"
                  name="movementType"
                  value="in"
                  checked={type === "in"}
                  onChange={() => {
                    setType("in");
                    setError(null);
                  }}
                  disabled={loading}
                />
                Stock IN
              </label>

              <label style={{ display: "flex", alignItems: "center", gap: "6px" }}>
                <input
                  type="radio"
                  name="movementType"
                  value="out"
                  checked={type === "out"}
                  onChange={() => {
                    setType("out");
                    setError(null);
                  }}
                  disabled={loading}
                />
                Stock OUT
              </label>
            </div>
          </div>

          <div className="form-group">
            <label>Quantity</label>

            <input
              type="text"
              inputMode="numeric"
              value={qty}
              onChange={(e) => {
                const value = e.target.value;

                if (/^\d*$/.test(value)) {
                  setQty(value);
                  setError(null);
                }
              }}
              placeholder="Enter quantity"
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label>Note</label>

            <textarea value={note} onChange={(e) => setNote(e.target.value)} placeholder="Enter adjustment reason..." rows={3} disabled={loading} />
          </div>

          <div className="product-detail-list">
            <div>
              <span>Current Stock</span>
              <strong>{product.stock}</strong>
            </div>

            <div>
              <span>Adjustment</span>
              <strong>{isValidQty ? `${type === "in" ? "+" : "-"}${parsedQty}` : "—"}</strong>
            </div>

            <div>
              <span>Stock After Adjustment</span>
              <strong className={afterAdjustment < 0 ? "stock-low" : "stock-normal"}>{isValidQty ? afterAdjustment : "—"}</strong>
            </div>
          </div>
        </div>

        <div className="modal-footer">
          <button type="button" className="secondary-button" onClick={onClose} disabled={loading}>
            Cancel
          </button>

          <button type="button" className="primary-button" onClick={handleSubmit} disabled={loading || !isValidQty}>
            {loading ? "Processing..." : "Submit Adjustment"}
          </button>
        </div>
      </div>
    </div>
  );
}
