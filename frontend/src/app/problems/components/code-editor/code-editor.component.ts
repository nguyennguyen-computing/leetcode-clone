import { Component, Input, Output, EventEmitter, OnInit, OnDestroy, ViewChild, ElementRef, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzSpinModule } from 'ng-zorro-antd/spin';
import { NzMessageService } from 'ng-zorro-antd/message';
import { Problem } from '../../models/problem.models';

declare const monaco: any;

export interface CodeSubmission {
    code: string;
    language: string;
    problemId: number;
}

export interface CodeExecutionResult {
    success: boolean;
    output?: string;
    error?: string;
    runtime?: number;
    memory?: number;
    testCasesPassed?: number;
    totalTestCases?: number;
}

@Component({
    selector: 'app-code-editor',
    standalone: true,
    imports: [
        CommonModule,
        FormsModule,
        NzSelectModule,
        NzButtonModule,
        NzIconModule,
        NzSpinModule
    ],
    template: `
    <div class="code-editor-container h-full flex flex-col">
      <!-- Editor Header -->
      <div class="editor-header flex items-center justify-between p-4 border-b border-gray-200">
        <div class="flex items-center space-x-4">
          <nz-select 
            [(ngModel)]="selectedLanguage" 
            (ngModelChange)="onLanguageChange($event)"
            class="w-32">
            <nz-option 
              *ngFor="let lang of supportedLanguages" 
              [nzValue]="lang.value" 
              [nzLabel]="lang.label">
            </nz-option>
          </nz-select>
          
          <button 
            nz-button 
            nzType="default" 
            (click)="resetCode()"
            [nzLoading]="isResetting">
            <span nz-icon nzType="reload" nzTheme="outline"></span>
            Reset
          </button>
        </div>
        
        <div class="flex items-center space-x-2">
          <button 
            nz-button 
            nzType="default" 
            (click)="runCode()"
            [nzLoading]="isRunning"
            [disabled]="!currentCode.trim()">
            <span nz-icon nzType="play-circle" nzTheme="outline"></span>
            Run Code
          </button>
          
          <button 
            nz-button 
            nzType="primary" 
            (click)="submitCode()"
            [nzLoading]="isSubmitting"
            [disabled]="!currentCode.trim()">
            <span nz-icon nzType="check-circle" nzTheme="outline"></span>
            Submit
          </button>
        </div>
      </div>
      
      <!-- Monaco Editor Container -->
      <div class="editor-content flex-1 relative">
        <div 
          #editorContainer 
          class="w-full h-full"
          [class.opacity-50]="isRunning || isSubmitting">
        </div>
        
        <!-- Loading Overlay -->
        <div 
          *ngIf="isRunning || isSubmitting" 
          class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 z-10">
          <nz-spin 
            [nzSpinning]="true" 
            [nzTip]="isRunning ? 'Running code...' : 'Submitting solution...'">
          </nz-spin>
        </div>
      </div>
      
      <!-- Results Panel -->
      <div 
        *ngIf="lastResult" 
        class="results-panel border-t border-gray-200 p-4 bg-gray-50 max-h-48 overflow-y-auto">
        <div class="flex items-center space-x-2 mb-2">
          <span 
            nz-icon 
            [nzType]="lastResult.success ? 'check-circle' : 'close-circle'"
            [class.text-green-500]="lastResult.success"
            [class.text-red-500]="!lastResult.success">
          </span>
          <span class="font-medium">
            {{ lastResult.success ? 'Success' : 'Failed' }}
          </span>
          
          <div *ngIf="lastResult.testCasesPassed !== undefined" class="text-sm text-gray-600">
            {{ lastResult.testCasesPassed }}/{{ lastResult.totalTestCases }} test cases passed
          </div>
        </div>
        
        <div *ngIf="lastResult.runtime || lastResult.memory" class="text-sm text-gray-600 mb-2">
          <span *ngIf="lastResult.runtime">Runtime: {{ lastResult.runtime }}ms</span>
          <span *ngIf="lastResult.memory" class="ml-4">Memory: {{ lastResult.memory }}KB</span>
        </div>
        
        <div *ngIf="lastResult.output" class="bg-white p-2 rounded border text-sm font-mono">
          <div class="text-gray-600 mb-1">Output:</div>
          <pre class="whitespace-pre-wrap">{{ lastResult.output }}</pre>
        </div>
        
        <div *ngIf="lastResult.error" class="bg-red-50 p-2 rounded border border-red-200 text-sm font-mono">
          <div class="text-red-600 mb-1">Error:</div>
          <pre class="whitespace-pre-wrap text-red-700">{{ lastResult.error }}</pre>
        </div>
      </div>
    </div>
  `,
    styles: [`
    .code-editor-container {
      min-height: 500px;
    }
    
    .editor-content {
      background: #1e1e1e;
    }
    
    .results-panel {
      font-size: 13px;
    }
    
    pre {
      margin: 0;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    }
  `]
})
export class CodeEditorComponent implements OnInit, AfterViewInit, OnDestroy {
    @Input() problem: Problem | null = null;
    @Output() codeRun = new EventEmitter<CodeSubmission>();
    @Output() codeSubmit = new EventEmitter<CodeSubmission>();

    @ViewChild('editorContainer', { static: true }) editorContainer!: ElementRef;

    private editor: any;
    private monacoLoaded = false;

