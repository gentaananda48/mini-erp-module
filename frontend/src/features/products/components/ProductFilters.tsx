export type ProductFilterType = "all" | "active" | "inactive" | "low-stock";

interface ProductFiltersProps {
  search: string;
  filter: ProductFilterType;
  onSearchChange: (value: string) => void;
  onFilterChange: (value: ProductFilterType) => void;
}

const filters: {
  label: string;
  value: ProductFilterType;
}[] = [
  { label: "All", value: "all" },
  { label: "Active", value: "active" },
  { label: "Inactive", value: "inactive" },
  { label: "Low Stock", value: "low-stock" },
];

export function ProductFilters({ search, filter, onSearchChange, onFilterChange }: ProductFiltersProps) {
  return (
    <div className="product-filters">
      <input type="search" className="search-input" placeholder="Search SKU or product name..." value={search} onChange={(event) => onSearchChange(event.target.value)} />

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
