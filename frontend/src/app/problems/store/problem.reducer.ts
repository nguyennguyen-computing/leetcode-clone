import { createReducer, on } from '@ngrx/store';
import { Problem, ProblemFilters } from '../models/problem.models';
import { CodeExecutionResult } from '../components/code-editor/code-editor.component';
import * as ProblemActions from './problem.actions';

export interface ProblemState {
  problems: Problem[];
  selectedProblem: Problem | null;
  filters: ProblemFilters;
  availableTags: string[];
  loading: boolean;
  error: string | null;
  total: number;
  currentPage: number;
  limit: number;
  // Execution state
  isRunning: boolean;
  isSubmitting: boolean;
  executionResult: CodeExecutionResult | null;
  executionError: string | null;
}

const initialFilters: ProblemFilters = {
  difficulty: [],
  tags: [],
  status: 'all',
  searchQuery: ''
};

export const initialState: ProblemState = {
  problems: [],
  selectedProblem: null,
  filters: initialFilters,
  availableTags: [],
  loading: false,
  error: null,
  total: 0,
  currentPage: 1,
  limit: 20,
  // Execution state
  isRunning: false,
  isSubmitting: false,
  executionResult: null,
  executionError: null
};

export const problemReducer = createReducer(
  initialState,

  // Load Problems
  on(ProblemActions.loadProblems, (state, { page = 1, limit = 20 }) => ({
    ...state,
    loading: true,
    error: null,
    currentPage: page,
    limit
  })),

  on(ProblemActions.loadProblemsSuccess, (state, { response }) => ({
    ...state,
    loading: false,
    problems: response.problems,
    total: response.total,
    currentPage: response.page,
    limit: response.limit
  })),

  on(ProblemActions.loadProblemsFailure, (state, { error }) => ({
    ...state,
    loading: false,
    error
  })),

  // Load Single Problem
  on(ProblemActions.loadProblem, ProblemActions.loadProblemBySlug, (state) => ({
    ...state,
    loading: true,
    error: null
  })),

  on(ProblemActions.loadProblemSuccess, (state, { problem }) => ({
    ...state,
    loading: false,
    selectedProblem: problem
  })),

  on(ProblemActions.loadProblemFailure, (state, { error }) => ({
    ...state,
    loading: false,
    error
  })),

  // Filters
  on(ProblemActions.updateFilters, (state, { filters }) => ({
    ...state,
    filters: { ...state.filters, ...filters }
  })),

  on(ProblemActions.resetFilters, (state) => ({
    ...state,
    filters: initialFilters
  })),

  on(ProblemActions.setSearchQuery, (state, { query }) => ({
    ...state,
    filters: { ...state.filters, searchQuery: query }
  })),

  // Tags
  on(ProblemActions.loadTags, (state) => ({
    ...state,
    loading: true,
    error: null
  })),

  on(ProblemActions.loadTagsSuccess, (state, { tags }) => ({
    ...state,
    loading: false,
    availableTags: tags
  })),

  on(ProblemActions.loadTagsFailure, (state, { error }) => ({
    ...state,
    loading: false,
    error
  })),

  // Clear Selected Problem
  on(ProblemActions.clearSelectedProblem, (state) => ({
    ...state,
    selectedProblem: null
  })),

  // Code Execution
  on(ProblemActions.runCode, (state) => ({
    ...state,
    isRunning: true,
    executionResult: null,
    executionError: null
  })),

  on(ProblemActions.runCodeSuccess, (state, { result }) => ({
    ...state,
    isRunning: false,
    executionResult: result
  })),

  on(ProblemActions.runCodeFailure, (state, { error }) => ({
    ...state,
    isRunning: false,
    executionError: error
  })),

  on(ProblemActions.submitCode, (state) => ({
    ...state,
    isSubmitting: true,
    executionResult: null,
    executionError: null
  })),

  on(ProblemActions.submitCodeSuccess, (state, { result }) => ({
    ...state,
    isSubmitting: false,
    executionResult: result
  })),

  on(ProblemActions.submitCodeFailure, (state, { error }) => ({
    ...state,
    isSubmitting: false,
    executionError: error
  })),

  on(ProblemActions.clearExecutionResult, (state) => ({
    ...state,
    executionResult: null,
    executionError: null
  }))
);