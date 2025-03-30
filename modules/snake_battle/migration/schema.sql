-- Ensure the uuid-ossp extension is enabled
-- You might need to run this manually in your database once: CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table (Passwordless, using UUID)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- Use UUID as primary key
    username VARCHAR(50) NOT NULL,
    elo INTEGER DEFAULT 1000 NOT NULL,   -- Default ELO rating
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    -- Removed email and password_hash
    -- Consider adding a UNIQUE constraint on username if desired, though it might not be strictly necessary for guest users
);

-- Rooms table (Unchanged)
CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    room_code VARCHAR(10) UNIQUE NOT NULL, -- A short, unique code for joining
    max_players INTEGER DEFAULT 2 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Room Members table (linking users (UUID) to rooms)
CREATE TABLE room_members (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Reference UUID
    room_id INTEGER NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_host BOOLEAN DEFAULT FALSE,
    -- Add any other room-specific player state if needed, e.g., snake color
    UNIQUE (user_id, room_id) -- A user can only be in a room once
);

-- Games table (to record match history and ELO changes, referencing UUID users)
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    room_id INTEGER NOT NULL REFERENCES rooms(id),
    winner_id UUID REFERENCES users(id), -- Reference UUID, Can be NULL
    loser_id UUID REFERENCES users(id),  -- Reference UUID, Can be NULL
    elo_change INTEGER,                  -- ELO points transferred
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP WITH TIME ZONE,
    game_data JSONB                      -- Optional: Store final game state or replay data
);

-- Indexes for faster lookups (Updated for UUID user_id)
CREATE INDEX idx_room_members_user_id ON room_members(user_id);
CREATE INDEX idx_room_members_room_id ON room_members(room_id);
CREATE INDEX idx_games_winner_id ON games(winner_id);
CREATE INDEX idx_games_loser_id ON games(loser_id);

-- Optional: Trigger to update user ELO after a game finishes (Updated for UUID)
-- CREATE OR REPLACE FUNCTION update_elo_after_game()
-- RETURNS TRIGGER AS $$
-- BEGIN
--   IF NEW.winner_id IS NOT NULL AND NEW.loser_id IS NOT NULL AND NEW.elo_change IS NOT NULL THEN
--     UPDATE users SET elo = elo + NEW.elo_change WHERE id = NEW.winner_id;
--     UPDATE users SET elo = elo - NEW.elo_change WHERE id = NEW.loser_id;
--   END IF;
--   RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER game_elo_update
-- AFTER INSERT ON games
-- FOR EACH ROW
-- EXECUTE FUNCTION update_elo_after_game();

-- Trigger to update updated_at timestamps (Unchanged function, applied to new users table)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE
ON users FOR EACH ROW EXECUTE PROCEDURE
update_updated_at_column();

CREATE TRIGGER update_rooms_updated_at BEFORE UPDATE
ON rooms FOR EACH ROW EXECUTE PROCEDURE
update_updated_at_column(); 