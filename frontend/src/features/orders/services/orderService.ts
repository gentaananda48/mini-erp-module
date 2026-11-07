import type { ApiResponse, PaginatedApiResponse, OrderListItem, OrderDetail, CreateOrderPayload } from "../../../types/api";

import { apiRequest } from "../../../lib/api";

export async function getOrders(page: number, limit: number): Promise<PaginatedApiResponse<OrderListItem>> {
  return apiRequest<PaginatedApiResponse<OrderListItem>>(`/orders?page=${page}&limit=${limit}`);
}

export async function getOrderById(id: number): Promise<ApiResponse<OrderDetail>> {
  return apiRequest<ApiResponse<OrderDetail>>(`/orders/${id}`);
}

export async function createOrder(payload: CreateOrderPayload, idempotencyKey: string): Promise<ApiResponse<OrderDetail>> {
  return apiRequest<ApiResponse<OrderDetail>>("/orders", {
    method: "POST",
    headers: {
      "Idempotency-Key": idempotencyKey,
    },
    body: JSON.stringify(payload),
  });
}

export async function cancelOrder(id: number): Promise<ApiResponse<{ id: number; status: string }>> {
  return apiRequest<ApiResponse<{ id: number; status: string }>>(`/orders/${id}/cancel`, {
    method: "POST",
  });
}
