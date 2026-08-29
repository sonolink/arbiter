-- +goose Up
CREATE TYPE forge AS ENUM ('github');
CREATE TYPE community AS ENUM ('discord');

CREATE TABLE links (
  id UUID PRIMARY KEY DEFAULT UUIDV7(),
  forge forge NOT NULL,
  forge_user_id TEXT COLLATE "C" NOT NULL, -- identifier should be byte-for-byte identical.
  community community NOT NULL,
  community_user_id TEXT COLLATE "C" NOT NULL, -- identifier should be byte-for-byte identical.
  created_at TIMESTAMPTZ GENERATED ALWAYS AS (uuid_extract_timestamp(id)) VIRTUAL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  UNIQUE (forge, forge_user_id, community),
  UNIQUE (community, community_user_id, forge)
);

-- +goose Down
DROP TABLE links;
DROP TYPE community;
DROP TYPE forge;
