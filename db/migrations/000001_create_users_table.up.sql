CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    identification_type   VARCHAR(20) NOT NULL,
    identification_number VARCHAR(50) NOT NULL,

    first_name VARCHAR(100) NOT NULL,
    last_name  VARCHAR(100) NOT NULL,

    birth_date DATE,

    address JSONB,

    phone_country_code VARCHAR(10) NOT NULL,
    phone_number       VARCHAR(30) NOT NULL,

    email VARCHAR(255),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_phone
    ON users(phone_country_code, phone_number);

CREATE INDEX idx_users_email
    ON users(email);
