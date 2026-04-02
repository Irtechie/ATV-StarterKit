import * as vscode from 'vscode';

export interface Layer3State {
    loadedExtensions: string[];
    copilotActive: boolean;
    copilotChatActive: boolean;
}

export function getLayer3State(): Layer3State {
    return {
        loadedExtensions: vscode.extensions.all
            .filter(e => e.isActive)
            .map(e => e.id),
        copilotActive:
            vscode.extensions.getExtension('github.copilot')?.isActive ?? false,
        copilotChatActive:
            vscode.extensions.getExtension('github.copilot-chat')?.isActive ?? false,
    };
}