    selectedLanguage = 'javascript';
    currentCode = '';
    isRunning = false;
    isSubmitting = false;
    isResetting = false;
    lastResult: CodeExecutionResult | null = null;

    supportedLanguages = [
        { value: 'javascript', label: 'JavaScript' },
        { value: 'python', label: 'Python' },
        { value: 'java', label: 'Java' }
    ];

    // Default code templates
    private defaultTemplates: { [key: string]: string } = {
        javascript: `/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number[]}
 */
var twoSum = function(nums, target) {
    // Write your solution here
    
};`,
        python: `def two_sum(nums, target):
    """
    :type nums: List[int]
    :type target: int
    :rtype: List[int]
    """
    # Write your solution here
    pass`,
        java: `class Solution {
    public int[] twoSum(int[] nums, int target) {
        // Write your solution here
        
    }
}`
    };

    constructor(private message: NzMessageService) { }

    ngOnInit() {
        this.loadMonacoEditor();
    }

    ngAfterViewInit() {
        if (this.monacoLoaded) {
            this.initializeEditor();
        }
    }

    ngOnDestroy() {
        if (this.editor) {
            this.editor.dispose();
        }
    }

    private async loadMonacoEditor() {
        if (typeof monaco !== 'undefined') {
            this.monacoLoaded = true;
            if (this.editorContainer) {
                this.initializeEditor();
            }
            return;
        }

        // Load Monaco Editor
        const script = document.createElement('script');
        script.src = 'assets/monaco-editor/min/vs/loader.js';
        script.onload = () => {
            (window as any).require.config({
                paths: {
                    vs: 'assets/monaco-editor/min/vs'
                }
            });

            (window as any).require(['vs/editor/editor.main'], () => {
                this.monacoLoaded = true;
                if (this.editorContainer) {
                    this.initializeEditor();
                }
            });
        };

        document.head.appendChild(script);
    }

    private initializeEditor() {
        if (!this.monacoLoaded || !this.editorContainer) {
            return;
        }

        // Set initial code
        this.currentCode = this.getTemplateCode();

        // Create editor
        this.editor = monaco.editor.create(this.editorContainer.nativeElement, {
            value: this.currentCode,
            language: this.getMonacoLanguage(this.selectedLanguage),
            theme: 'vs-dark',
            automaticLayout: true,
            minimap: { enabled: false },
            scrollBeyondLastLine: false,
            fontSize: 14,
            lineNumbers: 'on',
            roundedSelection: false,
            scrollbar: {
                vertical: 'visible',
                horizontal: 'visible'
            },
            suggestOnTriggerCharacters: true,
            quickSuggestions: true,
            wordBasedSuggestions: 'matchingDocuments'
        });

        // Listen for code changes
        this.editor.onDidChangeModelContent(() => {
            this.currentCode = this.editor.getValue();
        });
    }

    private getTemplateCode(): string {
        if (this.problem?.templateCode?.[this.selectedLanguage]) {
            return this.problem.templateCode[this.selectedLanguage];
        }
        return this.defaultTemplates[this.selectedLanguage] || '';
    }

    private getMonacoLanguage(language: string): string {
        const languageMap: { [key: string]: string } = {
            javascript: 'javascript',
            python: 'python',
            java: 'java'
        };
        return languageMap[language] || 'javascript';
    }

    onLanguageChange(language: string) {
        this.selectedLanguage = language;

        if (this.editor) {
            // Update editor language
            const model = this.editor.getModel();
            monaco.editor.setModelLanguage(model, this.getMonacoLanguage(language));

            // Load template code for new language
            const templateCode = this.getTemplateCode();
            this.editor.setValue(templateCode);
            this.currentCode = templateCode;
        }
    }

    resetCode() {
        this.isResetting = true;

        setTimeout(() => {
            const templateCode = this.getTemplateCode();
            if (this.editor) {
                this.editor.setValue(templateCode);
                this.currentCode = templateCode;
            }
            this.lastResult = null;
            this.isResetting = false;
            this.message.success('Code reset to template');
        }, 300);
    }

    runCode() {
        if (!this.currentCode.trim()) {
            this.message.warning('Please write some code first');
            return;
        }

        this.isRunning = true;
        this.lastResult = null;

        const submission: CodeSubmission = {
            code: this.currentCode,
            language: this.selectedLanguage,
            problemId: this.problem?.id || 0
        };

        this.codeRun.emit(submission);
    }

    submitCode() {
        if (!this.currentCode.trim()) {
            this.message.warning('Please write some code first');
            return;
        }

        this.isSubmitting = true;
        this.lastResult = null;

        const submission: CodeSubmission = {
            code: this.currentCode,
            language: this.selectedLanguage,
            problemId: this.problem?.id || 0
        };

        this.codeSubmit.emit(submission);
    }

    // Method to be called by parent component with execution results
    setExecutionResult(result: CodeExecutionResult) {
        this.lastResult = result;
        this.isRunning = false;
        this.isSubmitting = false;

        if (result.success) {
            this.message.success('Code executed successfully!');
        } else {
            this.message.error('Code execution failed');
        }
    }

    // Method to handle loading states from parent
    setRunningState(isRunning: boolean) {
        this.isRunning = isRunning;
    }

    setSubmittingState(isSubmitting: boolean) {
        this.isSubmitting = isSubmitting;
    }
}