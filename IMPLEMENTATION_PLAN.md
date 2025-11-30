# Help-Me-Budget: Implementation Roadmap

## 📊 Current State Analysis

### ✅ Completed
- Authentication (Supabase OAuth with Google/Facebook)
- Database schema (budgets, categories, accounts, transactions, budget_entries)
- Admin dashboard (user management, audit logs, session management)
- Basic infrastructure (Go API + SvelteKit frontend + PostgreSQL + Redis)

### ❌ Missing
- All core budget functionality (no budget endpoints or UI exist yet)
- The dashboard currently shows only static mockup data

---

## 🎯 Simplified User Journey

The workflow follows a natural progression that avoids overwhelming users:

### Phase 1: Get Started (First-Time Setup)
1. **Welcome & Quick Start** - Guide users through creating their first budget
2. **Add Accounts** - Link bank accounts, credit cards, cash
3. **Simple Budget Setup** - Create income & expenses with smart defaults

### Phase 2: Daily Use
4. **Dashboard Overview** - See current financial snapshot at a glance
5. **Track Transactions** - Quick entry or import from banks
6. **Monitor Progress** - Visual feedback on budget health

### Phase 3: Insights & Optimization
7. **Reports & Trends** - Understand spending patterns
8. **Adjust & Optimize** - Refine budget based on actual spending

---

## 🛠 Implementation Plan

### Sprint 1: Foundation - Accounts & Categories (Most Important)
**Goal:** Let users set up their financial accounts and categorize spending

#### Backend API (Go)
- [ ] `GET /api/accounts` - List all accounts
- [ ] `POST /api/accounts` - Create new account
- [ ] `PUT /api/accounts/:id` - Update account
- [ ] `DELETE /api/accounts/:id` - Delete account
- [ ] `GET /api/categories` - List categories (with hierarchy support)
- [ ] `POST /api/categories` - Create custom category
- [ ] Seed default categories (common expenses: groceries, rent, utilities, etc.)

#### Frontend UI (SvelteKit)
- [ ] **Onboarding flow** - Welcome screen → Add first account → Choose categories
- [ ] **Accounts page** - Card-style account display (like current dashboard mockup)
  - Add/edit/delete accounts
  - Show current balance
  - Account type icons (checking, savings, credit card, cash)
- [ ] **Categories page** - Visual category manager
  - Pre-populated common categories
  - Color & icon picker for personalization
  - Parent/child category support

**Why First:** Without accounts and categories, users can't organize anything. This is the foundation.

---

### Sprint 2: Budget Planning - The Core Experience
**Goal:** Simple budget creation that doesn't require finance expertise

#### Backend API (Go)
- [ ] `GET /api/budgets` - List user's budgets
- [ ] `POST /api/budgets` - Create new budget
- [ ] `GET /api/budgets/:id` - Get budget details with all entries
- [ ] `GET /api/budgets/:id/entries` - List budget entries
- [ ] `POST /api/budgets/:id/entries` - Add budget entry (income/expense)
- [ ] `PUT /api/budgets/:id/entries/:entry_id` - Update entry
- [ ] `DELETE /api/budgets/:id/entries/:entry_id` - Delete entry
- [ ] `GET /api/budgets/:id/projection` - Calculate projected cash flow (30/60/90 days)

#### Frontend UI (SvelteKit)
- [ ] **Budget Setup Wizard** - Step-by-step guided experience:
  1. "What's your monthly income?" (simple input)
  2. "What are your fixed expenses?" (rent, utilities, subscriptions)
  3. "What do you spend on variables?" (groceries, dining, entertainment)
  4. "Review your budget" (show income vs expenses with visual balance)

- [ ] **Budget Builder Page** - For advanced users who skip the wizard
  - Two-column layout: Income (left) | Expenses (right)
  - Quick-add entry buttons by frequency (monthly, weekly, one-time)
  - Real-time balance calculation: `Income - Expenses = $X left over`
  - Color-coded: Green (surplus), Yellow (tight), Red (deficit)

