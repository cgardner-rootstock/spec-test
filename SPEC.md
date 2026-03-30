# Rootstock Inventory Agent — Service Specification

## Overview

A Go REST API service that provides a Claude-powered inventory assistant for Rootstock ERP. The agent answers natural language questions about inventory stock levels and reorder recommendations by querying Salesforce/Rootstock data in real time.

**Motivation:** Salesforce has not disclosed AI pricing. This service hedges against unknown costs by keeping inventory AI capabilities under Rootstock's control.

## Scope

### In Scope (v1)

- **Check stock levels** — Return on-hand quantities for items across all sites, with per-site breakdowns and a total
- **Recommend reorder quantities** — Surface MRP-generated supply recommendations and outstanding demand for items

### Out of Scope (v1)

- Write operations (creating POs, adjusting inventory, etc.)
- Modules beyond Inventory (Sales Orders, Purchasing, etc.)
- Streaming responses (SSE)
- End-user authentication (per-user OAuth)
- Frontend UI (browser plugin, Salesforce embedded)

## Architecture

```
┌──────────┐       ┌──────────────────┐       ┌─────────────┐
│  Client  │──────▶│  Inventory Agent │──────▶│  Claude API │
│ curl/POST│◀──────│  (Go / Heroku)   │◀──────│  (Anthropic)│
└──────────┘       └────────┬─────────┘       └─────────────┘
                            │
                            ▼
                   ┌─────────────────┐
                   │ Salesforce API  │
                   │ (Rootstock ERP) │
                   └─────────────────┘
```

### Components

1. **REST API** — Accepts user messages, manages conversation state, returns agent responses
2. **Claude Integration** — Sends conversation history and tool definitions to the Anthropic API; processes tool-use responses
3. **Salesforce Client** — Executes SOQL queries against Rootstock objects via the Salesforce REST API
4. **Tool Definitions** — Custom tools that Claude can invoke to retrieve inventory data

## API

### POST /chat

Sends a message to the inventory agent within a conversation.

**Request:**

```json
{
  "conversation_id": "optional-uuid",
  "message": "What's the stock level for item 1001?"
}
```

- If `conversation_id` is omitted, a new conversation is created
- If provided, the message is appended to the existing conversation history

**Response:**

```json
{
  "conversation_id": "uuid",
  "response": "Item 1001 has 450 units on hand across all sites:\n- Site A: 200 units\n- Site B: 250 units"
}
```

### GET /conversations/{id}

Returns the full message history for a conversation.

### DELETE /conversations/{id}

Deletes a conversation and its history.

## System Prompt

The following system prompt is sent with every Claude API request to establish the agent's behavior:

```
You are an inventory assistant for Rootstock ERP. You help users check stock levels and understand reorder recommendations.

Rules:
- ALWAYS call a tool before answering inventory questions. Never guess quantities or dates.
- When asked about stock levels, use get_stock_levels. Always show the total across all sites AND the per-site breakdown unless the user asks about a specific site.
- When asked about reorder recommendations, use get_reorder_recommendations. Explain what MRP is suggesting in plain language — include the order type, quantity, and date required.
- If a tool returns no results, tell the user clearly that no data was found and suggest they verify the item number.
- If the user's request is ambiguous (e.g., partial item number), ask for clarification rather than guessing.
- Do not answer questions outside of inventory. Politely redirect the user.
- Keep responses concise and formatted for readability.
```

## Agent Tool Loop

When processing a user message, the service sends the conversation history to Claude. Claude may respond with text, tool calls, or both. The service handles this in a loop:

1. Send messages to Claude
2. If the response contains tool calls, execute them against Salesforce and append the results
3. Send the updated messages back to Claude
4. Repeat until Claude responds with text only (no tool calls)

**Max iterations:** 5 tool rounds per user message. If the limit is reached, the service appends a message telling Claude to respond with what it has so far. This prevents runaway loops and controls API costs.

## Agent Tools

Claude will be provided with the following tool definitions. When the model decides it needs data to answer a question, it will invoke one or more of these tools. The service executes the corresponding SOQL query and returns the results to Claude for interpretation.

### get_stock_levels

