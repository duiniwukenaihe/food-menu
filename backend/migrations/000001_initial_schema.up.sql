-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(50) DEFAULT 'user' NOT NULL,
    avatar VARCHAR(500),
    is_active BOOLEAN DEFAULT true NOT NULL,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    color VARCHAR(50),
    icon VARCHAR(100),
    is_active BOOLEAN DEFAULT true NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_by INTEGER REFERENCES users(id),
    updated_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_categories_is_active ON categories(is_active);
CREATE INDEX idx_categories_created_by ON categories(created_by);
CREATE INDEX idx_categories_updated_by ON categories(updated_by);
CREATE INDEX idx_categories_created_at ON categories(created_at);
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

-- Create contents table
CREATE TABLE IF NOT EXISTS contents (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    body TEXT,
    category_id INTEGER NOT NULL REFERENCES categories(id),
    author_id INTEGER NOT NULL REFERENCES users(id),
    tags VARCHAR(500),
    image_url VARCHAR(500),
    is_published BOOLEAN DEFAULT false NOT NULL,
    view_count INTEGER DEFAULT 0,
    updated_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_contents_title ON contents(title);
CREATE INDEX idx_contents_category_id ON contents(category_id);
CREATE INDEX idx_contents_author_id ON contents(author_id);
CREATE INDEX idx_contents_is_published ON contents(is_published);
CREATE INDEX idx_contents_view_count ON contents(view_count);
CREATE INDEX idx_contents_updated_by ON contents(updated_by);
CREATE INDEX idx_contents_created_at ON contents(created_at);
CREATE INDEX idx_contents_deleted_at ON contents(deleted_at);

-- Create dish_categories table
CREATE TABLE IF NOT EXISTS dish_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_dish_categories_name ON dish_categories(name);
CREATE INDEX idx_dish_categories_deleted_at ON dish_categories(deleted_at);

-- Create dishes table
CREATE TABLE IF NOT EXISTS dishes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES dish_categories(id),
    tags VARCHAR(500),
    is_active BOOLEAN DEFAULT true NOT NULL,
    is_seasonal BOOLEAN DEFAULT false NOT NULL,
    available_months VARCHAR(100),
    seasonal_note VARCHAR(500),
    ingredients TEXT,
    cooking_steps TEXT,
    prep_time INTEGER,
    cook_time INTEGER,
    servings INTEGER DEFAULT 1,
    difficulty VARCHAR(50),
    image_url VARCHAR(500),
    thumbnail_url VARCHAR(500),
    gallery_urls TEXT,
    video_url VARCHAR(500),
    created_by INTEGER REFERENCES users(id),
    updated_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_dishes_name ON dishes(name);
CREATE INDEX idx_dishes_category_id ON dishes(category_id);
CREATE INDEX idx_dishes_is_active ON dishes(is_active);
CREATE INDEX idx_dishes_is_seasonal ON dishes(is_seasonal);
CREATE INDEX idx_dishes_created_by ON dishes(created_by);
CREATE INDEX idx_dishes_updated_by ON dishes(updated_by);
CREATE INDEX idx_dishes_created_at ON dishes(created_at);
CREATE INDEX idx_dishes_deleted_at ON dishes(deleted_at);

-- Create dish_pairings table (pivot table for pairing rules)
CREATE TABLE IF NOT EXISTS dish_pairings (
    id SERIAL PRIMARY KEY,
    dish_id INTEGER NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    paired_dish_id INTEGER NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    pairing_type VARCHAR(50),
    score DECIMAL(3,2) DEFAULT 1.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_dish_pairings_dish_id ON dish_pairings(dish_id);
CREATE INDEX idx_dish_pairings_paired_dish_id ON dish_pairings(paired_dish_id);
CREATE UNIQUE INDEX idx_unique_pairing ON dish_pairings(dish_id, paired_dish_id, pairing_type);

-- Create media table
CREATE TABLE IF NOT EXISTS media (
    id SERIAL PRIMARY KEY,
    key VARCHAR(500) UNIQUE NOT NULL,
    url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    entity_type VARCHAR(50),
    entity_id INTEGER,
    is_primary BOOLEAN DEFAULT false,
    alt_text VARCHAR(255),
    created_by INTEGER REFERENCES users(id),
    updated_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_media_entity_type ON media(entity_type);
CREATE INDEX idx_media_entity_id ON media(entity_id);
CREATE INDEX idx_media_created_by ON media(created_by);
CREATE INDEX idx_media_updated_by ON media(updated_by);

-- Create recommendations table
CREATE TABLE IF NOT EXISTS recommendations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_id INTEGER REFERENCES contents(id) ON DELETE CASCADE,
    dish_id INTEGER REFERENCES dishes(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL,
    score DECIMAL(10,2),
    reason TEXT,
    algorithm VARCHAR(50),
    is_viewed BOOLEAN DEFAULT false NOT NULL,
    viewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_recommendations_user_id ON recommendations(user_id);
CREATE INDEX idx_recommendations_content_id ON recommendations(content_id);
CREATE INDEX idx_recommendations_dish_id ON recommendations(dish_id);
CREATE INDEX idx_recommendations_entity_type ON recommendations(entity_type);
CREATE INDEX idx_recommendations_score ON recommendations(score);
CREATE INDEX idx_recommendations_is_viewed ON recommendations(is_viewed);
CREATE INDEX idx_recommendations_created_at ON recommendations(created_at);
CREATE INDEX idx_recommendations_deleted_at ON recommendations(deleted_at);

-- Add constraint to ensure either content_id or dish_id is set
ALTER TABLE recommendations ADD CONSTRAINT chk_recommendation_entity CHECK (
    (content_id IS NOT NULL AND dish_id IS NULL AND entity_type = 'content') OR
    (content_id IS NULL AND dish_id IS NOT NULL AND entity_type = 'dish')
);
