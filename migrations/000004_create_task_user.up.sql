CREATE TABLE task_users (
                            id SERIAL PRIMARY KEY,
                            task_id INT NOT NULL,
                            user_id INT NOT NULL,
                            FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
                            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                            UNIQUE (task_id, user_id)
);
