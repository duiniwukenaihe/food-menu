-- Seed default users
INSERT INTO users (username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
VALUES 
    ('admin', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.tJxOHKZpW89TJUJGg7Vj1lFvWZ2jVe', 'Admin', 'User', 'admin', true, NOW(), NOW()),
    ('chef_demo', 'chef@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.tJxOHKZpW89TJUJGg7Vj1lFvWZ2jVe', 'Demo', 'Chef', 'user', true, NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- Seed dish categories
INSERT INTO dish_categories (name, description, icon, sort_order, is_active, created_at, updated_at)
VALUES
    ('Appetizers', 'Start your meal with these delicious appetizers', '🥗', 1, true, NOW(), NOW()),
    ('Main Course', 'Hearty and satisfying main dishes', '🍽️', 2, true, NOW(), NOW()),
    ('Desserts', 'Sweet treats to end your meal', '🍰', 3, true, NOW(), NOW()),
    ('Beverages', 'Refreshing drinks and cocktails', '🥤', 4, true, NOW(), NOW()),
    ('Breakfast', 'Start your day with a great breakfast', '🍳', 5, true, NOW(), NOW()),
    ('Salads', 'Fresh and healthy salad options', '🥗', 6, true, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- Seed content categories
INSERT INTO categories (name, description, color, icon, sort_order, is_active, created_at, updated_at)
VALUES
    ('Recipes', 'Detailed cooking recipes and guides', '#4CAF50', '📖', 1, true, NOW(), NOW()),
    ('Cooking Tips', 'Expert tips and techniques', '#2196F3', '💡', 2, true, NOW(), NOW()),
    ('Nutrition', 'Nutritional information and healthy eating', '#FF9800', '🥗', 3, true, NOW(), NOW()),
    ('Kitchen Tools', 'Reviews and guides for kitchen equipment', '#9C27B0', '🔪', 4, true, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- Seed sample dishes
INSERT INTO dishes (
    name, description, category_id, tags, is_active, is_seasonal, 
    ingredients, cooking_steps, prep_time, cook_time, servings, difficulty,
    image_url, created_by, created_at, updated_at
)
SELECT
    'Classic Margherita Pizza',
    'A traditional Italian pizza with fresh tomatoes, mozzarella, and basil',
    dc.id,
    'italian,vegetarian,cheese,pizza',
    true,
    false,
    '[
        {"name": "Pizza dough", "amount": "500g"},
        {"name": "San Marzano tomatoes", "amount": "400g"},
        {"name": "Fresh mozzarella", "amount": "250g"},
        {"name": "Fresh basil leaves", "amount": "10-12 leaves"},
        {"name": "Extra virgin olive oil", "amount": "2 tbsp"},
        {"name": "Salt", "amount": "to taste"}
    ]',
    'Step 1: Preheat oven to 475°F (245°C) with pizza stone inside.\nStep 2: Roll out pizza dough to desired thickness.\nStep 3: Spread crushed tomatoes evenly over dough.\nStep 4: Tear mozzarella and distribute over sauce.\nStep 5: Drizzle with olive oil and season with salt.\nStep 6: Bake for 10-12 minutes until crust is golden.\nStep 7: Add fresh basil leaves and serve immediately.',
    15,
    12,
    2,
    'medium',
    'https://images.unsplash.com/photo-1574071318508-1cdbab80d002',
    1,
    NOW(),
    NOW()
FROM dish_categories dc WHERE dc.name = 'Main Course'
ON CONFLICT DO NOTHING;

INSERT INTO dishes (
    name, description, category_id, tags, is_active, is_seasonal,
    available_months, seasonal_note,
    ingredients, cooking_steps, prep_time, cook_time, servings, difficulty,
    image_url, created_by, created_at, updated_at
)
SELECT
    'Summer Berry Salad',
    'A refreshing salad with mixed berries, arugula, and balsamic glaze',
    dc.id,
    'salad,summer,vegetarian,healthy,berries',
    true,
    true,
    '5,6,7,8',
    'Best enjoyed during summer months when berries are in season',
    '[
        {"name": "Mixed arugula", "amount": "4 cups"},
        {"name": "Strawberries", "amount": "1 cup, sliced"},
        {"name": "Blueberries", "amount": "1 cup"},
        {"name": "Raspberries", "amount": "1/2 cup"},
        {"name": "Goat cheese", "amount": "100g, crumbled"},
        {"name": "Candied pecans", "amount": "1/2 cup"},
        {"name": "Balsamic glaze", "amount": "3 tbsp"}
    ]',
    'Step 1: Wash and dry all berries thoroughly.\nStep 2: Place arugula in a large bowl.\nStep 3: Add strawberries, blueberries, and raspberries.\nStep 4: Sprinkle crumbled goat cheese over the top.\nStep 5: Add candied pecans.\nStep 6: Drizzle with balsamic glaze just before serving.\nStep 7: Toss gently and serve immediately.',
    10,
    0,
    4,
    'easy',
    'https://images.unsplash.com/photo-1540189549336-e6e99c3679fe',
    1,
    NOW(),
    NOW()
FROM dish_categories dc WHERE dc.name = 'Salads'
ON CONFLICT DO NOTHING;

INSERT INTO dishes (
    name, description, category_id, tags, is_active, is_seasonal,
    ingredients, cooking_steps, prep_time, cook_time, servings, difficulty,
    image_url, created_by, created_at, updated_at
)
SELECT
    'Chocolate Lava Cake',
    'Decadent chocolate cake with a molten center',
    dc.id,
    'dessert,chocolate,french,romantic',
    true,
    false,
    '[
        {"name": "Dark chocolate", "amount": "200g"},
        {"name": "Butter", "amount": "100g"},
        {"name": "Eggs", "amount": "4 large"},
        {"name": "Sugar", "amount": "100g"},
        {"name": "All-purpose flour", "amount": "50g"},
        {"name": "Vanilla extract", "amount": "1 tsp"},
        {"name": "Butter for ramekins", "amount": "as needed"},
        {"name": "Cocoa powder for dusting", "amount": "as needed"}
    ]',
    'Step 1: Preheat oven to 425°F (220°C).\nStep 2: Butter and dust 4 ramekins with cocoa powder.\nStep 3: Melt chocolate and butter together in a double boiler.\nStep 4: Whisk eggs and sugar until light and fluffy.\nStep 5: Fold melted chocolate into egg mixture.\nStep 6: Gently fold in flour and vanilla.\nStep 7: Divide batter among prepared ramekins.\nStep 8: Bake for 12-14 minutes until edges are set but center jiggles.\nStep 9: Let cool for 1 minute, then invert onto plates.\nStep 10: Serve immediately with vanilla ice cream or whipped cream.',
    20,
    14,
    4,
    'medium',
    'https://images.unsplash.com/photo-1624353365286-3f8d62daad51',
    1,
    NOW(),
    NOW()
FROM dish_categories dc WHERE dc.name = 'Desserts'
ON CONFLICT DO NOTHING;

INSERT INTO dishes (
    name, description, category_id, tags, is_active, is_seasonal,
    ingredients, cooking_steps, prep_time, cook_time, servings, difficulty,
    image_url, video_url, created_by, created_at, updated_at
)
SELECT
    'Classic Bruschetta',
    'Toasted bread topped with fresh tomatoes, garlic, and basil',
    dc.id,
    'italian,appetizer,vegetarian,quick',
    true,
    false,
    '[
        {"name": "Baguette", "amount": "1 loaf"},
        {"name": "Tomatoes", "amount": "4 large, diced"},
        {"name": "Garlic cloves", "amount": "3, minced"},
        {"name": "Fresh basil", "amount": "1/4 cup, chopped"},
        {"name": "Extra virgin olive oil", "amount": "4 tbsp"},
        {"name": "Balsamic vinegar", "amount": "1 tbsp"},
        {"name": "Salt and pepper", "amount": "to taste"}
    ]',
    'Step 1: Slice baguette into 1/2 inch thick slices.\nStep 2: Toast bread slices until golden brown.\nStep 3: Mix diced tomatoes, garlic, basil, olive oil, and balsamic vinegar in a bowl.\nStep 4: Season with salt and pepper.\nStep 5: Rub toasted bread with a cut garlic clove.\nStep 6: Top each slice with tomato mixture.\nStep 7: Drizzle with additional olive oil if desired.\nStep 8: Serve immediately.',
    15,
    5,
    6,
    'easy',
    'https://images.unsplash.com/photo-1572695157366-5e585ab2b69f',
    'https://www.youtube.com/watch?v=example',
    1,
    NOW(),
    NOW()
FROM dish_categories dc WHERE dc.name = 'Appetizers'
ON CONFLICT DO NOTHING;

INSERT INTO dishes (
    name, description, category_id, tags, is_active, is_seasonal,
    ingredients, cooking_steps, prep_time, cook_time, servings, difficulty,
    image_url, created_by, created_at, updated_at
)
SELECT
    'Fluffy Pancakes',
    'Light and fluffy American-style pancakes',
    dc.id,
    'breakfast,american,sweet,family-friendly',
    true,
    false,
    '[
        {"name": "All-purpose flour", "amount": "2 cups"},
        {"name": "Sugar", "amount": "2 tbsp"},
        {"name": "Baking powder", "amount": "2 tsp"},
        {"name": "Salt", "amount": "1/2 tsp"},
        {"name": "Milk", "amount": "1 3/4 cups"},
        {"name": "Eggs", "amount": "2 large"},
        {"name": "Butter, melted", "amount": "1/4 cup"},
        {"name": "Vanilla extract", "amount": "1 tsp"}
    ]',
    'Step 1: Mix flour, sugar, baking powder, and salt in a large bowl.\nStep 2: Whisk together milk, eggs, melted butter, and vanilla in another bowl.\nStep 3: Pour wet ingredients into dry ingredients and mix until just combined.\nStep 4: Let batter rest for 5 minutes.\nStep 5: Heat a griddle or non-stick pan over medium heat.\nStep 6: Pour 1/4 cup batter for each pancake.\nStep 7: Cook until bubbles form on surface, about 2-3 minutes.\nStep 8: Flip and cook until golden brown, another 2 minutes.\nStep 9: Serve hot with maple syrup and butter.',
    10,
    15,
    4,
    'easy',
    'https://images.unsplash.com/photo-1567620905732-2d1ec7ab7445',
    1,
    NOW(),
    NOW()
