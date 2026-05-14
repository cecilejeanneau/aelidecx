# Stack decision - Mermaid

## Decision tree

```mermaid
flowchart TD
    C["Project constraints"] --> H{Hosting?}

    H -->|"Paid cloud"| X1["Rejected (cost)"]
    H -->|"Free"| S{24/7 server?}

    S -->|"Personal PC"| X2["Not always on"]
    S -->|"Phone"| OK1["Always available"]

    OK1 --> RT{Android runtime?}
    RT --> TERMUX["Termux\n(single serious runtime)"]

    TERMUX --> LANG{Backend language?}
    LANG --> GO["Go\n10 MB RAM, single binary"]

    GO --> FW{Framework?}
    FW --> FIBER["Fiber\nExpress-like API"]

    FIBER --> DB{Database?}
    DB --> SQLITE["SQLite\nWAL + FTS5"]

    SQLITE --> FRONT{Frontend?}
    FRONT --> VANILLA["Vanilla JS\n0 KB, no build"]

    VANILLA --> EDIT{Rich editor?}
    EDIT --> TIPTAP["TipTap\nnative image paste"]

    TIPTAP --> FINAL["Final stack"]

    style X1 fill:#e94560,color:#fff
    style X2 fill:#e94560,color:#fff
    style FINAL fill:#4ecdc4,color:#111
```

## Detailed comparisons

### Backend language

```mermaid
quadrantChart
    title Backend language (phone)
    x-axis "Heavy RAM" --> "Light RAM"
    y-axis "Slow" --> "Fast"
    quadrant-1 "Ideal"
    quadrant-2 "Fast but heavy"
    quadrant-3 "Avoid"
    quadrant-4 "Light but slow"
    Go: [0.85, 0.9]
    Node.js: [0.5, 0.7]
    Python: [0.2, 0.3]
    Bun: [0.65, 0.85]
```

### Frontend framework

| Criteria | React | Vue 3 | Svelte | **Vanilla JS** |
|----------|-------|-------|--------|----------------|
| Bundle weight | 45 KB | 30 KB | 5 KB | **0 KB** |
| Build step | yes | yes | yes | **no** |
| HTML/JS separation | JSX | SFC | SFC | **full** |
| TipTap editor support | yes | yes | yes | **yes** |
| Learning curve | high | medium | low | **none** |

### Rich editor

| Criteria | TipTap | Quill | CKEditor | ProseMirror |
|----------|--------|-------|----------|-------------|
| Image paste | **native** | plugin | yes | low-level setup |
| CDN/ESM | **yes** | yes | partial | partial |
| Extensibility | **high** | medium | high | very high |
| Simple API | **yes** | yes | medium | complex |
| Free | **yes** | yes | freemium | yes |

## Final summary

```mermaid
mindmap
  root((my-notes))
    Runtime
      Termux
        Full Linux
        Free
        No root
    Backend
      Go
        10 MB RAM
        Single binary
      Fiber
        Express style
        Middleware support
    Data
      SQLite
        WAL mode
        FTS5 search
        Single backup file
    Frontend
      Vanilla JS
        app.js logic
        api.js HTTP
        editor.js TipTap
      TipTap
        Image paste
        Rich text
        CDN ESM
```