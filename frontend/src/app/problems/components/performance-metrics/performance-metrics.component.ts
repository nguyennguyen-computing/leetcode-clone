import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Subject, takeUntil } from 'rxjs';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzStatisticModule } from 'ng-zorro-antd/statistic';
import { NzProgressModule } from 'ng-zorro-antd/progress';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzSpinModule } from 'ng-zorro-antd/spin';
import { NzEmptyModule } from 'ng-zorro-antd/empty';
import { NzDividerModule } from 'ng-zorro-antd/divider';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { SubmissionService } from '../../services/submission.service';
import { SubmissionStats, SubmissionStatus } from '../../models/submission.models';

interface LanguageMetric {
  language: string;
  count: number;
  percentage: number;
  color: string;
}

interface StatusMetric {
  status: SubmissionStatus;
  count: number;
  percentage: number;
  color: string;
}

@Component({
  selector: 'app-performance-metrics',
  standalone: true,
  imports: [
    CommonModule,
    NzCardModule,
    NzStatisticModule,
    NzProgressModule,
    NzTagModule,
    NzSpinModule,
    NzEmptyModule,
    NzDividerModule,
    NzIconModule
  ],
  template: `
    <div class="performance-metrics">
      <nz-spin [nzSpinning]="loading">
        <div *ngIf="stats" class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Overall Statistics -->
          <nz-card nzTitle="Overall Statistics" class="h-fit">
            <div class="grid grid-cols-2 gap-4">
              <nz-statistic 
                nzTitle="Total Submissions" 
                [nzValue]="stats.totalSubmissions"
                [nzValueStyle]="{ color: '#1890ff' }">
                <ng-template #nzPrefix>
                  <span nz-icon nzType="file-text" nzTheme="outline"></span>
                </ng-template>
              </nz-statistic>
              
              <nz-statistic 
                nzTitle="Accepted" 
                [nzValue]="stats.acceptedSubmissions"
                [nzValueStyle]="{ color: '#52c41a' }">
                <ng-template #nzPrefix>
                  <span nz-icon nzType="check-circle" nzTheme="outline"></span>
                </ng-template>
              </nz-statistic>
            </div>
            
            <nz-divider></nz-divider>
            
            <div class="mb-4">
              <div class="flex justify-between items-center mb-2">
                <span class="text-sm font-medium">Acceptance Rate</span>
                <span class="text-sm font-bold" [style.color]="getAcceptanceRateColor()">
                  {{ stats.acceptanceRate.toFixed(1) }}%
                </span>
              </div>
              <nz-progress 
                [nzPercent]="stats.acceptanceRate" 
                [nzStrokeColor]="getAcceptanceRateColor()"
                [nzShowInfo]="false">
              </nz-progress>
            </div>
          </nz-card>

          <!-- Language Distribution -->
          <nz-card nzTitle="Language Distribution" class="h-fit">
            <div *ngIf="languageMetrics.length > 0; else noLanguageData">
              <div *ngFor="let metric of languageMetrics" class="mb-4 last:mb-0">
                <div class="flex justify-between items-center mb-2">
                  <div class="flex items-center">
                    <div 
                      class="w-3 h-3 rounded-full mr-2"
                      [style.background-color]="metric.color">
                    </div>
                    <span class="text-sm font-medium">{{ getLanguageLabel(metric.language) }}</span>
                  </div>
                  <div class="text-right">
                    <div class="text-sm font-bold">{{ metric.count }}</div>
                    <div class="text-xs text-gray-500">{{ metric.percentage.toFixed(1) }}%</div>
                  </div>
                </div>
                <nz-progress 
                  [nzPercent]="metric.percentage" 
                  [nzStrokeColor]="metric.color"
                  [nzShowInfo]="false"
                  nzSize="small">
                </nz-progress>
              </div>
            </div>
            
            <ng-template #noLanguageData>
              <nz-empty 
                nzNotFoundContent="No language data available"
                [nzNotFoundImage]="'simple'">
              </nz-empty>
            </ng-template>
          </nz-card>

          <!-- Status Distribution -->
          <nz-card nzTitle="Submission Status Distribution" class="lg:col-span-2">
            <div *ngIf="statusMetrics.length > 0; else noStatusData" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div *ngFor="let metric of statusMetrics" class="status-metric-card">
                <div class="flex items-center justify-between p-4 border rounded-lg" [style.border-color]="metric.color + '40'">
                  <div class="flex items-center">
                    <span 
                      nz-icon 
                      [nzType]="getStatusIcon(metric.status)"
                      class="text-lg mr-3"
                      [style.color]="metric.color">
                    </span>
                    <div>
                      <div class="font-medium text-sm">{{ getStatusShort(metric.status) }}</div>
                      <div class="text-xs text-gray-500">{{ metric.status }}</div>
                    </div>
                  </div>
                  <div class="text-right">
                    <div class="text-lg font-bold" [style.color]="metric.color">{{ metric.count }}</div>
                    <div class="text-xs text-gray-500">{{ metric.percentage.toFixed(1) }}%</div>
                  </div>
                </div>
              </div>
            </div>
            
            <ng-template #noStatusData>
              <nz-empty 
                nzNotFoundContent="No status data available"
                [nzNotFoundImage]="'simple'">
              </nz-empty>
            </ng-template>
          </nz-card>

          <!-- Performance Insights -->
          <nz-card nzTitle="Performance Insights" class="lg:col-span-2" *ngIf="stats.acceptedSubmissions > 0">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div class="text-center">
                <div class="text-2xl font-bold text-green-500 mb-2">
                  {{ getSuccessRate() }}%
                </div>
                <div class="text-sm text-gray-600">Success Rate</div>
                <div class="text-xs text-gray-500 mt-1">
                  First attempt success
                </div>
              </div>
              
              <div class="text-center">
                <div class="text-2xl font-bold text-blue-500 mb-2">
                  {{ getMostUsedLanguage() }}
                </div>
                <div class="text-sm text-gray-600">Preferred Language</div>
                <div class="text-xs text-gray-500 mt-1">
                  Most frequently used
                </div>
              </div>
              
              <div class="text-center">
                <div class="text-2xl font-bold text-purple-500 mb-2">
                  {{ getAverageAttempts() }}
                </div>
                <div class="text-sm text-gray-600">Avg. Attempts</div>
                <div class="text-xs text-gray-500 mt-1">
                  Per accepted solution
                </div>
              </div>
            </div>
          </nz-card>
        </div>

        <nz-empty 
          *ngIf="!loading && !stats"
          nzNotFoundContent="No performance data available"
          [nzNotFoundImage]="'simple'">
        </nz-empty>
      </nz-spin>
    </div>
  `,
  styles: [`
    .performance-metrics {
      width: 100%;
    }
    
    .status-metric-card {
      transition: transform 0.2s ease;
    }
    
    .status-metric-card:hover {
      transform: translateY(-2px);
    }
    
    :host ::ng-deep .ant-statistic-content {
      font-size: 24px;
    }
    
    :host ::ng-deep .ant-statistic-title {
      margin-bottom: 8px;
    }
  `]
})
export class PerformanceMetricsComponent implements OnInit, OnDestroy {
  @Input() problemId?: number; // If provided, show metrics for specific problem only

