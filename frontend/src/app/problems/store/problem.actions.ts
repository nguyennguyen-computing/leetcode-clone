import { createAction, props } from '@ngrx/store';
import { Problem, ProblemListResponse, ProblemFilters } from '../models/problem.models';
import { CodeSubmission, CodeExecutionResult } from '../components/code-editor/code-editor.component';

// Load Problems Actions
export const loadProblems = createAction(
  '[Problem] Load Problems',
  props<{ filters?: Partial<ProblemFilters>; page?: number; limit?: number }>()
);

export const loadProblemsSuccess = createAction(
  '[Problem] Load Problems Success',
  props<{ response: ProblemListResponse }>()
);

export const loadProblemsFailure = createAction(
  '[Problem] Load Problems Failure',
  props<{ error: string }>()
);

// Load Single Problem Actions
export const loadProblem = createAction(
  '[Problem] Load Problem',
  props<{ id: number }>()
);

export const loadProblemBySlug = createAction(
  '[Problem] Load Problem By Slug',
  props<{ slug: string }>()
);

export const loadProblemSuccess = createAction(
  '[Problem] Load Problem Success',
  props<{ problem: Problem }>()
);

export const loadProblemFailure = createAction(
  '[Problem] Load Problem Failure',
  props<{ error: string }>()
);

// Filter Actions
export const updateFilters = createAction(
  '[Problem] Update Filters',
  props<{ filters: Partial<ProblemFilters> }>()
);

export const resetFilters = createAction('[Problem] Reset Filters');

export const setSearchQuery = createAction(
  '[Problem] Set Search Query',
  props<{ query: string }>()
);

// Load Tags Actions
export const loadTags = createAction('[Problem] Load Tags');

export const loadTagsSuccess = createAction(
  '[Problem] Load Tags Success',
  props<{ tags: string[] }>()
);

export const loadTagsFailure = createAction(
  '[Problem] Load Tags Failure',
  props<{ error: string }>()
);

// Clear Selected Problem
export const clearSelectedProblem = createAction('[Problem] Clear Selected Problem');

// Code Execution Actions
export const runCode = createAction(
  '[Problem] Run Code',
  props<{ submission: CodeSubmission }>()
);

export const runCodeSuccess = createAction(
  '[Problem] Run Code Success',
  props<{ result: CodeExecutionResult }>()
);

export const runCodeFailure = createAction(
  '[Problem] Run Code Failure',
  props<{ error: string }>()
);

export const submitCode = createAction(
  '[Problem] Submit Code',
  props<{ submission: CodeSubmission }>()
);

export const submitCodeSuccess = createAction(
  '[Problem] Submit Code Success',
  props<{ result: CodeExecutionResult }>()
);

export const submitCodeFailure = createAction(
  '[Problem] Submit Code Failure',
  props<{ error: string }>()
);

// Clear Execution Results
export const clearExecutionResult = createAction('[Problem] Clear Execution Result');