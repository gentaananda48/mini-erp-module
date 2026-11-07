interface OrdersPageProps {
  onNavigateToProducts: () => void;
}

export function OrdersPage({ onNavigateToProducts }: OrdersPageProps) {
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

          <button type="button" className="primary-button">
            + Create Order
          </button>
        </div>
      </header>

      <div className="state-container">
        <h3>No orders yet</h3>
        <p>Orders will appear here once you create one.</p>
      </div>
    </main>
  );
}
