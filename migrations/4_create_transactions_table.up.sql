CREATE TABLE IF NOT EXISTS transactions (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(200) NOT NULL,
    type ENUM('income', 'outcome') NOT NULL,
    value INT NOT NULL,
    currency ENUM('BRL', 'USD') DEFAULT 'BRL',
    reference CHAR(7),
    status ENUM('pending', 'paid', 'canceled'),
    user_id INT NOT NULL,
    CONSTRAINT fk_transactions_user_id FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);