  private destroy$ = new Subject<void>();

  stats: SubmissionStats | null = null;
  loading = false;
  languageMetrics: LanguageMetric[] = [];
  statusMetrics: StatusMetric[] = [];

  private languageColors = ['#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1', '#13c2c2'];
  private statusColors: { [key in SubmissionStatus]: string } = {
    'Accepted': '#52c41a',
    'Wrong Answer': '#f5222d',
    'Time Limit Exceeded': '#faad14',
    'Memory Limit Exceeded': '#fa8c16',
    'Runtime Error': '#f5222d',
    'Compilation Error': '#ff4d4f',
    'Internal Error': '#8c8c8c'
  };

  constructor(private submissionService: SubmissionService) {}

  ngOnInit() {
    this.loadStats();
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  loadStats() {
    this.loading = true;
    
    this.submissionService
      .getSubmissionStats(this.problemId)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (stats) => {
          this.stats = stats;
          this.processLanguageMetrics();
          this.processStatusMetrics();
          this.loading = false;
        },
        error: (error) => {
          console.error('Error loading submission stats:', error);
          this.loading = false;
        }
      });
  }

  private processLanguageMetrics() {
    if (!this.stats) return;

    this.languageMetrics = Object.entries(this.stats.languageStats)
      .map(([language, count], index) => ({
        language,
        count,
        percentage: (count / this.stats!.totalSubmissions) * 100,
        color: this.languageColors[index % this.languageColors.length]
      }))
      .sort((a, b) => b.count - a.count);
  }

  private processStatusMetrics() {
    if (!this.stats) return;

    this.statusMetrics = Object.entries(this.stats.statusStats)
      .map(([status, count]) => ({
        status: status as SubmissionStatus,
        count,
        percentage: (count / this.stats!.totalSubmissions) * 100,
        color: this.statusColors[status as SubmissionStatus]
      }))
      .sort((a, b) => b.count - a.count);
  }

  getAcceptanceRateColor(): string {
    if (!this.stats) return '#1890ff';
    
    const rate = this.stats.acceptanceRate;
    if (rate >= 80) return '#52c41a';
    if (rate >= 60) return '#faad14';
    if (rate >= 40) return '#fa8c16';
    return '#f5222d';
  }

  getLanguageLabel(language: string): string {
    const labelMap: { [key: string]: string } = {
      javascript: 'JavaScript',
      python: 'Python',
      java: 'Java',
    };
    return labelMap[language] || language.charAt(0).toUpperCase() + language.slice(1);
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

  getStatusShort(status: SubmissionStatus): string {
    const shortMap: { [key in SubmissionStatus]: string } = {
      'Accepted': 'AC',
      'Wrong Answer': 'WA',
      'Time Limit Exceeded': 'TLE',
      'Memory Limit Exceeded': 'MLE',
      'Runtime Error': 'RE',
      'Compilation Error': 'CE',
      'Internal Error': 'IE'
    };
    return shortMap[status] || status;
  }

  getSuccessRate(): number {
    if (!this.stats || this.stats.totalSubmissions === 0) return 0;
    return Math.round(this.stats.acceptanceRate);
  }

  getMostUsedLanguage(): string {
    if (this.languageMetrics.length === 0) return 'N/A';
    return this.getLanguageLabel(this.languageMetrics[0].language);
  }

  getAverageAttempts(): string {
    if (!this.stats || this.stats.acceptedSubmissions === 0) return 'N/A';
    const avgAttempts = this.stats.totalSubmissions / this.stats.acceptedSubmissions;
    return avgAttempts.toFixed(1);
  }
}