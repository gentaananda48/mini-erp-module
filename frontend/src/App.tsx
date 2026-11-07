import { useState } from "react";
import { ProductsPage } from "./features/products/pages/ProductsPage";
import { OrdersPage } from "./features/orders/pages/OrdersPage";

type Page = "products" | "orders";

function App() {
  const [currentPage, setCurrentPage] = useState<Page>("products");

  return <>{currentPage === "products" ? <ProductsPage onNavigateToOrders={() => setCurrentPage("orders")} /> : <OrdersPage onNavigateToProducts={() => setCurrentPage("products")} />}</>;
}

export default App;
