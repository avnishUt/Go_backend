CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID REFERENCES restaurants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    channel TEXT NOT NULL,
    recipient TEXT NOT NULL,
    subject TEXT,
    body TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    attempts INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT notifications_channel_check CHECK (channel IN ('email', 'sms', 'push', 'in_app')),
    CONSTRAINT notifications_status_check CHECK (status IN ('pending', 'sent', 'failed'))
);

CREATE INDEX notifications_restaurant_id_idx ON notifications(restaurant_id);
CREATE INDEX notifications_status_idx ON notifications(status);
CREATE INDEX notifications_created_at_idx ON notifications(created_at);
