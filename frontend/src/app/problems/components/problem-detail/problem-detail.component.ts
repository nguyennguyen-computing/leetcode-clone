import { Component, inject, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { Store } from '@ngrx/store';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

import { NzCardModule } from 'ng-zorro-antd/card';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzDividerModule } from 'ng-zorro-antd/divider';
import { NzSpinModule } from 'ng-zorro-antd/spin';
import { NzAlertModule } from 'ng-zorro-antd/alert';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { NzDescriptionsModule } from 'ng-zorro-antd/descriptions';

import { selectSelectedProblem, selectProblemLoading, selectProblemError } from '../../store/problem.selectors';
import * as ProblemActions from '../../store/problem.actions';
import { Problem } from '../../models/problem.models';

@Component({
  selector: 'app-problem-detail',
  standalone: true,
  imports: [
    CommonModule,
    NzCardModule,
    NzButtonModule,
    NzTagModule,
    NzIconModule,
    NzDividerModule,
    NzSpinModule,
    NzAlertModule,
    NzTypographyModule,
    NzDescriptionsModule
  ],
  template: `
    <div class="min-h-screen bg-gray-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <!-- Back Button -->
        <div class="mb-6">
          <button 
            nz-button 
            nzType="default" 
            (click)="goBack()"
            class="flex items-center">
            <i nz-icon nzType="arrow-left" class="mr-2"></i>
            Back to Problems
          </button>
        </div>

        <!-- Loading State -->
        <div *ngIf="loading$ | async" class="text-center py-16">
          <nz-spin nzSize="large"></nz-spin>
          <p class="mt-4 text-gray-600">Loading problem...</p>
        </div>

        <!-- Error State -->
        <nz-alert
          *ngIf="error$ | async as error"
          nzType="error"
          [nzMessage]="error"
          nzShowIcon
          class="mb-6">
        </nz-alert>

        <!-- Problem Content -->
        <div *ngIf="problem$ | async as problem" class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <!-- Main Content -->
          <div class="lg:col-span-2">
            <nz-card>
              <!-- Problem Header -->
              <div class="mb-6">
                <div class="flex items-start justify-between mb-4">
                  <div>
                    <h1 class="text-3xl font-bold text-gray-900 mb-2">
                      {{ problem.id }}. {{ problem.title }}
                      <i 
                        *ngIf="problem.isSolved" 
                        nz-icon 
                        nzType="check-circle" 
                        class="text-green-500 ml-2">
                      </i>
                    </h1>
                    <div class="flex items-center space-x-4">
                      <nz-tag 
                        [nzColor]="getDifficultyColor(problem.difficulty)"
                        class="font-medium">
                        {{ problem.difficulty }}
                      </nz-tag>
                      <span class="text-gray-500">
                        Acceptance Rate: {{ problem.acceptanceRate || 0 }}%
                      </span>
                    </div>
                  </div>
                  <button 
                    nz-button 
                    nzType="primary" 
                    nzSize="large"
                    (click)="startSolving(problem)"
                    class="bg-gradient-to-r from-blue-500 to-purple-600 border-0">
                    <i nz-icon nzType="code" class="mr-2"></i>
                    Solve Problem
                  </button>
                </div>

                <!-- Tags -->
                <div class="flex flex-wrap gap-2">
                  <nz-tag 
                    *ngFor="let tag of problem.tags" 
                    nzColor="blue"
                    class="cursor-pointer hover:opacity-80">
                    {{ tag }}
                  </nz-tag>
                </div>
              </div>

              <nz-divider></nz-divider>

              <!-- Problem Description -->
              <div class="mb-8">
                <h2 class="text-xl font-semibold text-gray-900 mb-4">Description</h2>
                <div 
                  class="prose max-w-none text-gray-700 leading-relaxed"
                  [innerHTML]="formatDescription(problem.description)">
                </div>
              </div>

              <!-- Examples -->
              <div *ngIf="problem.examples?.length" class="mb-8">
                <h2 class="text-xl font-semibold text-gray-900 mb-4">Examples</h2>
                <div class="space-y-6">
                  <div 
                    *ngFor="let example of problem.examples; let i = index" 
                    class="bg-gray-50 rounded-lg p-4 border">
                    <h3 class="font-medium text-gray-900 mb-3">Example {{ i + 1 }}:</h3>
                    
                    <div class="space-y-3">
                      <div>
                        <span class="font-medium text-gray-700">Input:</span>
                        <pre class="mt-1 bg-white p-3 rounded border text-sm font-mono">{{ example.input }}</pre>
                      </div>
                      
                      <div>
                        <span class="font-medium text-gray-700">Output:</span>
                        <pre class="mt-1 bg-white p-3 rounded border text-sm font-mono">{{ example.output }}</pre>
                      </div>
                      
                      <div *ngIf="example.explanation">
                        <span class="font-medium text-gray-700">Explanation:</span>
                        <p class="mt-1 text-gray-600">{{ example.explanation }}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Constraints -->
              <div *ngIf="problem.constraints" class="mb-8">
                <h2 class="text-xl font-semibold text-gray-900 mb-4">Constraints</h2>
                <div 
                  class="bg-yellow-50 border border-yellow-200 rounded-lg p-4"
                  [innerHTML]="formatConstraints(problem.constraints)">
                </div>
              </div>
            </nz-card>
          </div>

          <!-- Sidebar -->
          <div class="lg:col-span-1">
            <div class="space-y-6">
              <!-- Problem Stats -->
              <nz-card nzTitle="Problem Stats" nzSize="small">
                <nz-descriptions nzBordered nzSize="small">
                  <nz-descriptions-item nzTitle="Difficulty">
                    <nz-tag [nzColor]="getDifficultyColor(problem.difficulty)">
                      {{ problem.difficulty }}
                    </nz-tag>
                  </nz-descriptions-item>
                  <nz-descriptions-item nzTitle="Acceptance Rate">
                    {{ problem.acceptanceRate || 0 }}%
                  </nz-descriptions-item>
                  <nz-descriptions-item nzTitle="Status">
                    <span *ngIf="problem.isSolved" class="text-green-600 font-medium">
                      <i nz-icon nzType="check-circle" class="mr-1"></i>
                      Solved
                    </span>
                    <span *ngIf="!problem.isSolved" class="text-gray-600">
                      <i nz-icon nzType="clock-circle" class="mr-1"></i>
                      Not Attempted
                    </span>
                  </nz-descriptions-item>
                </nz-descriptions>
              </nz-card>

              <!-- Related Topics -->
              <nz-card nzTitle="Related Topics" nzSize="small">
                <div class="flex flex-wrap gap-2">
                  <nz-tag 
                    *ngFor="let tag of problem.tags" 
                    nzColor="processing"
                    class="cursor-pointer hover:opacity-80">
                    {{ tag }}
                  </nz-tag>
                </div>
              </nz-card>

              <!-- Actions -->
              <nz-card nzTitle="Actions" nzSize="small">
                <div class="space-y-3">
                  <button 
                    nz-button 
                    nzType="primary" 
                    nzBlock
                    (click)="startSolving(problem)"
                    class="bg-gradient-to-r from-blue-500 to-purple-600 border-0">
                    <i nz-icon nzType="code" class="mr-2"></i>
                    Solve Problem
                  </button>
                  
                  <button 
                    nz-button 
                    nzType="default" 
                    nzBlock
                    (click)="viewSolutions(problem)">
                    <i nz-icon nzType="eye" class="mr-2"></i>
                    View Solutions
                  </button>
                  
                  <button 
                    nz-button 
                    nzType="default" 
                    nzBlock
                    (click)="viewDiscussion(problem)">
                    <i nz-icon nzType="message" class="mr-2"></i>
                    Discussion
                  </button>
                </div>
              </nz-card>
            </div>
          </div>
        </div>
      </div>
    </div>
  `
})
export class ProblemDetailComponent implements OnInit, OnDestroy {
  private readonly store = inject(Store);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly destroy$ = new Subject<void>();

