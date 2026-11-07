import { useEffect, useState } from "react";
import type { OrderListItem } from "../../../types/api";
import { getOrders } from "../services/orderService";

interface UseOrdersResult {
  orders: OrderListItem[];
  pagination: {
    page: number;
    limit: number;
    total: number;
  };
  loading: boolean;
  error: string | null;
  setPage: (page: number) => void;
  retry: () => void;
}

export function useOrders(limit = 10): UseOrdersResult {
  const [orders, setOrders] = useState<OrderListItem[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [retryCount, setRetryCount] = useState(0);

  useEffect(() => {
    let cancelled = false;

    async function fetchOrders() {
      setLoading(true);
      setError(null);

      try {
        const response = await getOrders(page, limit);

        if (cancelled) return;

        setOrders(response.data);
        setTotal(response.total);
      } catch (err) {
        if (cancelled) return;

        setError(err instanceof Error ? err.message : "Failed to load orders");
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    fetchOrders();

    return () => {
      cancelled = true;
    };
  }, [page, limit, retryCount]);

  function retry() {
    setRetryCount((current) => current + 1);
  }

  return {
    orders,
    pagination: {
      page,
      limit,
      total,
    },
    loading,
    error,
    setPage,
    retry,
  };
}
