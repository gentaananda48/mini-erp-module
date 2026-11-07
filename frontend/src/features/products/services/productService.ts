import { apiRequest } from "../../../lib/api";
import type { ApiResponse, PaginatedApiResponse, Product, CreateStockMovementPayload, StockMovement } from "../../../types/api";

export interface GetProductsParams {
  page?: number;
  limit?: number;
}

export interface GetStockMovementsParams {
  page?: number;
  limit?: number;
}

export interface CreateProductPayload {
  sku: string;
  name: string;
  price: number;
}

export interface UpdateProductPayload {
  sku: string;
  name: string;
  price: number;
}

export async function getProducts(params: GetProductsParams = {}): Promise<PaginatedApiResponse<Product>> {
  const searchParams = new URLSearchParams({
    page: String(params.page ?? 1),
    limit: String(params.limit ?? 10),
  });

  return apiRequest<PaginatedApiResponse<Product>>(`/products?${searchParams.toString()}`);
}

export async function createProduct(payload: CreateProductPayload): Promise<void> {
  await apiRequest<void>("/products", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function updateProduct(id: number, payload: UpdateProductPayload): Promise<Product> {
  return apiRequest<Product>(`/products/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export async function deleteProduct(id: number): Promise<void> {
  await apiRequest<void>(`/products/${id}`, {
    method: "DELETE",
  });
}

export async function getStockMovements(productId: number, params: GetStockMovementsParams = {}): Promise<PaginatedApiResponse<StockMovement>> {
  const searchParams = new URLSearchParams({
    page: String(params.page ?? 1),
    limit: String(params.limit ?? 5),
  });

  return apiRequest<PaginatedApiResponse<StockMovement>>(`/products/${productId}/movements?${searchParams.toString()}`);
}

export async function createStockMovement(productId: number, payload: CreateStockMovementPayload): Promise<ApiResponse<StockMovement>> {
  return apiRequest<ApiResponse<StockMovement>>(`/products/${productId}/stock`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}
