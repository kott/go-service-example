
CREATE TABLE IF NOT EXISTS articles (
id uuid NOT NULL DEFAULT uuid_generate_v4(),
title text,
body text,
created_at timestamptz not null,
updated_at timestamptz not null,
disabled_at timestamptz
);
