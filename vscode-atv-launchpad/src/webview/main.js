// ATV Launchpad Webview — renders dashboard panels from state updates.
// Communicates with extension host via acquireVsCodeApi().

(function () {
    const vscode = acquireVsCodeApi();

    // Tab switching
    document.querySelectorAll('.tab').forEach(tab => {
        tab.addEventListener('click', () => {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.panel').forEach(p => p.classList.add('hidden'));
            tab.classList.add('active');
            const panelId = 'panel-' + tab.dataset.tab;
            document.getElementById(panelId)?.classList.remove('hidden');
        });
    });

    // Listen for state updates from extension
    window.addEventListener('message', event => {
        const message = event.data;
        if (message.type === 'stateUpdate') {
            renderState(message.state, message.layer3);
        }
    });

    function renderState(state, layer3) {
        if (!state) { return; }
        renderMemory(state);
        renderContext(state);
        renderHealth(state, layer3);
        renderMoves(state);
        updateFooter(state);
    }

    function renderMemory(state) {
        const el = document.getElementById('panel-memory');
        if (!el) { return; }
        let html = '<h3 class="section-title">Repo Memory Artifacts</h3>';
        html += renderArtifactList('Brainstorms', state.brainstorms);
        html += renderArtifactList('Plans', state.plans);
        html += renderArtifactList('Solutions', state.solutions);

        const snap = state.launchpadSnapshot || {};
        const repo = snap.repoState || {};
        if (repo.memoryFileCount > 0) {
            html += '<h3 class="section-title">Copilot Memory (' + repo.memoryFileCount + ')</h3>';
        } else {
            html += '<h3 class="section-title">Copilot Memory</h3>';
            html += '<p class="placeholder">No .copilot-memory/ files yet</p>';
        }
        el.innerHTML = html;
    }

    function renderArtifactList(title, artifacts) {
        const items = artifacts || [];
        let html = '<h4>' + title + ' (' + items.length + ')</h4>';
        if (items.length === 0) {
            html += '<p class="placeholder">(empty)</p>';
        } else {
            html += '<ul class="artifact-list">';
            items.forEach(a => {
                const age = formatAge(a.modTime);
                html += '<li>' + escapeHtml(a.name) + '<span class="artifact-age">' + age + '</span></li>';
            });
            html += '</ul>';
        }
        return html;
    }

    function renderContext(state) {
        const el = document.getElementById('panel-context');
        if (!el) { return; }
        const ctx = state.contextEstimate || {};
        const snap = state.launchpadSnapshot || {};
        const repo = snap.repoState || {};

        let html = '<h3 class="section-title">Context Estimate</h3>';
        html += '<p class="stat">Instruction bytes: <span class="stat-value">' + (ctx.totalInstructionBytes || 0) + '</span></p>';
        html += '<p class="stat">Estimated tokens: <span class="stat-value">~' + (ctx.estimatedTokens || 0) + '</span></p>';

        html += '<h3 class="section-title">Capability Matrix</h3>';
        html += '<p class="stat">';
        html += '<span class="stat-value">' + (repo.installedAgents || 0) + '</span> agents  ';
        html += '<span class="stat-value">' + (repo.installedSkills || 0) + '</span> skills  ';
        html += '<span class="stat-value">' + (repo.instructionFileCount || 0) + '</span> instructions  ';
        html += '<span class="stat-value">' + (repo.promptFileCount || 0) + '</span> prompts';
        html += '</p>';
        html += '<p class="stat">';
        html += '<span class="stat-value">' + (repo.mcpServerCount || 0) + '</span> MCP servers  ';
        html += '<span class="stat-value">' + (repo.extensionRecommendationCount || 0) + '</span> extensions  ';
        html += '<span class="stat-value">' + (repo.gstackSkillCount || 0) + '</span> gstack skills';
        html += '</p>';

        html += '<h3 class="section-title">Copilot Config</h3>';
        html += statusLine(repo.hasCopilotInstructions, 'copilot-instructions.md');
        html += statusLine(repo.hasSetupSteps, 'copilot-setup-steps.yml');
        html += statusLine(repo.hasMCPConfig, 'MCP servers (' + (repo.mcpServerCount || 0) + ' configured)');

        el.innerHTML = html;
    }

    function renderHealth(state, layer3) {
        const el = document.getElementById('panel-health');
        if (!el) { return; }
        const snap = state.launchpadSnapshot || {};
        const outcome = snap.outcomeSummary || {};

        let html = '<h3 class="section-title">Install Intelligence</h3>';
        if (snap.hasManifest) {
            html += statusLine(true, 'Manifest: ' + (snap.manifestPath || ''));
            html += '<p class="stat">Outcomes: ';
            html += '<span style="color:var(--atv-success)">' + (outcome.done || 0) + ' done</span>  ';
            html += '<span style="color:var(--atv-warn)">' + (outcome.warning || 0) + ' warn</span>  ';
            html += '<span style="color:var(--atv-fail)">' + (outcome.failed || 0) + ' fail</span>  ';
            html += '<span style="color:var(--atv-dim)">' + (outcome.skipped || 0) + ' skip</span>';
            html += '</p>';
        } else {
            html += '<p class="status-warn">No manifest yet. Run atv-installer init --guided</p>';
        }

        // Drift entries
        const drift = state.driftEntries || [];
        html += '<h3 class="section-title">Install Drift</h3>';
        if (drift.length === 0) {
            html += '<p class="status-ok" style="padding-left:0">No drift detected</p>';
        } else {
            drift.forEach(d => {
                const cls = d.status === 'missing' ? 'status-fail' : 'status-warn';
                html += '<p class="' + cls + '">' + escapeHtml(d.path) + ' — ' + d.status + '</p>';
            });
        }

        // Layer 3 info
        if (layer3) {
            html += '<h3 class="section-title">VS Code Extensions</h3>';
            html += statusLine(layer3.copilotActive, 'GitHub Copilot');
            html += statusLine(layer3.copilotChatActive, 'GitHub Copilot Chat');
        }

        el.innerHTML = html;
    }

    function renderMoves(state) {
        const el = document.getElementById('panel-moves');
        if (!el) { return; }
        const snap = state.launchpadSnapshot || {};
        const recs = snap.recommendations || [];
        const sdkRecs = state.sdkRecommendations || [];

        let html = '<h3 class="section-title">Recommended Next Moves</h3>';

        if (sdkRecs.length > 0) {
            sdkRecs.forEach((rec, i) => {
                html += '<div class="rec-item" data-index="' + i + '">';
                html += '<span class="rec-title">' + (i + 1) + '. ' + escapeHtml(rec.title) + '</span>';
                if (rec.proposedAction) {
                    html += ' <span style="color:var(--atv-accent)">[actionable]</span>';
                }
                html += '<p class="rec-reason">' + escapeHtml(rec.reason) + '</p>';
                html += '</div>';
            });
        }

        if (recs.length === 0 && sdkRecs.length === 0) {
            html += '<p class="status-ok" style="padding-left:0">All clear — no recommended actions.</p>';
        } else {
            const offset = sdkRecs.length;
            recs.forEach((rec, i) => {
                html += '<div class="rec-item" data-index="' + (offset + i) + '">';
                html += '<span class="rec-title">' + (offset + i + 1) + '. ' + escapeHtml(rec.title) + '</span>';
                html += '<span class="rec-priority">P' + rec.priority + '</span>';
                html += '<p class="rec-reason">' + escapeHtml(rec.reason) + '</p>';
                html += '</div>';
            });
        }

        el.innerHTML = html;
    }

    function updateFooter(state) {
        const lastEvent = document.getElementById('last-event');
        const indicator = document.getElementById('status-indicator');
        if (lastEvent && state.lastFSEvent) {
            const ago = timeSince(new Date(state.lastFSEvent));
            lastEvent.textContent = 'Last FS event: ' + ago + ' ago';
        }
        if (indicator) {
            indicator.textContent = state.sdkOnline ? '🟢 Online' : '⚪ Offline';
        }
    }

    // Helpers
    function statusLine(ok, label) {
        const cls = ok ? 'status-ok' : 'status-missing';
        return '<p class="' + cls + '">' + escapeHtml(label) + '</p>';
    }

    function escapeHtml(str) {
        if (!str) { return ''; }
        const div = document.createElement('div');
        div.appendChild(document.createTextNode(str));
        return div.innerHTML;
    }

    function formatAge(isoDate) {
        if (!isoDate) { return ''; }
        return timeSince(new Date(isoDate));
    }

    function timeSince(date) {
        const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
        if (seconds < 60) { return seconds + 's'; }
        const minutes = Math.floor(seconds / 60);
        if (minutes < 60) { return minutes + 'm'; }
        const hours = Math.floor(minutes / 60);
        if (hours < 24) { return hours + 'h'; }
        return Math.floor(hours / 24) + 'd';
    }
})();
