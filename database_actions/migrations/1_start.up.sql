CREATE TABLE IF NOT EXISTS breeds (
    id INT AUTO_INCREMENT PRIMARY KEY,
    species VARCHAR(50),
    pet_size VARCHAR(50),
    name VARCHAR(100),
    average_male_adult_weight INT,
    average_female_adult_weight INT
);