-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create movies table
CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration INTEGER NOT NULL, -- in minutes
    genre VARCHAR(100) NOT NULL,
    release_date DATE NOT NULL,
    poster_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create halls table
CREATE TABLE halls (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    capacity INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create show_times table
CREATE TABLE show_times (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER REFERENCES movies(id) ON DELETE CASCADE,
    hall_id INTEGER REFERENCES halls(id) ON DELETE CASCADE,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create seats table
CREATE TABLE seats (
    id SERIAL PRIMARY KEY,
    hall_id INTEGER REFERENCES halls(id) ON DELETE CASCADE,
    show_time_id INTEGER REFERENCES show_times(id) ON DELETE CASCADE,
    row_number VARCHAR(2) NOT NULL,
    seat_number INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(hall_id, row_number, seat_number)
);

-- Create tickets table
CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    show_time_id INTEGER REFERENCES show_times(id) ON DELETE CASCADE,
    seat_id INTEGER REFERENCES seats(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'reserved',
    booking_code VARCHAR(50) UNIQUE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(show_time_id, seat_id)
);

-- Create bookings table
CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    show_time_id INTEGER REFERENCES show_times(id) ON DELETE CASCADE,
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    booked_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create booking_seats table
CREATE TABLE booking_seats (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER REFERENCES bookings(id) ON DELETE CASCADE,
    seat_id INTEGER REFERENCES seats(id) ON DELETE CASCADE,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_movies_release_date ON movies(release_date);
CREATE INDEX idx_show_times_start_time ON show_times(start_time);
CREATE INDEX idx_tickets_booking_code ON tickets(booking_code);
CREATE INDEX idx_tickets_user_id ON tickets(user_id);

-- Insert sample movie data
INSERT INTO movies (title, description, duration, genre, release_date, poster_url) VALUES
('The Shawshank Redemption', 'Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.', 142, 'Drama', '1994-09-23', 'https://example.com/shawshank.jpg'),
('The Godfather', 'The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.', 175, 'Crime', '1972-03-24', 'https://example.com/godfather.jpg'),
('Inception', 'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.', 148, 'Sci-Fi', '2010-07-16', 'https://example.com/inception.jpg');

-- Insert sample halls
INSERT INTO halls (name, capacity) VALUES
('Hall A', 100),
('Hall B', 150),
('Hall C', 200);

-- Insert sample seats for Hall A
INSERT INTO seats (hall_id, row_number, seat_number, status) 
SELECT 1, chr(64 + row_num), seat_num, 'AVAILABLE'
FROM generate_series(1, 10) row_num,
     generate_series(1, 10) seat_num;

-- Insert sample seats for Hall B
INSERT INTO seats (hall_id, row_number, seat_number, status)
SELECT 2, chr(64 + row_num), seat_num, 'AVAILABLE'
FROM generate_series(1, 10) row_num,
     generate_series(1, 15) seat_num;

-- Insert sample seats for Hall C
INSERT INTO seats (hall_id, row_number, seat_number, status)
SELECT 3, chr(64 + row_num), seat_num, 'AVAILABLE'
FROM generate_series(1, 10) row_num,
     generate_series(1, 20) seat_num;

-- Insert sample showtimes
INSERT INTO show_times (movie_id, hall_id, start_time, end_time, price) VALUES
(1, 1, CURRENT_TIMESTAMP + interval '1 day', CURRENT_TIMESTAMP + interval '1 day 3 hours', 15.00),
(2, 2, CURRENT_TIMESTAMP + interval '1 day', CURRENT_TIMESTAMP + interval '1 day 3 hours', 18.00),
(3, 3, CURRENT_TIMESTAMP + interval '1 day', CURRENT_TIMESTAMP + interval '1 day 3 hours', 20.00);

-- Insert sample user
INSERT INTO users (email, password, name, phone) VALUES
('user@example.com', '$2a$10$3QxDjD1ylgPnRgQLhBrTaeGxZHhN9.DpBTPv6L9pB.HxQv2RHllMa', 'Test User', '1234567890');

-- Add triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_halls_updated_at
    BEFORE UPDATE ON halls
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_movies_updated_at
    BEFORE UPDATE ON movies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_show_times_updated_at
    BEFORE UPDATE ON show_times
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_seats_updated_at
    BEFORE UPDATE ON seats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tickets_updated_at
    BEFORE UPDATE ON tickets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookings_updated_at
    BEFORE UPDATE ON bookings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_booking_seats_updated_at
    BEFORE UPDATE ON booking_seats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();