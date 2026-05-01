Questions: 
1. CQRS segregation is needed because we need atomic like operation on critical system like banking? Query reads state and can't modify

  1. Why CQRS?
  Not quite. Atomicity is a database concern — PostgreSQL
  transactions handle that. CQRS is about keeping two
  concerns with different shapes from fighting each other 
  in the same code path.
  Writes are complex: they need business rule validation,
  aggregate loading, concurrency checks. Reads are simple:
  just give me the current balance, maybe in a specific
  shape for a UI. If you put both through the same object,
  the write model's complexity bleeds into reads (slow,
  over-engineered) and read requirements distort the write
  model.

  CQRS says: don't fight it, just split them. Each side  can evolve and be optimized independently.


2. are state changes stored as events? meaning all events in the past are parsed in order to get the current state?

  2. Events as the only state — yes, that's correct, but
  that's Event Sourcing, not CQRS. They're two separate    
  patterns that happen to work very well together in this
  project. Important distinction:                          
                                 
  - CQRS — separate read path from write path
  - Event Sourcing — store events instead of current state;
   rebuild state by replaying them                         
                                                           
  You can use one without the other. This project uses     
  both.           


## Q&A — Concepts Clarified During Learning

  ### CQRS                                                 
   
  **Q: Is CQRS needed for atomic operations on critical    
  systems like banking?**
                                                           
  Not quite. Atomicity is a database concern — PostgreSQL  
  transactions handle that.
  CQRS is about keeping two concerns with different shapes 
  from fighting each other
  in the same code path. Writes are complex (validation,
  aggregate loading, concurrency
  checks). Reads are simple (just return current state).
  CQRS says: split them so each                            
  side can evolve and be optimized independently.
                                                           
  **Q: Are state changes stored as events, meaning all past
   events are replayed to get current state?**

  Yes — but that's **Event Sourcing**, not CQRS. They're   
  two separate patterns used together here:
  - **CQRS** — separate read path from write path          
  - **Event Sourcing** — store events instead of current
  state; rebuild state by replaying them

  **CQRS flow:**
  Command → CommandHandler → loads Aggregate → Aggregate
  enforces rules                                        
          → produces Event → Event saved to store

  Query → QueryHandler → replays Events from store →       
  returns current state
        (never modifies anything)                          
                  
  ---

  ### Handlers

  **Q: Are handlers like internal APIs that commands or    
  queries call to do actions?**
                                                           
  Slightly inverted. The handler *is* the entry point — it
  receives the command or query
  and orchestrates the response. It's the boundary between
  the outside world (HTTP) and                             
  your domain logic.
                                                           
  Handler orchestration flow:
  1. Receive the command
  2. Load the account (replay events from the store)       
  3. Call business logic on the aggregate
  4. Save the resulting new event to the store             
                                                           
  The handler does not contain business logic itself — it
  only calls it. Rules like                                
  "can't withdraw more than balance" live in the aggregate,
   not the handler.                                        
   
  ---                                                      
                  
  ### Aggregates

  **Q: Does the aggregate hold the business logic?**       
   
  Yes. The aggregate is the guardian of business rules —   
  nothing changes state without
  going through it. This ensures rules can't be bypassed by
   a second code path.                                     
  This is the core DDD idea: the aggregate is the single
  authority over its own state.                            
                  
  ---                                                      
                  
  ### Performance

  **Q: If state is derived by replaying events every time, 
  isn't that compute intensive as events pile up?**
                                                           
  Yes, that's a real problem. Two solutions exist:

  - **Snapshots** — every N events, save the current state.
   On load, start from the
    latest snapshot and replay only events after it.       
  Instead of 10,000 replays, at most N.                    
  - **Projections** — maintain a pre-computed table (e.g.
  `account_balance`) updated                               
    every time a new event is saved. Reads query the table
  directly — no replay at all.                             
   
  For this project: get naive replay working first.        
  Snapshots and projections are listed
  as bonus features in the spec — understand the problem   
  before optimizing it.

