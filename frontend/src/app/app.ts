import { Component, signal, inject, OnInit } from '@angular/core';
import { RouterOutlet, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { Store } from '@ngrx/store';

import { NzLayoutModule } from 'ng-zorro-antd/layout';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzIconModule } from 'ng-zorro-antd/icon';

import * as AuthActions from './auth/store/auth.actions';
import { selectCurrentUser, selectIsAuthenticated } from './auth/store/auth.selectors';

@Component({
  selector: 'app-root',
  imports: [
    CommonModule,
    RouterOutlet,
    RouterLink,
    NzLayoutModule,
    NzButtonModule,
    NzIconModule
  ],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App implements OnInit {
  protected readonly title = signal('LeetCode Clone');
  private readonly store = inject(Store);

  currentUser$ = this.store.select(selectCurrentUser);
  isAuthenticated$ = this.store.select(selectIsAuthenticated);

  ngOnInit(): void {
    // Initialize authentication state from localStorage
    this.store.dispatch(AuthActions.initializeAuth());
  }

  logout(): void {
    this.store.dispatch(AuthActions.logout());
  }
}
