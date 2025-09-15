import { createAction, props } from '@ngrx/store';
import { 
  LoginRequest, 
  RegisterRequest, 
  AuthResponse, 
  User,
  PasswordResetRequest,
  PasswordResetConfirmRequest
} from '../models/auth.models';

// Login Actions
export const login = createAction(
  '[Auth] Login',
  props<{ credentials: LoginRequest }>()
);

export const loginSuccess = createAction(
  '[Auth] Login Success',
  props<{ response: AuthResponse }>()
);

export const loginFailure = createAction(
  '[Auth] Login Failure',
  props<{ error: string }>()
);

// Register Actions
export const register = createAction(
  '[Auth] Register',
  props<{ userData: RegisterRequest }>()
);

export const registerSuccess = createAction(
  '[Auth] Register Success',
  props<{ response: AuthResponse }>()
);

export const registerFailure = createAction(
  '[Auth] Register Failure',
  props<{ error: string }>()
);

// Logout Actions
export const logout = createAction('[Auth] Logout');

export const logoutSuccess = createAction('[Auth] Logout Success');

// Password Reset Actions
export const requestPasswordReset = createAction(
  '[Auth] Request Password Reset',
  props<{ request: PasswordResetRequest }>()
);

export const requestPasswordResetSuccess = createAction(
  '[Auth] Request Password Reset Success'
);

export const requestPasswordResetFailure = createAction(
  '[Auth] Request Password Reset Failure',
  props<{ error: string }>()
);

export const confirmPasswordReset = createAction(
  '[Auth] Confirm Password Reset',
  props<{ request: PasswordResetConfirmRequest }>()
);

export const confirmPasswordResetSuccess = createAction(
  '[Auth] Confirm Password Reset Success'
);

export const confirmPasswordResetFailure = createAction(
  '[Auth] Confirm Password Reset Failure',
  props<{ error: string }>()
);

// Initialize Auth State
export const initializeAuth = createAction('[Auth] Initialize');

export const setAuthenticatedUser = createAction(
  '[Auth] Set Authenticated User',
  props<{ user: User; token: string }>()
);

// Clear Error
export const clearError = createAction('[Auth] Clear Error');