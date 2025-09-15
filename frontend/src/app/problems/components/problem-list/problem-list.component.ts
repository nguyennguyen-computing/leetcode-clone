import { Component, inject, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { Store } from '@ngrx/store';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

import { NzCardModule } from 'ng-zorro-antd/card';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzStatisticModule } from 'ng-zorro-antd/statistic';
import { NzGridModule } from 'ng-zorro-antd/grid';
import { NzEmptyModule } from 'ng-zorro-antd/empty';
import { NzSpinModule } from 'ng-zorro-antd/spin';
import { NzPaginationModule } from 'ng-zorro-antd/pagination';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzLayoutModule } from 'ng-zorro-antd/layout';

import { selectCurrentUser } from '../../../auth/store/auth.selectors';
import { 
  selectProblems, 
  selectProblemLoading, 
  selectProblemError,
  selectProblemStats,
  selectProblemPagination
} from '../../store/problem.selectors';
import * as ProblemActions from '../../store/problem.actions';
import { ProblemFilterComponent } from '../problem-filter/problem-filter.component';
import { Problem } from '../../models/problem.models';



@Component({
  selector: 'app-problem-list',
  standalone: true,
  imports: [
    CommonModule,
    NzCardModule,
    NzButtonModule,
    NzTagModule,
    NzIconModule,
    NzStatisticModule,
    NzGridModule,
    NzEmptyModule,
    NzSpinModule,
    NzPaginationModule,
    NzInputModule,
    NzSelectModule,
    NzLayoutModule,
    ProblemFilterComponent
  ],
  template: `
    <div class="min-h-screen bg-gray-50">
      <!-- Hero Section -->
      <div class="bg-gradient-to-r from-blue-600 to-purple-700 text-white">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <div class="text-center">
            <h1 class="text-4xl font-bold mb-4">
              Welcome back, {{ (currentUser$ | async)?.username }}!
            </h1>
            <p class="text-xl text-blue-100 mb-8">
              Continue your coding journey with our curated problems
            </p>
            
            <!-- Stats -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-2xl mx-auto">
              <nz-card class="bg-white/10 backdrop-blur-sm border-white/20 text-center">
                <nz-statistic 
                  [nzValue]="0" 
                  nzTitle="Problems Solved" 
                  [nzValueStyle]="{ color: '#fff', fontSize: '24px', fontWeight: 'bold' }">
                </nz-statistic>
              </nz-card>
              <nz-card class="bg-white/10 backdrop-blur-sm border-white/20 text-center">
                <nz-statistic 
                  [nzValue]="0" 
                  nzTitle="Submissions" 
                  [nzValueStyle]="{ color: '#fff', fontSize: '24px', fontWeight: 'bold' }">
                </nz-statistic>
              </nz-card>
              <nz-card class="bg-white/10 backdrop-blur-sm border-white/20 text-center">
                <nz-statistic 
                  [nzValue]="0" 
                  nzTitle="Acceptance Rate" 
                  nzSuffix="%" 
                  [nzValueStyle]="{ color: '#fff', fontSize: '24px', fontWeight: 'bold' }">
                </nz-statistic>
              </nz-card>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <!-- Quick Actions -->
        <div class="mb-12">
          <h2 class="text-2xl font-bold text-gray-900 mb-6">Quick Start</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <nz-card class="hover:shadow-lg transition-shadow cursor-pointer border-l-4 border-l-green-500">
              <div class="flex items-center space-x-4">
                <div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
                  <i nz-icon nzType="play-circle" class="text-green-600 text-xl"></i>
                </div>
                <div>
                  <h3 class="font-semibold text-gray-900">Easy Problems</h3>
                  <p class="text-sm text-gray-600">Start with basics</p>
                </div>
              </div>
            </nz-card>
            
            <nz-card class="hover:shadow-lg transition-shadow cursor-pointer border-l-4 border-l-yellow-500">
              <div class="flex items-center space-x-4">
                <div class="w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center">
                  <i nz-icon nzType="fire" class="text-yellow-600 text-xl"></i>
                </div>
                <div>
                  <h3 class="font-semibold text-gray-900">Medium Problems</h3>
                  <p class="text-sm text-gray-600">Level up your skills</p>
                </div>
              </div>
            </nz-card>
            
            <nz-card class="hover:shadow-lg transition-shadow cursor-pointer border-l-4 border-l-red-500">
              <div class="flex items-center space-x-4">
                <div class="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center">
                  <i nz-icon nzType="thunderbolt" class="text-red-600 text-xl"></i>
                </div>
                <div>
                  <h3 class="font-semibold text-gray-900">Hard Problems</h3>
                  <p class="text-sm text-gray-600">Challenge yourself</p>
                </div>
              </div>
            </nz-card>
            
            <nz-card class="hover:shadow-lg transition-shadow cursor-pointer border-l-4 border-l-purple-500">
              <div class="flex items-center space-x-4">
                <div class="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
                  <i nz-icon nzType="trophy" class="text-purple-600 text-xl"></i>
                </div>
                <div>
                  <h3 class="font-semibold text-gray-900">Contests</h3>
                  <p class="text-sm text-gray-600">Compete with others</p>
                </div>
              </div>
            </nz-card>
          </div>
        </div>

        <!-- Problems Section -->
        <div class="mb-12">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-2xl font-bold text-gray-900">All Problems</h2>
            <div class="flex space-x-4">
              <button nz-button nzType="default">
                <i nz-icon nzType="filter" class="mr-2"></i>
                Filter
              </button>
              <button nz-button nzType="default">
                <i nz-icon nzType="sort-ascending" class="mr-2"></i>
                Sort
              </button>
            </div>
          </div>

          <!-- Empty State -->
          <nz-card class="text-center py-16">
            <nz-empty 
              nzNotFoundImage="https://gw.alipayobjects.com/zos/antfincdn/ZHrcdLPrvN/empty.svg"
              nzNotFoundContent="Problems will be available soon">
              <div nz-empty-footer>
                <p class="text-gray-600 mb-4">
                  We're working hard to bring you an amazing collection of coding problems.
                </p>
                <p class="text-sm text-gray-500 mb-6">
                  In the meantime, you can explore the authentication system we just built!
                </p>
                <button nz-button nzType="primary" class="bg-gradient-to-r from-blue-500 to-purple-600 border-0">
                  <i nz-icon nzType="bell" class="mr-2"></i>
                  Notify Me When Ready
                </button>
              </div>
            </nz-empty>
          </nz-card>
        </div>

        <!-- Recent Activity -->
        <div>
          <h2 class="text-2xl font-bold text-gray-900 mb-6">Recent Activity</h2>
          <nz-card>
            <div class="text-center py-8">
              <i nz-icon nzType="history" class="text-4xl text-gray-400 mb-4"></i>
              <p class="text-gray-600">No recent activity yet</p>
              <p class="text-sm text-gray-500">Start solving problems to see your progress here</p>
            </div>
          </nz-card>
        </div>
      </div>
    </div>
  `
})
export class ProblemListComponent {
  private readonly store = inject(Store);
  
  currentUser$ = this.store.select(selectCurrentUser);
}