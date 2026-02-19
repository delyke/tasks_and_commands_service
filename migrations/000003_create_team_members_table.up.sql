-- Team members table (many-to-many with role)
CREATE TABLE IF NOT EXISTS team_members (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    team_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    role ENUM('owner', 'admin', 'member') NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE INDEX idx_team_members_team_user (team_id, user_id),
    INDEX idx_team_members_user_id (user_id),
    INDEX idx_team_members_role (role),

    CONSTRAINT fk_team_members_team_id
        FOREIGN KEY (team_id) REFERENCES teams(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_team_members_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
