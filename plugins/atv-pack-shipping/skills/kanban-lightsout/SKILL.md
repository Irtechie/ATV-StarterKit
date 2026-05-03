---
name: kanban-lightsout
description: "Compatibility forwarder for the renamed kanban-work skill. Use when the user says 'lightsout' or invokes the old kanban-lightsout skill name."
argument-hint: "[path to kanban manifest, or blank to find latest]"
---

# Kanban Lightsout Compatibility Forwarder

`kanban-lightsout` has been renamed to `kanban-work` because the workflow includes HITL pauses and is not fully hands-off.

Forward this invocation to `kanban-work` with the same arguments.

<input> #$ARGUMENTS </input>
