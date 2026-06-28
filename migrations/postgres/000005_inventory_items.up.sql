CREATE TABLE inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES restaurant_branches(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    unit TEXT NOT NULL DEFAULT 'piece',
    current_stock NUMERIC(14,3) NOT NULL DEFAULT 0,
    minimum_stock NUMERIC(14,3) NOT NULL DEFAULT 0,
    reorder_level NUMERIC(14,3) NOT NULL DEFAULT 0,
    allow_negative_stock BOOLEAN NOT NULL DEFAULT false,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT inventory_items_status_check CHECK (status IN ('active', 'inactive')),
    CONSTRAINT inventory_items_stock_check CHECK (allow_negative_stock OR current_stock >= 0),
    CONSTRAINT inventory_items_scope_unique UNIQUE (restaurant_id, branch_id, product_id)
);

CREATE INDEX inventory_items_restaurant_id_idx ON inventory_items(restaurant_id);
CREATE INDEX inventory_items_branch_id_idx ON inventory_items(branch_id);
CREATE INDEX inventory_items_product_id_idx ON inventory_items(product_id);
