# Sprint 3: Testing Guide

## Prerequisites

1. **Start the database and Redis:**
   ```bash
   cd database && make up
   ```

2. **Start the Go API server:**
   ```bash
   cd api && go run ./cmd/server
   ```
   Server runs on http://localhost:3000

3. **Start the SvelteKit frontend:**
   ```bash
   cd frontend/help-me-budget && npm run dev
   ```
   Frontend runs on http://localhost:5173

4. **Login** to the application using Supabase OAuth (Google/Facebook)

## Testing Transactions Feature

### 1. Navigate to Transactions Page
- Click "Transactions" in the left sidebar
- URL: http://localhost:5173/dashboard/transactions

### 2. Test Empty State
- If no transactions exist, you should see:
  - 💸 icon
  - "No transactions yet" message
  - "Add Your First Transaction" button

### 3. Test Create Transaction (Expense)
1. Click "Add Transaction" button (top-right or empty state)
2. Fill out the form:
   - Type: **Expense**
   - Amount: **50.00**
   - Account: Select an account (create one in /dashboard/accounts if needed)
   - Category: Select a category (optional)
   - Date: Today (default)
   - Description: "Groceries at Whole Foods"
   - Notes: "Weekly shopping"
3. Click "Add Transaction"
4. Verify:
   - ✅ Modal closes
   - ✅ New transaction appears in list
   - ✅ Amount shows in red with minus sign (-$50.00)
   - ✅ Description displays correctly
   - ✅ Category and Account show properly

### 4. Test Create Transaction (Income)
1. Click "Add Transaction" button again
2. Fill out the form:
   - Type: **Income**
   - Amount: **1500.00**
   - Account: Select an account
   - Category: Select income category (e.g., "Salary")
   - Date: Today
   - Description: "Monthly paycheck"
3. Click "Add Transaction"
4. Verify:
   - ✅ Amount shows in green with plus sign (+$1,500.00)
   - ✅ Transaction appears at top of list (newest first)

### 5. Test Filters

**Filter by Account:**
1. Select an account from "Account" dropdown
2. Click "Apply Filters"
3. Verify only transactions from that account appear

**Filter by Category:**
1. Select a category from "Category" dropdown
2. Click "Apply Filters"
3. Verify only transactions in that category appear

**Filter by Date Range:**
1. Set "From Date" to 7 days ago
2. Set "To Date" to today
3. Click "Apply Filters"
4. Verify only transactions within date range appear

**Clear Filters:**
1. Click "Clear" button
2. Verify all transactions appear again
3. Verify URL resets to /dashboard/transactions

### 6. Test Edit Transaction
1. Hover over a transaction row
2. Click the blue edit icon (pencil)
3. Modify the transaction:
   - Change amount to **75.00**
   - Update description
4. Click "Save Changes"
5. Verify:
   - ✅ Modal closes
   - ✅ Transaction updates in list
   - ✅ Amount and description reflect changes

### 7. Test Delete Transaction
1. Click the red delete icon (trash) on a transaction
2. Verify delete confirmation modal shows:
   - ✅ Transaction details displayed
   - ✅ "Are you sure?" warning message
3. Click "Delete"
4. Verify:
   - ✅ Modal closes
   - ✅ Transaction removed from list

5. Click delete on another transaction
6. Click "Cancel" instead
7. Verify:
   - ✅ Modal closes
   - ✅ Transaction still in list

### 8. Test Form Validation

**Create Modal:**
1. Open create modal
2. Try submitting without filling required fields
3. Verify browser validation triggers:
   - ✅ Amount is required
   - ✅ Account is required
   - ✅ Date is required

**Invalid Amount:**
1. Enter negative or zero amount
2. Verify validation prevents submission

### 9. Test Table Display
1. Add multiple transactions (5-10)
2. Verify:
   - ✅ Transactions ordered by date (newest first)
   - ✅ Visual indicators: green dot (income), red dot (expense)
   - ✅ Color-coded amounts: green +, red -
   - ✅ Table scrolls horizontally on mobile
   - ✅ Hover effect on rows

### 10. Test Keyboard Navigation
1. Tab through the create modal form
2. Verify all fields are keyboard accessible
3. Press Escape to close modal (if implemented)

## Backend API Testing (Optional)

### Using curl or Postman:

**Get Transactions:**
```bash
curl -H "X-API-Key: YOUR_API_KEY" \
     -H "X-User-ID: YOUR_USER_ID" \
     http://localhost:3000/api/transactions
```

**Create Transaction:**
```bash
curl -X POST \
     -H "X-API-Key: YOUR_API_KEY" \
     -H "X-User-ID: YOUR_USER_ID" \
     -H "Content-Type: application/json" \
     -d '{
       "account_id": "ACCOUNT_UUID",
       "amount": 50.00,
       "transaction_type": "expense",
       "description": "Test transaction",
       "transaction_date": "2024-12-01"
     }' \
     http://localhost:3000/api/transactions
```

**Filter Transactions:**
```bash
curl -H "X-API-Key: YOUR_API_KEY" \
     -H "X-User-ID: YOUR_USER_ID" \
     "http://localhost:3000/api/transactions?start_date=2024-11-01&end_date=2024-12-01"
```

## Known Issues / Future Enhancements

### Not Implemented Yet (Future Sprints):
- CSV import for bulk transaction upload
- Automatic matching to budget entries
- Bulk operations (categorize multiple)
- Transaction attachments (receipts)
- Recurring transaction templates
- Advanced search (by description text)

### Accessibility Warnings:
- Modal backdrop click handlers need keyboard event handlers
- These are warnings, not blocking errors
- Will be addressed in accessibility improvements

## Success Criteria

All tests should pass with:
- ✅ No compilation errors (Go or TypeScript)
- ✅ Transactions create, update, delete successfully
- ✅ Filters work correctly
- ✅ UI displays properly on desktop and mobile
- ✅ Form validation prevents invalid data
- ✅ Visual indicators match transaction type
- ✅ Navigation works seamlessly

## Troubleshooting

**"Failed to load transactions":**
- Check that Go API server is running on port 3000
- Verify API_SECRET_KEY matches in both .env files
- Check browser console for errors

**"Unauthorized" error:**
- Ensure you're logged in via Supabase
- Check that getLocalUserId() returns valid user ID
- Verify API key is set correctly

**Empty accounts/categories dropdown:**
- Navigate to /dashboard/accounts and create an account
- Navigate to /dashboard/categories and seed or create categories

**Filters not working:**
- Check URL query parameters are being set
- Verify backend receives filter parameters
- Check network tab for API request/response
