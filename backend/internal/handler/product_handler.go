package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"backend/internal/dto"
	"backend/internal/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) GetProducts(
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

    req := dto.ProductListRequest{
        Search: query.Get("search"),
        Page:   page,
        Limit:  limit,
    }

    result, err := h.service.GetProducts(r.Context(), req)
    if err != nil {
        writeError(
            w,
            http.StatusInternalServerError,
            "Failed to retrieve products",
            "INTERNAL_SERVER_ERROR",
        )
        return
    }

    writeListResponse(
        w,
        "Products retrieved successfully",
        result.Data,
        result.Page,
        result.Limit,
        result.Total,
    )
}


func (h *ProductHandler) CreateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.CreateProductRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"Invalid Request Body",
			"INVALID_REQUEST_BODY",
		)
		return
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(
			w,
			http.StatusBadRequest,
			"Request body must contain only one JSON object",
			"INVALID_REQUEST_BODY",
		)
		return
	}

	result, err := h.service.CreateProduct(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrSKURequired),
			errors.Is(err, service.ErrNameRequired),
			errors.Is(err, service.ErrPriceRequired),
			errors.Is(err, service.ErrInvalidPrice):

			writeError(
				w,
				http.StatusBadRequest,
				err.Error(),
				"INVALID_REQUEST_BODY",
			)
			return

		case errors.Is(err, service.ErrSKUExists):

			writeError(
				w,
				http.StatusConflict,
				err.Error(),
				"SKU_EXISTS",
			)
			return

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"failed to create product",
				"INTERNAL_SERVER_ERROR",
			)
			return
		}
	}

	writeResponse(w, http.StatusCreated, "Product created successfully", result)
}

func (h *ProductHandler) UpdateProduct(
    w http.ResponseWriter,
    r *http.Request,
) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

    if err != nil || id <= 0 {
        writeError(
            w,
            http.StatusBadRequest,
            "invalid product id",
            "INVALID_PRODUCT_ID",
        )
        return
    }

    var req dto.UpdateProductRequest

    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&req); err != nil {
        writeError(
            w,
            http.StatusBadRequest,
            "Invalid Request Body",
            "INVALID_REQUEST_BODY",
        )
        return
    }

    result, err := h.service.UpdateProduct(
        r.Context(),
        id,
        req,
    )

    if err != nil {
        switch {
        case errors.Is(err, service.ErrProductNotFound):
            writeError(
                w,
                http.StatusNotFound,
                err.Error(),
                "PRODUCT_NOT_FOUND",
            )

        case errors.Is(err, service.ErrSKURequired),
            errors.Is(err, service.ErrNameRequired),
            errors.Is(err, service.ErrPriceRequired),
            errors.Is(err, service.ErrInvalidPrice):
            writeError(
                w,
                http.StatusBadRequest,
                err.Error(),
                "INVALID_REQUEST_BODY",
            )

        default:
            writeError(
                w,
                http.StatusInternalServerError,
                "failed to update product",
                "INTERNAL_SERVER_ERROR",
            )
        }

        return
    }

    writeResponse(w, http.StatusOK, "Product updated successfully", result)
}

func (h *ProductHandler) DeleteProduct(
    w http.ResponseWriter,
    r *http.Request,
) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

    if err != nil || id <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid product id",
			"INVALID_PRODUCT_ID",
		)
        return
    }

    err = h.service.DeleteProduct(r.Context(), id)

    if err != nil {
        switch {
        case errors.Is(err, service.ErrProductNotFound):
            writeError(
                w,
                http.StatusNotFound,
                err.Error(),
                "PRODUCT_NOT_FOUND",
            )

        case errors.Is(err, service.ErrProductHasHistory):
            writeError(
                w,
                http.StatusConflict,
                err.Error(),
                "PRODUCT_HAS_HISTORY",
            )

        default:
            writeError(
                w,
                http.StatusInternalServerError,
                "failed to delete product",
                "INTERNAL_SERVER_ERROR",
            )
        }
        return
    }

    w.WriteHeader(http.StatusNoContent)
}