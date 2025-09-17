import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subject, takeUntil, debounceTime, distinctUntilChanged } from 'rxjs';
import { NzTableModule } from 'ng-zorro-antd/table';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzEmptyModule } from 'ng-zorro-antd/empty';
import { NzSpinModule } from 'ng-zorro-antd/spin';
import { NzModalModule } from 'ng-zorro-antd/modal';
import { NzMessageService } from 'ng-zorro-antd/message';
import { SubmissionService } from '../../services/submission.service';
import { 
  Submission, 
  SubmissionFilters, 
  SubmissionStatus,
  DetailedSubmissionResult 
} from '../../models/submission.models';
import { SubmissionResultComponent } from '../submission-result/submission-result.component';

@Component({
  selector: 'app-submission-history',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    NzTableModule,
    NzCardModule,
    NzSelectModule,
    NzDatePickerModule,
    NzButtonModule,
    NzIconModule,
    NzTagModule,
    NzEmptyModule,
    NzSpinModule,
    NzModalModule,
    SubmissionResultComponent
  ],
  template: `
    <nz-card nzTitle="Submission History" [nzExtra]="filterTemplate">
      <ng-template #filterTemplate>
        <div class="flex items-center space-x-2">
          <nz-select 
            [(ngModel)]="filters.status" 
            (ngModelChange)="onFilterChange()"
            nzPlaceHolder="Status"
            class="w-32">
            <nz-option nzValue="all" nzLabel="All Status"></nz-option>
            <nz-option 
              *ngFor="let status of statusOptions" 
              [nzValue]="status" 
              [nzLabel]="status">
            </nz-option>
          </nz-select>

          <nz-select 
            [(ngModel)]="filters.language" 
            (ngModelChange)="onFilterChange()"
            nzPlaceHolder="Language"
            class="w-32">
            <nz-option nzValue="all" nzLabel="All Languages"></nz-option>
            <nz-option 
              *ngFor="let lang of languageOptions" 
              [nzValue]="lang.value" 
              [nzLabel]="lang.label">
            </nz-option>
          </nz-select>

          <nz-range-picker 
            [(ngModel)]="dateRange"
            (ngModelChange)="onDateRangeChange($event)"
            nzFormat="yyyy-MM-dd"
            class="w-64">
          </nz-range-picker>

          <button 
            nz-button 
            nzType="default" 
            (click)="resetFilters()"
            nzSize="small">
            <span nz-icon nzType="reload"></span>
            Reset
          </button>
        </div>
      </ng-template>

      <nz-spin [nzSpinning]="loading">
        <nz-table 
          #submissionTable
          [nzData]="submissions"
          [nzTotal]="total"
          [nzPageSize]="pageSize"
          [nzPageIndex]="currentPage"
          [nzShowSizeChanger]="true"
          [nzPageSizeOptions]="[10, 20, 50]"
          (nzPageIndexChange)="onPageChange($event)"
          (nzPageSizeChange)="onPageSizeChange($event)"
          nzShowPagination
          nzSize="middle">
          
          <thead>
            <tr>
              <th nzWidth="60px">#</th>
              <th *ngIf="!problemId">Problem</th>
              <th nzWidth="100px">Status</th>
              <th nzWidth="100px">Language</th>
              <th nzWidth="80px">Runtime</th>
              <th nzWidth="80px">Memory</th>
              <th nzWidth="100px">Test Cases</th>
              <th nzWidth="150px">Submitted</th>
              <th nzWidth="100px">Actions</th>
            </tr>
          </thead>
          
          <tbody>
            <tr *ngFor="let submission of submissionTable.data; let i = index">
              <td>{{ submission.id }}</td>
              
              <td *ngIf="!problemId">
                <div *ngIf="submission.problem">
                  <div class="font-medium">{{ submission.problem.title }}</div>
                  <nz-tag 
                    [nzColor]="getDifficultyColor(submission.problem.difficulty)"
                    nzSize="small">
                    {{ submission.problem.difficulty }}
                  </nz-tag>
                </div>
              </td>
              
              <td>
                <nz-tag 
                  [nzColor]="getStatusColor(submission.status)"
                  class="text-xs">
                  <span nz-icon [nzType]="getStatusIcon(submission.status)" class="mr-1"></span>
                  {{ getStatusShort(submission.status) }}
                </nz-tag>
              </td>
              
              <td>
                <span class="text-sm font-mono">{{ submission.language }}</span>
              </td>
              
              <td>
                <span *ngIf="submission.runtimeMs" class="text-sm">
                  {{ submission.runtimeMs }}ms
                </span>
                <span *ngIf="!submission.runtimeMs" class="text-gray-400">-</span>
              </td>
              
              <td>
                <span *ngIf="submission.memoryKb" class="text-sm">
                  {{ formatMemory(submission.memoryKb) }}
                </span>
                <span *ngIf="!submission.memoryKb" class="text-gray-400">-</span>
              </td>
              
              <td>
                <span class="text-sm">
                  {{ submission.testCasesPassed }}/{{ submission.totalTestCases }}
                </span>
              </td>
              
              <td>
                <span class="text-sm text-gray-600">
                  {{ formatDate(submission.submittedAt) }}
                </span>
              </td>
              
              <td>
                <button 
                  nz-button 
                  nzType="link" 
                  nzSize="small"
                  (click)="viewSubmissionDetails(submission)"
                  [nzLoading]="loadingDetails === submission.id">
                  <span nz-icon nzType="eye"></span>
                  View
                </button>
              </td>
            </tr>
          </tbody>
        </nz-table>

        <nz-empty 
          *ngIf="!loading && submissions.length === 0"
          nzNotFoundContent="No submissions found"
          [nzNotFoundImage]="'simple'">
        </nz-empty>
      </nz-spin>
    </nz-card>

    <!-- Submission Details Modal -->
    <nz-modal
      [(nzVisible)]="detailsModalVisible"
      nzTitle="Submission Details"
      [nzFooter]="null"
      nzWidth="80%"
      (nzOnCancel)="closeDetailsModal()">
      
      <ng-container *nzModalContent>
        <div *ngIf="selectedSubmissionDetails" class="max-h-96 overflow-y-auto">
          <app-submission-result [result]="selectedSubmissionDetails"></app-submission-result>
        </div>
        
        <nz-spin *ngIf="!selectedSubmissionDetails" nzTip="Loading submission details...">
          <div class="h-32"></div>
        </nz-spin>
      </ng-container>
    </nz-modal>
  `,
  styles: [`
    :host ::ng-deep .ant-table-tbody > tr > td {
      padding: 8px 16px;
    }
    
    :host ::ng-deep .ant-tag {
      margin: 0;
    }
    
    .filter-section {
      margin-bottom: 16px;
    }
  `]
})
export class SubmissionHistoryComponent implements OnInit, OnDestroy {
  @Input() problemId?: number; // If provided, show submissions for specific problem only