FROM dish_categories dc WHERE dc.name = 'Breakfast'
ON CONFLICT DO NOTHING;

-- Seed dish pairings
INSERT INTO dish_pairings (dish_id, paired_dish_id, pairing_type, score, created_at)
SELECT 
    d1.id,
    d2.id,
    'complements',
    0.95,
    NOW()
FROM dishes d1
CROSS JOIN dishes d2
WHERE d1.name = 'Classic Bruschetta' AND d2.name = 'Classic Margherita Pizza'
ON CONFLICT DO NOTHING;

INSERT INTO dish_pairings (dish_id, paired_dish_id, pairing_type, score, created_at)
SELECT 
    d1.id,
    d2.id,
    'complements',
    0.90,
    NOW()
FROM dishes d1
CROSS JOIN dishes d2
WHERE d1.name = 'Chocolate Lava Cake' AND d2.name = 'Classic Margherita Pizza'
ON CONFLICT DO NOTHING;

INSERT INTO dish_pairings (dish_id, paired_dish_id, pairing_type, score, created_at)
SELECT 
    d1.id,
    d2.id,
    'complements',
    0.85,
    NOW()
FROM dishes d1
CROSS JOIN dishes d2
WHERE d1.name = 'Summer Berry Salad' AND d2.name = 'Fluffy Pancakes'
ON CONFLICT DO NOTHING;
