CREATE TABLE IF NOT EXISTS vehicles (
    id serial PRIMARY KEY,
    make VARCHAR (50),
    model VARCHAR (50),
    vin VARCHAR (50) UNIQUE
);