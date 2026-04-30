# Event-Sourced Bank Account (Go)

## Project Overview

Build a minimal event-sourced banking system to learn Event Sourcing, CQRS, Event-Driven Architecture, and Domain-Driven Design concepts.

## What It Does

Simple bank account system where you can:
- Create account
- Deposit money
- Withdraw money
- Check balance

**The Key Difference:** Instead of storing current balance, you store **events** and rebuild state from them.

## Core Concepts Demonstrated

### 1. Event Sourcing
- Store events: `AccountCreated`, `MoneyDeposited`, `MoneyWithdrawn`
- Rebuild account state by replaying events
- Never update, only append

### 2. CQRS (Command Query Responsibility Segregation)
- **Command side:** Handle commands (`CreateAccount`, `Deposit`, `Withdraw`)
- **Query side:** Read current balance from event stream
- Separate write model from read model

### 3. Event-Driven Architecture
- Commands generate events
- Events are the source of truth
- Other systems could subscribe to events

### 4. Domain-Driven Design
- **Aggregate:** `Account` (enforces business rules)
- **Commands:** Actions users want to take
- **Events:** Things that happened
- **Value Objects:** `Money`, `AccountID`

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL (store events)
- **Framework:** Standard library or Gin (simple REST API)
- **No external event sourcing libs** (learn by building it yourself)

## Project Structure

```
bank-account-es/
├── cmd/
│   └── api/
│       └── main.go          # Entry point
├── domain/
│   ├── account.go           # Account aggregate
│   ├── commands.go          # Command definitions
│   └── events.go            # Event definitions
├── eventstore/
│   └── postgres.go          # Event storage
├── handlers/
│   ├── command_handler.go   # Process commands
│   └── query_handler.go     # Answer queries
└── api/
    └── routes.go            # HTTP endpoints
```

## Implementation Phases

### Phase 1: Core Event Sourcing (Days 1-2)

**Step 1: Define Events**
- Create event structure with ID, AggregateID, Type, Data, Timestamp
- Define specific events: AccountCreated, MoneyDeposited, MoneyWithdrawn

**Step 2: Event Store (PostgreSQL)**
- Create events table with columns: id, aggregate_id, event_type, data (JSONB), timestamp
- Implement SaveEvents and GetEvents functions

**Step 3: Account Aggregate**
- Create Account struct with ID, Balance, and uncommitted Events
- Implement business logic methods (Deposit, Withdraw)
- Implement applyEvent method to update state from events
- Ensure business rules are enforced (no negative deposits, etc.)

**Step 4: Rebuild from Events**
- Implement LoadAccount function that replays all events
- Apply each event in order to rebuild current state
- No stored state, only events

### Phase 2: CQRS (Days 3-4)

**Command Side (Writes)**
- Define command structures: CreateAccountCommand, DepositCommand, WithdrawCommand
- Implement command handlers that:
  - Load account from events
  - Execute business logic
  - Save new events

**Query Side (Reads)**
- Implement query handlers for:
  - GetBalance: rebuild account and return balance
  - GetAccountHistory: return all events for an account
- Keep read and write logic completely separate

### Phase 3: REST API (Day 5)

**Endpoints:**
- POST /accounts - Create account
- POST /accounts/:id/deposit - Deposit money
- POST /accounts/:id/withdraw - Withdraw money
- GET /accounts/:id/balance - Get current balance
- GET /accounts/:id/events - View complete event history

**Implementation:**
- Use Gin or standard library for HTTP routing
- Map HTTP requests to commands
- Return appropriate responses and error codes

## What You'll Learn

✅ Event sourcing fundamentals
✅ CQRS pattern (separate reads/writes)
✅ Event-driven thinking
✅ Domain modeling (DDD)
✅ Go project structure
✅ PostgreSQL with Go
✅ Building REST APIs in Go

## Bonus Features (Optional)

### Projections
- Maintain a separate `account_balance` table
- Update automatically when events are saved
- Query this for fast balance lookups instead of replaying events

### Snapshots
- Store account state at intervals (every 100 events)
- Replay only events since last snapshot
- Optimizes performance for accounts with many events

### Event Versioning
- Handle schema changes in event definitions
- Support migrating old events to new formats
- Practice backwards compatibility

## Key Architectural Decisions

### Why Event Sourcing?
- Complete audit trail of all changes
- Can rebuild state at any point in time
- Natural fit for financial transactions
- Learn the pattern properly before using it in larger systems

### Why CQRS?
- Separates write complexity from read simplicity
- Allows different optimization strategies for reads vs writes
- Mirrors real-world eYou architecture

### Why PostgreSQL?
- JSONB support for flexible event data
- ACID transactions for consistency
- Familiar tool to most developers

## Testing Strategy

1. **Unit Tests:** Test aggregate business logic (deposit/withdraw rules)
2. **Integration Tests:** Test event store save/load
3. **API Tests:** Test HTTP endpoints end-to-end
4. **Event Replay Tests:** Verify rebuilding state from events works correctly

## Common Pitfalls to Avoid

- Don't update events after saving (they're immutable)
- Don't skip events when replaying (order matters)
- Don't put business logic in event handlers (keep it in aggregates)
- Don't forget to handle concurrency (two simultaneous withdrawals)

## Success Criteria

You've successfully completed this project when you can:
1. Create an account and see the AccountCreated event
2. Deposit money and see balance increase
3. Withdraw money and see balance decrease
4. View complete event history for an account
5. Delete the account_balance and rebuild it from events
6. Explain the difference between commands, events, and queries

## Next Steps After Completion

Once you finish this project:
1. Add more complex business rules (overdraft limits, transaction fees)
2. Build the mini-eYou clone using these same patterns
3. Study how eYou implements event sourcing in their actual codebase
4. Be ready to discuss this architecture in interviews