-- Constraints para tabla schools

ALTER TABLE schools ADD CONSTRAINT schools_code_unique UNIQUE (code);
ALTER TABLE schools ADD CONSTRAINT schools_subscription_tier_check CHECK (subscription_tier IN ('free', 'basic', 'premium', 'enterprise'));
