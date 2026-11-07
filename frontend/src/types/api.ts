export interface ApiResponse<T> {
  httpCode: number;
  message: string;
  data: T;
}

export interface PaginatedApiResponse<T> extends ApiResponse<T[]> {
  page: number;
  limit: number;
  total: number;
}

export interface Product {
  id: number;
  sku: string;
  name: string;
  price: number;
  stock: number;
  is_active: boolean;
  created_at: string;
  last_movement_at: string | null;
}

export interface CreateStockMovementPayload {
  type: "in" | "out";
  qty: number;
  note: string;
}

export interface StockMovement {
  id: number;
  productId: number;
  type: "in" | "out";
  qty: number;
  note: string;
  available: number;
  createdAt: string;
}
