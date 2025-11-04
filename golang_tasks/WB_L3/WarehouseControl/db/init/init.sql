CREATE TABLE IF NOT EXISTS items (
  id SERIAL PRIMARY KEY,
  sku TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
  location TEXT,
  note TEXT,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS items_audit (
  id SERIAL PRIMARY KEY,
  item_id INTEGER,
  operation TEXT NOT NULL, -- INSERT, UPDATE, DELETE
  changed_by TEXT,
  changed_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  old_data JSONB,
  new_data JSONB
);

CREATE OR REPLACE FUNCTION items_audit_trigger_fn() RETURNS trigger AS $$
DECLARE
  actor text;
  old_row jsonb;
  new_row jsonb;
BEGIN
  BEGIN
    actor := current_setting('myapp.current_user', true);
  EXCEPTION WHEN others THEN
    actor := NULL;
  END;

  IF TG_OP = 'INSERT' THEN
    new_row := to_jsonb(NEW);
    INSERT INTO items_audit(item_id, operation, changed_by, changed_at, old_data, new_data)
    VALUES (NEW.id, 'INSERT', actor, now(), NULL, new_row);
    RETURN NEW;
  ELSIF TG_OP = 'UPDATE' THEN
    old_row := to_jsonb(OLD);
    new_row := to_jsonb(NEW);
    INSERT INTO items_audit(item_id, operation, changed_by, changed_at, old_data, new_data)
    VALUES (NEW.id, 'UPDATE', actor, now(), old_row, new_row);
    NEW.updated_at := now();
    RETURN NEW;
  ELSIF TG_OP = 'DELETE' THEN
    old_row := to_jsonb(OLD);
    INSERT INTO items_audit(item_id, operation, changed_by, changed_at, old_data, new_data)
    VALUES (OLD.id, 'DELETE', actor, now(), old_row, NULL);
    RETURN OLD;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_items_audit ON items;
CREATE TRIGGER trg_items_audit
AFTER INSERT OR UPDATE OR DELETE ON items
FOR EACH ROW EXECUTE FUNCTION items_audit_trigger_fn();

CREATE INDEX IF NOT EXISTS idx_items_audit_item_id ON items_audit (item_id);
CREATE INDEX IF NOT EXISTS idx_items_audit_changed_at ON items_audit (changed_at);
