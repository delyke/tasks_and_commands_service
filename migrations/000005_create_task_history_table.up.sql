-- Task history table (audit log)
CREATE TABLE IF NOT EXISTS task_history (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    task_id BIGINT UNSIGNED NOT NULL,
    changed_by BIGINT UNSIGNED NOT NULL,
    field_name VARCHAR(50) NOT NULL,
    old_value TEXT,
    new_value TEXT,
    changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    INDEX idx_task_history_task_id (task_id),
    INDEX idx_task_history_changed_by (changed_by),
    INDEX idx_task_history_changed_at (changed_at),
    INDEX idx_task_history_task_changed (task_id, changed_at),

    CONSTRAINT fk_task_history_task_id
        FOREIGN KEY (task_id) REFERENCES tasks(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_task_history_changed_by
        FOREIGN KEY (changed_by) REFERENCES users(id)
        ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
