# AGENTS.md

These rules apply to every code change made in this repository.

- Do not write comments in code. Code must be self-explanatory through naming.
- Do not hardcode values. Extract every literal (string, number, color, spacing, URL, etc.) into a named variable or constant.
- Every component and piece of code must be scalable rather than tied to fixed dimensions or a single viewport. Derive sizing, spacing, typography, imagery, and layout from the available container or viewport dimensions, and verify behavior across smaller and larger widths and heights.
- Do not create a new UI component without first searching the codebase for an existing one that already serves the purpose. Always reuse an existing component when one exists. If no suitable component exists, ask before creating a new one.
- Every UI component and every UI addition must be spatially symmetrical relative to the overall UI. For any element, opposing spacing values (top margin vs. bottom margin, left padding vs. right padding, etc.) must be equal unless there is an explicit, stated reason for asymmetry.
- Commit messages must be concise: a single short line, never a multi-paragraph changelog.

# Definition of DONE

A task or a code change is not done unless it passes the below:

## Verification is mandatory

These rules are requirements, not suggestions. A change is complete if and only if it has been checked against every rule above.

- After every single change, explicitly verify that change against each rule above, one by one.
- A rule violation means the change fails verification. Do not present it as done. Rewrite it until it passes every rule.
- The only exception is a deliberate, explicitly stated reason given at the time of the change (e.g. the symmetry rule's own "unless there is an explicit, stated reason for asymmetry" clause). Silence is not an exception.

when finishing a task report each single rule and whether you have checked against it or not (in form of a table). If you haven't, then re-do the check. If the code doesn't fulfill the requirement check, re-do the entire task from scratch until all requirements pass.
