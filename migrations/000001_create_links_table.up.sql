CREATE TABLE links (
                       id BIGSERIAL PRIMARY KEY,
                       short_code VARCHAR(10) UNIQUE NOT NULL,
                       original_url TEXT NOT NULL,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                       click_count BIGINT NOT NULL DEFAULT 0
);