package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"backend/internal/dto"
	"backend/internal/service"
)

type StockHandler struct {
	service *service.StockService
}

func NewStockHandler(s *service.StockService) *StockHandler {
	return &StockHandler{
		service: s,
	}
}

func (h *StockHandler) CreateMovement(
	w http.ResponseWriter,
	r *http.Request,
) {
	productID, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil || productID <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"Invalid product ID",
			"INVALID_PRODUCT_ID",
		)
		return
	}

	var req dto.CreateStockMovementRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
			"INVALID_REQUEST_BODY",
		)
		return
	}

	result, err := h.service.CreateMovement(
		r.Context(),
		productID,
		req,
	)

	if err != nil {
		switch {	
		case errors.Is(err, service.ErrInvalidMovementType):
			writeError(
				w,
				http.StatusBadRequest,
				"Invalid movement type. Allowed values are in or out",
				"INVALID_MOVEMENT_TYPE",
			)

		case errors.Is(err, service.ErrInvalidQuantity):
			writeError(
				w,
				http.StatusBadRequest,
				"Quantity must be greater than zero",
				"INVALID_QUANTITY",
			)

		case errors.Is(err, service.ErrProductNotFound):
			writeError(
				w,
				http.StatusNotFound,
				"Product not found",
				"PRODUCT_NOT_FOUND",
			)

		default:
			var insufficientErr *service.InsufficientStockError

			if errors.As(err, &insufficientErr) {
				writeJSON(w, http.StatusConflict, struct {
					CodeHttp int    `json:"codeHttp"`
					Message  string `json:"message"`
					Error    struct {
						Code      string `json:"code"`
						Available int64  `json:"available"`
					} `json:"error"`
				}{
					CodeHttp: http.StatusConflict,
					Message:  "Insufficient stock",
					Error: struct {
						Code      string `json:"code"`
						Available int64  `json:"available"`
					}{
						Code:      "INSUFFICIENT_STOCK",
						Available: insufficientErr.Available,
					},
				})

				return
			}

			writeError(
				w,
				http.StatusInternalServerError,
				"Failed to create stock movement",
				"INTERNAL_SERVER_ERROR",
			)
		}

		return
	}

	writeResponse(
		w,
		http.StatusCreated,
		"Stock movement created successfully",
		result,
	)
}

func (h *StockHandler) GetMovements(
	w http.ResponseWriter,
	r *http.Request,
) {
	productID, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil || productID <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"Invalid product ID",
			"INVALID_PRODUCT_ID",
		)
		return
	}

	query := r.URL.Query()

	page := 1
	limit := 10

	if value := query.Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)

		if err != nil || parsed < 1 {
			writeError(
				w,
				http.StatusBadRequest,
				"Invalid page parameter",
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
				"Invalid limit parameter",
				"INVALID_LIMIT",
			)
			return
		}

		limit = parsed
	}

	req := dto.StockMovementListRequest{
		Page:  page,
		Limit: limit,
	}

	result, err := h.service.GetMovements(
		r.Context(),
		productID,
		req,
	)

	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"Product not found",
				"PRODUCT_NOT_FOUND",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve stock movements",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	writeListResponse(
		w,
		"Stock movements retrieved successfully",
		result.Data,
		result.Page,
		result.Limit,
		result.Total,
	)
}