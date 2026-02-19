-- Task comments table
CREATE TABLE IF NOT EXISTS task_comments (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    task_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    INDEX idx_task_comments_task_id (task_id),
    INDEX idx_task_comments_user_id (user_id),
    INDEX idx_task_comments_created_at (created_at),

    CONSTRAINT fk_task_comments_task_id
        FOREIGN KEY (task_id) REFERENCES tasks(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_task_comments_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
