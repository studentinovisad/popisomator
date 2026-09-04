CREATE TYPE notification_kind AS ENUM ('item_request', 'item_expiry');

CREATE TABLE notifications (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  recipient_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  kind notification_kind NOT NULL,
  read BOOLEAN NOT NULL DEFAULT false,
  UNIQUE (id, kind)
);

CREATE INDEX idx_notifications_recipient_id ON notifications(recipient_id);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

CREATE TABLE notifdesc_item_request (
  notification_id BIGINT PRIMARY KEY REFERENCES notifications(id) ON DELETE CASCADE,
  kind notification_kind GENERATED ALWAYS AS ('item_request') STORED,
  user_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,

  FOREIGN KEY(notification_id, kind)
    REFERENCES notifications(id, kind) ON DELETE CASCADE,
  FOREIGN KEY(user_id, item_id)
    REFERENCES item_requests(user_id, item_id) ON DELETE CASCADE
);

CREATE TYPE notifdesc_expiry_type AS ENUM ('expiring_soon', 'expired');

CREATE TABLE notifdesc_item_expiry (
  notification_id BIGINT PRIMARY KEY REFERENCES notifications(id) ON DELETE CASCADE,
  kind notification_kind GENERATED ALWAYS AS ('item_expiry') STORED,
  item_id BIGINT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
  expiry_type notifdesc_expiry_type NOT NULL,

  FOREIGN KEY(notification_id, kind)
    REFERENCES notifications(id, kind) ON DELETE CASCADE
);
