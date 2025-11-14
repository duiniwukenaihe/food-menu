# Database Schema Documentation

## Overview

This document describes the database schema for the application, covering users, dishes, media assets, recommendations, and related entities.

## Technology Stack

- **Database**: PostgreSQL 15+
- **ORM**: GORM v1.25.4
- **Migration Tool**: golang-migrate v4
- **Schema Design**: Relational with foreign keys and indexes

## Table Structure

### 1. Users Table

Stores user credentials and profile information.

```sql
CREATE TABLE users (
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
```

**Indexes:**
- `idx_users_username` on `username`
- `idx_users_email` on `email`
- `idx_users_role` on `role`
- `idx_users_is_active` on `is_active`
- `idx_users_created_at` on `created_at`
- `idx_users_deleted_at` on `deleted_at` (soft delete support)

**Fields:**
- `id`: Primary key, auto-incrementing
- `username`: Unique username for login
- `email`: Unique email address
- `password_hash`: Bcrypt hashed password
- `first_name`, `last_name`: User profile information
- `role`: User role (e.g., 'user', 'admin')
- `avatar`: URL to user profile picture
- `is_active`: Account status flag
- `last_login`: Timestamp of last successful login
- `created_at`, `updated_at`: Audit timestamps
- `deleted_at`: Soft delete timestamp (NULL if not deleted)

---

### 2. Categories Table

Content categories for articles and blog posts.

```sql
CREATE TABLE categories (
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
```

**Indexes:**
- `idx_categories_name` on `name`
- `idx_categories_is_active` on `is_active`
- `idx_categories_created_by` on `created_by`
- `idx_categories_updated_by` on `updated_by`

**Relationships:**
- `created_by`: Foreign key to `users.id`
- `updated_by`: Foreign key to `users.id`

---

### 3. Contents Table

Stores articles, blog posts, and other content.

```sql
CREATE TABLE contents (
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
```

**Indexes:**
- `idx_contents_title` on `title`
- `idx_contents_category_id` on `category_id`
- `idx_contents_author_id` on `author_id`
- `idx_contents_is_published` on `is_published`
- `idx_contents_view_count` on `view_count`

**Relationships:**
- `category_id`: Foreign key to `categories.id`
- `author_id`: Foreign key to `users.id`
- `updated_by`: Foreign key to `users.id`

---

### 4. Dish Categories Table

Categories specific to dishes/recipes.

```sql
CREATE TABLE dish_categories (
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
```

**Indexes:**
- `idx_dish_categories_name` on `name`
- `idx_dish_categories_deleted_at` on `deleted_at`

---

### 5. Dishes Table

Main table for recipe/dish information with comprehensive details.

```sql
CREATE TABLE dishes (
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
```

**Indexes:**
- `idx_dishes_name` on `name`
- `idx_dishes_category_id` on `category_id`
- `idx_dishes_is_active` on `is_active`
- `idx_dishes_is_seasonal` on `is_seasonal`
- `idx_dishes_created_by` on `created_by`
- `idx_dishes_updated_by` on `updated_by`
- `idx_dishes_created_at` on `created_at`

**Fields:**
- **Basic Information:**
  - `name`: Dish name
  - `description`: Brief description
  - `category_id`: Link to dish category
  - `tags`: Comma-separated tags for filtering

- **Seasonal Configuration:**
  - `is_seasonal`: Flag for seasonal availability
  - `available_months`: Months when available (e.g., "1,2,3,11,12")
  - `seasonal_note`: Note about seasonality

- **Recipe Details:**
  - `ingredients`: JSON array of ingredients with amounts
  - `cooking_steps`: Detailed cooking instructions
  - `prep_time`: Preparation time in minutes
  - `cook_time`: Cooking time in minutes
  - `servings`: Number of servings
  - `difficulty`: Difficulty level (easy/medium/hard)

- **Media:**
  - `image_url`: Main image URL
  - `thumbnail_url`: Thumbnail image URL
  - `gallery_urls`: JSON array of additional images
  - `video_url`: Cooking video URL

- **Audit Fields:**
  - `created_by`, `updated_by`: User IDs who created/updated
  - `created_at`, `updated_at`: Timestamps
  - `deleted_at`: Soft delete timestamp

**Relationships:**
- `category_id`: Foreign key to `dish_categories.id`
- `created_by`: Foreign key to `users.id`
- `updated_by`: Foreign key to `users.id`

---

### 6. Dish Pairings Table

Pivot table for dish pairing recommendations.

```sql
CREATE TABLE dish_pairings (
    id SERIAL PRIMARY KEY,
    dish_id INTEGER NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    paired_dish_id INTEGER NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    pairing_type VARCHAR(50),
    score DECIMAL(3,2) DEFAULT 1.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
```

