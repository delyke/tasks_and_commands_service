-- Tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    team_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status ENUM('todo', 'in_progress', 'done') NOT NULL DEFAULT 'todo',
    priority ENUM('low', 'medium', 'high') NOT NULL DEFAULT 'medium',
    assignee_id BIGINT UNSIGNED,
    created_by BIGINT UNSIGNED NOT NULL,
    due_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    INDEX idx_tasks_team_id (team_id),
    INDEX idx_tasks_assignee_id (assignee_id),
    INDEX idx_tasks_created_by (created_by),
    INDEX idx_tasks_status (status),
    INDEX idx_tasks_due_date (due_date),
    -- Composite index for common query patterns
    INDEX idx_tasks_team_status (team_id, status),
    INDEX idx_tasks_team_assignee (team_id, assignee_id),
    INDEX idx_tasks_team_status_created (team_id, status, created_at),

    CONSTRAINT fk_tasks_team_id
        FOREIGN KEY (team_id) REFERENCES teams(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_tasks_assignee_id
        FOREIGN KEY (assignee_id) REFERENCES users(id)
        ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_tasks_created_by
        FOREIGN KEY (created_by) REFERENCES users(id)
        ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
