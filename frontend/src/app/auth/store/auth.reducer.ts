import { createReducer, on } from '@ngrx/store';
import { AuthState, initialAuthState } from '../models/auth.models';
import * as AuthActions from './auth.actions';

export const authReducer = createReducer(
  initialAuthState,

  // Login
  on(AuthActions.login, (state) => ({
    ...state,
    loading: true,
    error: null
  })),

  on(AuthActions.loginSuccess, (state, { response }) => ({
    ...state,
    user: response.user,
    token: response.token,
    isAuthenticated: true,
    loading: false,
    error: null
  })),

  on(AuthActions.loginFailure, (state, { error }) => ({
    ...state,
    user: null,
    token: null,
    isAuthenticated: false,
    loading: false,
    error
  })),

  // Register
  on(AuthActions.register, (state) => ({
    ...state,
    loading: true,
    error: null
  })),

  on(AuthActions.registerSuccess, (state, { response }) => ({
    ...state,
    user: response.user,
    token: response.token,
    isAuthenticated: true,
    loading: false,
    error: null
  })),

  on(AuthActions.registerFailure, (state, { error }) => ({
    ...state,
    user: null,
    token: null,
    isAuthenticated: false,
    loading: false,
    error
  })),

  // Logout
  on(AuthActions.logout, (state) => ({
    ...state,
    loading: true
  })),

  on(AuthActions.logoutSuccess, (state) => ({
    ...initialAuthState
  })),

  // Password Reset
  on(AuthActions.requestPasswordReset, (state) => ({
    ...state,
    loading: true,
    error: null
  })),

  on(AuthActions.requestPasswordResetSuccess, (state) => ({
    ...state,
    loading: false,
    error: null
  })),

  on(AuthActions.requestPasswordResetFailure, (state, { error }) => ({
    ...state,
    loading: false,
    error
  })),

  on(AuthActions.confirmPasswordReset, (state) => ({
    ...state,
    loading: true,
    error: null
  })),

  on(AuthActions.confirmPasswordResetSuccess, (state) => ({
    ...state,
    loading: false,
    error: null
  })),

  on(AuthActions.confirmPasswordResetFailure, (state, { error }) => ({
    ...state,
    loading: false,
    error
  })),

  // Initialize Auth
  on(AuthActions.setAuthenticatedUser, (state, { user, token }) => ({
    ...state,
    user,
    token,
    isAuthenticated: true,
    loading: false,
    error: null
  })),

  // Clear Error
  on(AuthActions.clearError, (state) => ({
    ...state,
    error: null
  }))
);