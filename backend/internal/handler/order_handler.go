package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"backend/internal/dto"
	"backend/internal/repository"
	"backend/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

type OrderErrorDetail struct {
	Code      string `json:"code"`
	SKU       string `json:"sku,omitempty"`
	Available int64  `json:"available,omitempty"`
}

type OrderErrorResponse struct {
	CodeHttp int               `json:"codeHttp"`
	Message  string            `json:"message"`
	Error    OrderErrorDetail  `json:"error"`
}

func (h *OrderHandler) CreateOrder(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.CreateOrderRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON request", "INVALID_REQUEST")
		return
	}

	// Reject trailing JSON values.
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "Invalid JSON request", "INVALID_REQUEST")
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")

	result, err := h.service.CreateOrder(
		r.Context(),
		req,
		idempotencyKey,
	)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrCustomerRequired):
			writeError(w, http.StatusBadRequest, err.Error(), "CUSTOMER_REQUIRED")

		case errors.Is(err, service.ErrItemsRequired):
			writeError(w, http.StatusBadRequest, err.Error(), "ITEMS_REQUIRED")

		case errors.Is(err, service.ErrSKURequired):
			writeError(w, http.StatusBadRequest, err.Error(), "SKU_REQUIRED")

		case errors.Is(err, service.ErrInvalidQuantity):
			writeError(w, http.StatusBadRequest, err.Error(), "INVALID_QUANTITY")

		case errors.Is(err, service.ErrDuplicateSKU):
			writeError(w, http.StatusBadRequest, err.Error(), "DUPLICATE_SKU")

		case errors.Is(err, service.ErrIdempotencyKeyRequired):
			writeError(w, http.StatusBadRequest, err.Error(), "IDEMPOTENCY_KEY_REQUIRED")

		case errors.Is(err, repository.ErrIdempotencyKeyConflict):
			writeError(w, http.StatusConflict, err.Error(), "IDEMPOTENCY_KEY_CONFLICT")

		default:
			var stockErr *repository.InsufficientOrderStockError

			if errors.As(err, &stockErr) {
				writeOrderError(
					w,
					http.StatusConflict,
					"Insufficient stock",
					OrderErrorDetail{
						Code:      "INSUFFICIENT_STOCK",
						SKU:       stockErr.SKU,
						Available: stockErr.Available,
					},
				)
				return
			}

			var skuErr *repository.UnknownSKUError

			if errors.As(err, &skuErr) {
				writeError(
					w,
					http.StatusBadRequest,
					"Unknown SKU: "+skuErr.SKU,
					"UNKNOWN_SKU",
				)
				return
			}

			writeError(
				w,
				http.StatusInternalServerError,
				"Internal server error",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	writeResponse(
		w,
		http.StatusCreated,
		"Order created successfully",
		result,
	)
}

func writeOrderError(
	w http.ResponseWriter,
	status int,
	message string,
	detail OrderErrorDetail,
) {
	writeJSON(w, status, OrderErrorResponse{
		CodeHttp: status,
		Message:  message,
		Error:    detail,
	})
}


func (h *OrderHandler) CancelOrder(
    w http.ResponseWriter,
    r *http.Request,
) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

    if err != nil || id <= 0 {
        writeError(
            w,
            http.StatusBadRequest,
            "Invalid order ID",
            "INVALID_ORDER_ID",
        )
        return
    }

    err = h.service.CancelOrder(r.Context(), id)

    if err != nil {
        switch {
        case errors.Is(err, service.ErrOrderNotFound):
            writeError(
                w,
                http.StatusNotFound,
                "Order not found",
                "ORDER_NOT_FOUND",
            )

        case errors.Is(err, service.ErrOrderAlreadyCancelled):
            writeError(
                w,
                http.StatusConflict,
                "Order has already been cancelled",
                "ORDER_ALREADY_CANCELLED",
            )

        case errors.Is(err, service.ErrOrderCannotBeCancelled):
            writeError(
                w,
                http.StatusConflict,
                "Order cannot be cancelled",
                "ORDER_CANNOT_BE_CANCELLED",
            )

        default:
            writeError(
                w,
                http.StatusInternalServerError,
                "Internal server error",
                "INTERNAL_SERVER_ERROR",
            )
        }

        return
    }

    writeResponse(
        w,
        http.StatusOK,
        "Order cancelled successfully",
        map[string]any{
            "id":     id,
            "status": "CANCELLED",
        },
    )
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
    if err != nil || id <= 0 {
        writeError(
            w,
            http.StatusBadRequest,
            "Invalid order ID",
            "INVALID_ORDER_ID",
        )
        return
    }

    result, err := h.service.GetOrderByID(r.Context(), id)

    if errors.Is(err, service.ErrOrderNotFound) {
        writeError(
            w,
            http.StatusNotFound,
            "Order not found",
            "ORDER_NOT_FOUND",
        )
        return
    }

    if err != nil {
        writeError(
            w,
            http.StatusInternalServerError,
            "Failed to retrieve order",
            "INTERNAL_SERVER_ERROR",
        )
        return
    }

    writeResponse(
        w,
        http.StatusOK,
        "Order retrieved successfully",
        result,
    )
}

func (h *OrderHandler) GetOrders(
    w http.ResponseWriter,
    r *http.Request,
) {
    query := r.URL.Query()

    page := 1
    limit := 10

    if value := query.Get("page"); value != "" {
        parsed, err := strconv.Atoi(value)
        if err != nil || parsed < 1 {
            writeError(
                w,
                http.StatusBadRequest,
                "Invalid page",
                "INVALID_PAGE",
            )
            return
        }
        page = parsed
    }

    if value := query.Get("limit"); value != "" {
        parsed, err := strconv.Atoi(value)
        if err != nil || parsed < 1 || parsed > 100 {
            writeError(
                w,
                http.StatusBadRequest,
                "Limit must be between 1 and 100",
                "INVALID_LIMIT",
            )
            return
        }
        limit = parsed
    }

    result, err := h.service.GetOrders(r.Context(), page, limit)
    if err != nil {
        writeError(
            w,
            http.StatusInternalServerError,
            "Failed to retrieve orders",
            "INTERNAL_SERVER_ERROR",
        )
        return
    }

	writeListResponse(
		w,
		"Orders retrieved successfully",
		result.Data,
		result.Page,
		result.Limit,
		result.Total,
	)
}