Returns on-hand inventory for an item across all sites, with per-site breakdown and total.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| item_number | string | yes | The item number to look up |
| site | string | no | Optional site to filter to a single site |

**Data Sources:**

- `rstk__icitemsite__c` — aggregated on-hand qty by site (`icitemsite_qtyoh__c`)
- `rstk__icitem__c` — item master for item details and unit of measure (`icitem_invuom__c`)
- `rstk__sysite__c` — site name and description

**SOQL (example):**

```sql
SELECT
  icitemsite_icitem__r.Name,
  icitemsite_icitem__r.rstk__icitem_invuom__c,
  icitemsite_site__r.Name,
  icitemsite_site__r.rstk__sysite_descr__c,
  rstk__icitemsite_qtyoh__c,
  rstk__icitemsite_qtycons__c
FROM rstk__icitemsite__c
WHERE icitemsite_icitem__r.Name = :itemNumber
```

### get_reorder_recommendations

Returns MRP supply recommendations and outstanding demand for an item, giving Claude the data needed to explain what MRP suggests and why.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| item_number | string | yes | The item number to look up |
| site | string | no | Optional site filter |

**Data Sources:**

- `rstk__mrpsup__c` — MRP-generated supply recommendations (`mrpsup_qtyoutstdg__c`, `mrpsup_dtereqd__c`, `mrpsup_ordtype__c`)
- `rstk__icixr__c` — consolidated supply/demand index (`icixr_dtereqd__c`, `icixr_ordtype__c`)
- `rstk__icitem__c` — item master for safety stock (`icitem_sspolqty__c`), policy qty (`icitem_mrppolqty__c`), and planning policy (`icitem_mrpplanpol__c`)
- `rstk__icitemsite__c` — current on-hand for context

**SOQL (example):**

```sql
SELECT
  rstk__mrpsup_ordno__c,
  rstk__mrpsup_ordtype__c,
  rstk__mrpsup_qtyoutstdg__c,
  rstk__mrpsup_dtereqd__c,
  rstk__mrpsup_sts__c
FROM rstk__mrpsup__c
WHERE rstk__mrpsup_item__r.Name = :itemNumber
ORDER BY rstk__mrpsup_dtereqd__c ASC
```

## Authentication

### Salesforce

- **Method:** OAuth 2.0 Connected App (JWT Bearer or Username-Password flow)
- **User:** Dedicated integration user with read-only access to inventory objects
- **Permissions:** Read access to: `icitem__c`, `icitemsite__c`, `iclocitem__c`, `sysite__c`, `mrpsup__c`, `icixr__c`, `icreplenish__c`

### Claude API

- **Method:** Anthropic API key
- **Model:** `claude-sonnet-4-6` (default), `claude-opus-4-6` (for complex multi-step analysis)

### Client Authentication

- None for v1 (localhost/internal testing only)
- Future: API key or OAuth

## Conversation Management

- Conversations are stored in-memory for v1 (map of conversation ID to message history)
- Each conversation maintains the full Claude message history (user messages, assistant messages, tool calls, tool results)
- Conversations are created on first message and can be explicitly deleted
- Future: persistent storage (PostgreSQL on Heroku)

## Data Model

### Conversation

```go
type Conversation struct {
    ID        string
    Messages  []anthropic.Message
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Deployment

- **Platform:** Heroku
- **Dyno:** Single web dyno for v1
- **Config vars:**
  - `ANTHROPIC_API_KEY` — Claude API key
  - `SF_CLIENT_ID` — Salesforce connected app client ID
  - `SF_CLIENT_SECRET` — Salesforce connected app client secret
  - `SF_USERNAME` — Integration user username
  - `SF_LOGIN_URL` — Salesforce login URL (login.salesforce.com or test.salesforce.com)

## Future Considerations

- Write operations (create POs, adjust inventory)
- Additional modules (Sales Orders, Purchasing, Manufacturing)
- Per-user OAuth (user context instead of integration user)
- SSE streaming for real-time token delivery
- Persistent conversation storage
- Frontend UI (Salesforce embedded component or browser plugin)
- Usage/consumption trend analysis using `icusage__c` history
- Location-level detail from `iclocitem__c` (bin-level balances)
