import { Component, inject, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { Store } from '@ngrx/store';
import { Subject } from 'rxjs';
import { takeUntil, filter } from 'rxjs/operators';

import { NzLayoutModule } from 'ng-zorro-antd/layout';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzDividerModule } from 'ng-zorro-antd/divider';
import { NzSpinModule } from 'ng-zorro-antd/spin';
import { NzAlertModule } from 'ng-zorro-antd/alert';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { NzTabsModule } from 'ng-zorro-antd/tabs';
import { NzMessageService } from 'ng-zorro-antd/message';

import { CodeEditorComponent, CodeSubmission, CodeExecutionResult } from '../code-editor/code-editor.component';
import { SubmissionHistoryComponent } from '../submission-history/submission-history.component';
import { PerformanceMetricsComponent } from '../performance-metrics/performance-metrics.component';
import { SubmissionErrorComponent } from '../submission-error/submission-error.component';
import { 
  selectSelectedProblem, 
  selectProblemLoading, 
  selectProblemError,
  selectExecutionState
} from '../../store/problem.selectors';
import * as ProblemActions from '../../store/problem.actions';
import { Problem } from '../../models/problem.models';

@Component({
  selector: 'app-problem-solve',
  standalone: true,
  imports: [
    CommonModule,
    NzLayoutModule,
    NzCardModule,
    NzButtonModule,
    NzTagModule,
    NzIconModule,
    NzDividerModule,
    NzSpinModule,
    NzAlertModule,
    NzTypographyModule,
    NzTabsModule,
    CodeEditorComponent,
    SubmissionHistoryComponent,
    PerformanceMetricsComponent,
    SubmissionErrorComponent
  ],
  template: `
    <div class="min-h-screen bg-gray-50">
      <!-- Header -->
      <div class="bg-white border-b border-gray-200 px-4 py-3">
        <div class="max-w-7xl mx-auto flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <button 
              nz-button 
              nzType="text" 
              (click)="goBack()"
              class="flex items-center">
              <i nz-icon nzType="arrow-left" class="mr-2"></i>
              Back
            </button>
            
            <div *ngIf="problem$ | async as problem" class="flex items-center space-x-3">
              <h1 class="text-lg font-semibold text-gray-900">
                {{ problem.id }}. {{ problem.title }}
              </h1>
              <nz-tag 
                [nzColor]="getDifficultyColor(problem.difficulty)"
                class="font-medium">
                {{ problem.difficulty }}
              </nz-tag>
            </div>
          </div>
          
          <div class="flex items-center space-x-2">
            <button 
              nz-button 
              nzType="default"
              (click)="toggleDescription()"
              [nzType]="showDescription ? 'primary' : 'default'">
              <i nz-icon nzType="eye" class="mr-1"></i>
              Description
            </button>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div *ngIf="loading$ | async" class="flex items-center justify-center h-96">
        <nz-spin nzSize="large"></nz-spin>
        <p class="ml-4 text-gray-600">Loading problem...</p>
      </div>

      <!-- Error State -->
      <div *ngIf="error$ | async as error" class="max-w-7xl mx-auto px-4 py-8">
        <nz-alert
          nzType="error"
          [nzMessage]="error"
          nzShowIcon>
        </nz-alert>
      </div>

      <!-- Main Content -->
      <div *ngIf="problem$ | async as problem" class="max-w-7xl mx-auto h-screen flex">
        <!-- Problem Description Panel -->
        <div 
          *ngIf="showDescription" 
          class="w-1/2 bg-white border-r border-gray-200 overflow-y-auto">
          <div class="p-6">
            <!-- Problem Header -->
            <div class="mb-6">
              <h2 class="text-2xl font-bold text-gray-900 mb-3">
                {{ problem.title }}
              </h2>
              <div class="flex items-center space-x-4 mb-4">
                <nz-tag 
                  [nzColor]="getDifficultyColor(problem.difficulty)"
                  class="font-medium">
                  {{ problem.difficulty }}
                </nz-tag>
                <span class="text-gray-500 text-sm">
                  Acceptance Rate: {{ problem.acceptanceRate || 0 }}%
                </span>
              </div>
              
              <!-- Tags -->
              <div class="flex flex-wrap gap-2 mb-4">
                <nz-tag 
                  *ngFor="let tag of problem.tags" 
                  nzColor="blue"
                  class="text-xs">
                  {{ tag }}
                </nz-tag>
              </div>
            </div>

            <nz-divider></nz-divider>

            <!-- Problem Description -->
            <div class="mb-6">
              <h3 class="text-lg font-semibold text-gray-900 mb-3">Description</h3>
              <div 
                class="prose max-w-none text-gray-700 text-sm leading-relaxed"
                [innerHTML]="formatDescription(problem.description)">
              </div>
            </div>

            <!-- Examples -->
            <div *ngIf="problem.examples?.length" class="mb-6">
              <h3 class="text-lg font-semibold text-gray-900 mb-3">Examples</h3>
              <div class="space-y-4">
                <div 
                  *ngFor="let example of problem.examples; let i = index" 
                  class="bg-gray-50 rounded-lg p-4 border">
                  <h4 class="font-medium text-gray-900 mb-2 text-sm">Example {{ i + 1 }}:</h4>
                  
                  <div class="space-y-2">
                    <div>
                      <span class="font-medium text-gray-700 text-xs">Input:</span>
                      <pre class="mt-1 bg-white p-2 rounded border text-xs font-mono">{{ example.input }}</pre>
                    </div>
                    
                    <div>
                      <span class="font-medium text-gray-700 text-xs">Output:</span>
                      <pre class="mt-1 bg-white p-2 rounded border text-xs font-mono">{{ example.output }}</pre>
                    </div>
                    
                    <div *ngIf="example.explanation">
                      <span class="font-medium text-gray-700 text-xs">Explanation:</span>
                      <p class="mt-1 text-gray-600 text-xs">{{ example.explanation }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Constraints -->
            <div *ngIf="problem.constraints">
              <h3 class="text-lg font-semibold text-gray-900 mb-3">Constraints</h3>
              <div 
                class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 text-sm"
                [innerHTML]="formatConstraints(problem.constraints)">
              </div>
            </div>
          </div>
        </div>

        <!-- Code Editor and Results Panel -->
        <div [class]="showDescription ? 'w-1/2' : 'w-full'" class="flex flex-col">
          <!-- Code Editor -->
          <div class="flex-1">
            <app-code-editor
              #codeEditor
              [problem]="problem"
              (codeRun)="onCodeRun($event)"
              (codeSubmit)="onCodeSubmit($event)"
              class="h-full block">
            </app-code-editor>
          </div>

          <!-- Results and Submission Interface -->
          <div class="h-1/3 bg-white border-t border-gray-200 overflow-hidden">
            <nz-tabset nzSize="small" class="h-full">
              <!-- Real-time Results Tab -->
              <nz-tab nzTitle="Results">
                <div class="p-4 h-full overflow-y-auto">
                  <!-- Show execution results or errors -->
                  <div *ngIf="lastExecutionResult">
                    <div *ngIf="lastExecutionResult.success; else errorTemplate">
                      <!-- Success Result -->
                      <div class="bg-green-50 border border-green-200 rounded-lg p-4 mb-4">
                        <div class="flex items-center mb-2">
                          <i nz-icon nzType="check-circle" nzTheme="fill" class="text-green-500 mr-2"></i>
                          <span class="font-medium text-green-800">Success!</span>
                        </div>
                        
                        <div class="grid grid-cols-3 gap-4 text-sm">
                          <div *ngIf="lastExecutionResult.runtime">
                            <span class="text-gray-600">Runtime:</span>
                            <span class="font-medium ml-1">{{ lastExecutionResult.runtime }}ms</span>
                          </div>
                          <div *ngIf="lastExecutionResult.memory">
                            <span class="text-gray-600">Memory:</span>
                            <span class="font-medium ml-1">{{ lastExecutionResult.memory }}KB</span>
                          </div>
                          <div *ngIf="lastExecutionResult.testCasesPassed !== undefined">
                            <span class="text-gray-600">Test Cases:</span>
                            <span class="font-medium ml-1">{{ lastExecutionResult.testCasesPassed }}/{{ lastExecutionResult.totalTestCases }}</span>
                          </div>
                        </div>
                        
                        <div *ngIf="lastExecutionResult.output" class="mt-3">
                          <div class="text-sm text-gray-600 mb-1">Output:</div>
                          <pre class="bg-white p-2 rounded border text-xs font-mono">{{ lastExecutionResult.output }}</pre>
                        </div>
                      </div>
                    </div>
                    
                    <ng-template #errorTemplate>
                      <!-- Error Result -->
                      <app-submission-error
                        [status]="getSubmissionStatus(lastExecutionResult)"
                        [errorMessage]="lastExecutionResult.error"
                        [testCasesPassed]="lastExecutionResult.testCasesPassed || 0"
                        [totalTestCases]="lastExecutionResult.totalTestCases || 0">
                      </app-submission-error>
                    </ng-template>
                  </div>
                  
                  <!-- No results yet -->
                  <div *ngIf="!lastExecutionResult" class="text-center text-gray-500 py-8">
                    <i nz-icon nzType="play-circle" nzTheme="outline" class="text-2xl mb-2"></i>
                    <p>Run your code to see results here</p>
                  </div>
                </div>
              </nz-tab>

              <!-- Submission History Tab -->
              <nz-tab nzTitle="Submissions">
                <div class="h-full overflow-hidden">
                  <app-submission-history 
                    [problemId]="(problem$ | async)?.id"
                    class="h-full block">
                  </app-submission-history>
                </div>
              </nz-tab>

              <!-- Performance Metrics Tab -->
              <nz-tab nzTitle="Analytics">
                <div class="p-4 h-full overflow-y-auto">
                  <app-performance-metrics 
                    [problemId]="(problem$ | async)?.id">
                  </app-performance-metrics>
                </div>
              </nz-tab>
            </nz-tabset>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    :host {
      display: block;
      height: 100vh;
    }
    
    .prose code {
      background-color: #f3f4f6;
      padding: 2px 4px;
      border-radius: 4px;
      font-size: 0.875rem;
    }
  `]
})
export class ProblemSolveComponent implements OnInit, OnDestroy {
  @ViewChild('codeEditor') codeEditor!: CodeEditorComponent;
  
