import { useEffect, useState } from "react";

import type { OrderListItem, OrderDetail } from "../../../types/api";

import { getOrderById, cancelOrder } from "../services/orderService";

interface OrderDetailModalProps {
  order: OrderListItem;
  onClose: () => void;
  onSuccess: () => void;
}

export function OrderDetailModal({ order, onClose, onSuccess }: OrderDetailModalProps) {
  const [detail, setDetail] = useState<OrderDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function fetchDetail() {
      try {
        setLoading(true);
        setError(null);

        const response = await getOrderById(order.id);

        if (!cancelled) {
          setDetail(response.data);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "Failed to load order details");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    fetchDetail();

    return () => {
      cancelled = true;
    };
  }, [order.id]);

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

  async function handleCancel() {
    if (!detail || detail.status !== "CREATED") {
      return;
    }

    const confirmed = window.confirm(`Are you sure you want to cancel order #${detail.id}?`);

    if (!confirmed) {
      return;
    }

    try {
      setSubmitting(true);
      setError(null);

      await cancelOrder(detail.id);

      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to cancel order");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-container" onClick={(event) => event.stopPropagation()}>
        <div className="modal-header">
          <div>
            <h2>Order Details</h2>
            <p>Order #{order.id}</p>
          </div>

          <button type="button" className="modal-close" onClick={onClose} disabled={submitting}>
            &times;
          </button>
        </div>

        <div className="modal-body">
          {loading && (
            <div className="state-container">
              <p>Loading order details...</p>
            </div>
          )}

          {!loading && error && !detail && <div className="form-error">{error}</div>}

          {!loading && detail && (
            <>
              <div className="product-detail-list">
                <div>
                  <span>Customer</span>
                  <strong>{detail.customer}</strong>
                </div>

                <div>
                  <span>Status</span>
                  <strong>{detail.status}</strong>
                </div>

                <div>
                  <span>Created At</span>
                  <strong>{formatDate(detail.createdAt)}</strong>
                </div>

                <div>
                  <span>Order Total</span>
                  <strong>{formatCurrency(detail.total)}</strong>
                </div>
              </div>

              <div className="order-items-section">
                <h3>Order Items</h3>

                <div className="order-items-table-wrapper">
                  <table className="order-items-table">
                    <thead>
                      <tr>
                        <th>SKU</th>
                        <th className="text-center">Qty</th>
                        <th className="text-right">Unit Price</th>
                        <th className="text-right">Subtotal</th>
                      </tr>
                    </thead>

                    <tbody>
                      {detail.items.map((item, index) => (
                        <tr key={`${item.sku}-${index}`}>
                          <td className="sku-cell">{item.sku}</td>
                          <td className="text-center">{item.qty}</td>
                          <td className="text-right">{formatCurrency(item.unitPrice)}</td>
                          <td className="text-right">{formatCurrency(item.subtotal)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {error && <div className="form-error">{error}</div>}
            </>
          )}
        </div>

        <div className="modal-footer">
          <button type="button" className="outline-button" onClick={onClose} disabled={submitting}>
            Close
          </button>

          {!loading && detail?.status === "CREATED" && (
            <button type="button" className="danger-button" onClick={handleCancel} disabled={submitting}>
              {submitting ? "Cancelling..." : "Cancel Order"}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
