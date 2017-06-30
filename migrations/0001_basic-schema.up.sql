-- done
create table meal_type(
    id serial primary key,
    type varchar(50) unique
);

create table weekday(
    id serial primary key,
    day varchar(50) unique
);

create table house(
     id serial primary key,
     name varchar(355) not null,
     grocery_day varchar(50) references weekday(day) not null,
     household_number int not null
);

-- done
create table user_info(
    id serial primary key,
    email varchar(355) unique not null,
    password varchar(355) not null,
    username varchar(50) not null
);

create table recipe(
     id serial primary key,
     name varchar(355) not null,
     serves_for int,
     image varchar(50)
);

create table ingredient(
    id serial primary key,
    name varchar(355) not null,
    carb_per_100g float8 not null,
    protein_per_100g float8 not null,
    fat_per_100g float8 not null,
    fiber_per_100g float8 not null,
    calories_per_100g float8 not null
);

create table unit(
     id serial primary key,
     name varchar(50) unique not null
 );

create table step(
    recipe_id int4 references recipe(id) on delete cascade,
    id int not null,
    text text not null,
    primary key(recipe_id, id),
    unique(recipe_id, id)
);

create table step_ingredient(
    recipe_id int4,
    step_id int4,
    ingredient_id int4 references ingredient(id),
    unit_id int4 references unit(id) not null,
    amount float8 not null,
    primary key(recipe_id, step_id, ingredient_id),
    foreign key(recipe_id, step_id) references step(recipe_id,id) on delete cascade
);

create table item_in_storage(
   house_id int4 references house(id) on delete cascade,
   ingredient_id int4 references ingredient(id),
   amount float8 not null,
   unit_id int4 references unit(id) not null,
   primary key(house_id, ingredient_id)
);

create table user_recipe(
    user_id int4 references user_info(id) on delete cascade,
    recipe_id int4 references recipe(id) on delete cascade,
    primary key (user_id, recipe_id)
);

create table house_recipe(
    house_id int4 references house(id) on delete cascade,
    recipe_id int4 references recipe(id) on delete cascade,
    primary key (house_id, recipe_id)
);

create table schedule(
    house_id int4 references house(id) on delete cascade,
    week_id int4 references weekday(id) not null,
    type_id int4 references meal_type(id) not null,
    recipe_id int4 references recipe(id),
    primary key(house_id, week_id, type_id)
);

 create table ownership(
     own_type serial primary key,
     description varchar(50) not null
 );

create table member_of(
    user_id int4 references user_info(id) on delete cascade,
    house_id int4 references house(id) on delete cascade,
    own_type int4 references ownership(own_type) on delete cascade,
    primary key(user_id, house_id, own_type)
);

create table recipe_type(
    recipe_id int4 references recipe(id) on delete cascade,
    type_id int4 references meal_type(id) not null,
    primary key(recipe_id, type_id)
);

insert into user_info(username, email, password)
values 
('guli', 'gulipek5@gmail.com', 'password');

insert into user_info(username, email, password)
values 
('digo', 'rod.dearaujo@gmail.com', 'password1');

insert into user_info(username, email, password)
values 
('jane', 'mynameisjane@gmail.com', 'password2');

insert into user_info(username, email, password)
values 
('joe', 'iamjoe@gmail.com', 'password3');

insert into weekday(day)
values
('Monday');

insert into weekday(day)
values
('Tuesday');

insert into weekday(day)
values
('Wednesday');

insert into weekday(day)
values
('Thursday');

insert into weekday(day)
values
('Friday');

insert into weekday(day)
values
('Saturday');

insert into weekday(day)
values
('Sunday');

insert into meal_type(type)
values
('Breakfast');

insert into meal_type(type)
values
('Snack');

insert into meal_type(type)
values
('Lunch');

insert into meal_type(type)
values
('Dinner');

insert into house(name, grocery_day, household_number)
values
('My Lovely Home', 'Friday', 2);

insert into house(name, grocery_day, household_number)
values
('My Lovely Home', 'Saturday', 4);

insert into house(name, grocery_day, household_number)
values
('Spaceship', 'Sunday', 1);

insert into recipe(name, serves_for)
values
('Baked Potato', 4);

insert into recipe(name, serves_for)
values
('Beans with rice', 6);

insert into recipe(name, serves_for)
values
('No Flour Pancake', 2);

insert into recipe(name, serves_for)
values
('Roast Chicken', 4);

insert into recipe_type(recipe_id, type_id)
values
(2, 4);

insert into recipe_type(recipe_id, type_id)
values
(4, 4);

insert into recipe_type(recipe_id, type_id)
values
(2, 3);

insert into recipe_type(recipe_id, type_id)
values
(3, 1);

