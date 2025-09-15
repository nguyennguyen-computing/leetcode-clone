import { Injectable, inject } from '@angular/core';
import { Router } from '@angular/router';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { of } from 'rxjs';
import { map, exhaustMap, catchError, tap } from 'rxjs/operators';

import { AuthService } from '../services/auth.service';
import * as AuthActions from './auth.actions';

@Injectable()
export class AuthEffects {
  private readonly actions$ = inject(Actions);
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  login$ = createEffect(() =>
    this.actions$.pipe(
      ofType(AuthActions.login),
      exhaustMap(action =>
        this.authService.login(action.credentials).pipe(
          map(response => AuthActions.loginSuccess({ response })),
          catchError(error => {
            const errorMessage = error.error?.message || 'Login failed';
            return of(AuthActions.loginFailure({ error: errorMessage }));
          })
        )
      )
    )
  );

  loginSuccess$ = createEffect(() =>
    this.actions$.pipe(
      ofType(AuthActions.loginSuccess),
      tap(() => {
        this.router.navigate(['/problems']);
      })
    ),
    { dispatch: false }
  );

  register$ = createEffect(() =>
    this.actions$.pipe(
      ofType(AuthActions.register),
      exhaustMap(action =>
        this.authService.register(action.userData).pipe(
          map(response => AuthActions.registerSuccess({ response })),
          catchError(error => {
            const errorMessage = error.error?.message || 'Registration failed';
            return of(AuthActions.registerFailure({ error: errorMessage }));
          })
        )
      )
    )
  );

  registerSuccess$ = createEffect(() =>
    this.actions$.pipe(
      ofType(AuthActions.registerSuccess),
      tap(() => {
        this.router.navigate(['/problems']);
      })
    ),
    { dispatch: false }
  );

  logout$ = createEffect(() =>
    this.actions$.pipe(
      ofType(AuthActions.logout),
      tap(() => {
        this.authService.logout();
      }),
      map(() => AuthActions.logoutSuccess())
    )
  );

  logoutSuccess$ = createEffect(() =>
    this.actions$.pipe(
      ofType(AuthActions.logoutSuccess),
      tap(() => {
        this.router.navigate(['/login']);
      })
    ),
    { dispatch: false }
  );

  requestPasswordReset$ = createEffect(() =>
    this.actions$.pipe(
      ofType(AuthActions.requestPasswordReset),
      exhaustMap(action =>
        this.authService.requestPasswordReset(action.request).pipe(
          map(() => AuthActions.requestPasswordResetSuccess()),
          catchError(error => {
            const errorMessage = error.error?.message || 'Password reset request failed';
            return of(AuthActions.requestPasswordResetFailure({ error: errorMessage }));
          })
        )
      )
    )
  );

  confirmPasswordReset$ = createEffect(() =>
    this.actions$.pipe(
      ofType(AuthActions.confirmPasswordReset),
      exhaustMap(action =>
        this.authService.confirmPasswordReset(action.request).pipe(
          map(() => AuthActions.confirmPasswordResetSuccess()),
          catchError(error => {
            const errorMessage = error.error?.message || 'Password reset failed';
            return of(AuthActions.confirmPasswordResetFailure({ error: errorMessage }));
          })
        )
      )
    )
  );

  confirmPasswordResetSuccess$ = createEffect(() =>
    this.actions$.pipe(
      ofType(AuthActions.confirmPasswordResetSuccess),
      tap(() => {
        this.router.navigate(['/login']);
      })
    ),
    { dispatch: false }
  );

  initializeAuth$ = createEffect(() =>
    this.actions$.pipe(
      ofType(AuthActions.initializeAuth),
      map(() => {
        const token = this.authService.getToken();
        const user = this.authService.getCurrentUser();
        
        if (token && user && this.authService.isAuthenticated()) {
          return AuthActions.setAuthenticatedUser({ user, token });
        } else {
          return AuthActions.logoutSuccess();
        }
      })
    )
  );
}