CREATE TABLE ads (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    category VARCHAR(100) NOT NULL,
    contact_phone VARCHAR(12),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_ads_category ON ads(category);
CREATE INDEX idx_ads_created_at ON ads(created_at DESC);