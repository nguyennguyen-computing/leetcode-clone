import { Routes } from '@angular/router';
import { authGuard, guestGuard } from './auth/guards/auth.guard';

export const routes: Routes = [
  {
    path: '',
    redirectTo: '/problems',
    pathMatch: 'full'
  },
  {
    path: 'login',
    loadComponent: () => import('./auth/components/login/login.component').then(m => m.LoginComponent),
    canActivate: [guestGuard]
  },
  {
    path: 'register',
    loadComponent: () => import('./auth/components/register/register.component').then(m => m.RegisterComponent),
    canActivate: [guestGuard]
  },
  {
    path: 'problems',
    loadComponent: () => import('./problems/components/problem-list/problem-list.component').then(m => m.ProblemListComponent),
    canActivate: [authGuard]
  },
  {
    path: 'problems/:id',
    loadComponent: () => import('./problems/components/problem-detail/problem-detail.component').then(m => m.ProblemDetailComponent),
    canActivate: [authGuard]
  },
  {
    path: 'problems/:id/solve',
    loadComponent: () => import('./problems/components/problem-solve/problem-solve.component').then(m => m.ProblemSolveComponent),
    canActivate: [authGuard]
  },
  {
    path: '**',
    redirectTo: '/problems'
  }
];