  private readonly store = inject(Store);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly message = inject(NzMessageService);
  private readonly destroy$ = new Subject<void>();

  problem$ = this.store.select(selectSelectedProblem);
  loading$ = this.store.select(selectProblemLoading);
  error$ = this.store.select(selectProblemError);
  executionState$ = this.store.select(selectExecutionState);

  showDescription = true;
  lastExecutionResult: CodeExecutionResult | null = null;

  ngOnInit() {
    // Load problem based on route params
    this.route.params.pipe(
      takeUntil(this.destroy$)
    ).subscribe(params => {
      const id = params['id'];
      if (id) {
        this.store.dispatch(ProblemActions.loadProblem({ id: +id }));
      }
    });

    // Listen to execution state changes
    this.executionState$.pipe(
      takeUntil(this.destroy$),
      filter(state => state.result !== null || state.error !== null)
    ).subscribe(state => {
      if (this.codeEditor) {
        if (state.result) {
          this.lastExecutionResult = state.result;
          this.codeEditor.setExecutionResult(state.result);
          
          // Show success message for accepted submissions
          if (state.result.success && state.result.testCasesPassed === state.result.totalTestCases) {
            this.message.success('Solution accepted! 🎉');
          }
        }
        
        if (state.error) {
          this.message.error(state.error);
          this.codeEditor.setRunningState(false);
          this.codeEditor.setSubmittingState(false);
        }
      }
    });
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
    this.store.dispatch(ProblemActions.clearSelectedProblem());
    this.store.dispatch(ProblemActions.clearExecutionResult());
  }

