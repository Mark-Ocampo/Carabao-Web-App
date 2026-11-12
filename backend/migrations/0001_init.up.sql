-- CARA schema (PostgreSQL). All money is integer centavos (BIGINT).

CREATE TABLE users (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                TEXT NOT NULL UNIQUE,
    password_hash        TEXT NOT NULL,
    full_name            TEXT NOT NULL,
    emergency_min_months INT NOT NULL DEFAULT 3 CHECK (emergency_min_months >= 0),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Cookie sessions: store only a hash of the session token
CREATE TABLE sessions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE accounts (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name          TEXT NOT NULL,
    type          TEXT NOT NULL CHECK (type IN ('cash', 'bank', 'ewallet', 'credit')),
    balance_cents BIGINT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- kind is the wants vs. needs tag, used for expense categories only
CREATE TABLE categories (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name    TEXT NOT NULL,
    type    TEXT NOT NULL CHECK (type IN ('income', 'expense')),
    kind    TEXT CHECK (kind IN ('need', 'want')),
    UNIQUE (user_id, name, type),
    CHECK (type = 'expense' OR kind IS NULL)
);

CREATE TABLE transactions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id   UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    category_id  UUID REFERENCES categories(id) ON DELETE SET NULL,
    type         TEXT NOT NULL CHECK (type IN ('income', 'expense')),
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    note         TEXT,
    occurred_on  DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- month is stored as the first day of the month, e.g. 2026-11-01
CREATE TABLE budgets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON