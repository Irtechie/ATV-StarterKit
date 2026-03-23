---
name: lfg
description: Full autonomous engineering workflow
---

## Arguments
[feature description]

CRITICAL: You MUST execute every step below IN ORDER. Do NOT skip any step. Do NOT jump ahead to coding or implementation. The plan phase (steps 2-3) MUST be completed and verified BEFORE any work begins. Violating this order produces bad output.

1. `/ce-plan $ARGUMENTS`

   GATE: STOP. Verify that `/ce-plan` produced a plan file in `docs/plans/`. If no plan file was created, run `/ce-plan $ARGUMENTS` again. Do NOT proceed to step 2 until a written plan exists.

2. `/deepen-plan`

   GATE: STOP. Confirm the plan has been deepened and updated. The plan file in `docs/plans/` should now contain additional detail. Do NOT proceed to step 3 without a deepened plan.

3. `/ce-work`

   GATE: STOP. Verify that implementation work was performed - files were created or modified beyond the plan. Do NOT proceed to step 4 if no code changes were made.

4. `/ce-review`

5. `/resolve_todo_parallel`

6. `/test-browser`

7. `/feature-video`

8. `/ce-compound`

9. Output `<promise>DONE</promise>` after the video is in the PR and the solution has been compounded.

Start with step 1 now. Remember: plan FIRST, then work. Never skip the plan.
