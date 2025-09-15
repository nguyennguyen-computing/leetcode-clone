import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Store } from '@ngrx/store';

import { NzCardModule } from 'ng-zorro-antd/card';
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox';
import { NzRadioModule } from 'ng-zorro-antd/radio';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzDividerModule } from 'ng-zorro-antd/divider';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzInputModule } from 'ng-zorro-antd/input';

import { selectProblemFilters, selectAvailableTags } from '../../store/problem.selectors';
import * as ProblemActions from '../../store/problem.actions';
import { ProblemFilters } from '../../models/problem.models';

@Component({
    selector: 'app-problem-filter',
    standalone: true,
    imports: [
        CommonModule,
        FormsModule,
        NzCardModule,
        NzCheckboxModule,
        NzRadioModule,
        NzButtonModule,
        NzIconModule,
        NzDividerModule,
        NzTagModule,
        NzInputModule
    ],
    template: `
    <nz-card nzTitle="Filters" class="sticky top-4">
      <div class="space-y-6">
        <!-- Search -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Search</label>
          <nz-input-group nzSuffixIcon="search">
            <input 
              nz-input 
              placeholder="Search problems..." 
              [(ngModel)]="searchQuery"
              (ngModelChange)="onSearchChange($event)"
              class="w-full"
            />
          </nz-input-group>
        </div>

        <nz-divider></nz-divider>

        <!-- Status Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">Status</label>
          <nz-radio-group 
            [(ngModel)]="currentFilters.status" 
            (ngModelChange)="onStatusChange($event)"
            class="flex flex-col space-y-2">
            <label nz-radio nzValue="all" class="flex items-center">
              <span class="ml-2">All Problems</span>
            </label>
            <label nz-radio nzValue="solved" class="flex items-center">
              <span class="ml-2 text-green-600">Solved</span>
            </label>
            <label nz-radio nzValue="unsolved" class="flex items-center">
              <span class="ml-2 text-gray-600">Unsolved</span>
            </label>
          </nz-radio-group>
        </div>

        <nz-divider></nz-divider>

        <!-- Difficulty Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">Difficulty</label>
          <div class="space-y-2">
            <label 
              nz-checkbox 
              [(ngModel)]="difficultyChecked.Easy"
              (ngModelChange)="onDifficultyChange('Easy', $event)"
              class="flex items-center">
              <span class="ml-2 text-green-600 font-medium">Easy</span>
            </label>
            <label 
              nz-checkbox 
              [(ngModel)]="difficultyChecked.Medium"
              (ngModelChange)="onDifficultyChange('Medium', $event)"
              class="flex items-center">
              <span class="ml-2 text-yellow-600 font-medium">Medium</span>
            </label>
            <label 
              nz-checkbox 
              [(ngModel)]="difficultyChecked.Hard"
              (ngModelChange)="onDifficultyChange('Hard', $event)"
              class="flex items-center">
              <span class="ml-2 text-red-600 font-medium">Hard</span>
            </label>
          </div>
        </div>

        <nz-divider></nz-divider>

        <!-- Tags Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">Topics</label>
          <div class="max-h-48 overflow-y-auto space-y-2">
            <label 
              *ngFor="let tag of availableTags$ | async" 
              nz-checkbox 
              [ngModel]="isTagSelected(tag)"
              (ngModelChange)="onTagChange(tag, $event)"
              class="flex items-center">
              <span class="ml-2 text-sm">{{ tag }}</span>
            </label>
          </div>
        </div>

        <nz-divider></nz-divider>

        <!-- Clear Filters -->
        <button 
          nz-button 
          nzType="default" 
          nzBlock
          (click)="clearFilters()"
          class="w-full">
          <i nz-icon nzType="clear" class="mr-2"></i>
          Clear All Filters
        </button>
      </div>
    </nz-card>
  `
})
export class ProblemFilterComponent implements OnInit {
    private readonly store = inject(Store);

    currentFilters$ = this.store.select(selectProblemFilters);
    availableTags$ = this.store.select(selectAvailableTags);

    currentFilters: ProblemFilters = {
        difficulty: [],
        tags: [],
        status: 'all',
        searchQuery: ''
    };

    difficultyChecked = {
        Easy: false,
        Medium: false,
        Hard: false
    };

    searchQuery = '';

    ngOnInit() {
        this.store.dispatch(ProblemActions.loadTags());

        this.currentFilters$.subscribe(filters => {
            this.currentFilters = { ...filters };
            this.searchQuery = filters.searchQuery;
            this.updateDifficultyChecked();
        });
    }

    onSearchChange(query: string) {
        this.store.dispatch(ProblemActions.setSearchQuery({ query }));
    }

    onStatusChange(status: 'all' | 'solved' | 'unsolved') {
        this.store.dispatch(ProblemActions.updateFilters({
            filters: { status }
        }));
    }

    onDifficultyChange(difficulty: 'Easy' | 'Medium' | 'Hard', checked: boolean) {
        const currentDifficulties = [...this.currentFilters.difficulty];

        if (checked) {
            if (!currentDifficulties.includes(difficulty)) {
                currentDifficulties.push(difficulty);
            }
        } else {
            const index = currentDifficulties.indexOf(difficulty);
            if (index > -1) {
                currentDifficulties.splice(index, 1);
            }
        }

        this.store.dispatch(ProblemActions.updateFilters({
            filters: { difficulty: currentDifficulties }
        }));
    }

    onTagChange(tag: string, checked: boolean) {
        const currentTags = [...this.currentFilters.tags];

        if (checked) {
            if (!currentTags.includes(tag)) {
                currentTags.push(tag);
            }
        } else {
            const index = currentTags.indexOf(tag);
            if (index > -1) {
                currentTags.splice(index, 1);
            }
        }

        this.store.dispatch(ProblemActions.updateFilters({
            filters: { tags: currentTags }
        }));
    }

    isTagSelected(tag: string): boolean {
        return this.currentFilters.tags.includes(tag);
    }

    clearFilters() {
        this.store.dispatch(ProblemActions.resetFilters());
    }

    private updateDifficultyChecked() {
        this.difficultyChecked = {
            Easy: this.currentFilters.difficulty.includes('Easy'),
            Medium: this.currentFilters.difficulty.includes('Medium'),
            Hard: this.currentFilters.difficulty.includes('Hard')
        };
    }
}