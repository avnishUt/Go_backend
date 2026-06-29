CREATE TABLE inventory_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES restaurant_branches(id) ON DELETE SET NULL,
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    quantity NUMERIC(14,3) NOT NULL,
    stock_before NUMERIC(14,3) NOT NULL,
    stock_after NUMERIC(14,3) NOT NULL,
    reason TEXT,
    reference_type TEXT,
    reference_id UUID,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT inventory_transactions_type_check CHECK (type IN ('stock_in', 'stock_out', 'adjustment', 'wastage', 'return')),
    CONSTRAINT inventory_transactions_quantity_check CHECK (quantity > 0)
);

CREATE INDEX inventory_transactions_restaurant_id_idx ON inventory_transactions(restaurant_id);
CREATE INDEX inventory_transactions_inventory_item_id_idx ON inventory_transactions(inventory_item_id);
CREATE INDEX inventory_transactions_created_at_idx ON inventory_transactions(created_at);