  goBack() {
    this.router.navigate(['/problems']);
  }

  toggleDescription() {
    this.showDescription = !this.showDescription;
  }

  onCodeRun(submission: CodeSubmission) {
    this.store.dispatch(ProblemActions.runCode({ submission }));
  }

  onCodeSubmit(submission: CodeSubmission) {
    this.store.dispatch(ProblemActions.submitCode({ submission }));
  }

  getDifficultyColor(difficulty: string): string {
    switch (difficulty) {
      case 'Easy': return 'green';
      case 'Medium': return 'orange';
      case 'Hard': return 'red';
      default: return 'default';
    }
  }

  formatDescription(description: string): string {
    return description
      .replace(/\n/g, '<br>')
      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.*?)\*/g, '<em>$1</em>')
      .replace(/`(.*?)`/g, '<code>$1</code>');
  }

  formatConstraints(constraints: string): string {
    return constraints
      .split('\n')
      .filter(line => line.trim())
      .map(line => `<li class="text-sm text-gray-700">${line.trim()}</li>`)
      .join('');
  }

  getSubmissionStatus(result: CodeExecutionResult): 'Accepted' | 'Wrong Answer' | 'Runtime Error' | 'Time Limit Exceeded' | 'Memory Limit Exceeded' | 'Compilation Error' | 'Internal Error' {
    if (!result.success) {
      if (result.error) {
        if (result.error.includes('timeout') || result.error.includes('time limit')) {
          return 'Time Limit Exceeded';
        }
        if (result.error.includes('memory') || result.error.includes('out of memory')) {
          return 'Memory Limit Exceeded';
        }
        if (result.error.includes('compilation') || result.error.includes('syntax')) {
          return 'Compilation Error';
        }
        return 'Runtime Error';
      }
      return 'Wrong Answer';
    }
    return 'Accepted';
  }
}