**Indexes:**
- `idx_dish_pairings_dish_id` on `dish_id`
- `idx_dish_pairings_paired_dish_id` on `paired_dish_id`
- `idx_unique_pairing` unique index on `(dish_id, paired_dish_id, pairing_type)`

**Fields:**
- `dish_id`: Primary dish ID
- `paired_dish_id`: ID of paired dish
- `pairing_type`: Type of pairing (e.g., 'complements', 'alternative', 'beverage')
- `score`: Strength of pairing (0.00 to 1.00)

**Relationships:**
- `dish_id`: Foreign key to `dishes.id` with CASCADE delete
- `paired_dish_id`: Foreign key to `dishes.id` with CASCADE delete

---

### 7. Media Table

Stores metadata for uploaded media files.

```sql
CREATE TABLE media (
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
```

**Indexes:**
- `idx_media_entity_type` on `entity_type`
- `idx_media_entity_id` on `entity_id`
- `idx_media_created_by` on `created_by`

**Fields:**
- `key`: Unique storage key/path
- `url`: Public URL to access media
- `file_name`: Original filename
- `content_type`: MIME type
- `size`: File size in bytes
- `entity_type`: Type of entity (e.g., 'dish', 'user', 'content')
- `entity_id`: ID of associated entity
- `is_primary`: Flag for primary/featured image
- `alt_text`: Alt text for accessibility
- `created_by`, `updated_by`: Audit user IDs

---

### 8. Recommendations Table

Stores personalized recommendations for users.

```sql
CREATE TABLE recommendations (
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
    deleted_at TIMESTAMP,
    CONSTRAINT chk_recommendation_entity CHECK (
        (content_id IS NOT NULL AND dish_id IS NULL AND entity_type = 'content') OR
        (content_id IS NULL AND dish_id IS NOT NULL AND entity_type = 'dish')
    )
);
```

**Indexes:**
- `idx_recommendations_user_id` on `user_id`
- `idx_recommendations_content_id` on `content_id`
- `idx_recommendations_dish_id` on `dish_id`
- `idx_recommendations_entity_type` on `entity_type`
- `idx_recommendations_score` on `score`
- `idx_recommendations_is_viewed` on `is_viewed`

**Fields:**
- `user_id`: User receiving recommendation
- `content_id`: Recommended content (if applicable)
- `dish_id`: Recommended dish (if applicable)
- `entity_type`: Type of recommendation ('content' or 'dish')
- `score`: Recommendation score/confidence
- `reason`: Human-readable explanation
- `algorithm`: Algorithm used for recommendation
- `is_viewed`: Whether user has viewed recommendation
- `viewed_at`: When recommendation was viewed

**Constraints:**
- `chk_recommendation_entity`: Ensures either content_id or dish_id is set, matching entity_type

**Relationships:**
- `user_id`: Foreign key to `users.id` with CASCADE delete
- `content_id`: Foreign key to `contents.id` with CASCADE delete
- `dish_id`: Foreign key to `dishes.id` with CASCADE delete

---

## Entity Relationships

### ER Diagram (Textual)

```
users (1) ----< (N) contents [author]
users (1) ----< (N) dishes [creator]
users (1) ----< (N) categories [creator]
users (1) ----< (N) media [uploader]
users (1) ----< (N) recommendations [recipient]

categories (1) ----< (N) contents
dish_categories (1) ----< (N) dishes

dishes (1) ----< (N) dish_pairings [primary dish]
dishes (1) ----< (N) dish_pairings [paired dish]
dishes (1) ----< (N) recommendations
contents (1) ----< (N) recommendations

media (N) ----< (1) entity [polymorphic: dish, user, content]
```

---

## Migrations

Migrations are located in `/backend/migrations/` and managed using `golang-migrate`.

### Available Migrations

1. **000001_initial_schema** - Creates all base tables with indexes and constraints
2. **000002_seed_data** - Seeds initial data including:
   - Default admin and demo users
   - Dish categories (Appetizers, Main Course, Desserts, etc.)
   - Content categories (Recipes, Cooking Tips, etc.)
   - Sample dishes with full details
   - Sample dish pairings

### Running Migrations

```bash
# Apply all pending migrations
make db-migrate

# Reset database and re-run migrations
make db-reset
```

---

## Indexes Strategy

Indexes are strategically placed to optimize common query patterns:

1. **Foreign Keys**: All foreign key columns are indexed
2. **Filter Fields**: Columns used in WHERE clauses (is_active, is_seasonal, etc.)
3. **Search Fields**: Text columns used in searches (name, username, email)
4. **Sort Fields**: Columns used in ORDER BY (created_at, view_count, score)
5. **Unique Constraints**: username, email, category names

