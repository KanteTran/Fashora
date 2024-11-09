-- User Phones table (Main User Table)
CREATE TABLE user_phones (
     phone_id VARCHAR(255) PRIMARY KEY,                -- Primary key
     user_name VARCHAR(100),                           -- User name, nullable
     birthday DATE,                                    -- Birthday, nullable
     address VARCHAR(255),                             -- Address, nullable
     device_id VARCHAR(100),                           -- Device identifier, nullable
     gender INT CHECK (gender IN (0, 1, 2)),           -- Gender: 0 (male), 1 (female), 2 (other), nullable
     password_hash VARCHAR(255),                       -- Hashed password, nullable
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,   -- Automatically sets time on creation
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP                      -- Automatically updates time on modification -- Index on soft delete field
);

CREATE TABLE tokens (
    phone_id VARCHAR(255) PRIMARY KEY,                -- Định nghĩa phone_id là PRIMARY KEY
    token VARCHAR(255) UNIQUE NOT NULL,               -- JWT token, phải là duy nhất
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Tự động đặt thời gian tạo
    expired_time TIMESTAMP NOT NULL,                  -- Thời gian hết hạn của token
    CONSTRAINT fk_phone_id FOREIGN KEY (phone_id) REFERENCES user_phones(phone_id) -- Ràng buộc FOREIGN KEY
);

-- Pose Images table
CREATE TABLE pose_images (
                             pose_id VARCHAR(255) PRIMARY KEY,                  -- Unique identifier for the pose image
                             phone_id VARCHAR(255) REFERENCES user_phones(phone_id),  -- Foreign key to user_phones
                             image_url VARCHAR(255),                            -- URL to the item image (stored in S3)
                             status VARCHAR(50)                                 -- e.g., "available", "unavailable"
);

-- Wardrobe Items table
CREATE TABLE wardrobe_items (
                                item_id VARCHAR(255) PRIMARY KEY,                  -- Unique identifier for the wardrobe item
                                phone_id VARCHAR(255) REFERENCES user_phones(phone_id),  -- Foreign key to user_phones
                                image_url VARCHAR(255),                            -- URL to the item image (stored in S3)
                                status INT CHECK (status IN (0, 1))                -- 1: available, 0: unavailable
);

-- Try-On Sessions table
CREATE TABLE try_on_sessions (
                                 session_id VARCHAR(255) PRIMARY KEY,               -- Unique identifier for the try-on session
                                 store_id VARCHAR(255),                             -- Store ID, nullable
                                 phone_id VARCHAR(255) REFERENCES user_phones(phone_id),  -- Foreign key to user_phones
                                 clothing_image_url VARCHAR(255),                   -- URL to the clothing image selected for try-on
                                 pose_image_url VARCHAR(255),                       -- URL to the pose image of the user
                                 mask_url VARCHAR(255),                             -- URL to the mask image drawn by the user
                                 try_on_result_url VARCHAR(255),                    -- URL to the try-on result image
                                 status VARCHAR(50),                                -- e.g., "in_progress", "completed"
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- Automatically sets time on creation
                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Shop table
CREATE TABLE user_shop (
                           phone_id VARCHAR(255) PRIMARY KEY,                 -- Primary key
                           store_name VARCHAR(100),                           -- Store name, nullable
                           address VARCHAR(255),                              -- Address, nullable
                           password VARCHAR(255),                             -- Password (should be hashed), nullable
                           status INT CHECK (status IN (0, 1))                -- 1: active, 0: inactive
);

-- Store Items table
CREATE TABLE store_items (
                             item_id VARCHAR(255) PRIMARY KEY,                  -- Unique identifier for the item
                             store_id VARCHAR(255) REFERENCES user_shop(phone_id), -- Foreign key to user_shop
                             image_url VARCHAR(255),                            -- URL to the item image (stored in S3)
                             status VARCHAR(50),                                -- e.g., "available", "unavailable"
                             name_item VARCHAR(100)                             -- Name of the item
);

-- Count Items table
CREATE TABLE count_items (
                             phone_id VARCHAR(255) PRIMARY KEY REFERENCES user_phones(phone_id), -- Foreign key to user_phones
                             count_wardrobe_items INT,                       -- Number of items in wardrobe
                             count_pose_items INT                            -- Number of pose images
);