  problem$ = this.store.select(selectSelectedProblem);
  loading$ = this.store.select(selectProblemLoading);
  error$ = this.store.select(selectProblemError);

  ngOnInit() {
    this.route.params.pipe(
      takeUntil(this.destroy$)
    ).subscribe(params => {
      const id = params['id'];
      const slug = params['slug'];
      
      if (id) {
        this.store.dispatch(ProblemActions.loadProblem({ id: +id }));
      } else if (slug) {
        this.store.dispatch(ProblemActions.loadProblemBySlug({ slug }));
      }
    });
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
    this.store.dispatch(ProblemActions.clearSelectedProblem());
  }

  goBack() {
    this.router.navigate(['/problems']);
  }

  startSolving(problem: Problem) {
    this.router.navigate(['/problems', problem.id, 'solve']);
  }

  viewSolutions(problem: Problem) {
    this.router.navigate(['/problems', problem.id, 'solutions']);
  }

  viewDiscussion(problem: Problem) {
    this.router.navigate(['/problems', problem.id, 'discussion']);
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
    // Convert newlines to <br> tags and handle basic formatting
    return description
      .replace(/\n/g, '<br>')
      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.*?)\*/g, '<em>$1</em>')
      .replace(/`(.*?)`/g, '<code class="bg-gray-100 px-1 py-0.5 rounded text-sm">$1</code>');
  }

  formatConstraints(constraints: string): string {
    // Format constraints with bullet points
    return constraints
      .split('\n')
      .filter(line => line.trim())
      .map(line => `<li class="text-sm text-gray-700">${line.trim()}</li>`)
      .join('');
  }
}