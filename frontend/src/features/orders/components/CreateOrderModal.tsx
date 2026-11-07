import { useEffect, useState } from "react";

import type { Product } from "../../../types/api";
import { getProducts } from "../../products/services/productService";
import { createOrder } from "../services/orderService";

interface CreateOrderModalProps {
  onClose: () => void;
  onSuccess: () => void;
}

interface OrderItemInput {
  sku: string;
  qty: string;
}

function generateIdempotencyKey() {
  return crypto.randomUUID();
}

export function CreateOrderModal({ onClose, onSuccess }: CreateOrderModalProps) {
  const [customer, setCustomer] = useState("");
  const [items, setItems] = useState<OrderItemInput[]>([{ sku: "", qty: "1" }]);

  const [products, setProducts] = useState<Product[]>([]);
  const [productsLoading, setProductsLoading] = useState(true);
  const [productsError, setProductsError] = useState<string | null>(null);

  const [idempotencyKey, setIdempotencyKey] = useState(generateIdempotencyKey);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function fetchAllProducts() {
      try {
        setProductsLoading(true);
        setProductsError(null);

        const limit = 100;

        const firstResponse = await getProducts({
          page: 1,
          limit,
        });

        let allProducts = [...firstResponse.data];

        const totalPages = Math.ceil(firstResponse.total / firstResponse.limit);

        for (let page = 2; page <= totalPages; page++) {
          const response = await getProducts({
            page,
            limit: firstResponse.limit,
          });

          allProducts = [...allProducts, ...response.data];
        }

        if (!cancelled) {
          setProducts(allProducts.filter((product) => product.is_active));
        }
      } catch (err) {
        if (!cancelled) {
          setProductsError(err instanceof Error ? err.message : "Failed to fetch products.");
        }
      } finally {
        if (!cancelled) {
          setProductsLoading(false);
        }
      }
    }

    void fetchAllProducts();

    return () => {
      cancelled = true;
    };
  }, []);

  function updateCustomer(value: string) {
    setCustomer(value);
    setIdempotencyKey(generateIdempotencyKey());
  }

  function updateItem(index: number, field: keyof OrderItemInput, value: string) {
    setItems((current) => current.map((item, itemIndex) => (itemIndex === index ? { ...item, [field]: value } : item)));

    setIdempotencyKey(generateIdempotencyKey());
  }

  function addItem() {
    setItems((current) => [...current, { sku: "", qty: "1" }]);

    setIdempotencyKey(generateIdempotencyKey());
  }

  function removeItem(index: number) {
    setItems((current) => current.filter((_, itemIndex) => itemIndex !== index));

    setIdempotencyKey(generateIdempotencyKey());
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError(null);

    if (productsLoading || productsError) {
      setError("Product list is not ready. Please try again.");
      return;
    }

    const normalizedCustomer = customer.trim();

    if (!normalizedCustomer) {
      setError("Customer name is required.");
      return;
    }

    if (items.length === 0) {
      setError("At least one order item is required.");
      return;
    }

    const normalizedItems = items.map((item) => ({
      sku: item.sku.trim(),
      qty: Number(item.qty),
    }));

    if (normalizedItems.some((item) => !item.sku)) {
      setError("Please select a product for every order item.");
      return;
    }

    if (normalizedItems.some((item) => !Number.isSafeInteger(item.qty) || item.qty <= 0)) {
      setError("Quantity must be a positive integer.");
      return;
    }

    const skuList = normalizedItems.map((item) => item.sku.toLowerCase());

    if (new Set(skuList).size !== skuList.length) {
      setError("Duplicate SKU is not allowed.");
      return;
    }

    try {
      setLoading(true);

      await createOrder(
        {
          customer: normalizedCustomer,
          items: normalizedItems,
        },
        idempotencyKey,
      );

      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create order");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="modal-overlay" onClick={loading ? undefined : onClose}>
      <div className="modal-container" onClick={(event) => event.stopPropagation()}>
        <form onSubmit={handleSubmit}>
          <div className="modal-header">
            <div>
              <h2>Create Order</h2>
              <p>Create a new customer order.</p>
            </div>

            <button type="button" className="modal-close" onClick={onClose} disabled={loading}>
              &times;
            </button>
          </div>

          <div className="modal-body">
            <div className="form-group">
              <label htmlFor="order-customer">Customer</label>

              <input id="order-customer" type="text" value={customer} onChange={(event) => updateCustomer(event.target.value)} placeholder="Enter customer name" disabled={loading} />
            </div>

            <div className="order-items-header">
              <h3>Order Items</h3>

              <button type="button" className="outline-button" onClick={addItem} disabled={loading}>
                + Add Item
              </button>
            </div>

            {productsError && <div className="form-error">Failed to load products: {productsError}</div>}

            {items.map((item, index) => {
              const selectedProduct = products.find((product) => product.sku === item.sku);

              return (
                <div className="order-item-row" key={index}>
                  <div className="form-group">
                    <label htmlFor={`order-product-${index}`}>Product</label>

                    <select id={`order-product-${index}`} value={item.sku} onChange={(event) => updateItem(index, "sku", event.target.value)} disabled={loading || productsLoading || !!productsError} required>
                      <option value="">{productsLoading ? "Loading products..." : products.length === 0 ? "No active products available" : "Select product"}</option>

                      {products.map((product) => {
                        const alreadySelected = items.some((selectedItem, selectedIndex) => selectedIndex !== index && selectedItem.sku === product.sku);

                        return (
                          <option key={product.id} value={product.sku} disabled={alreadySelected}>
                            {product.name} ({product.sku})
                          </option>
                        );
                      })}
                    </select>

                    {selectedProduct && <small>Available stock: {selectedProduct.stock}</small>}
                  </div>

                  <div className="form-group">
                    <label htmlFor={`order-qty-${index}`}>Quantity</label>

                    <input
                      id={`order-qty-${index}`}
                      type="text"
                      inputMode="numeric"
                      value={item.qty}
                      onChange={(event) => {
                        const value = event.target.value;

                        if (/^\d*$/.test(value)) {
                          updateItem(index, "qty", value);
                        }
                      }}
                      disabled={loading}
                    />
                  </div>

                  <button type="button" className="remove-item-button" onClick={() => removeItem(index)} disabled={loading || items.length === 1} title="Remove item">
                    &times;
                  </button>
                </div>
              );
            })}

            {error && <div className="form-error">{error}</div>}
          </div>

          <div className="modal-footer">
            <button type="button" className="outline-button" onClick={onClose} disabled={loading}>
              Cancel
            </button>

            <button type="submit" className="primary-button" disabled={loading || productsLoading || !!productsError || products.length === 0}>
              {loading ? "Creating..." : "Create Order"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
