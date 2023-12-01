CREATE TABLE IF NOT EXISTS skeleton_plans (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    description VARCHAR(200),
    status ENUM('enabled', 'disabled') DEFAULT 'enabled',
    user_id INT NOT NULL,
    CONSTRAINT fk_skeleton_plans_user_id FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
)