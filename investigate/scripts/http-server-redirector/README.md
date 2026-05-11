---

### 1. The "Standard Move" (301)

Use this to confirm that the agent correctly switches to **GET** and **drops the body** (as per RFC 7231).

```bash
curl -v -L \
     -X POST \
     -d '{"test": "this-should-disappear"}' \
     -H "X-Auth-Token: Secret123" \
     "http://localhost:8081?status=301"

```

* **Success at 8082:** Method: **GET**, Body: **[EMPTY]**.

---

### 2. The "Temporary Preservation" (307)

This is the primary fix for your ticket. Use this to verify that the agent **rewinds the body** and keeps the **POST** method.

```bash
curl -v -L \
     -X POST \
     -d '{"test": "this-must-stay"}' \
     -H "X-Auth-Token: Secret123" \
     "http://localhost:8081?status=307"

```

* **Success at 8082:** Method: **POST**, Body: **{"test": "this-must-stay"}**.

---

### 3. The "Permanent Preservation" (308)

Use this to ensure your Go logic treats 308 exactly like 307.

```bash
curl -v -L \
     -X POST \
     -d '{"test": "permanent-stay"}' \
     -H "X-Auth-Token: Secret123" \
     "http://localhost:8081?status=308"

```

* **Success at 8082:** Method: **POST**, Body: **{"test": "permanent-stay"}**.

---
