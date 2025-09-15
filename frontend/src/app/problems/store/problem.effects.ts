import { Injectable, inject } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { Store } from '@ngrx/store';
import { of } from 'rxjs';
import { map, catchError, switchMap, withLatestFrom } from 'rxjs/operators';

import { ProblemService } from '../services/problem.service';
import { ExecutionService } from '../services/execution.service';
import * as ProblemActions from './problem.actions';
import { selectProblemFilters } from './problem.selectors';

@Injectable()
export class ProblemEffects {
  private readonly actions$ = inject(Actions);
  private readonly problemService = inject(ProblemService);
  private readonly executionService = inject(ExecutionService);
  private readonly store = inject(Store);

  loadProblems$ = createEffect(() =>
    this.actions$.pipe(
      ofType(ProblemActions.loadProblems),
      withLatestFrom(this.store.select(selectProblemFilters)),
      switchMap(([action, currentFilters]) => {
        const filters = action.filters || currentFilters;
        return this.problemService.getProblems(filters, action.page, action.limit).pipe(
          map(response => ProblemActions.loadProblemsSuccess({ response })),
          catchError(error => of(ProblemActions.loadProblemsFailure({ 
            error: error.message || 'Failed to load problems' 
          })))
        );
      })
    )
  );

  loadProblem$ = createEffect(() =>
    this.actions$.pipe(
      ofType(ProblemActions.loadProblem),
      switchMap(action =>
        this.problemService.getProblem(action.id).pipe(
          map(problem => ProblemActions.loadProblemSuccess({ problem })),
          catchError(error => of(ProblemActions.loadProblemFailure({ 
            error: error.message || 'Failed to load problem' 
          })))
        )
      )
    )
  );

  loadProblemBySlug$ = createEffect(() =>
    this.actions$.pipe(
      ofType(ProblemActions.loadProblemBySlug),
      switchMap(action =>
        this.problemService.getProblemBySlug(action.slug).pipe(
          map(problem => ProblemActions.loadProblemSuccess({ problem })),
          catchError(error => of(ProblemActions.loadProblemFailure({ 
            error: error.message || 'Failed to load problem' 
          })))
        )
      )
    )
  );

  loadTags$ = createEffect(() =>
    this.actions$.pipe(
      ofType(ProblemActions.loadTags),
      switchMap(() =>
        this.problemService.getAvailableTags().pipe(
          map(tags => ProblemActions.loadTagsSuccess({ tags })),
          catchError(error => of(ProblemActions.loadTagsFailure({ 
            error: error.message || 'Failed to load tags' 
          })))
        )
      )
    )
  );

  // Reload problems when filters change
  updateFilters$ = createEffect(() =>
    this.actions$.pipe(
      ofType(ProblemActions.updateFilters),
      map(() => ProblemActions.loadProblems({ page: 1 }))
    )
  );

  resetFilters$ = createEffect(() =>
    this.actions$.pipe(
      ofType(ProblemActions.resetFilters),
      map(() => ProblemActions.loadProblems({ page: 1 }))
    )
  );

  // Code Execution Effects
  runCode$ = createEffect(() =>
    this.actions$.pipe(
      ofType(ProblemActions.runCode),
      switchMap(action =>
        this.executionService.runCode(action.submission).pipe(
          map(result => ProblemActions.runCodeSuccess({ result })),
          catchError(error => of(ProblemActions.runCodeFailure({ 
            error: error.message || 'Failed to run code' 
          })))
        )
      )
    )
  );

  submitCode$ = createEffect(() =>
    this.actions$.pipe(
      ofType(ProblemActions.submitCode),
      switchMap(action =>
        this.executionService.submitCode(action.submission).pipe(
          map(result => ProblemActions.submitCodeSuccess({ result })),
          catchError(error => of(ProblemActions.submitCodeFailure({ 
            error: error.message || 'Failed to submit code' 
          })))
        )
      )
    )
  );
}