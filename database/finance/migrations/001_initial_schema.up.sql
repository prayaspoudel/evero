-- Finance Module - Initial Schema
-- Version: 001
-- Description: Create core tables for finance module including accounts, transactions, ledgers, and payments

-- ============================================================================
-- Core Tables
-- ============================================================================

-- Chart of Accounts
CREATE TABLE IF NOT EXISTS finance_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_code VARCHAR(50) UNIQUE NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) NOT NULL, -- ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE
    parent_account_id UUID REFERENCES finance_accounts(id) ON DELETE SET NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    is_active BOOLEAN DEFAULT true,
    is_system BOOLEAN DEFAULT false,
    description TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- General Ledger
CREATE TABLE IF NOT EXISTS finance_ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_number VARCHAR(50) UNIQUE NOT NULL,
    entry_date DATE NOT NULL,
    reference_type VARCHAR(50), -- INVOICE, PAYMENT, ADJUSTMENT, JOURNAL
    reference_id UUID,
    description TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, POSTED, REVERSED
    posted_at TIMESTAMP WITH TIME ZONE,
    posted_by UUID,
    reversed_at TIMESTAMP WITH TIME ZONE,
    reversed_by UUID,
    reversal_entry_id UUID REFERENCES finance_ledger_entries(id),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Ledger Entry Lines (Double-entry bookkeeping)
CREATE TABLE IF NOT EXISTS finance_ledger_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ledger_entry_id UUID NOT NULL REFERENCES finance_ledger_entries(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES finance_accounts(id),
    debit_amount DECIMAL(19, 4) DEFAULT 0,
    credit_amount DECIMAL(19, 4) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    exchange_rate DECIMAL(19, 6) DEFAULT 1,
    description TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_debit_or_credit CHECK (
        (debit_amount > 0 AND credit_amount = 0) OR 
        (credit_amount > 0 AND debit_amount = 0)
    )
);

-- Transactions
CREATE TABLE IF NOT EXISTS finance_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_number VARCHAR(50) UNIQUE NOT NULL,
    transaction_type VARCHAR(50) NOT NULL, -- PAYMENT, RECEIPT, TRANSFER, REFUND
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    amount DECIMAL(19, 4) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    from_account_id UUID REFERENCES finance_accounts(id),
    to_account_id UUID REFERENCES finance_accounts(id),
    payment_method VARCHAR(50), -- CASH, CARD, BANK_TRANSFER, CHECK
    reference_number VARCHAR(100),
    status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, COMPLETED, FAILED, CANCELLED
    description TEXT,
    ledger_entry_id UUID REFERENCES finance_ledger_entries(id),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Invoices
CREATE TABLE IF NOT EXISTS finance_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id UUID NOT NULL,
    invoice_date DATE NOT NULL,
    due_date DATE NOT NULL,
    subtotal DECIMAL(19, 4) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(19, 4) DEFAULT 0,
    discount_amount DECIMAL(19, 4) DEFAULT 0,
    total_amount DECIMAL(19, 4) NOT NULL DEFAULT 0,
    paid_amount DECIMAL(19, 4) DEFAULT 0,
    balance_due DECIMAL(19, 4) NOT NULL DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(20) DEFAULT 'DRAFT', -- DRAFT, SENT, PARTIAL, PAID, OVERDUE, CANCELLED
    payment_terms VARCHAR(50), -- NET_30, NET_60, DUE_ON_RECEIPT
    notes TEXT,
    ledger_entry_id UUID REFERENCES finance_ledger_entries(id),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Invoice Line Items
CREATE TABLE IF NOT EXISTS finance_invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES finance_invoices(id) ON DELETE CASCADE,
    item_description TEXT NOT NULL,
    quantity DECIMAL(19, 4) NOT NULL DEFAULT 1,
    unit_price DECIMAL(19, 4) NOT NULL,
    tax_rate DECIMAL(5, 2) DEFAULT 0,
    discount_rate DECIMAL(5, 2) DEFAULT 0,
    line_total DECIMAL(19, 4) NOT NULL,
    account_id UUID REFERENCES finance_accounts(id),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Payments
CREATE TABLE IF NOT EXISTS finance_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_number VARCHAR(50) UNIQUE NOT NULL,
    invoice_id UUID REFERENCES finance_invoices(id),
    transaction_id UUID REFERENCES finance_transactions(id),
    payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    amount DECIMAL(19, 4) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    payment_method VARCHAR(50) NOT NULL,
    reference_number VARCHAR(100),
    status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, COMPLETED, FAILED, REFUNDED
    notes TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Budget Management
