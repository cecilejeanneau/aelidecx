# Data model - Mermaid

```mermaid
erDiagram
    notes {
        INTEGER id PK "AUTOINCREMENT"
        TEXT title "NOT NULL, DEFAULT ''"
        TEXT content "NOT NULL, DEFAULT '' (HTML)"
        TEXT folder "NOT NULL, DEFAULT 'inbox'"
        DATETIME created_at "DEFAULT CURRENT_TIMESTAMP"
        DATETIME updated_at "DEFAULT CURRENT_TIMESTAMP"
    }

    notes_fts {
        INTEGER rowid FK "-> notes.id"
        TEXT title "indexed copy"
        TEXT content "indexed copy"
    }

    uploads {
        TEXT filename "timestamp_nano.ext"
        TEXT format ".png .jpg .jpeg .gif .webp .svg"
    }

    notes ||--|| notes_fts : "trigger sync"
    notes ||--o{ uploads : "referenced in HTML content"
```

## Allowed values for folder

| Value | Description |
|-------|-------------|
| `inbox` | Unsorted notes |
| `concepts` | Core knowledge |
| `practices` | Learned best practices |
| `ui-ux` | UI/UX captures and analysis |
| `postmortems` | Errors and lessons learned |
| `references` | Links, resources, tools |

## FTS5 triggers

```mermaid
flowchart LR
    INSERT[INSERT note] --> AI["Trigger notes_ai<br/>-> INSERT into FTS"]
    UPDATE[UPDATE note] --> AU["Trigger notes_au<br/>-> DELETE old + INSERT new into FTS"]
    DELETE[DELETE note] --> AD["Trigger notes_ad<br/>-> DELETE from FTS"]
```