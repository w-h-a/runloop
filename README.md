# runloop                                                  
                                                                                                                         
## Problem

When LLM calls are expensive, non-deterministic, and slow, that means lost progress on crashes, opaque failures, and no durable record of what happened or what it cost. Runloop separates *building* an agent from *running* agent workloads as managed infrastructure.

## Solution

Every state transition — LLM call, tool call, result, failure — is an immutable event appended to a durable log. Current state is derived by replaying events, never stored as a mutable snapshot. Runs are first-class resources with identity, lifecycle, and a query API.

## Architecture

```mermaid
graph TB
    Client[REST Client] -->|POST /runs, GET /runs/:id| API[HTTP Handlers]
    API --> RunSvc[Run Service]
    RunSvc --> Executor[Executor]
    RunSvc --> ES[Event Store]

    Executor --> LLM[LLM Client]
    Executor --> TD[Tool Dispatcher]
    Executor --> ES

    LLM -->|HTTP| OpenAI[OpenAI-compatible API]
    TD -->|HTTP| Tools[Tool Endpoints]
    ES -->|append events| DB[(SQLite)]
```

## Usage

Coming soon...