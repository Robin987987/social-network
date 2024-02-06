CREATE TABLE IF NOT EXISTS group_invitations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_id INT NOT NULL,
    join_user_id INT NOT NULL,
    invite_user_id INT,
    status TEXT NOT NULL CHECK( status IN ('pending', 'accepted', 'declined')) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (join_user_id) REFERENCES users(id)
    FOREIGN KEY (invite_user_id) REFERENCES users(id)
);