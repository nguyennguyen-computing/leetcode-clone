import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Problem, ProblemListResponse, ProblemFilters } from '../models/problem.models';

@Injectable({
  providedIn: 'root'
})
export class ProblemService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = '/api/v1';

  getProblems(filters?: Partial<ProblemFilters>, page = 1, limit = 20): Observable<ProblemListResponse> {
    let params = new HttpParams()
      .set('page', page.toString())
      .set('limit', limit.toString());

    if (filters?.difficulty?.length) {
      params = params.set('difficulty', filters.difficulty.join(','));
    }

    if (filters?.tags?.length) {
      params = params.set('tags', filters.tags.join(','));
    }

    if (filters?.status && filters.status !== 'all') {
      params = params.set('status', filters.status);
    }

    if (filters?.searchQuery) {
      params = params.set('search', filters.searchQuery);
    }

    return this.http.get<ProblemListResponse>(`${this.apiUrl}/problems`, { params });
  }

  getProblem(id: number): Observable<Problem> {
    return this.http.get<Problem>(`${this.apiUrl}/problems/${id}`);
  }

  getProblemBySlug(slug: string): Observable<Problem> {
    return this.http.get<Problem>(`${this.apiUrl}/problems/slug/${slug}`);
  }

  getAvailableTags(): Observable<string[]> {
    return this.http.get<string[]>(`${this.apiUrl}/problems/tags`);
  }
}