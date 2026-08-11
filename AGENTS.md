# AGENTS.md

These rules apply to every code change made in this repository.

- Do not write comments in code. Code must be self-explanatory through naming.
- Do not hardcode values. Extract every literal (string, number, color, spacing, URL, etc.) into a named variable or constant.
- Do not create a new UI component without first searching the codebase for an existing one that already serves the purpose. Always reuse an existing component when one exists. If no suitable component exists, ask before creating a new one.
- Every UI component and every UI addition must be spatially symmetrical relative to the overall UI. For any element, opposing spacing values (top margin vs. bottom margin, left padding vs. right padding, etc.) must be equal unless there is an explicit, stated reason for asymmetry.
