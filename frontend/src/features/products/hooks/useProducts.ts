import { useEffect, useState } from "react";

import type { Product } from "../../../types/api";
import { getProducts } from "../services/productService";

interface Pagination {
  page: number;
  limit: number;
  total: number;
}

interface UseProductsResult {
  products: Product[];
  pagination: Pagination | null;
  loading: boolean;
  error: string | null;
  page: number;
  setPage: (page: number) => void;
  retry: () => void;
}

export function useProducts(limit = 10): UseProductsResult {
  const [products, setProducts] = useState<Product[]>([]);
  const [pagination, setPagination] = useState<Pagination | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [retryCount, setRetryCount] = useState(0);

  const retry = () => {
    setRetryCount((prev) => prev + 1);
  };
  useEffect(() => {
    let cancelled = false;
    async function fetchProducts() {
      setLoading(true);
      setError(null);

      try {
        const response = await getProducts({ page, limit });

        if (!cancelled) {
          setProducts(response.data);
          setPagination({
            page: response.page,
            limit: response.limit,
            total: response.total,
          });
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "Failed to fetch products.");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void fetchProducts();

    return () => {
      cancelled = true;
    };
  }, [page, limit, retryCount]);

  return {
    products,
    pagination,
    loading,
    error,
    page,
    setPage,
    retry,
  };
}
