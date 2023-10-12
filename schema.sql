CREATE TABLE IF NOT EXISTS tax_transaction
(
    id                  SERIAL      PRIMARY KEY,
    transaction_date    INT(11)     NOT NULL,
    deposit_rp          BIGINT(20)  DEFAULT 0,
    withdraw_rp         BIGINT(20)  DEFAULT 0,
    fee                 BIGINT(20)  DEFAULT 0,
    upline_bonus        BIGINT(20)  DEFAULT 0,
    remain              BIGINT(20)  DEFAULT 0,
    ppn                 BIGINT(20)  DEFAULT 0
)

CREATE INDEX IF NOT EXISTS tax_transaction_transaction_date_index
    ON tax_transaction (transaction_date);