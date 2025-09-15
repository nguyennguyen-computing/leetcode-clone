import { Component, inject, OnInit, OnDestroy } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule, AbstractControl } from '@angular/forms';
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
  selector: 'app-register',
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
    <div class="min-h-screen bg-gradient-to-br from-purple-50 via-white to-blue-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full space-y-8">
        <div>
          <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Create your account
          </h2>
          <p class="mt-2 text-center text-sm text-gray-600">
            Or
            <a routerLink="/login" class="font-medium text-blue-600 hover:text-blue-500">
              sign in to your existing account
            </a>
          </p>
        </div>

        <nz-card class="shadow-lg">
          <form nz-form [formGroup]="registerForm" (ngSubmit)="onSubmit()" nzLayout="vertical">
            
            <nz-alert 
              *ngIf="error$ | async as error" 
              [nzMessage]="error" 
              nzType="error" 
              nzShowIcon
              class="mb-4">
            </nz-alert>

            <nz-form-item>
              <nz-form-label [nzRequired]="true" nzFor="username">Username</nz-form-label>
              <nz-form-control nzErrorTip="Username must be 3-50 characters and contain only letters, numbers, and underscores">
                <input 
                  nz-input 
                  id="username"
                  formControlName="username" 
                  type="text" 
                  placeholder="Enter your username"
                  size="large"
                />
              </nz-form-control>
            </nz-form-item>

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
              <nz-form-control nzErrorTip="Password must be at least 6 characters long">
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

            <nz-form-item>
              <nz-form-label [nzRequired]="true" nzFor="confirmPassword">Confirm Password</nz-form-label>
              <nz-form-control nzErrorTip="Passwords do not match">
                <input 
                  nz-input 
                  id="confirmPassword"
                  formControlName="confirmPassword" 
                  type="password" 
                  placeholder="Confirm your password"
                  size="large"
                />
              </nz-form-control>
            </nz-form-item>

            <nz-form-item class="mb-0">
              <nz-form-control>
                <button 
                  nz-button 
                  nzType="primary" 
                  nzSize="large"
                  [nzLoading]="loading$ | async"
                  [disabled]="registerForm.invalid"
                  class="w-full"
                  type="submit">
                  Create Account
                </button>
              </nz-form-control>
            </nz-form-item>
          </form>
        </nz-card>
      </div>
    </div>
  `
})
export class RegisterComponent implements OnInit, OnDestroy {
  private readonly fb = inject(FormBuilder);
  private readonly store = inject(Store);
  private readonly destroy$ = new Subject<void>();

  registerForm!: FormGroup;
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
    this.registerForm = this.fb.group({
      username: ['', [
        Validators.required, 
        Validators.minLength(3), 
        Validators.maxLength(50),
        this.usernameValidator
      ]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      confirmPassword: ['', [Validators.required]]
    }, { validators: this.passwordMatchValidator });
  }

  private clearErrorOnFormChange(): void {
    this.registerForm.valueChanges
      .pipe(takeUntil(this.destroy$))
      .subscribe(() => {
        this.store.dispatch(AuthActions.clearError());
      });
  }

  private usernameValidator(control: AbstractControl) {
    const value = control.value;
    if (!value) return null;
    
    const usernameRegex = /^[a-zA-Z0-9_]+$/;
    return usernameRegex.test(value) ? null : { invalidUsername: true };
  }

  private passwordMatchValidator(form: AbstractControl) {
    const password = form.get('password');
    const confirmPassword = form.get('confirmPassword');
    
    if (!password || !confirmPassword) return null;
    
    return password.value === confirmPassword.value ? null : { passwordMismatch: true };
  }

  onSubmit(): void {
    if (this.registerForm.valid) {
      const { confirmPassword, ...userData } = this.registerForm.value;
      this.store.dispatch(AuthActions.register({ userData }));
    } else {
      this.markFormGroupTouched();
    }
  }

  private markFormGroupTouched(): void {
    Object.keys(this.registerForm.controls).forEach(key => {
      const control = this.registerForm.get(key);
      control?.markAsTouched();
    });
  }
}