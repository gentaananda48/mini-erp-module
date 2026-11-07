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

export interface CreateOrderItemPayload {
  sku: string;
  qty: number;
}

export interface CreateOrderPayload {
  customer: string;
  items: CreateOrderItemPayload[];
}

export interface OrderItem {
  sku: string;
  qty: number;
  unitPrice: number;
  subtotal: number;
}

export interface Order {
  id: number;
  customer: string;
  status: string;
  total: number;
  items: OrderItem[];
  createdAt: string;
}

export interface CancelOrderResponse {
  id: number;
  status: string;
}

export interface OrderListItem {
  id: number;
  customer: string;
  status: "CREATED" | "CANCELLED";
  total: number;
  createdAt: string;
}

export interface OrderItem {
  sku: string;
  qty: number;
  unitPrice: number;
  subtotal: number;
}

export interface OrderDetail {
  id: number;
  customer: string;
  status: "CREATED" | "CANCELLED";
  total: number;
  items: OrderItem[];
  createdAt: string;
}

export interface CreateOrderPayload {
  customer: string;
  items: {
    sku: string;
    qty: number;
  }[];
}
