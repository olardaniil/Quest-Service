
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(20) UNIQUE CHECK (username <> ''),
    balance INTEGER DEFAULT 0 CHECK (balance >= 0)
);

CREATE TABLE quests (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    cost INTEGER DEFAULT 0 CHECK (cost >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE quests_complete (
    user_id INTEGER,
    quest_id INTEGER,
    completed_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (quest_id) REFERENCES quests(id)
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    quest_id INTEGER NOT NULL,
    name VARCHAR(150) NOT NULL,
    is_reusable BOOLEAN DEFAULT FALSE,
    cost INTEGER DEFAULT 0 CHECK ( cost >= 0 ),
    deleted_at TIMESTAMP,
    FOREIGN KEY (quest_id) REFERENCES quests(id)
);

CREATE TABLE tasks_complete (
    user_id INTEGER,
    task_id INTEGER,
    completed_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

