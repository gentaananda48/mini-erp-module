import { useState } from "react";
import { createProduct } from "../services/productService";

interface AddProductModalProps {
  onClose: () => void;
  onSuccess: () => void;
}

export function AddProductModal({ onClose, onSuccess }: AddProductModalProps) {
  const [sku, setSku] = useState("");
  const [name, setName] = useState("");
  const [price, setPrice] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit() {
    const parsedPrice = Number(price);

    if (!sku.trim() || !name.trim() || !price.trim()) {
      setError("SKU, product name, and price are required.");
      return;
    }

    if (!Number.isSafeInteger(parsedPrice) || parsedPrice <= 0) {
      setError("Price must be a positive whole number.");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      await createProduct({
        sku: sku.trim(),
        name: name.trim(),
        price: parsedPrice,
      });

      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create product.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="modal-overlay">
      <div className="modal-container">
        <div className="modal-header">
          <h2>Add Product</h2>

          <button type="button" className="modal-close" onClick={onClose} disabled={loading}>
            ×
          </button>
        </div>

        <div className="modal-body">
          {error && <div className="form-error">{error}</div>}

          <div className="form-group">
            <label>SKU</label>
            <input value={sku} onChange={(e) => setSku(e.target.value)} placeholder="e.g. SKU-015" disabled={loading} />
          </div>

          <div className="form-group">
            <label>Product Name</label>
            <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Enter product name" disabled={loading} />
          </div>

          <div className="form-group">
            <label>Price (IDR)</label>
            <input
              type="text"
              inputMode="numeric"
              value={price}
              onChange={(e) => {
                const value = e.target.value;

                if (/^\d*$/.test(value)) {
                  setPrice(value);
                }
              }}
              placeholder="Enter product price"
              disabled={loading}
            />
          </div>
        </div>

        <div className="modal-footer">
          <button type="button" className="secondary-button" onClick={onClose} disabled={loading}>
            Cancel
          </button>

          <button type="button" className="primary-button" onClick={handleSubmit} disabled={loading}>
            {loading ? "Creating..." : "Create Product"}
          </button>
        </div>
      </div>
    </div>
  );
}
