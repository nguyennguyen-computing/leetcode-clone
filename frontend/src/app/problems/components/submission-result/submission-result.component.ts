import { Component, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzStatisticModule } from 'ng-zorro-antd/statistic';
import { NzProgressModule } from 'ng-zorro-antd/progress';
import { NzCollapseModule } from 'ng-zorro-antd/collapse';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { NzDividerModule } from 'ng-zorro-antd/divider';
import { DetailedSubmissionResult, TestCaseResult, SubmissionStatus } from '../../models/submission.models';

@Component({
  selector: 'app-submission-result',
  standalone: true,
  imports: [
    CommonModule,
    NzCardModule,
    NzIconModule,
    NzTagModule,
    NzStatisticModule,
    NzProgressModule,
    NzCollapseModule,
    NzTypographyModule,
    NzDividerModule
  ],
  template: `
    <div class="submission-result" *ngIf="result">
      <!-- Overall Status Card -->
      <nz-card 
        [nzTitle]="'Submission Result'" 
        class="mb-4"
        [nzExtra]="statusTemplate">
        
        <ng-template #statusTemplate>
          <nz-tag 
            [nzColor]="getStatusColor(result.submission.status)"
            class="text-sm px-3 py-1">
            <span nz-icon [nzType]="getStatusIcon(result.submission.status)" class="mr-1"></span>
            {{ result.submission.status }}
          </nz-tag>
        </ng-template>

        <!-- Performance Metrics -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
          <nz-statistic 
            nzTitle="Runtime" 
            [nzValue]="result.overallStats.runtimeMs" 
            nzSuffix="ms"
            [nzValueStyle]="{ color: getPerformanceColor('runtime', result.overallStats.runtimeMs) }">
          </nz-statistic>
          
          <nz-statistic 
            nzTitle="Memory" 
            [nzValue]="result.overallStats.memoryKb" 
            nzSuffix="KB"
            [nzValueStyle]="{ color: getPerformanceColor('memory', result.overallStats.memoryKb) }">
          </nz-statistic>
          
          <nz-statistic 
            nzTitle="Test Cases" 
            [nzValue]="result.overallStats.testCasesPassed" 
            [nzSuffix]="'/' + result.overallStats.totalTestCases">
          </nz-statistic>
        </div>

        <!-- Progress Bar -->
        <div class="mb-4">
          <div class="flex justify-between items-center mb-2">
            <span class="text-sm font-medium">Test Cases Progress</span>
            <span class="text-sm text-gray-500">
              {{ result.overallStats.testCasesPassed }}/{{ result.overallStats.totalTestCases }}
            </span>
          </div>
          <nz-progress 
            [nzPercent]="getTestCaseProgress()" 
            [nzStatus]="getProgressStatus()"
            [nzStrokeColor]="getProgressColor()">
          </nz-progress>
        </div>

        <!-- Error Message (if any) -->
        <div *ngIf="result.submission.errorMessage" class="mb-4">
          <nz-divider nzText="Error Details" nzOrientation="left"></nz-divider>
          <div class="bg-red-50 border border-red-200 rounded-lg p-4">
            <div class="flex items-start">
              <span nz-icon nzType="exclamation-circle" nzTheme="fill" class="text-red-500 mr-2 mt-1"></span>
              <div class="flex-1">
                <h4 class="text-red-800 font-medium mb-2">Execution Error</h4>
                <pre class="text-sm text-red-700 whitespace-pre-wrap font-mono">{{ result.submission.errorMessage }}</pre>
              </div>
            </div>
          </div>
        </div>
      </nz-card>

      <!-- Test Cases Details -->
      <nz-card nzTitle="Test Cases Details" *ngIf="result.testCaseResults.length > 0">
        <nz-collapse nzAccordion>
          <nz-collapse-panel 
            *ngFor="let testCase of result.testCaseResults; let i = index"
            [nzHeader]="getTestCaseHeader(i, testCase)"
            [nzActive]="!testCase.passed"
            [nzExtra]="testCaseStatusTemplate">
            
            <ng-template #testCaseStatusTemplate>
              <span 
                nz-icon 
                [nzType]="testCase.passed ? 'check-circle' : 'close-circle'"
                [class.text-green-500]="testCase.passed"
                [class.text-red-500]="!testCase.passed">
              </span>
            </ng-template>

            <div class="test-case-details">
              <!-- Input -->
              <div class="mb-4">
                <h5 class="text-sm font-medium text-gray-700 mb-2">Input:</h5>
                <div class="bg-gray-50 border rounded p-3">
                  <pre class="text-sm font-mono whitespace-pre-wrap">{{ testCase.input }}</pre>
                </div>
              </div>

              <!-- Expected Output -->
              <div class="mb-4">
                <h5 class="text-sm font-medium text-gray-700 mb-2">Expected Output:</h5>
                <div class="bg-gray-50 border rounded p-3">
                  <pre class="text-sm font-mono whitespace-pre-wrap">{{ testCase.expectedOutput }}</pre>
                </div>
              </div>

              <!-- Actual Output (if available) -->
              <div *ngIf="testCase.actualOutput !== undefined" class="mb-4">
                <h5 class="text-sm font-medium text-gray-700 mb-2">Your Output:</h5>
                <div 
                  class="border rounded p-3"
                  [class.bg-green-50]="testCase.passed"
                  [class.border-green-200]="testCase.passed"
                  [class.bg-red-50]="!testCase.passed"
                  [class.border-red-200]="!testCase.passed">
                  <pre class="text-sm font-mono whitespace-pre-wrap">{{ testCase.actualOutput }}</pre>
                </div>
              </div>

              <!-- Test Case Error -->
              <div *ngIf="testCase.error" class="mb-4">
                <h5 class="text-sm font-medium text-red-700 mb-2">Error:</h5>
                <div class="bg-red-50 border border-red-200 rounded p-3">
                  <pre class="text-sm font-mono whitespace-pre-wrap text-red-700">{{ testCase.error }}</pre>
                </div>
              </div>

              <!-- Performance for this test case -->
              <div *ngIf="testCase.runtimeMs || testCase.memoryKb" class="grid grid-cols-2 gap-4">
                <div *ngIf="testCase.runtimeMs">
                  <span class="text-xs text-gray-500">Runtime: </span>
                  <span class="text-xs font-medium">{{ testCase.runtimeMs }}ms</span>
                </div>
                <div *ngIf="testCase.memoryKb">
                  <span class="text-xs text-gray-500">Memory: </span>
                  <span class="text-xs font-medium">{{ testCase.memoryKb }}KB</span>
                </div>
              </div>
            </div>
          </nz-collapse-panel>
        </nz-collapse>
      </nz-card>
    </div>
  `,
  styles: [`
    .submission-result {
      max-width: 100%;
    }
    
    pre {
      margin: 0;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 12px;
      line-height: 1.4;
    }
    
    .test-case-details h5 {
      margin-bottom: 8px;
    }
    
    :host ::ng-deep .ant-collapse-header {
      align-items: center;
    }
    
    :host ::ng-deep .ant-statistic-content {
      font-size: 20px;
    }
  `]
})
export class SubmissionResultComponent implements OnInit {
  @Input() result: DetailedSubmissionResult | null = null;

