import { Component, inject, OnInit, OnDestroy } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { Store } from '@ngrx/store';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { CommonModule } from '@angular/common';

import { NzFormModule } from 'ng-zorro-antd/form';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzAlertModule } from 'ng-zorro-antd/alert';
import { NzIconModule } from 'ng-zorro-antd/icon';

import * as AuthActions from '../../store/auth.actions';
import { selectAuthLoading, selectAuthError } from '../../store/auth.selectors';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterLink,
    NzFormModule,
    NzInputModule,
    NzButtonModule,
    NzCardModule,
    NzAlertModule,
    NzIconModule
  ],
  template: `
    <div class="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full space-y-8">
        <div>
          <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Sign in to your account
          </h2>
          <p class="mt-2 text-center text-sm text-gray-600">
            Or
            <a routerLink="/register" class="font-medium text-blue-600 hover:text-blue-500">
              create a new account
            </a>
          </p>
        </div>

        <nz-card class="shadow-lg">
          <form nz-form [formGroup]="loginForm" (ngSubmit)="onSubmit()" nzLayout="vertical">
            
            <nz-alert 
              *ngIf="error$ | async as error" 
              [nzMessage]="error" 
              nzType="error" 
              nzShowIcon
              class="mb-4">
            </nz-alert>

            <nz-form-item>
              <nz-form-label [nzRequired]="true" nzFor="email">Email</nz-form-label>
              <nz-form-control nzErrorTip="Please enter a valid email address">
                <input 
                  nz-input 
                  id="email"
                  formControlName="email" 
                  type="email" 
                  placeholder="Enter your email"
                  size="large"
                />
              </nz-form-control>
            </nz-form-item>

            <nz-form-item>
              <nz-form-label [nzRequired]="true" nzFor="password">Password</nz-form-label>
              <nz-form-control nzErrorTip="Please enter your password">
                <input 
                  nz-input 
                  id="password"
                  formControlName="password" 
                  type="password" 
                  placeholder="Enter your password"
                  size="large"
                />
              </nz-form-control>
            </nz-form-item>

            <div class="flex items-center justify-between mb-4">
              <div class="text-sm">
                <a routerLink="/forgot-password" class="font-medium text-blue-600 hover:text-blue-500">
                  Forgot your password?
                </a>
              </div>
            </div>

            <nz-form-item class="mb-0">
              <nz-form-control>
                <button 
                  nz-button 
                  nzType="primary" 
                  nzSize="large"
                  [nzLoading]="loading$ | async"
                  [disabled]="loginForm.invalid"
                  class="w-full"
                  type="submit">
                  Sign in
                </button>
              </nz-form-control>
            </nz-form-item>
          </form>
        </nz-card>
      </div>
    </div>
  `
})
export class LoginComponent implements OnInit, OnDestroy {
  private readonly fb = inject(FormBuilder);
  private readonly store = inject(Store);
  private readonly destroy$ = new Subject<void>();

  loginForm!: FormGroup;
  loading$ = this.store.select(selectAuthLoading);
  error$ = this.store.select(selectAuthError);

  ngOnInit(): void {
    this.initializeForm();
    this.clearErrorOnFormChange();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private initializeForm(): void {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  private clearErrorOnFormChange(): void {
    this.loginForm.valueChanges
      .pipe(takeUntil(this.destroy$))
      .subscribe(() => {
        this.store.dispatch(AuthActions.clearError());
      });
  }

  onSubmit(): void {
    if (this.loginForm.valid) {
      const credentials = this.loginForm.value;
      this.store.dispatch(AuthActions.login({ credentials }));
    } else {
      this.markFormGroupTouched();
    }
  }

  private markFormGroupTouched(): void {
    Object.keys(this.loginForm.controls).forEach(key => {
      const control = this.loginForm.get(key);
      control?.markAsTouched();
    });
  }
}