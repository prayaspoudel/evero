# Finance Module

## Overview

The Finance module provides comprehensive financial management capabilities for the Evero platform, including general ledger, accounts payable/receivable, budgeting, and financial reporting.

## Features

### General Ledger
- Chart of accounts management
- Double-entry bookkeeping
- Journal entries
- Account reconciliation
- Period closing
- Multi-currency support

### Accounts Receivable
- Customer invoicing
- Payment tracking
- Aging reports
- Credit management
- Collection management
- Automatic reminders

### Accounts Payable
- Vendor management
- Bill processing
- Payment scheduling
- Expense tracking
- Purchase order integration

### Budgeting & Forecasting
- Budget creation and management
- Budget vs actual analysis
- Variance reporting
- Cash flow forecasting
- Financial planning

### Financial Reporting
- Balance sheet
- Income statement
- Cash flow statement
- Trial balance
- Custom reports
- Consolidation

## Architecture

```
app/finance/
├── entities/           # Domain models
├── models/            # Database models
├── repositories/      # Data access layer
├── usecases/         # Business logic
├── controllers/      # HTTP handlers
├── middleware/       # Finance-specific middleware
└── routes.go         # Route definitions
```

## Database Schema

### Core Tables

1. **finance_accounts** - Chart of accounts
2. **finance_ledger_entries** - General ledger header
3. **finance_ledger_lines** - General ledger detail
4. **finance_transactions** - Financial transactions
5. **finance_invoices** - Customer invoices
6. **finance_invoice_items** - Invoice line items
7. **finance_payments** - Payment records
8. **finance_budgets** - Budget management
9. **finance_periods** - Financial periods
10. **finance_reconciliations** - Account reconciliation

See [database/finance/migrations/001_initial_schema.up.sql](../../database/finance/migrations/001_initial_schema.up.sql) for complete schema.

## API Endpoints

### Chart of Accounts
- `GET /api/v1/finance/accounts` - List accounts
- `POST /api/v1/finance/accounts` - Create account
- `GET /api/v1/finance/accounts/:id` - Get account details
- `PUT /api/v1/finance/accounts/:id` - Update account
- `DELETE /api/v1/finance/accounts/:id` - Deactivate account

### General Ledger
- `GET /api/v1/finance/ledger` - List ledger entries
- `POST /api/v1/finance/ledger` - Create journal entry
- `GET /api/v1/finance/ledger/:id` - Get entry details
- `POST /api/v1/finance/ledger/:id/post` - Post entry
- `POST /api/v1/finance/ledger/:id/reverse` - Reverse entry

### Invoices
- `GET /api/v1/finance/invoices` - List invoices
- `POST /api/v1/finance/invoices` - Create invoice
- `GET /api/v1/finance/invoices/:id` - Get invoice
- `PUT /api/v1/finance/invoices/:id` - Update invoice
- `POST /api/v1/finance/invoices/:id/send` - Send invoice
- `POST /api/v1/finance/invoices/:id/pay` - Record payment
- `GET /api/v1/finance/invoices/:id/pdf` - Download PDF

### Payments
- `GET /api/v1/finance/payments` - List payments
- `POST /api/v1/finance/payments` - Record payment
- `GET /api/v1/finance/payments/:id` - Get payment
- `POST /api/v1/finance/payments/:id/refund` - Process refund

### Transactions
- `GET /api/v1/finance/transactions` - List transactions
- `POST /api/v1/finance/transactions` - Create transaction
- `GET /api/v1/finance/transactions/:id` - Get transaction
- `PUT /api/v1/finance/transactions/:id` - Update transaction

### Budgets
- `GET /api/v1/finance/budgets` - List budgets
- `POST /api/v1/finance/budgets` - Create budget
- `GET /api/v1/finance/budgets/:id` - Get budget
- `PUT /api/v1/finance/budgets/:id` - Update budget
- `GET /api/v1/finance/budgets/:id/variance` - Get variance report

### Reports
- `GET /api/v1/finance/reports/balance-sheet` - Balance sheet
- `GET /api/v1/finance/reports/income-statement` - Income statement
- `GET /api/v1/finance/reports/cash-flow` - Cash flow statement
- `GET /api/v1/finance/reports/trial-balance` - Trial balance
- `GET /api/v1/finance/reports/aging` - Aging report
- `GET /api/v1/finance/reports/budget-variance` - Budget variance

### Period Management
- `GET /api/v1/finance/periods` - List periods
- `POST /api/v1/finance/periods` - Create period
- `POST /api/v1/finance/periods/:id/close` - Close period
- `POST /api/v1/finance/periods/:id/reopen` - Reopen period

### Reconciliation
- `GET /api/v1/finance/reconciliations` - List reconciliations
- `POST /api/v1/finance/reconciliations` - Start reconciliation
- `PUT /api/v1/finance/reconciliations/:id` - Update reconciliation
- `POST /api/v1/finance/reconciliations/:id/complete` - Complete reconciliation

## Configuration

### Environment Variables
```env
FINANCE_DB_HOST=localhost
FINANCE_DB_PORT=5432
FINANCE_DB_NAME=evero
FINANCE_PORT=3003
FINANCE_DEFAULT_CURRENCY=USD
FINANCE_FISCAL_YEAR_START=01-01
FINANCE_AUTO_POST_ENTRIES=false
```