- [ ] **Budget Entry Form** - Smart, context-aware form
  - Category dropdown (filtered by income/expense)
  - Amount input with currency formatting
  - Frequency selector with clear examples:
    - Monthly: "Like rent or salary"
    - Weekly: "Like groceries or gas"
    - Fortnightly: "Every two weeks"
    - Once-off: "One-time expense or income"
  - Date picker (when does this start?)
  - Optional: End date (for temporary expenses)

**Why Second:** This is the heart of the application - planning finances. Must be intuitive and non-intimidating.

---

### Sprint 3: Transaction Tracking - Reality Check
**Goal:** Enter or import actual transactions to compare against budget

#### Backend API (Go)
- [ ] `GET /api/transactions` - List transactions (with filters: date range, account, category)
- [ ] `POST /api/transactions` - Create transaction
- [ ] `PUT /api/transactions/:id` - Update transaction
- [ ] `DELETE /api/transactions/:id` - Delete transaction
- [ ] `POST /api/transactions/:id/categorize` - Assign category
- [ ] `GET /api/transactions/unmatched` - Get unmatched transactions (for review)

#### Frontend UI (SvelteKit)
- [ ] **Transaction List Page** - Clean, scannable list
  - Filters: Date range, account, category, amount range
  - Visual indicators: Income (green), Expense (red)
  - Quick actions: Categorize, edit, delete
  - Bulk select for categorization

- [ ] **Quick Add Transaction** - Fast entry modal/drawer
  - Amount (with +/- for income/expense)
  - Category dropdown
  - Account selector
  - Date (defaults to today)
  - Optional: Description/notes

- [ ] **Transaction Import** (CSV upload - simple version)
  - Upload CSV from bank export
  - Map columns (Date, Amount, Description)
  - Preview before import
  - Auto-categorize based on description matching

**Why Third:** Users need to see how their actual spending compares to their plan. This creates accountability.

---

### Sprint 4: Dashboard - The Financial Snapshot
**Goal:** Replace static mockup with real data insights

#### Backend API (Go)
- [ ] `GET /api/dashboard/summary` - Overview metrics:
  - Total balance across all accounts
  - This month: Income vs Expenses vs Budget
  - Upcoming bills (next 7/30 days)
  - Budget health score (0-100)

- [ ] `GET /api/dashboard/recent-activity` - Recent transactions (last 10)
- [ ] `GET /api/dashboard/spending-by-category` - For charts

#### Frontend UI (SvelteKit)
- [ ] **Dashboard redesign** with real data:
  - **Top**: Total net worth + month-over-month change
  - **Budget Progress Bar**: Visual of income vs expenses this month
  - **Upcoming Bills**: List of budget entries due soon
  - **Recent Transactions**: Last 10 transactions
  - **Spending Breakdown**: Pie/bar chart by category
  - **Quick Actions**: Add transaction, view budget, add account

**Why Fourth:** Dashboard needs real data to be useful. Build it after core features exist.

---

### Sprint 5: Intelligent Matching - Automation
**Goal:** Auto-link transactions to budget entries to save time

#### Backend API (Go)
- [ ] Transaction matching algorithm:
  - Match by description patterns (Netflix → "Netflix subscription" budget entry)
  - Match by amount tolerance (±$1-2 for slight variations)
  - Match by category + frequency + timing
  - Assign confidence levels: `auto_high`, `auto_low`, `unmatched`

- [ ] `GET /api/transactions/suggested-matches/:id` - Get match suggestions
- [ ] `POST /api/transactions/:id/link` - Link transaction to budget entry
- [ ] `POST /api/budget-entries/:id/matching-rules` - Define matching rules (JSONB)

#### Frontend UI (SvelteKit)
- [ ] **Matching Review Page** - For unmatched/low-confidence transactions
  - Show suggested matches with confidence %
  - One-click approve or manual search
  - "Teach" mode: Link + create matching rule for future

