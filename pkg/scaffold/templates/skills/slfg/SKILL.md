---
name: slfg
description: Full autonomous engineering workflow using swarm mode for parallel execution
---

## Arguments
[feature description]

Swarm-enabled LFG. Run these steps in order, parallelizing where indicated. Do not stop between steps — complete every step through to the end.

## Sequential Phase

1. `/ce-plan $ARGUMENTS`
2. `/deepen-plan`
3. `/ce-work` — **Use swarm mode**: Make a Task list and launch an army of agent swarm subagents to build the plan

## Parallel Phase

After work completes, launch steps 5 and 6 as **parallel swarm agents** (both only need code to be written):

4. `/ce-review` — spawn as background Task agent
5. `/test-browser` — spawn as background Task agent

Wait for both to complete before continuing.

## Finalize Phase

6. `/resolve_todo_parallel` — resolve any findings from the review
7. `/feature-video` — record the final walkthrough and add to PR
8. `/ce-compound` — document the solution for future sessions
9. Output `<promise>DONE</promise>` after the video is in the PR and the solution has been compounded

Start with step 1 now.
