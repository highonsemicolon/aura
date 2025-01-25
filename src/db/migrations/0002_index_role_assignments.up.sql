-- CREATE INDEX idx_user_id (user_id);
-- CREATE INDEX idx_resource_id (resource_id);
CREATE INDEX idx_user_resource ON role_assignments(user_id, resource_id);

