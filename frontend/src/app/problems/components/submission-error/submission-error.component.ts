import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NzAlertModule } from 'ng-zorro-antd/alert';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzCollapseModule } from 'ng-zorro-antd/collapse';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { NzDividerModule } from 'ng-zorro-antd/divider';
import { SubmissionStatus } from '../../models/submission.models';

interface ErrorSuggestion {
  title: string;
  description: string;
  code?: string;
}

@Component({
  selector: 'app-submission-error',
  standalone: true,
  imports: [
    CommonModule,
    NzAlertModule,
    NzIconModule,
    NzCollapseModule,
    NzTagModule,
    NzTypographyModule,
    NzDividerModule
  ],
  template: `
    <div class="submission-error" *ngIf="errorMessage || status !== 'Accepted'">
      <!-- Main Error Alert -->
      <nz-alert
        [nzType]="getAlertType()"
        [nzMessage]="getErrorTitle()"
        [nzDescription]="getErrorDescription()"
        nzShowIcon
        class="mb-4">
      </nz-alert>

      <!-- Detailed Error Message -->
      <div *ngIf="errorMessage" class="error-details mb-4">
        <nz-collapse>
          <nz-collapse-panel 
            nzHeader="Error Details" 
            [nzActive]="true"
            [nzExtra]="errorIconTemplate">
            
            <ng-template #errorIconTemplate>
              <span nz-icon nzType="bug" nzTheme="outline" class="text-red-500"></span>
            </ng-template>

            <div class="bg-red-50 border border-red-200 rounded-lg p-4">
              <pre class="text-sm font-mono text-red-700 whitespace-pre-wrap m-0">{{ errorMessage }}</pre>
            </div>
          </nz-collapse-panel>
        </nz-collapse>
      </div>

      <!-- Error-specific suggestions -->
      <div *ngIf="suggestions.length > 0" class="suggestions">
        <h4 class="text-sm font-medium text-gray-700 mb-3 flex items-center">
          <span nz-icon nzType="bulb" nzTheme="outline" class="mr-2 text-yellow-500"></span>
          Suggestions to Fix This Error
        </h4>
        
        <nz-collapse>
          <nz-collapse-panel 
            *ngFor="let suggestion of suggestions; let i = index"
            [nzHeader]="suggestion.title"
            [nzExtra]="suggestionIconTemplate">
            
            <ng-template #suggestionIconTemplate>
              <span nz-icon nzType="right-circle" nzTheme="outline" class="text-blue-500"></span>
            </ng-template>

            <div class="suggestion-content">
              <p class="text-sm text-gray-700 mb-3">{{ suggestion.description }}</p>
              
              <div *ngIf="suggestion.code" class="code-example">
                <div class="text-xs text-gray-500 mb-2">Example:</div>
                <div class="bg-gray-50 border rounded p-3">
                  <pre class="text-sm font-mono text-gray-800 m-0">{{ suggestion.code }}</pre>
                </div>
              </div>
            </div>
          </nz-collapse-panel>
        </nz-collapse>
      </div>

      <!-- Common debugging tips -->
      <div class="debugging-tips mt-4">
        <nz-divider nzText="Debugging Tips" nzOrientation="left"></nz-divider>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <div class="tip-card p-3 bg-blue-50 border border-blue-200 rounded">
            <div class="flex items-start">
              <span nz-icon nzType="check-circle" nzTheme="outline" class="text-blue-500 mr-2 mt-1"></span>
              <div>
                <div class="text-sm font-medium text-blue-800">Test with Examples</div>
                <div class="text-xs text-blue-600">Run your code with the provided examples first</div>
              </div>
            </div>
          </div>
          
          <div class="tip-card p-3 bg-green-50 border border-green-200 rounded">
            <div class="flex items-start">
              <span nz-icon nzType="eye" nzTheme="outline" class="text-green-500 mr-2 mt-1"></span>
              <div>
                <div class="text-sm font-medium text-green-800">Add Debug Prints</div>
                <div class="text-xs text-green-600">Use console.log() to trace your logic</div>
              </div>
            </div>
          </div>
          
          <div class="tip-card p-3 bg-yellow-50 border border-yellow-200 rounded">
            <div class="flex items-start">
              <span nz-icon nzType="clock-circle" nzTheme="outline" class="text-yellow-500 mr-2 mt-1"></span>
              <div>
                <div class="text-sm font-medium text-yellow-800">Check Edge Cases</div>
                <div class="text-xs text-yellow-600">Empty arrays, single elements, large inputs</div>
              </div>
            </div>
          </div>
          
          <div class="tip-card p-3 bg-purple-50 border border-purple-200 rounded">
            <div class="flex items-start">
              <span nz-icon nzType="code" nzTheme="outline" class="text-purple-500 mr-2 mt-1"></span>
              <div>
                <div class="text-sm font-medium text-purple-800">Review Constraints</div>
                <div class="text-xs text-purple-600">Ensure your solution handles all constraints</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .submission-error {
      max-width: 100%;
    }
    
    pre {
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      line-height: 1.4;
    }
    
    .tip-card {
      transition: transform 0.2s ease;
    }
    
    .tip-card:hover {
      transform: translateY(-1px);
    }
    
    .suggestion-content {
      font-size: 13px;
    }
    
    :host ::ng-deep .ant-collapse-header {
      align-items: center;
    }
  `]
})
export class SubmissionErrorComponent {
  @Input() status: SubmissionStatus = 'Internal Error';
  @Input() errorMessage?: string;
  @Input() testCasesPassed: number = 0;
  @Input() totalTestCases: number = 0;

