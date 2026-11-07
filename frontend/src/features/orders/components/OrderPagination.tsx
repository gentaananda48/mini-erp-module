interface OrderPaginationProps {
  page: number;
  limit: number;
  total: number;
  onPageChange: (page: number) => void;
  disabled?: boolean;
}

export function OrderPagination({ page, limit, total, onPageChange, disabled = false }: OrderPaginationProps) {
  const totalPages = Math.ceil(total / limit);

  if (totalPages <= 1) {
    return null;
  }

  const startPage = Math.max(1, Math.min(page - 2, totalPages - 4));

  const endPage = Math.min(totalPages, startPage + 4);

  const pages = Array.from({ length: endPage - startPage + 1 }, (_, index) => startPage + index);

  const startItem = (page - 1) * limit + 1;
  const endItem = Math.min(page * limit, total);

  return (
    <div className="pagination-container">
      <div className="pagination-info">
        Showing {startItem}–{endItem} of {total} orders
      </div>

      <div className="pagination-controls">
        <button type="button" disabled={disabled || page === 1} onClick={() => onPageChange(page - 1)} className="pagination-button">
          Previous
        </button>

        {pages.map((pageNumber) => (
          <button key={pageNumber} type="button" disabled={disabled} onClick={() => onPageChange(pageNumber)} className={`pagination-button ${pageNumber === page ? "current" : ""}`} aria-current={pageNumber === page ? "page" : undefined}>
            {pageNumber}
          </button>
        ))}

        <button type="button" disabled={disabled || page === totalPages} onClick={() => onPageChange(page + 1)} className="pagination-button">
          Next
        </button>
      </div>
    </div>
  );
}
