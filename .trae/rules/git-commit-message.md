---
alwaysApply: true
scene: git_message
---

---
alwaysApply: true
scene: git_message
---

Generate commit messages using the following rules:

1. Use Conventional Commit types in **UPPERCASE** only:
   - FEAT
   - FIX
   - REFACTOR
   - PERF
   - TEST
   - DOCS
   - STYLE
   - CHORE
   - BUILD
   - CI
   - REVERT

2. The commit title format must always be:

   TYPE:: short summary

   Example:
   - FEAT:: add debt payment endpoint
   - FIX:: resolve nil pointer in debt filter
   - DOCS:: update swagger for debt APIs

3. The summary must:
   - Use imperative mood (add, fix, update, remove, refactor, etc.).
   - Start with a lowercase word after `::`.
   - Be concise (preferably under 72 characters).
   - Do not end with a period.

4. If additional details are needed, use a blank line after the title followed by bullet points beginning with `-`.

   Example:

   FEAT:: add debt payment endpoint

   - implement pay debt usecase
   - update repository to set paid_at automatically
   - add handler and swagger documentation

5. Avoid:
   - Emojis
   - Issue numbers unless explicitly requested
   - Personal pronouns (I, we, my)
   - Unnecessary filler words

6. Use present tense and describe what the commit introduces or changes, not what was done in the past.

7. Prefer one logical change per commit message. If multiple related changes are included, summarize them in the title and list the main changes as bullet points.