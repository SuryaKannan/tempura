# Tempura

A lightweight durable execution system inspired by Temporal. Single Go binary + Python SDK. No cloud, no cluster, no infrastructure.

## What it is

Tempura is the SQLite to Temporal's Postgres — a deliberate simplification for a specific class of problem: **sequential, expensive, crash-sensitive pipelines on a single machine**.

One line in your Dockerfile. That's the entire deployment story.

## Architecture

```
Python Application (@batter/@fry decorators)
        │
        │ HTTP
        ▼
Tempura Go Server (REST API)
        │
        │ read/write
        ▼
SQLite (execution state + step outputs)
```

- **Go binary** — API server, state ledger, heartbeat watcher. Embeds SQLite. No external dependencies.
- **Python SDK** — decorator-based. Thin HTTP client. Drives execution, handles checkpointing transparently.
- **SQLite** — stores execution state and step inputs/outputs. Lives next to your app.

## Execution Model

Resume from last checkpoint — NOT Temporal's replay model.

```
NORMAL RUN:
@batter registers → @fry step1 ✅ → @fry step2 ✅ → @fry step3 ✅ → DONE

CRASH:
@batter registers → @fry step1 ✅ → @fry step2 💥 PROCESS DIES

RESUME (process restarted externally):
@batter → queries server → step1 ✅ SKIP → step2 ⏳ RESUME → step3 ⏳ → DONE
```

On resume, each `@fry` checks the Go server before executing. If a step has a stored output, return it immediately without re-executing. The process must be restarted externally (Docker restart policy, cron, etc).

## Decorators

- `@batter` — marks a full workflow. Registers execution with Go server on start. Queries state on resume.
- `@fry` — marks an individual activity. Checks stored state before executing. Stores result after completing.

```python
@batter
def training_pipeline():
    dataset = fetch_dataset()
    processed = preprocess(dataset)
    model = train(processed)
    metrics = evaluate(model)
    export_onnx(model)

@fry
def fetch_dataset(): ...

@fry
def preprocess(dataset): ...
```

## Output Storage

Step outputs are stored transparently:

- JSON-serialisable → stored directly in SQLite
- Non-serialisable (numpy arrays, dataframes, etc) → pickled to `~/.tempura/pickles/{execution_id}/{step_name}.pkl`, path stored in SQLite

Configurable via `TEMPURA_STORAGE_PATH`.

## SDK Safeguards

The SDK detects and loudly fails on common misuse rather than silently producing wrong results:

- **Step sequence fingerprinting** — records step order on first run, aborts on resume if sequence diverges
- **Loop detection** — warns if `@fry` is called inside a for/while loop
- **Input hashing** — warns if step inputs differ from original run on resume
- **Idempotency flag** — `@fry(idempotent=False)` forces explicit handling of side-effectful steps rather than silent re-execution

## Limitations

Tempura is intentionally constrained. If you need any of these, use Temporal:

- Sequential execution only — no parallel activities
- No branching logic
- No cross-workflow communication or signals
- No dynamic step generation (steps inside loops)
- No distributed workers
- Single machine only

Workflows must be linear and deterministic — same step sequence every run.

## Sweet Spot Use Cases

- Local ML training pipelines (fetch → preprocess → train → evaluate → export)
- LLM/eval pipelines where each step is an expensive API call
- ETL pipelines on a single machine
- Long-running background jobs that need crash recovery

Not a fit for: GitHub Actions (ephemeral runners), parallel workloads, event-driven workflows.

## Go Server Responsibilities

- REST API for execution registration, step checkpointing, state queries
- SQLite reads/writes
- Heartbeat watcher (background goroutine) — marks stale RUNNING executions as FAILED
- Graceful shutdown — flushes state before exit

## What This Is Not

- Not a Temporal replacement for production distributed systems
- Not a scheduler or cron system
- Not a workflow visualisation tool
- Not suitable for workflows requiring parallelism or branching
