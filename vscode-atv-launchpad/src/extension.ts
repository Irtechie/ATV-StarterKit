import * as vscode from 'vscode';
import * as path from 'path';
import { StatePoller } from './statePoller';
import { getLayer3State } from './layer3';

let panel: vscode.WebviewPanel | undefined;
let poller: StatePoller | undefined;

export function activate(context: vscode.ExtensionContext) {
    const openCmd = vscode.commands.registerCommand('atvLaunchpad.open', () => {
        if (panel) {
            panel.reveal(vscode.ViewColumn.Beside);
            return;
        }

        panel = vscode.window.createWebviewPanel(
            'atvLaunchpad',
            'ATV Launchpad',
            vscode.ViewColumn.Beside,
            {
                enableScripts: true,
                retainContextWhenHidden: true,
                localResourceRoots: [
                    vscode.Uri.file(path.join(context.extensionPath, 'src', 'webview'))
                ]
            }
        );

        const webviewPath = path.join(context.extensionPath, 'src', 'webview');
        const styleUri = panel.webview.asWebviewUri(
            vscode.Uri.file(path.join(webviewPath, 'style.css'))
        );
        const scriptUri = panel.webview.asWebviewUri(
            vscode.Uri.file(path.join(webviewPath, 'main.js'))
        );

        panel.webview.html = getWebviewContent(styleUri, scriptUri);

        // Start state polling
        const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
        if (workspaceRoot) {
            poller = new StatePoller(workspaceRoot);
            poller.onStateChange((state) => {
                const layer3 = getLayer3State();
                panel?.webview.postMessage({
                    type: 'stateUpdate',
                    state,
                    layer3
                });
            });
            poller.start(context);
        }

        // Handle messages from webview
        panel.webview.onDidReceiveMessage(
            (message) => {
                switch (message.type) {
                    case 'approveAction':
                        vscode.window.showInformationMessage(
                            `Action approved: ${message.command}`
                        );
                        break;
                    case 'refresh':
                        poller?.refresh();
                        break;
                }
            },
            undefined,
            context.subscriptions
        );

        panel.onDidDispose(() => {
            panel = undefined;
            poller?.dispose();
            poller = undefined;
        });
    });

    context.subscriptions.push(openCmd);

    // Auto-open if .atv/ exists
    const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
    if (workspaceRoot) {
        const atvPath = path.join(workspaceRoot, '.atv');
        vscode.workspace.fs.stat(vscode.Uri.file(atvPath)).then(
            () => {
                // .atv/ exists — show a status bar item
                const statusItem = vscode.window.createStatusBarItem(
                    vscode.StatusBarAlignment.Right, 100
                );
                statusItem.text = '$(dashboard) ATV';
                statusItem.tooltip = 'Open ATV Launchpad';
                statusItem.command = 'atvLaunchpad.open';
                statusItem.show();
                context.subscriptions.push(statusItem);
            },
            () => { /* .atv/ doesn't exist, do nothing */ }
        );
    }
}

export function deactivate() {
    poller?.dispose();
    panel?.dispose();
}

function getWebviewContent(styleUri: vscode.Uri, scriptUri: vscode.Uri): string {
    return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src ${styleUri}; script-src ${scriptUri};">
    <link href="${styleUri}" rel="stylesheet">
    <title>ATV Launchpad</title>
</head>
<body>
    <div id="app">
        <header class="header">
            <h1>⚡ ATV Launchpad ⚡</h1>
            <span class="subtitle">Live dashboard · event-driven</span>
        </header>
        <nav class="tab-bar">
            <button class="tab active" data-tab="memory">1: Memory</button>
            <button class="tab" data-tab="context">2: Context</button>
            <button class="tab" data-tab="health">3: Health</button>
            <button class="tab" data-tab="moves">4: Moves</button>
        </nav>
        <main id="content">
            <div class="panel" id="panel-memory">
                <p class="placeholder">Waiting for state...</p>
            </div>
            <div class="panel hidden" id="panel-context"></div>
            <div class="panel hidden" id="panel-health"></div>
            <div class="panel hidden" id="panel-moves"></div>
        </main>
        <footer class="footer">
            <span id="last-event">Last FS event: never</span>
            <span id="status-indicator">⚪ Waiting</span>
        </footer>
    </div>
    <script src="${scriptUri}"></script>
</body>
</html>`;
}
