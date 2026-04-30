# Conecto Pipeline Engine

## 🧠 Purpose

This file captures the **current architecture and conventions** of Conecto, a dynamic, streaming, multi-connector pipeline engine written in Go.
It is intended to give **AI tools (Copilot, ChatGPT)** and developers full context without re-explaining the system.

---

# 🏗️ High-Level Architecture

```text
Connector (REST/DB/File/Kafka)
        ↓
ConnectorEngine (retry + pagination)
        ↓
Event stream (chan Event)
        ↓
Transformer (optional, batch-based)
        ↓
SinkEngine (batching + retry)
        ↓
Sink Adapter (Postgres, etc.)
```

---

# 📦 Core Data Model

## Event

```go
type Event struct {
	Payload []byte        // raw data (JSON or other)
	Cursor  core.Cursor   // checkpoint AFTER this event
}
```

---

## Cursor (Engine-level state)

```go
type Cursor map[string]string
```

* Generic and serializable
* Stored externally (DB / state store)
* Used for checkpointing + resume
* MUST NOT contain API-specific structs

---

## Batch

```go
type Batch struct {
	Events []Event
	Cursor Cursor
}
```

* Cursor represents **next position after batch**
* Used by engine to continue pagination

---

# 🔁 Cursor Model (CRITICAL)

## Separation

| Layer        | Type              |
| ------------ | ----------------- |
| Engine       | `Cursor` (map)    |
| Provider/API | `PageCursor`      |
| Connector    | maps between both |

---

## Rules

* Engine owns cursor lifecycle
* Connector translates cursor ↔ API pagination
* Provider only understands its own cursor type

---

# 🔌 Connector System

## Interface

```go
type Connector interface {
	Open(ctx context.Context, state Cursor) error
	FetchBatch(ctx context.Context, state Cursor) (Batch, error)
	Close() error
}
```

---

## Responsibilities

* Fetch external data (REST, DB, Kafka, etc.)
* Handle pagination
* Map:

  * Engine Cursor → API cursor
  * API cursor → Engine Cursor
* Return batches of events

---

## MUST guarantees

* Must eventually terminate
* Must advance cursor
* Must return empty batch OR nil cursor when done

---

# ⚙️ ConnectorEngine

## Responsibilities

* Retry fetches (with backoff)
* Maintain cursor progression
* Stream events into channel

---

## Loop behavior

```go
for {
	batch := FetchBatch()

	if len(batch.Events) == 0 || batch.Cursor == nil {
		return nil
	}

	for _, ev := range batch.Events {
		out <- ev
	}

	current = batch.Cursor
}
```

---

# 🔄 Transformer Layer

## Philosophy

* Optional
* Stateless
* Batch-oriented
* Works on `[]byte` by default

---

## Modes

### 1. Raw mode (default)

* No decoding
* Fast
* Works across all sinks

---

### 2. Structured mode (optional)

* Uses Codec (JSON, etc.)
* Decode → modify → encode

---

## Interface

```go
type Transformer interface {
	Transform(ctx context.Context, batch []Event) ([]Event, error)
}
```

---

# 🧩 JSON Field Selection

Instead of full decoding, use:

* gjson or similar
* JSON path extraction

Purpose:

* lightweight projection
* avoid schema coupling
* improve performance

---

# 🗄️ Sink System

## SinkEngine responsibilities

* Read from channel
* Batch events
* Apply transformer
* Retry batch processing

---

## Behavior

```go
for ev := range in {
	batch = append(batch, ev)

	if len(batch) >= BatchSize {
		process(batch)
	}
}

flush remaining
```

---

# 🔌 Adapter (Sink-specific logic)

## Responsibilities

* Decode payload (`[]byte`) → structured data
* Build DB queries
* Apply schema mapping

---

## Example

```go
type Adapter interface {
	Decode([]byte) (map[string]interface{}, error)
	BuildUpsertQuery(...)
}
```

---

## Rule

> Adapter is the ONLY place that understands schema

---

# 🔁 Retry Strategy

| Layer  | Scope       |
| ------ | ----------- |
| Source | API fetch   |
| Sink   | batch write |

---

## Rules

* Retry source on transient errors
* Retry sink per batch (NOT entire pipeline)
* Avoid global retries (duplicates risk)

---

# ⚠️ Concurrency Model

## Pipeline

```go
go ConnectorEngine → writes events
go SinkEngine      → consumes events
```

---

## Critical rules

* Source MUST close channel
* Sink MUST drain channel
* Pipeline MUST handle:

  * success
  * error
  * cancellation

---

## Completion signaling

Use explicit signal:

```go
doneCh <- struct{}{}
```

NOT only:

```go
close(doneCh)
```

---

# 🚨 Common Failure Modes

## Deadlocks

* Source blocked on channel (sink not consuming)
* Sink blocked in process()
* Channel never closed

---

## Infinite loops

* Cursor not advancing
* Missing termination condition

---

## Pipeline hangs

* No success signal
* Select waiting forever

---

# 🧠 Design Principles

1. **Byte-first pipeline**

   * `[]byte` is the transport layer

2. **Late decoding**

   * Decode only in sink (or explicitly in transformer)

3. **Separation of concerns**

   * Connector = fetch
   * Transformer = shape
   * Sink = persist

4. **Determinism**

   * Cursor-driven progression
   * Replay-safe

5. **Extensibility**

   * New connectors without engine changes
   * New sinks without transformer changes

---

# 🚀 Current Capabilities

* Streaming ingestion
* REST pagination with cursor
* Batch processing
* Retry (source + sink)
* Pluggable connectors
* Pluggable transformers
* Postgres sink with adapter

---

# 🚀 Future Directions

* DAG execution (fan-out / fan-in)
* Worker pools in sink
* Multi-sink pipelines
* Schema registry
* Exactly-once semantics
* Declarative pipeline config (already in progress)
* Auto schema → DB mapping

---

# 🧾 One-line Summary

> Conecto is a streaming, cursor-driven pipeline engine where connectors fetch data, transformers optionally shape it, and sinks own schema and persistence, with strict separation of concerns and backpressure-aware execution.

---