CREATE TABLE IF NOT EXISTS finance_budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_name VARCHAR(255) NOT NULL,
    fiscal_year INTEGER NOT NULL,
    account_id UUID NOT NULL REFERENCES finance_accounts(id),
    period VARCHAR(20) NOT NULL, -- MONTHLY, QUARTERLY, ANNUAL
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    budgeted_amount DECIMAL(19, 4) NOT NULL,
    actual_amount DECIMAL(19, 4) DEFAULT 0,
    variance DECIMAL(19, 4) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'ACTIVE', -- ACTIVE, CLOSED, EXCEEDED
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Financial Periods (for closing)
CREATE TABLE IF NOT EXISTS finance_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'OPEN', -- OPEN, CLOSED, LOCKED
    closed_at TIMESTAMP WITH TIME ZONE,
    closed_by UUID,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(start_date, end_date)
);

-- Reconciliation
CREATE TABLE IF NOT EXISTS finance_reconciliations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES finance_accounts(id),
    reconciliation_date DATE NOT NULL,
    statement_balance DECIMAL(19, 4) NOT NULL,
    book_balance DECIMAL(19, 4) NOT NULL,
    difference DECIMAL(19, 4) NOT NULL,
    status VARCHAR(20) DEFAULT 'IN_PROGRESS', -- IN_PROGRESS, RECONCILED, DISCREPANCY
    reconciled_by UUID,
    reconciled_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- Indexes for Performance
-- ============================================================================

-- Accounts
CREATE INDEX idx_finance_accounts_type ON finance_accounts(account_type);
CREATE INDEX idx_finance_accounts_parent ON finance_accounts(parent_account_id);
CREATE INDEX idx_finance_accounts_code ON finance_accounts(account_code);
CREATE INDEX idx_finance_accounts_active ON finance_accounts(is_active);

-- Ledger
CREATE INDEX idx_finance_ledger_entries_date ON finance_ledger_entries(entry_date);
CREATE INDEX idx_finance_ledger_entries_status ON finance_ledger_entries(status);
CREATE INDEX idx_finance_ledger_entries_reference ON finance_ledger_entries(reference_type, reference_id);
CREATE INDEX idx_finance_ledger_lines_entry ON finance_ledger_lines(ledger_entry_id);
CREATE INDEX idx_finance_ledger_lines_account ON finance_ledger_lines(account_id);

-- Transactions
CREATE INDEX idx_finance_transactions_date ON finance_transactions(transaction_date);
CREATE INDEX idx_finance_transactions_type ON finance_transactions(transaction_type);
CREATE INDEX idx_finance_transactions_status ON finance_transactions(status);
CREATE INDEX idx_finance_transactions_from_account ON finance_transactions(from_account_id);
CREATE INDEX idx_finance_transactions_to_account ON finance_transactions(to_account_id);

-- Invoices
CREATE INDEX idx_finance_invoices_customer ON finance_invoices(customer_id);
CREATE INDEX idx_finance_invoices_date ON finance_invoices(invoice_date);
CREATE INDEX idx_finance_invoices_due_date ON finance_invoices(due_date);
CREATE INDEX idx_finance_invoices_status ON finance_invoices(status);
CREATE INDEX idx_finance_invoice_items_invoice ON finance_invoice_items(invoice_id);

-- Payments
CREATE INDEX idx_finance_payments_invoice ON finance_payments(invoice_id);
CREATE INDEX idx_finance_payments_transaction ON finance_payments(transaction_id);
CREATE INDEX idx_finance_payments_date ON finance_payments(payment_date);
CREATE INDEX idx_finance_payments_status ON finance_payments(status);

-- Budgets
CREATE INDEX idx_finance_budgets_account ON finance_budgets(account_id);
CREATE INDEX idx_finance_budgets_fiscal_year ON finance_budgets(fiscal_year);
CREATE INDEX idx_finance_budgets_dates ON finance_budgets(start_date, end_date);

-- Periods
CREATE INDEX idx_finance_periods_dates ON finance_periods(start_date, end_date);
CREATE INDEX idx_finance_periods_status ON finance_periods(status);

-- Reconciliations
CREATE INDEX idx_finance_reconciliations_account ON finance_reconciliations(account_id);
CREATE INDEX idx_finance_reconciliations_date ON finance_reconciliations(reconciliation_date);
CREATE INDEX idx_finance_reconciliations_status ON finance_reconciliations(status);

-- ============================================================================
-- Comments
-- ============================================================================

COMMENT ON TABLE finance_accounts IS 'Chart of accounts - defines all financial accounts';
COMMENT ON TABLE finance_ledger_entries IS 'General ledger entries - header for journal entries';
COMMENT ON TABLE finance_ledger_lines IS 'General ledger lines - detail lines for double-entry bookkeeping';
COMMENT ON TABLE finance_transactions IS 'Financial transactions - payments, receipts, transfers';
COMMENT ON TABLE finance_invoices IS 'Customer invoices';
COMMENT ON TABLE finance_invoice_items IS 'Line items for invoices';
COMMENT ON TABLE finance_payments IS 'Payment records linked to invoices and transactions';
COMMENT ON TABLE finance_budgets IS 'Budget management and tracking';
COMMENT ON TABLE finance_periods IS 'Financial periods for period closing';
COMMENT ON TABLE finance_reconciliations IS 'Bank and account reconciliation records';
