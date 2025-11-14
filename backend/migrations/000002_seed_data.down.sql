DELETE FROM dish_pairings WHERE (dish_id, paired_dish_id) IN (
    SELECT d1.id, d2.id
    FROM dishes d1
    JOIN dishes d2 ON 1=1
    WHERE d1.name IN ('Classic Bruschetta', 'Chocolate Lava Cake', 'Summer Berry Salad')
);

DELETE FROM dishes WHERE name IN (
    'Classic Margherita Pizza',
    'Summer Berry Salad',
    'Chocolate Lava Cake',
    'Classic Bruschetta',
    'Fluffy Pancakes'
);

DELETE FROM categories WHERE name IN (
    'Recipes', 'Cooking Tips', 'Nutrition', 'Kitchen Tools'
);

DELETE FROM dish_categories WHERE name IN (
    'Appetizers', 'Main Course', 'Desserts', 'Beverages', 'Breakfast', 'Salads'
);

DELETE FROM users WHERE username IN ('admin', 'chef_demo');