---

## Data Types & Constraints

### Standard Field Types

- **IDs**: `SERIAL` (auto-incrementing integer)
- **Strings**: 
  - Short fields (names, etc.): `VARCHAR(255)`
  - Medium fields (descriptions): `VARCHAR(500)`
  - Long fields (body, instructions): `TEXT`
- **Booleans**: `BOOLEAN` with defaults
- **Timestamps**: `TIMESTAMP` with `CURRENT_TIMESTAMP` default
- **Numbers**: 
  - Counts/durations: `INTEGER`
  - Scores/ratings: `DECIMAL(10,2)`

### Soft Deletes

Most tables support soft deletion via `deleted_at` timestamp:
- NULL = active record
- Non-NULL = soft deleted record

Tables with soft delete:
- users
- categories
- contents
- dish_categories
- dishes
- recommendations

---

## Repository Pattern

The application uses a repository pattern for database access with the following features:

### Repository Interfaces

Located in `/backend/internal/repository/`, providing:

- **Context-Aware Queries**: All operations support `context.Context`
- **Transaction Support**: Repositories can operate within transactions
- **CRUD Operations**: Create, Read, Update, Delete primitives
- **Specialized Queries**: Domain-specific query methods

### Example Usage

```go
// Basic usage
dishRepo := repository.NewDishRepository(db)
dish, err := dishRepo.FindByID(1)

// With context
ctx := context.Background()
dish, err = dishRepo.WithContext(ctx).FindByID(1)

// Within transaction
err := txManager.WithTransaction(ctx, func(txCtx context.Context) error {
    tx := txCtx.Value("tx").(*gorm.DB)
    dishRepo := dishRepo.WithTx(tx)
    return dishRepo.Create(newDish)
})
```

---

## Query Optimization

### N+1 Prevention

Use `Preload` for related entities:

```go
db.Preload("Category").Preload("CreatedBy").Find(&dishes)
```

### Pagination

All list queries support offset/limit pagination:

```go
dishes, total, err := dishRepo.List(filters, offset, limit)
```

### Full-Text Search

Search uses LOWER() with LIKE for case-insensitive matching across multiple fields.

---

## Seed Data

Default seed data includes:

### Users
- **admin** - Admin user (role: admin)
- **chef_demo** - Demo chef user (role: user)

**Note:** Default password hash is for "password123" - change in production!

### Dish Categories
- Appetizers, Main Course, Desserts, Beverages, Breakfast, Salads

### Sample Dishes
1. **Classic Margherita Pizza** - Italian main course
2. **Summer Berry Salad** - Seasonal salad (May-Aug)
3. **Chocolate Lava Cake** - Dessert
4. **Classic Bruschetta** - Appetizer with video
5. **Fluffy Pancakes** - Breakfast

All dishes include:
- Complete ingredient lists (JSON format)
- Step-by-step cooking instructions
- Prep/cook times and servings
- Difficulty ratings
- Image URLs

### Dish Pairings
- Bruschetta + Margherita Pizza (complements)
- Lava Cake + Margherita Pizza (dessert pairing)
- Berry Salad + Pancakes (breakfast pairing)

---

## Best Practices

1. **Always use transactions** for multi-table operations
2. **Use soft deletes** instead of hard deletes where possible
3. **Index foreign keys** and frequently queried columns
4. **Store JSON data as TEXT** (ingredients, gallery_urls) for flexibility
5. **Use audit fields** (created_by, updated_by) for tracking
6. **Validate constraints** at both database and application layers
7. **Use prepared statements** (handled automatically by GORM)
8. **Implement connection pooling** (configured in database.go)

---

## Schema Evolution

When making schema changes:

1. Create a new migration file in `/backend/migrations/`
2. Follow naming convention: `{version}_{description}.up.sql` and `.down.sql`
3. Test both up and down migrations
4. Update this documentation
5. Run migrations in staging before production

---

## Backup & Recovery

### Backup
```bash
pg_dump -U postgres appdb > backup.sql
```

### Restore
```bash
psql -U postgres appdb < backup.sql
```

---

## Performance Considerations

- **Connection Pool**: Max 25 connections (configurable in `database.go`)
- **Indexes**: All foreign keys and common filter fields indexed
- **Query Optimization**: Use EXPLAIN ANALYZE for slow queries
- **Caching**: Consider caching for frequently accessed, rarely changed data
- **Pagination**: Always use limit/offset for large result sets

---

## Contact & Support

For schema-related questions or issues, please refer to:
- Database Administrator
- Backend Development Team
- Architecture Documentation

Last Updated: 2024
