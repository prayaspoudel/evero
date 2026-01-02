-- Finance Module - Rollback Initial Schema
-- Version: 001
-- Description: Drop all finance module tables

-- Drop indexes first
DROP INDEX IF EXISTS idx_finance_reconciliations_status;
DROP INDEX IF EXISTS idx_finance_reconciliations_date;
DROP INDEX IF EXISTS idx_finance_reconciliations_account;
DROP INDEX IF EXISTS idx_finance_periods_status;
DROP INDEX IF EXISTS idx_finance_periods_dates;
DROP INDEX IF EXISTS idx_finance_budgets_dates;
DROP INDEX IF EXISTS idx_finance_budgets_fiscal_year;
DROP INDEX IF EXISTS idx_finance_budgets_account;
DROP INDEX IF EXISTS idx_finance_payments_status;
DROP INDEX IF EXISTS idx_finance_payments_date;
DROP INDEX IF EXISTS idx_finance_payments_transaction;
DROP INDEX IF EXISTS idx_finance_payments_invoice;
DROP INDEX IF EXISTS idx_finance_invoice_items_invoice;
DROP INDEX IF EXISTS idx_finance_invoices_status;
DROP INDEX IF EXISTS idx_finance_invoices_due_date;
DROP INDEX IF EXISTS idx_finance_invoices_date;
DROP INDEX IF EXISTS idx_finance_invoices_customer;
DROP INDEX IF EXISTS idx_finance_transactions_to_account;
DROP INDEX IF EXISTS idx_finance_transactions_from_account;
DROP INDEX IF EXISTS idx_finance_transactions_status;
DROP INDEX IF EXISTS idx_finance_transactions_type;
DROP INDEX IF EXISTS idx_finance_transactions_date;
DROP INDEX IF EXISTS idx_finance_ledger_lines_account;
DROP INDEX IF EXISTS idx_finance_ledger_lines_entry;
DROP INDEX IF EXISTS idx_finance_ledger_entries_reference;
DROP INDEX IF EXISTS idx_finance_ledger_entries_status;
DROP INDEX IF EXISTS idx_finance_ledger_entries_date;
DROP INDEX IF EXISTS idx_finance_accounts_active;
DROP INDEX IF EXISTS idx_finance_accounts_code;
DROP INDEX IF EXISTS idx_finance_accounts_parent;
DROP INDEX IF EXISTS idx_finance_accounts_type;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS finance_reconciliations;
DROP TABLE IF EXISTS finance_periods;
DROP TABLE IF EXISTS finance_budgets;
DROP TABLE IF EXISTS finance_payments;
DROP TABLE IF EXISTS finance_invoice_items;
DROP TABLE IF EXISTS finance_invoices;
DROP TABLE IF EXISTS finance_transactions;
DROP TABLE IF EXISTS finance_ledger_lines;
DROP TABLE IF EXISTS finance_ledger_entries;
DROP TABLE IF EXISTS finance_accounts;
