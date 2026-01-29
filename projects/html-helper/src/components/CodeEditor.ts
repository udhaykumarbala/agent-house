import { EditorView, basicSetup } from 'codemirror';
import { html } from '@codemirror/lang-html';
import { oneDark } from '@codemirror/theme-one-dark';

export class CodeEditor {
  private container: HTMLElement;
  private editor: EditorView | null = null;
  private onChange: ((code: string) => void) | null = null;

  constructor(container: HTMLElement) {
    this.container = container;
  }

  initialize(initialCode: string = '', onChange: (code: string) => void): void {
    this.onChange = onChange;

    this.editor = new EditorView({
      doc: initialCode,
      extensions: [
        basicSetup,
        html(),
        oneDark,
        EditorView.updateListener.of((update) => {
          if (update.docChanged && this.onChange) {
            this.onChange(update.state.doc.toString());
          }
        }),
      ],
      parent: this.container,
    });

    // Style the editor container
    this.container.style.height = '100%';
    this.container.style.overflow = 'auto';
  }

  getCode(): string {
    return this.editor?.state.doc.toString() || '';
  }

  setCode(code: string): void {
    if (this.editor) {
      this.editor.dispatch({
        changes: {
          from: 0,
          to: this.editor.state.doc.length,
          insert: code,
        },
      });
    }
  }

  destroy(): void {
    if (this.editor) {
      this.editor.destroy();
      this.editor = null;
    }
  }

  focus(): void {
    if (this.editor) {
      this.editor.focus();
    }
  }
}