### Config Files
- `config/finance/local.json`
- `config/finance/development.json`
- `config/finance/production.json`

## Double-Entry Bookkeeping

### Account Types
- **ASSET** - Debit normal balance
- **LIABILITY** - Credit normal balance
- **EQUITY** - Credit normal balance
- **REVENUE** - Credit normal balance
- **EXPENSE** - Debit normal balance

### Journal Entry Rules
- Every transaction must have equal debits and credits
- At least one debit and one credit per entry
- Total debits = Total credits (balanced entry)

### Example Entry
```json
{
  "entry_date": "2026-01-02",
  "description": "Customer payment received",
  "lines": [
    {
      "account_code": "1000",  // Cash
      "debit_amount": 1000.00
    },
    {
      "account_code": "1200",  // Accounts Receivable
      "credit_amount": 1000.00
    }
  ]
}
```

## Financial Periods

### Period Management
- Define fiscal year periods (monthly/quarterly/annual)
- Enforce posting to open periods only
- Period closing process
- Year-end closing
- Audit trail for period activities

### Closing Process
1. Verify all transactions are posted
2. Run reconciliations
3. Review period reports
4. Close the period (no further posting)
5. Optional: Lock period (prevent reopening)

## Multi-Currency Support

### Features
- Multi-currency transactions
- Exchange rate management
- Automatic currency conversion
- Realized/unrealized gains/losses
- Base currency reporting

### Currency Handling
```json
{
  "amount": 1000.00,
  "currency": "EUR",
  "exchange_rate": 1.18,
  "base_amount": 1180.00,  // In USD
  "base_currency": "USD"
}
```

## Integration

### Healthcare Module
- Patient billing integration
- Insurance claim payments
- Automatic revenue recognition
- Payment allocation

### Insurance Module
- Premium collection
- Claim payment processing
- Commission payments
- Policy accounting

### Access Module
- User authentication
- Role-based permissions
- Audit logging
- Multi-company support

## Security & Compliance

### Financial Controls
- Segregation of duties
- Approval workflows
- Audit trail for all transactions
- No direct data deletion (soft delete)
- Immutable ledger entries

### Compliance
- GAAP/IFRS standards
- Tax reporting requirements
- SOX compliance
- Audit support
- Data retention policies

## Automation

### Scheduled Tasks
- Automatic invoice generation
- Recurring transactions
- Payment reminders
- Budget monitoring alerts
- Period-end processes

### Workflows
- Invoice approval workflow
- Payment approval workflow
- Journal entry approval
- Expense approval

## Build & Deploy

### Build
```bash
make build-finance
```

### Setup
```bash
make setup-finance
```

### Run Migrations
```bash
make migrate-finance
```

### Deploy
```bash
make deploy-finance
```

## Testing

```bash
make test-finance
```

### Test Data
- Sample chart of accounts
- Test transactions
- Sample invoices
- Budget templates

## Monitoring

### Key Metrics
- Revenue
- Expenses
- Profit margin
- Cash balance
- Accounts receivable aging
- Accounts payable aging
- Budget variance

### Alerts
- Negative cash balance
- Overdue invoices
- Budget exceeded
- Failed transactions
- Reconciliation discrepancies

## Reporting

### Standard Reports
1. **Balance Sheet** - Financial position
2. **Income Statement** - Profitability
3. **Cash Flow Statement** - Cash movements
4. **Trial Balance** - Account balances
5. **General Ledger** - Detailed transactions
6. **Aging Report** - AR/AP aging
7. **Budget Variance** - Budget vs actual

### Custom Reports
- Configurable report builder
- Scheduled report generation
- Export formats (PDF, Excel, CSV)
- Dashboard widgets

## Best Practices

### Chart of Accounts
- Use consistent numbering scheme
- Group accounts logically
- Document account purposes
- Regular review and cleanup

### Journal Entries
- Clear descriptions
- Supporting documentation
- Timely posting
- Regular reconciliation

### Invoicing
- Consistent numbering
- Clear terms and conditions
- Timely delivery
- Follow-up on overdue

## Development

### Getting Started

1. Setup the module:
   ```bash
   make setup-finance
   ```

2. Run migrations:
   ```bash
   make migrate-finance
   ```

3. Build and run:
   ```bash
   make build-finance
   ./bin/finance
   ```

### Code Structure

- Follow clean architecture principles
- Use dependency injection
- Implement repository pattern
- Maintain data integrity
- Write comprehensive tests

## Documentation

- [API Documentation](./API.md)
- [Database Schema](./SCHEMA.md)
- [Accounting Guide](./ACCOUNTING.md)
- [Integration Guide](./INTEGRATION.md)
- [Compliance Guide](./COMPLIANCE.md)

## Troubleshooting

### Common Issues

**Unbalanced Entry**
- Verify total debits = total credits
- Check decimal precision
- Review account types

**Posting to Closed Period**
- Verify period status
- Reopen period if needed
- Create adjusting entry in current period

**Currency Conversion**
- Update exchange rates
- Verify currency codes
- Check conversion logic

## Support

For issues or questions:
1. Check module documentation
2. Review error logs
3. Verify configuration
4. Contact finance module team

---

**Last Updated:** January 2, 2026  
**Version:** 1.0.0  
**Status:** In Development