insert into recipe_type(recipe_id, type_id)
values
(3, 2);

insert into recipe_type(recipe_id, type_id)
values
(1, 4);

insert into recipe_type(recipe_id, type_id)
values
(1, 3);

insert into ingredient(name, carb_per_100g, protein_per_100g, fat_per_100g, fiber_per_100g, calories_per_100g)
values
('potato', 17, 0, 0.1, 2.2, 77);

insert into ingredient(name, carb_per_100g, protein_per_100g, fat_per_100g, fiber_per_100g, calories_per_100g)
values
('milk', 5, 3.4, 1, 0, 42);

insert into ingredient(name, carb_per_100g, protein_per_100g, fat_per_100g, fiber_per_100g, calories_per_100g)
values
('parmesan cheese', 4.1, 38, 29, 0, 431);

insert into ingredient(name, carb_per_100g, protein_per_100g, fat_per_100g, fiber_per_100g, calories_per_100g)
values
('chicken breast', 0, 31, 3.6, 0, 165);

insert into ingredient(name, carb_per_100g, protein_per_100g, fat_per_100g, fiber_per_100g, calories_per_100g)
values
('oregano', 69, 9, 4.3, 43, 265);

insert into ingredient(name, carb_per_100g, protein_per_100g, fat_per_100g, fiber_per_100g, calories_per_100g)
values
('banana', 23, 1.1, 0.3, 2.6, 89);

insert into ingredient(name, carb_per_100g, protein_per_100g, fat_per_100g, fiber_per_100g, calories_per_100g)
values
('egg', 1.1, 13, 11, 0, 155);

insert into ingredient(name, carb_per_100g, protein_per_100g, fat_per_100g, fiber_per_100g, calories_per_100g)
values
('baking powder', 28, 0, 0, 0.2, 53);

insert into ingredient(id, name, carb_per_100g, protein_per_100g, fat_per_100g, fiber_per_100g, calories_per_100g)
values
(0,'empty', 0, 0, 0, 0, 0);

insert into step(recipe_id, id, text)
values 
(1, 1, 'peel and cut the potatoes into an inch thick disks');

insert into step(recipe_id, id, text)
values 
(1, 2, 'mix the milk and parmesan together');

insert into step(recipe_id, id, text)
values 
(1, 3, 'mix everything together and put them in the oven at 425 degrees for 45 minutes');

insert into unit(name)
values
('kg');

insert into unit(name)
values
('pound');

insert into unit(name)
values
('lb');

insert into unit(name)
values
('tbsp');

insert into unit(name)
values
('tsp');

insert into unit(name)
values
('ounce');

insert into unit(name)
values
('quantity');

insert into unit(name)
values
('litre');

insert into unit(name)
values
('cup');

insert into unit(name)
values
('grams');

insert into step_ingredient(recipe_id, step_id, ingredient_id, unit_id, amount)
values
(1,1,1,8,4);

insert into step_ingredient(recipe_id, step_id, ingredient_id, unit_id, amount)
values
(1,2,2,9,0.25);

insert into step_ingredient(recipe_id, step_id, ingredient_id, unit_id, amount)
values
(1,2,3,10,1);

insert into house_recipe(house_id, recipe_id)
values 
(1, 1);

insert into house_recipe(house_id, recipe_id)
values 
(1, 4);

insert into house_recipe(house_id, recipe_id)
values 
(3, 2);

insert into item_in_storage(house_id, ingredient_id, amount, unit_id)
values
(1, 1, 5, 8);

insert into item_in_storage(house_id, ingredient_id, amount, unit_id)
values
(3, 2, 0, 2);

insert into ownership(description)
values
('owner');

insert into ownership(description)
values
('resident');

insert into ownership(description)
values
('not allowed');

insert into member_of(user_id, house_id, own_type)
values 
(1, 1, 1);

insert into member_of(user_id, house_id, own_type)
values 
(2, 1, 2);

insert into member_of(user_id, house_id, own_type)
values 
(4, 1, 3);

insert into user_recipe(user_id, recipe_id)
values
(1,1);

insert into user_recipe(user_id, recipe_id)
values
(2,2);

insert into user_recipe(user_id, recipe_id)
values
(2,3);

insert into user_recipe(user_id, recipe_id)
values
(4,4);

insert into schedule(house_id, week_id, type_id, recipe_id)
values 
(1,1,2,1);

insert into schedule(house_id, week_id, type_id, recipe_id)
values 
(2,6,1,3);

insert into schedule(house_id, week_id, type_id, recipe_id)
values 
(2,2,1,3);

insert into schedule(house_id, week_id, type_id, recipe_id)
values 
(2,3,3,4);