- [ ] **Budget Entry Matching Rules** - Optional advanced feature
  - Define patterns: "Contains: AMZN" → Amazon Shopping category
  - Amount range: $9.99-$10.99 → Netflix subscription

**Why Fifth:** Automation reduces manual work, but only valuable after users have budgets and transactions.

---

### Sprint 6: Reports & Insights - Learn & Improve
**Goal:** Help users understand their financial behavior

#### Backend API (Go)
- [ ] `GET /api/reports/spending-trends` - Month-over-month spending by category
- [ ] `GET /api/reports/budget-variance` - Budget vs actual analysis
- [ ] `GET /api/reports/cash-flow-projection` - Future balance forecast (based on budget entries)
- [ ] `GET /api/reports/top-expenses` - Biggest spending categories

#### Frontend UI (SvelteKit)
- [ ] **Reports Page** - Insight cards:
  - **Spending Trends**: Line chart showing category spending over time
  - **Budget Performance**: Where you're over/under budget
  - **Cash Flow Forecast**: Projected balance for next 90 days
  - **Expense Breakdown**: What's eating your budget?

- [ ] **Export Options**: PDF report, CSV data

**Why Last:** Reports require historical data. Build after users have been tracking for a while.

---

## 🎨 UX/UI Principles for Simplicity

### Design Philosophy
1. **Progressive Disclosure** - Show simple options first, advanced features on demand
2. **Smart Defaults** - Pre-fill common values (today's date, primary account, etc.)
3. **Visual Feedback** - Use colors/icons to communicate budget health at a glance
4. **Contextual Help** - Inline tooltips and examples (not separate help docs)
5. **Mobile-First** - Budget tracking happens on-the-go

### Key UI Patterns
- **Wizard for beginners** - Step-by-step guided setup
- **Quick-add everywhere** - Floating action button for fast transaction entry
- **Drag-and-drop** - For categorizing transactions
- **Inline editing** - Click to edit, no separate forms
- **Real-time calculations** - Show balance/surplus as users type

### Color System for Budget Health
- 🟢 **Green**: Surplus (income > expenses)
- 🟡 **Yellow**: Tight (within 5% of budget)
- 🔴 **Red**: Deficit (overspending)
- 🔵 **Blue**: Neutral/informational

---

## 📦 Optional Future Enhancements (Post-MVP)

### Nice-to-Haves (Later)
- [ ] Bank API integration (Plaid/Yodlee) for automatic transaction sync
- [ ] Bill reminders & notifications
- [ ] Shared budgets (household/partner budgeting)
- [ ] Goals & savings tracking ("Save $5000 for vacation")
- [ ] Debt payoff calculator
- [ ] Mobile app (React Native / Flutter)
- [ ] AI-powered insights ("You spent 30% more on dining this month")
- [ ] Recurring transaction templates
- [ ] Budget templates (starter budgets for common situations)

---

## 🚀 Getting Started: Recommended First Steps

1. **Sprint 1 Week 1**: Implement account management (backend + frontend)
2. **Sprint 1 Week 2**: Build category system with smart defaults
3. **Sprint 2 Week 1**: Create budget planning API + simple wizard UI
4. **Sprint 2 Week 2**: Budget builder page with entry management
5. **Sprint 3**: Transaction tracking (manual entry first, CSV import second)

**Estimated Timeline:** 6-8 weeks for MVP (Sprints 1-4), with Sprints 5-6 as polish/automation.

---

## 💡 Key Success Metrics

Track these to know if the app is helping users:
- Time to create first budget: < 5 minutes
- Daily active users (checking dashboard)
- Budget adherence rate (staying within planned amounts)
- Transaction categorization accuracy
- User retention (coming back month over month)

---

## Notes

This plan prioritizes **simplicity over features**. The goal is to make budgeting feel manageable, not overwhelming. Each sprint delivers value users can immediately use, building confidence progressively.
