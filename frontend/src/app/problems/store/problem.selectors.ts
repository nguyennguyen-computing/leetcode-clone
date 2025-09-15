import { createFeatureSelector, createSelector } from '@ngrx/store';
import { ProblemState } from './problem.reducer';

export const selectProblemState = createFeatureSelector<ProblemState>('problems');

export const selectProblems = createSelector(
  selectProblemState,
  (state) => state.problems
);

export const selectSelectedProblem = createSelector(
  selectProblemState,
  (state) => state.selectedProblem
);

export const selectProblemFilters = createSelector(
  selectProblemState,
  (state) => state.filters
);

export const selectAvailableTags = createSelector(
  selectProblemState,
  (state) => state.availableTags
);

export const selectProblemLoading = createSelector(
  selectProblemState,
  (state) => state.loading
);

export const selectProblemError = createSelector(
  selectProblemState,
  (state) => state.error
);

export const selectProblemPagination = createSelector(
  selectProblemState,
  (state) => ({
    total: state.total,
    currentPage: state.currentPage,
    limit: state.limit,
    totalPages: Math.ceil(state.total / state.limit)
  })
);

export const selectFilteredProblems = createSelector(
  selectProblems,
  selectProblemFilters,
  (problems, filters) => {
    let filtered = [...problems];

    // Apply difficulty filter
    if (filters.difficulty.length > 0) {
      filtered = filtered.filter(problem => 
        filters.difficulty.includes(problem.difficulty)
      );
    }

    // Apply tags filter
    if (filters.tags.length > 0) {
      filtered = filtered.filter(problem =>
        filters.tags.some(tag => problem.tags.includes(tag))
      );
    }

    // Apply status filter
    if (filters.status !== 'all') {
      filtered = filtered.filter(problem => {
        if (filters.status === 'solved') {
          return problem.isSolved === true;
        } else if (filters.status === 'unsolved') {
          return problem.isSolved !== true;
        }
        return true;
      });
    }

    // Apply search query
    if (filters.searchQuery) {
      const query = filters.searchQuery.toLowerCase();
      filtered = filtered.filter(problem =>
        problem.title.toLowerCase().includes(query) ||
        problem.description.toLowerCase().includes(query) ||
        problem.tags.some(tag => tag.toLowerCase().includes(query))
      );
    }

    return filtered;
  }
);

export const selectProblemStats = createSelector(
  selectProblems,
  (problems) => {
    const total = problems.length;
    const solved = problems.filter(p => p.isSolved).length;
    const easy = problems.filter(p => p.difficulty === 'Easy').length;
    const medium = problems.filter(p => p.difficulty === 'Medium').length;
    const hard = problems.filter(p => p.difficulty === 'Hard').length;
    const easySolved = problems.filter(p => p.difficulty === 'Easy' && p.isSolved).length;
    const mediumSolved = problems.filter(p => p.difficulty === 'Medium' && p.isSolved).length;
    const hardSolved = problems.filter(p => p.difficulty === 'Hard' && p.isSolved).length;

    return {
      total,
      solved,
      acceptanceRate: total > 0 ? Math.round((solved / total) * 100) : 0,
      byDifficulty: {
        easy: { total: easy, solved: easySolved },
        medium: { total: medium, solved: mediumSolved },
        hard: { total: hard, solved: hardSolved }
      }
    };
  }
);

// Execution Selectors
export const selectIsRunning = createSelector(
  selectProblemState,
  (state) => state.isRunning
);

export const selectIsSubmitting = createSelector(
  selectProblemState,
  (state) => state.isSubmitting
);

export const selectExecutionResult = createSelector(
  selectProblemState,
  (state) => state.executionResult
);

export const selectExecutionError = createSelector(
  selectProblemState,
  (state) => state.executionError
);

export const selectExecutionState = createSelector(
  selectProblemState,
  (state) => ({
    isRunning: state.isRunning,
    isSubmitting: state.isSubmitting,
    result: state.executionResult,
    error: state.executionError
  })
);