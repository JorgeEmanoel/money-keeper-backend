CREATE TABLE IF NOT EXISTS skeletons (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(200) NOT NULL,
    direction ENUM('income', 'outcome') NOT NULL,
    frequency ENUM('monthly', 'random'),
    value INT NOT NULL,
    currency ENUM('BRL', 'USD') DEFAULT 'BRL',
    plan_id INT NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_skeletons_plan_id FOREIGN KEY (plan_id) REFERENCES plans(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    CONSTRAINT fk_skeletons_user_id FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
)
