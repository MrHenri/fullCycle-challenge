create table if not exists accounts (
    id uuid PRIMARY KEY,
    number TEXT,
    balance REAL
);

create table if not exists transfers (
    id uuid PRIMARY KEY,
    from_account TEXT,
    to_account TEXT,
    amount REAL
);