# 🐃 CARA
### *Cashflow Analysis, Resource Allocation, Budgeting, and Optimization*

> **CARA** is a personal finance web app that helps users understand, manage, and improve their financial life. It turns your own income, expenses, budgets, and goals into clear answers using transparent calculations — available on any device through your browser, with a local-first mobile app to follow.

Inspired by the resilience, strength, and reliability of the **carabao** — the national animal of the Philippines and a symbol of hard work and endurance — this application is built to help users carry and manage the weight of their financial responsibilities in a smarter, more sustainable way.

> 🚧 **Status:** CARA is in early planning and development. Everything below describes the intended design; features are not yet released.

---

## 📌 Table of Contents

- [Overview](#-overview)
- [Core Objectives](#-core-objectives)
- [Product Direction](#-product-direction)
- [Privacy & Security](#-privacy--security)
- [How CARA Calculates](#-how-cara-calculates)
- [Key Features](#-key-features)
- [Planned Features](#-planned-features--future-development)
- [Tech Stack](#️-tech-stack)
- [Project Structure](#-project-structure)
- [Design System & UX Guidelines](#-design-system--ux-guidelines)
- [Platform Support](#-platform-support)
- [Roadmap](#-roadmap)
- [Getting Started](#-getting-started)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🧭 Overview

Unlike traditional budgeting apps that only record expenses after the fact, **CARA turns your own numbers into clear answers** — how much you can safely spend today, how many months your emergency fund would last, and whether a purchase fits your budget and goals.

Every result is computed from your own income, expenses, budgets, and goals using simple, visible formulas. CARA shows the formula behind each number, so there is no black box.

Core capabilities:

- ✅ Expense & income tracking
- ✅ Quick Log — fast transaction entry
- ✅ Budgets with wants vs. needs categories
- ✅ "Can I afford this?" calculator with transparent, rule-based checks
- ✅ Goal management & emergency fund monitoring
- ✅ Monthly financial reports and cashflow forecasting from your own history
- ✅ Export — Excel, PDF, and charts, with optional delivery to your email
- ✅ Secure account that works on any device with a browser

---

## 🎯 Core Objectives

| Objective | Description |
|---|---|
| 💡 Financial Awareness | Help users understand where their money goes and why |
| 🧮 Transparent Insights | Give clear answers from simple formulas the user can see and verify |
| ⚖️ Wants vs. Needs | Help users separate essential and discretionary spending |
| 🏆 Goal Support | Guide users toward short-term and long-term financial goals |
| 🛡️ Financial Safety | Monitor the emergency fund and warn about risk |
| ⚡ Real-Time Decisions | Support spending decisions at the moment of purchase |
| 🔒 Trustworthy Security | Protect users' financial data with strong, proven practices |

---

## 🧭 Product Direction

CARA is being built **website first**, then extended to mobile.

| Stage | What it is |
|---|---|
| 🌐 **Website** *(first release)* | Cloud-based web app. Your data is stored securely on CARA's servers and available from any device with a browser. |
| 📱 **Mobile app** *(later)* | A local-first Android app with instant, offline Quick Log, which can optionally sync with your website account. |

#### Why Website First
- Nothing to install, so it is the fastest way to reach users and get feedback
- The server can generate reports, exports, and email summaries
- Forgotten passwords can be reset without losing data
- One codebase and one deployment to maintain while the product takes shape

#### Why Mobile Comes Later
- Logging a quick expense is fastest from a phone, even with weak signal
- A local-first mobile app gives users an offline option and more control over where data lives
- It is built on top of a website that already works

---

## 🔐 Privacy & Security

CARA is built on a foundational principle of modern financial systems:

> **Your financial data belongs to you — always.**
> CARA exists to assist, not to control or exploit it.

CARA's website stores your data in the cloud, so protecting it is a core responsibility.

#### 🛡️ How We Protect Your Data
- All traffic is encrypted in transit (HTTPS/TLS)
- Data is encrypted at rest in the database
- Passwords are hashed with a modern algorithm (argon2id or bcrypt) and never stored in plain text
- Sessions use secure, HTTP-only cookies
- Optional two-factor authentication
- Every request is checked so users can only access their own data
- Sign-in attempts are rate-limited to slow down attacks
- Financial details are kept out of application logs
- Regular, encrypted database backups

#### 👤 Your Control
- Export all your data at any time (CSV / JSON)
- Delete your account and all associated data permanently
- Collect only the information needed to run the app

#### 🔐 Our Commitments
- No selling of financial data
- No targeted advertising based on spending behavior
- No hidden trackers or profiling
- No dark patterns influencing spending decisions
- Clear, plain-language explanation of what is stored and why

#### ⚠️ Honest Trade-Offs
- **CARA's servers can read your data.** This is what lets CARA generate reports, exports, and email summaries, but it also means CARA must protect the database carefully. Users who want the data to stay on their own device should use the planned local-first mobile app.
- **An internet connection is required** to use the website.
- **Some trusted services handle data on CARA's behalf,** such as hosting and email delivery, and will be listed openly.
- **Exported files and emailed reports leave CARA's protection** once they are sent to another app, inbox, or drive.
- **No system is risk-free.** CARA follows security best practices and will notify users promptly if a breach ever affects their data.

#### 🇵🇭 Compliance
CARA is being designed with the Philippines' **Data Privacy Act of 2012 (RA 10173)** in mind, including consent, purpose limitation, data-subject rights, and security measures. Formal compliance steps (such as any National Privacy Commission requirements) will be confirmed before public launch.

---

## 🧮 How CARA Calculates

CARA's insights come from simple, explainable arithmetic on your own entries. Every figure can show its formula.

| Insight | Calculation |
|---|---|
| Net cashflow | Total income − total expenses for the period |
| Savings rate | (Income − expenses) ÷ income |
| Budget remaining | Budget limit − spent so far in the category |
| Safe to spend today | Remaining budget ÷ days left in the period |
| Budget alert | Warn at set thresholds (e.g., 80% and 100% of a limit) |
| Wants vs. needs split | Spending in "needs" categories vs. "wants" categories, as a % of expenses |
| Emergency fund coverage | Fund balance ÷ average monthly essential expenses = months covered |
| Goal projection | (Target − saved) ÷ monthly contribution = months to reach the goal |
| Spending run-rate | Spent so far ÷ days elapsed × days in the period |
| Cashflow forecast | Starting balance + scheduled income − scheduled bills, plus a trailing average of everyday spending |

> 💰 All amounts are stored as exact values (integer centavos), never floating-point numbers, so totals are always accurate.

#### 🛒 "Can I Afford This?" Check
Enter the price of an item and CARA runs a set of rules against your numbers:
1. Does it fit within the remaining budget for its category?
2. Would it reduce your emergency fund coverage below your chosen minimum (e.g., 3 months)?
3. How many days would it delay your nearest goal?
4. Does it push your projected end-of-month balance below zero?

The result is a clear pass / caution / not-now verdict with the reason for each rule, so the user always sees *why*.

---

## ✨ Key Features

### ⚡ Quick Log
Fast transaction entry designed for real life:
- Add income or an expense in a few taps, from anywhere in the app
- Smart suggestions based on your own history (recent categories and amounts)
- Keyboard-friendly on desktop, thumb-friendly on mobile browsers

---

### 💰 Income & Expense Management
- Multiple income sources
- Fixed and variable expenses
- Savings tracking
- Full financial overview

---

### 📊 Dashboard & Spending Tracker
- Daily, weekly, monthly analytics
- Category breakdowns and charts
- Wants vs. needs split, based on how you tag each category
- Safe-to-spend-today at a glance

---

### 🎯 Budgets
- Monthly budgets per category
- Remaining amount and safe-to-spend calculation
- Threshold alerts as you approach a limit

---

### 🛒 Affordability Check
- Manual price entry
- Rule-based pass / caution / not-now verdict
- Shows budget, emergency fund, and goal impact with the formula behind each

---

### 🎯 Goal Management
- Savings goals tracking
- Progress and projected completion date
- Adjustable monthly contribution

---

### 🛟 Emergency Fund Tracker
- Fund balance monitoring
- Months-of-coverage indicator
- Warnings when coverage drops below your target

---

### 🔔 Smart Alerts
- Budget threshold warnings
- Overspending warnings based on run-rate
- Upcoming bill reminders
- In-app alerts, with optional email notifications

---

### 💡 Money Tips
- Short financial tips shown when a rule is triggered (for example, a tip about building an emergency fund when coverage is low)
- Written content built into the app

---

### 📘 Monthly Reports
- Income vs. expenses
- Savings rate and net cashflow
- Category breakdowns and month-over-month changes
- Optional monthly report sent to your email

---

### 📤 Export & Share
- Download reports as Excel (.xlsx) and PDF, and charts as images
- Send a report to your own email with one click
- Optional save to your own Google Drive
- CSV / JSON export of all data

---

### 📅 Cashflow Timeline
- Bill tracking
- Income scheduling
- Balance forecasting from scheduled items and your own spending history

---

### 👤 Account & Security Settings
- Email or Google sign-in
- Optional two-factor authentication
- Data export and permanent account deletion
- Active session management

---

## 🗓️ Planned Features — Future Development

- **Local-first Android app** with instant offline Quick Log and optional sync with your website account
- Installable web app (PWA) for a home-screen shortcut
- Smart onboarding templates (Filipino users)
- Bill & subscription tracking improvements
- Credit / installment monitoring
- Automatic transaction capture from payment SMS on Android (GCash, Maya, BDO, BPI), processed on-device and subject to Google Play's SMS permission policy
- Shared family budgeting
- Optional end-to-end encrypted storage for users who want even CARA's servers to be unable to read their data
- Optional AI insights, **off by default**, which would send only the data needed for each request to an AI provider and would be clearly disclosed
- Word (.docx) report export

---

## 🏗️ Tech Stack

### Frontend (Web)
| Technology | Purpose |
|---|---|
| Next.js (React, TypeScript) | Web application and routing |
| Charting library | Spending and cashflow charts |
| Generated API client (from OpenAPI spec) | Type-safe calls to the backend |

### Backend
| Technology | Purpose |
|---|---|
| Go | REST API, business logic, and calculations |
| OpenAPI specification | Single source of truth for the API and generated types |
| Background jobs | Report generation, exports, and scheduled emails |
| Transactional email service | Alerts and report delivery |

### Database & Storage
| Technology | Purpose |
|---|---|
| PostgreSQL | Accounts, transactions, categories, budgets, goals |
| Database migrations | Versioned schema changes |
| Cloud object storage | Generated export files |

### Security
| Technology | Purpose |
|---|---|
| argon2id / bcrypt | Password hashing |
| Secure, HTTP-only cookie sessions | Authentication |
| TOTP two-factor authentication | Optional extra sign-in protection |
| Google OAuth | Optional Google sign-in and Drive delivery |

### Development Tools
| Technology | Purpose |
|---|---|
| Postman *or* Bruno | API testing and documentation |
| Docker *(where available)* or local PostgreSQL with pgAdmin 4 | Local development environment |
| Automated tests and CI | Verify calculations and API behavior on every change |

### Mobile *(later)*
| Technology | Purpose |
|---|---|
| Flutter *or* React Native *(final choice pending)* | Local-first Android app |
| SQLite with SQLCipher | Encrypted on-device database |
| Android Keystore / secure storage | Protects the database encryption key |

### Structure Rules
- **`app/` holds routing only.** Pages import from `features/` and render; no business logic inside.
- **Features do not import from each other.** Shared code moves up to `components/`, `lib/`, or `hooks/`.
- **Each feature exposes an `index.ts`,** so other code imports from the feature, not its internals.
- **Backend layers go one way:** handler → service → repository. Formulas live only in the service layer.
- **Money is integer centavos** end to end, formatted only at the display edge.
- **The OpenAPI spec comes first.** Web and mobile API clients are generated from it.
- **Mobile logic exists twice by design.** Formulas run on-device for offline use, so they live in one separate module and share test cases with the Go backend.
- **Offline entries use client-generated UUIDs,** so they never collide when synced.

---

## 🎨 Design System & UX Guidelines

CARA's design is **simple and direct**: a clean, calm fintech UI with an emerald green primary theme, minimal cognitive load, and a friendly tone. The fundamental job of the interface is to help people finish what they came to do with as little effort and confusion as possible.

Core principles:
- Show the numbers and the formulas, never a black box
- Responsive on phones, tablets, and desktops
- Accessible by default (keyboard navigation, screen-reader support, readable contrast)
- Follow familiar patterns, and save originality for what makes CARA unique

### 🧠 Laws of UX

Based on the [Laws of UX](https://lawsofux.com) collection.

#### Seeing and grouping
| Law | Meaning | How CARA applies it |
|---|---|---|
| Proximity | Things close together look related | Group a budget's name, amount, and bar together |
| Similarity | Things that look alike are grouped | Same style for all category chips |
| Common Region | Things in one boundary are grouped | Cards for budgets, goals, and transactions |
| Uniform Connectedness | Linked things look more related | Connect a goal to its progress bar |
| Prägnanz | People read shapes in the simplest way | Simple charts, no decoration |
| Aesthetic-Usability Effect | Attractive designs feel easier to use | Calm, polished visuals build trust |
| Von Restorff Effect | The odd one out is remembered | One accent color for the main action |

#### Memory and mental effort
| Law | Meaning | How CARA applies it |
|---|---|---|
| Miller's Law / Chunking | Memory holds about 7 ± 2 items | Short lists, grouped sections |
| Working Memory | Don't make users remember across screens | Show the amount and category on the confirm step |
| Cognitive Load | Mental effort should stay low | Defaults, few fields, plain language |
| Serial Position Effect | First and last items are remembered best | Put key actions first and last in menus |
| Peak-End Rule | Experiences are judged by the peak and the end | Make saving and "goal reached" moments feel good |
| Zeigarnik Effect | Unfinished tasks stay in mind | Show setup and goal progress |

#### Decisions and behavior
| Law | Meaning | How CARA applies it |
|---|---|---|
| Hick's Law | More choices, slower decisions | Quick Log shows only a few suggested categories first |
| Choice Overload | Too many options cause stress | Limit options per screen |
| Goal-Gradient Effect | Effort increases near a goal | Progress bars for goals and budgets |
| Flow | Balanced challenge and no interruptions | Keep Quick Log uninterrupted |

#### Interaction and expectations
| Law | Meaning | How CARA applies it |
|---|---|---|
| Fitts's Law | Big, close targets are faster to hit | Large Quick Log button in the thumb zone |
| Jakob's Law | Users expect your app to work like others | Use standard navigation, forms, and icons |
| Mental Model | Users bring their own idea of how it works | Use everyday words: "Safe to spend today" |
| Doherty Threshold | Responses under ~400 ms keep users engaged | Optimistic updates, skeleton loaders |
| Paradox of the Active User | People skip instructions | Learn by doing, hints in context |

#### Design principles
| Law | Meaning | How CARA applies it |
|---|---|---|
| Postel's Law | Be liberal in what you accept, conservative in what you output | Accept `1250`, `1,250`, or `₱1,250`; always display `₱1,250.00` |
| Tesler's Law | Complexity can't be removed, only moved | The system handles the math so the user doesn't |
| Occam's Razor | Prefer the simplest solution | Fewest fields and steps |
| Pareto Principle | ~80% of value comes from ~20% of features | Polish Quick Log and the dashboard first |
| Parkinson's Law | Tasks expand to fill the time given | Short forms get finished faster |

#### Jakob's Law in practice
- Logo top-left goes home; navigation at the top (desktop) or a bottom tab bar (mobile)
- Labels above inputs, primary button at the bottom, "Forgot password?" under the password field
- Standard icons: magnifier = search, gear = settings, trash = delete, three dots = more
- Transaction lists grouped by date, "+" to add, progress bars for budgets
- Standard gestures: pull to refresh, swipe back, tap outside to close
- If a pattern must change, make it clearly better and guide the user

### 🔤 Typography

| Rule | Guideline |
|---|---|
| Typefaces | Use 1 or 2 families; vary with weight and size |
| Type scale | Fixed sizes (for example 12, 14, 16, 20, 24, 32) reused everywhere |
| Body size | At least 16px |
| Line length | About 45–75 characters |
| Line height | About 1.4–1.6 for body text, tighter for headings |
| Contrast | At least 4.5:1 for normal text |
| Alignment | Left-align body text; avoid justified or long centered text |
| Letter spacing | Default for body; slightly loose for all-caps labels |
| All caps | Only for short labels |
| Numbers | Tabular figures (`font-variant-numeric: tabular-nums`), right-aligned, so amounts line up |

### 🌈 Color

- **60-30-10 rule:** 60% neutrals (backgrounds), 30% secondary (cards, navigation), 10% accent (key actions)
- **Palette:** emerald green primary, one neutral scale, and three status colors (success, warning, error)
- **Build shades** for each color (for example 50–900) and reuse them
- **Meaning stays consistent:** green = positive/success, amber = warning, red = error/danger, blue = information/links
- **Never rely on color alone:** pair color with a sign, icon, or label (for example "−₱500" or "Over budget" with an icon)
- **Contrast:** 4.5:1 for normal text, 3:1 for large text
- **Use red carefully.** Use amber for early warnings (matching the 80% alert) and reserve red for real overspending, so money never feels alarming
- **Wants vs. needs:** two distinct, non-alarming colors
- **Dark mode:** slightly desaturated colors, no pure black backgrounds

### ⬜ Shape & Spacing

**Nested radius rule:** inner radius = outer radius − padding. This keeps curves concentric so the gap looks even.

```css
.card {
  --radius: 24px;
  --pad: 8px;
  border-radius: var(--radius);
  padding: var(--pad);
}
.card .inner {
  border-radius: calc(var(--radius) - var(--pad));  /* 16px */
}
```

- Use a small radius scale (for example 4, 8, 12, 16, 24, full) and reuse it
- Same radius for similar elements (all buttons match, all cards match)
- Pills (`9999px`) for chips and badges, used consistently
- Use a spacing scale (4, 8, 12, 16, 24, 32) instead of arbitrary values

### 🎞️ Motion

Motion must serve a function: feedback, orientation, focus, continuity, or perceived speed.

| Rule | Guideline |
|---|---|
| Duration | About 150–300 ms for most UI; 100–150 ms for hover and toggles |
| Easing | Ease-out for entering, ease-in for leaving; avoid linear |
| Properties | Animate `transform` and `opacity` only |
| Reduced motion | Respect `prefers-reduced-motion` |
| Blocking | Never make users wait for an animation |
| Consistency | Same action, same timing and easing everywhere |

CARA-specific:
- Briefly animate changing numbers (like safe-to-spend), but never delay the final value
- Soft fill for goal and budget progress bars
- Calm alerts: no shaking or flashing around money
- Fade or slide in the "Can I afford this?" verdict

### 🧩 Components

| Component | Use it when | Avoid when |
|---|---|---|
| **FAB** | One primary action dominates (Quick Log "+") | Several equal actions |
| **FAB speed-dial** | 2–4 related actions (expense, income, transfer) | More than 4 options |
| **Bottom sheet** | Focused mobile tasks (category picker, Quick Log form) | Long, complex flows |
| **Modal / dialog** | Critical or destructive decisions | Routine tasks |
| **Three-dot menu** | Secondary per-row actions (edit, duplicate, delete) | The primary action |
| **Toast / snackbar** | Brief confirmation, with Undo ("Transaction saved") | Critical errors |
| **Banner** | Ongoing status ("Budget 85% used", "You're offline") | Temporary confirmations |
| **Skeleton** | Loading lists, cards, and charts | Waits under ~300 ms |
| **Spinner / progress bar** | Short unknown waits / measurable progress | Full-page loads |
| **Empty state** | No data yet; explain and offer the next step | Never leave a screen blank |
| **Tooltip** | Explaining an icon or formula on hover or focus | Essential info (touch can't hover) |
| **Tabs / segmented control** | Daily / Weekly / Monthly views | Sequential steps |
| **Chips** | Category filters, wants vs. needs tags | Long labels |
| **Swipe actions** | Edit or delete a row on mobile | As the only way to reach an action |

Quick rules:
- **Toast vs. banner vs. modal:** toast = "it worked", banner = "this is ongoing", modal = "you must decide"
- **Prefer Undo** over "Are you sure?" for most destructive actions
- **Toasts:** 4–6 seconds, at most 2–3 stacked, announced to screen readers (`aria-live`)
- **Bottom sheets:** drag handle, visible close option, dismissible by tapping outside
- **Skeletons:** match the real layout so nothing jumps when content loads

### ✨ Small UX Details

**Touch and input**
- Touch targets at least 44×44px with enough spacing
- Right keyboard for each field (`inputmode="decimal"` for amounts, `type="email"` for emails)
- Autofocus the amount field in Quick Log
- Never clear user input after an error
- Show/hide password toggle; support password managers

**Feedback and states**
- Every interactive element has hover, focus, active, disabled, and loading states
- Visible focus ring for keyboard users
- Disable the submit button while saving to prevent double submissions
- Optimistic updates: show the saved item immediately, roll back only on failure

**Forms**
- Labels above inputs, not placeholder-only
- Inline validation after leaving a field, with specific messages ("Enter an amount greater than 0")
- Sensible defaults: today's date, last-used category and account
- Format amounts on blur: `₱1,250.00`
- Mark optional fields when most are required

**Content and clarity**
- Specific button labels ("Save expense", not "Submit")
- Plain-language errors that say what happened and how to fix it
- Name what will be lost in destructive confirmations ("Delete 'Groceries' budget?")
- Show the formula next to calculated numbers

**Layout and flow**
- Primary action in the thumb zone on mobile
- Preserve scroll position and filters when going back
- Reserve space for images, charts, and data to avoid layout shift
- Sticky actions for long forms
- Respect system dark mode and reduced-motion settings

**Finance-specific**
- Right-align numbers and use tabular figures
- Always show the currency and use consistent decimals
- Pair positive and negative amounts with a sign or icon, not just color

---

## 📱 Platform Support

| Platform | Status |
|---|---|
| 🌐 Website | **MVP — Priority Release** |
| 🤖 Android (local-first app) | Phase 3 |
| 🍎 iOS | Planned |

---

## 🚀 Roadmap

### Phase 1 — Website MVP
- [ ] Accounts and secure sign-in
- [ ] Transactions (add, edit, delete) and categories
- [ ] Quick Log
- [ ] Budgets with safe-to-spend and threshold alerts
- [ ] Dashboard and charts
- [ ] Goal tracking
- [ ] Emergency fund tracker
- [ ] Data export and permanent account deletion

### Phase 2 — Insights & Reports
- [ ] "Can I afford this?" rule-based check
- [ ] Wants vs. needs breakdown
- [ ] Cashflow timeline and forecast
- [ ] Monthly reports
- [ ] Excel / PDF / chart exports
- [ ] Email delivery and optional Google Drive save
- [ ] Two-factor authentication

### Phase 3 — Mobile
- [ ] Choose mobile framework (Flutter or React Native)
- [ ] Local-first Android app with offline Quick Log
- [ ] Optional sync with the website account
- [ ] Installable web app (PWA)

### Phase 4 — Expansion
- [ ] iOS app
- [ ] Payment SMS capture on Android
- [ ] Shared family budgets
- [ ] Optional end-to-end encrypted storage
- [ ] Optional, off-by-default AI insights

---

## 🛠️ Getting Started

CARA is still in early development and there is no public build yet. Setup and run instructions will be added here once the first MVP milestone is ready.

---

## 🤝 Contributing

Contribution guidelines are coming soon. In the meantime, ideas and feedback are welcome through the repository's issue tracker.

---

## 💡 Why CARA?

Built for resilience — like the **carabao**.

- 🐃 **Strength** in financial discipline
- 🛡️ **Stability** in uncertain times
- 💪 **Long-term financial endurance**

---

## 📄 License

License information will be defined in a future release.

---

<p align="center">
  Made with ❤️ in the Philippines 🇵🇭
  <br/>
  <strong>Your money. Your data. Your future.</strong>
</p>