  private destroy$ = new Subject<void>();

  submissions: Submission[] = [];
  loading = false;
  loadingDetails: number | null = null;
  total = 0;
  currentPage = 1;
  pageSize = 20;

  filters: SubmissionFilters = {
    status: 'all',
    language: 'all'
  };

  dateRange: [Date, Date] | null = null;

  statusOptions: SubmissionStatus[] = [
    'Accepted',
    'Wrong Answer',
    'Time Limit Exceeded',
    'Memory Limit Exceeded',
    'Runtime Error',
    'Compilation Error',
    'Internal Error'
  ];

  languageOptions = [
    { value: 'javascript', label: 'JavaScript' },
    { value: 'python', label: 'Python' },
    { value: 'java', label: 'Java' }
  ];

  // Modal state
  detailsModalVisible = false;
  selectedSubmissionDetails: DetailedSubmissionResult | null = null;

  constructor(
    private submissionService: SubmissionService,
    private message: NzMessageService
  ) {}

  ngOnInit() {
    this.loadSubmissions();
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  loadSubmissions() {
    this.loading = true;
    
    const filters = { ...this.filters };
    if (this.problemId) {
      filters.problemId = this.problemId;
    }

    this.submissionService
      .getSubmissions(filters, this.currentPage, this.pageSize)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (response) => {
          this.submissions = response.submissions;
          this.total = response.total;
          this.loading = false;
        },
        error: (error) => {
          console.error('Error loading submissions:', error);
          this.message.error('Failed to load submissions');
          this.loading = false;
        }
      });
  }

  onFilterChange() {
    this.currentPage = 1;
    this.loadSubmissions();
  }

  onDateRangeChange(dates: [Date, Date] | null) {
    this.dateRange = dates;
    if (dates) {
      this.filters.dateRange = {
        start: dates[0].toISOString().split('T')[0],
        end: dates[1].toISOString().split('T')[0]
      };
    } else {
      delete this.filters.dateRange;
    }
    this.onFilterChange();
  }

  resetFilters() {
    this.filters = {
      status: 'all',
      language: 'all'
    };
    this.dateRange = null;
    this.onFilterChange();
  }

  onPageChange(page: number) {
    this.currentPage = page;
    this.loadSubmissions();
  }

  onPageSizeChange(size: number) {
    this.pageSize = size;
    this.currentPage = 1;
    this.loadSubmissions();
  }

  viewSubmissionDetails(submission: Submission) {
    this.loadingDetails = submission.id;
    this.detailsModalVisible = true;
    this.selectedSubmissionDetails = null;

    this.submissionService
      .getDetailedSubmissionResult(submission.id)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (details) => {
          this.selectedSubmissionDetails = details;
          this.loadingDetails = null;
        },
        error: (error) => {
          console.error('Error loading submission details:', error);
          this.message.error('Failed to load submission details');
          this.loadingDetails = null;
          this.detailsModalVisible = false;
        }
      });
  }

  closeDetailsModal() {
    this.detailsModalVisible = false;
    this.selectedSubmissionDetails = null;
    this.loadingDetails = null;
  }

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

  getDifficultyColor(difficulty: 'Easy' | 'Medium' | 'Hard'): string {
    const colorMap = {
      'Easy': 'green',
      'Medium': 'orange',
      'Hard': 'red'
    };
    return colorMap[difficulty] || 'default';
  }

  formatMemory(memoryKb: number): string {
    if (memoryKb < 1024) {
      return `${memoryKb}KB`;
    }
    return `${(memoryKb / 1024).toFixed(1)}MB`;
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / (1000 * 60));
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    
    return date.toLocaleDateString();
  }
}