  get suggestions(): ErrorSuggestion[] {
    return this.getErrorSuggestions();
  }

  getAlertType(): 'error' | 'warning' | 'info' {
    switch (this.status) {
      case 'Time Limit Exceeded':
      case 'Memory Limit Exceeded':
        return 'warning';
      case 'Wrong Answer':
        return 'info';
      default:
        return 'error';
    }
  }

  getErrorTitle(): string {
    const titleMap: { [key in SubmissionStatus]: string } = {
      'Wrong Answer': 'Wrong Answer',
      'Time Limit Exceeded': 'Time Limit Exceeded',
      'Memory Limit Exceeded': 'Memory Limit Exceeded',
      'Runtime Error': 'Runtime Error',
      'Compilation Error': 'Compilation Error',
      'Internal Error': 'Internal Error',
      'Accepted': 'Accepted'
    };
    return titleMap[this.status] || 'Submission Failed';
  }

  getErrorDescription(): string {
    const descriptionMap: { [key in SubmissionStatus]: string } = {
      'Wrong Answer': `Your solution produced incorrect output. ${this.testCasesPassed}/${this.totalTestCases} test cases passed.`,
      'Time Limit Exceeded': 'Your solution took too long to execute. Consider optimizing your algorithm.',
      'Memory Limit Exceeded': 'Your solution used too much memory. Try to reduce memory usage.',
      'Runtime Error': 'Your code crashed during execution. Check for null pointer exceptions, array bounds, etc.',
      'Compilation Error': 'Your code failed to compile. Check for syntax errors.',
      'Internal Error': 'An internal server error occurred. Please try again.',
      'Accepted': 'Your solution is correct!'
    };
    return descriptionMap[this.status] || 'An error occurred during submission.';
  }

  private getErrorSuggestions(): ErrorSuggestion[] {
    switch (this.status) {
      case 'Wrong Answer':
        return [
          {
            title: 'Check Your Logic',
            description: 'Review your algorithm step by step. Make sure you understand the problem requirements correctly.',
          },
          {
            title: 'Test with Examples',
            description: 'Run your code with the provided examples and trace through the execution manually.',
          },
          {
            title: 'Handle Edge Cases',
            description: 'Consider special cases like empty inputs, single elements, or boundary values.',
            code: `// Example: Check for empty array
if (nums.length === 0) {
    return [];
}`
          }
        ];

      case 'Time Limit Exceeded':
        return [
          {
            title: 'Optimize Time Complexity',
            description: 'Your algorithm might be too slow. Look for ways to reduce time complexity (e.g., O(n²) to O(n log n)).',
          },
          {
            title: 'Avoid Nested Loops',
            description: 'Multiple nested loops can cause timeout. Consider using hash maps or other data structures.',
            code: `// Instead of nested loops:
// for (let i = 0; i < n; i++) {
//     for (let j = 0; j < n; j++) { ... }
// }

// Use a hash map:
const map = new Map();
for (let i = 0; i < n; i++) {
    // O(1) lookup instead of O(n)
}`
          },
          {
            title: 'Check for Infinite Loops',
            description: 'Make sure your loops have proper termination conditions.',
          }
        ];

      case 'Memory Limit Exceeded':
        return [
          {
            title: 'Reduce Space Usage',
            description: 'Try to solve the problem with less memory. Reuse variables and avoid creating unnecessary data structures.',
          },
          {
            title: 'Use In-Place Operations',
            description: 'Modify the input array directly instead of creating new arrays when possible.',
            code: `// Instead of creating new array:
// const result = nums.map(x => x * 2);

// Modify in-place:
for (let i = 0; i < nums.length; i++) {
    nums[i] *= 2;
}`
          }
        ];

      case 'Runtime Error':
        return [
          {
            title: 'Check Array Bounds',
            description: 'Make sure you\'re not accessing array elements outside the valid range.',
            code: `// Always check bounds:
if (i >= 0 && i < arr.length) {
    return arr[i];
}`
          },
          {
            title: 'Handle Null/Undefined',
            description: 'Check for null or undefined values before using them.',
            code: `// Check for null:
if (node !== null && node.val !== undefined) {
    // Safe to use node.val
}`
          },
          {
            title: 'Avoid Division by Zero',
            description: 'Check denominators before division operations.',
            code: `// Safe division:
if (denominator !== 0) {
    result = numerator / denominator;
}`
          }
        ];

      case 'Compilation Error':
        return [
          {
            title: 'Check Syntax',
            description: 'Look for missing semicolons, brackets, or parentheses.',
          },
          {
            title: 'Variable Declarations',
            description: 'Make sure all variables are properly declared before use.',
            code: `// Correct variable declaration:
let result = 0;
const arr = [1, 2, 3];`
          },
          {
            title: 'Function Syntax',
            description: 'Ensure your function signature matches the expected format.',
          }
        ];

      default:
        return [
          {
            title: 'Try Again',
            description: 'This might be a temporary issue. Please try submitting your solution again.',
          }
        ];
    }
  }
}