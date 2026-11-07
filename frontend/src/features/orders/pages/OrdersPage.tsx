import { useEffect, useMemo, useState } from "react";

import { OrderTable } from "../components/OrderTable";
import { OrderFilters, type OrderFilterType } from "../components/OrderFilters";
import { OrderSummary } from "../components/OrderSummary";
import { OrderPagination } from "../components/OrderPagination";

import { useDebounce } from "../../products/hooks/useDebounce";
import { useOrders } from "../hooks/useOrders";
import { OrderDetailModal } from "../components/OrderDetailModal";
import type { OrderListItem } from "../../../types/api";
import { CreateOrderModal } from "../components/CreateOrderModal";

interface OrdersPageProps {
  onNavigateToProducts: () => void;
}

export function OrdersPage({ onNavigateToProducts }: OrdersPageProps) {
  const { orders, pagination, loading, error, setPage, retry } = useOrders(5);

  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState<OrderFilterType>("all");
  const [selectedOrder, setSelectedOrder] = useState<OrderListItem | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  const debouncedSearch = useDebounce(search, 300);

  useEffect(() => {
    setPage(1);
  }, [debouncedSearch, filter, setPage]);

  const filteredOrders = useMemo(() => {
    const keyword = debouncedSearch.trim().toLowerCase();

    return orders.filter((order) => {
      const matchesSearch = String(order.id).includes(keyword) || order.customer.toLowerCase().includes(keyword);

      const matchesFilter = (() => {
        switch (filter) {
          case "created":
            return order.status === "CREATED";

          case "cancelled":
            return order.status === "CANCELLED";

          case "all":
          default:
            return true;
        }
      })();

      return matchesSearch && matchesFilter;
    });
  }, [orders, debouncedSearch, filter]);

  return (
    <main className="page-container">
      <header className="page-header">
        <div>
          <h1>Orders</h1>
        </div>

        <div className="page-header-actions">
          <button type="button" className="outline-button" onClick={onNavigateToProducts}>
            Products
          </button>

          <button type="button" className="primary-button" onClick={() => setIsCreateModalOpen(true)}>
            + Create Order
          </button>
        </div>
      </header>

      {!loading && !error && pagination && <OrderSummary orders={orders} totalOrders={pagination.total} />}

      <OrderFilters search={search} filter={filter} onSearchChange={setSearch} onFilterChange={setFilter} />

      {loading && (
        <div className="state-container">
          <p>Loading orders...</p>
        </div>
      )}

      {!loading && error && (
        <div className="state-container error-state">
          <h3>Something went wrong</h3>
          <button type="button" className="primary-button" onClick={retry}>
            Retry
          </button>
        </div>
      )}

      {!loading && !error && filteredOrders.length === 0 && (
        <div className="state-container">
          <h3>No orders found</h3>
          <p>No orders match your current search or filter.</p>
        </div>
      )}

      {!loading && !error && filteredOrders.length > 0 && (
        <>
          <div className="results-info">
            Showing {filteredOrders.length} of {orders.length} loaded orders
            {pagination && ` (${pagination.total} total in database)`}
          </div>

          <OrderTable orders={filteredOrders} onDetail={setSelectedOrder} />

          {pagination && <OrderPagination page={pagination.page} limit={pagination.limit} total={pagination.total} onPageChange={setPage} disabled={loading} />}
        </>
      )}

      {selectedOrder && (
        <OrderDetailModal
          order={selectedOrder}
          onClose={() => setSelectedOrder(null)}
          onSuccess={() => {
            setSelectedOrder(null);
            retry();
          }}
        />
      )}

      {isCreateModalOpen && (
        <CreateOrderModal
          onClose={() => setIsCreateModalOpen(false)}
          onSuccess={() => {
            setIsCreateModalOpen(false);
            setPage(1);
            retry();
          }}
        />
      )}
    </main>
  );
}
