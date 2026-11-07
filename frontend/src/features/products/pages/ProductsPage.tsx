import { useEffect, useMemo, useState } from "react";

import { ProductTable } from "../components/ProductTable";
import { ProductFilters, type ProductFilterType } from "../components/ProductFilters";
import { ProductSummary } from "../components/ProductSummary";
import { ProductPagination } from "../components/ProductPagination";

import { useDebounce } from "../hooks/useDebounce";
import { useProducts } from "../hooks/useProducts";
import { ProductDetailModal } from "../components/ProductDetailModal";
import type { Product } from "../../../types/api";
import { AddProductModal } from "../components/AddProductModal";
import { StockAdjustmentModal } from "../components/StockAdjustmentModal";

interface ProductsPageProps {
  onNavigateToOrders: () => void;
}

export function ProductsPage({ onNavigateToOrders }: ProductsPageProps) {
  const { products, pagination, loading, error, setPage, retry } = useProducts(5);

  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState<ProductFilterType>("all");
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [adjustmentProduct, setAdjustmentProduct] = useState<Product | null>(null);

  const debouncedSearch = useDebounce(search, 300);

  useEffect(() => {
    setPage(1);
  }, [debouncedSearch, filter, setPage]);

  const filteredProducts = useMemo(() => {
    const keyword = debouncedSearch.trim().toLowerCase();

    return products.filter((product) => {
      const matchesSearch = product.sku.toLowerCase().includes(keyword) || product.name.toLowerCase().includes(keyword);

      const matchesFilter = (() => {
        switch (filter) {
          case "active":
            return product.is_active;

          case "inactive":
            return !product.is_active;

          case "low-stock":
            return product.stock <= 10;

          case "all":
          default:
            return true;
        }
      })();

      return matchesSearch && matchesFilter;
    });
  }, [products, debouncedSearch, filter]);

  return (
    <main className="page-container">
      <header className="page-header">
        <div>
          <h1>Products</h1>
        </div>

        <div className="page-header-actions">
          <button type="button" className="outline-button" onClick={onNavigateToOrders}>
            Orders
          </button>

          <button type="button" className="primary-button" onClick={() => setIsAddModalOpen(true)}>
            + Add Product
          </button>
        </div>
      </header>

      {!loading && !error && pagination && <ProductSummary products={products} totalProducts={pagination.total} />}

      <ProductFilters search={search} filter={filter} onSearchChange={setSearch} onFilterChange={setFilter} />

      {loading && (
        <div className="state-container">
          <p>Loading products...</p>
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

      {!loading && !error && filteredProducts.length === 0 && (
        <div className="state-container">
          <h3>No products found</h3>
          <p>No products match your current search or filter.</p>
        </div>
      )}

      {!loading && !error && filteredProducts.length > 0 && (
        <>
          <div className="results-info">
            Showing {filteredProducts.length} of {products.length} loaded products
            {pagination && ` (${pagination.total} total in database)`}
          </div>

          <ProductTable products={filteredProducts} onDetail={setSelectedProduct} onAdjustment={setAdjustmentProduct} />

          {pagination && <ProductPagination page={pagination.page} limit={pagination.limit} total={pagination.total} onPageChange={setPage} disabled={loading} />}
        </>
      )}

      {selectedProduct && (
        <ProductDetailModal
          product={selectedProduct}
          onClose={() => setSelectedProduct(null)}
          onSuccess={() => {
            setSelectedProduct(null);
            retry();
          }}
        />
      )}

      {isAddModalOpen && (
        <AddProductModal
          onClose={() => setIsAddModalOpen(false)}
          onSuccess={() => {
            setIsAddModalOpen(false);
            setPage(1);
            retry();
          }}
        />
      )}

      {adjustmentProduct && (
        <StockAdjustmentModal
          product={adjustmentProduct}
          onClose={() => setAdjustmentProduct(null)}
          onSuccess={() => {
            setAdjustmentProduct(null);
            retry();
          }}
        />
      )}
    </main>
  );
}
