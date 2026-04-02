import * as vscode from 'vscode';
import * as fs from 'fs';
import * as path from 'path';

export class StatePoller {
    private watcher: vscode.FileSystemWatcher | undefined;
    private callback: ((state: unknown) => void) | undefined;
    private stateFilePath: string;

    constructor(workspaceRoot: string) {
        this.stateFilePath = path.join(workspaceRoot, '.atv', 'launchpad-state.json');
    }

    onStateChange(cb: (state: unknown) => void): void {
        this.callback = cb;
    }

    start(context: vscode.ExtensionContext): void {
        this.readAndNotify();

        const dir = path.dirname(this.stateFilePath);
        this.watcher = vscode.workspace.createFileSystemWatcher(
            new vscode.RelativePattern(dir, 'launchpad-state.json')
        );

        this.watcher.onDidChange(() => this.readAndNotify());
        this.watcher.onDidCreate(() => this.readAndNotify());
        context.subscriptions.push(this.watcher);
    }

    refresh(): void {
        this.readAndNotify();
    }

    dispose(): void {
        this.watcher?.dispose();
    }

    private readAndNotify(): void {
        try {
            const content = fs.readFileSync(this.stateFilePath, 'utf8');
            const state: unknown = JSON.parse(content);
            this.callback?.(state);
        } catch {
            // File missing or invalid JSON — standalone fallback
        }
    }
}