  ngOnInit() {}

  getStatusColor(status: SubmissionStatus): string {
    const colorMap: { [key in SubmissionStatus]: string } = {
      'Accepted': 'green',
      'Wrong Answer': 'red',
      'Time Limit Exceeded': 'orange',
      'Memory Limit Exceeded': 'orange',
      'Runtime Error': 'red',
      'Compilation Error': 'red',
      'Internal Error': 'red'
    };
    return colorMap[status] || 'default';
  }

  getStatusIcon(status: SubmissionStatus): string {
    const iconMap: { [key in SubmissionStatus]: string } = {
      'Accepted': 'check-circle',
      'Wrong Answer': 'close-circle',
      'Time Limit Exceeded': 'clock-circle',
      'Memory Limit Exceeded': 'warning',
      'Runtime Error': 'exclamation-circle',
      'Compilation Error': 'code',
      'Internal Error': 'question-circle'
    };
    return iconMap[status] || 'question-circle';
  }

  getPerformanceColor(metric: 'runtime' | 'memory', value: number): string {
    // These thresholds would ideally come from the backend or be configurable
    if (metric === 'runtime') {
      if (value < 100) return '#52c41a'; // green
      if (value < 500) return '#faad14'; // orange
      return '#f5222d'; // red
    } else {
      if (value < 10000) return '#52c41a'; // green
      if (value < 50000) return '#faad14'; // orange
      return '#f5222d'; // red
    }
  }

  getTestCaseProgress(): number {
    if (!this.result) return 0;
    return Math.round((this.result.overallStats.testCasesPassed / this.result.overallStats.totalTestCases) * 100);
  }

  getProgressStatus(): 'success' | 'exception' | 'normal' {
    if (!this.result) return 'normal';
    
    const progress = this.getTestCaseProgress();
    if (progress === 100) return 'success';
    if (this.result.submission.status !== 'Accepted') return 'exception';
    return 'normal';
  }

  getProgressColor(): string {
    const status = this.getProgressStatus();
    if (status === 'success') return '#52c41a';
    if (status === 'exception') return '#f5222d';
    return '#1890ff';
  }

  getTestCaseHeader(index: number, testCase: TestCaseResult): string {
    const status = testCase.passed ? 'Passed' : 'Failed';
    return `Test Case ${index + 1} - ${status}`;
  }
}