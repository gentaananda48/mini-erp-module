export type OrderFilterType = "all" | "created" | "cancelled";

interface OrderFiltersProps {
  search: string;
  filter: OrderFilterType;
  onSearchChange: (value: string) => void;
  onFilterChange: (value: OrderFilterType) => void;
}

const filters: {
  label: string;
  value: OrderFilterType;
}[] = [
  { label: "All", value: "all" },
  { label: "Created", value: "created" },
  { label: "Cancelled", value: "cancelled" },
];

export function OrderFilters({ search, filter, onSearchChange, onFilterChange }: OrderFiltersProps) {
  return (
    <div className="product-filters">
      <input type="search" className="search-input" placeholder="Search by order ID or customer..." value={search} onChange={(event) => onSearchChange(event.target.value)} />

      <div className="filter-group">
        {filters.map((item) => (
          <button key={item.value} type="button" className={`filter-button ${filter === item.value ? "selected" : ""}`} onClick={() => onFilterChange(item.value)}>
            {item.label}
          </button>
        ))}
      </div>
    </div>
  );
}
