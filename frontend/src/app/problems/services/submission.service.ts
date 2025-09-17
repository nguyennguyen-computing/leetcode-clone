import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { 
  Submission, 
  SubmissionFilters, 
  SubmissionListResponse, 
  SubmissionStats,
  DetailedSubmissionResult 
} from '../models/submission.models';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class SubmissionService {
  private readonly apiUrl = `${environment.apiUrl}/submissions`;

  constructor(private http: HttpClient) {}

  getSubmissions(
    filters: Partial<SubmissionFilters> = {},
    page: number = 1,
    limit: number = 20
  ): Observable<SubmissionListResponse> {
    let params = new HttpParams()
      .set('page', page.toString())
      .set('limit', limit.toString());

    if (filters.status && filters.status !== 'all') {
      params = params.set('status', filters.status);
    }

    if (filters.language && filters.language !== 'all') {
      params = params.set('language', filters.language);
    }

    if (filters.problemId) {
      params = params.set('problemId', filters.problemId.toString());
    }

    if (filters.dateRange) {
      params = params.set('startDate', filters.dateRange.start);
      params = params.set('endDate', filters.dateRange.end);
    }

    return this.http.get<SubmissionListResponse>(this.apiUrl, { params });
  }

  getSubmission(id: number): Observable<Submission> {
    return this.http.get<Submission>(`${this.apiUrl}/${id}`);
  }

  getDetailedSubmissionResult(id: number): Observable<DetailedSubmissionResult> {
    return this.http.get<DetailedSubmissionResult>(`${this.apiUrl}/${id}/details`);
  }

  getSubmissionStats(problemId?: number): Observable<SubmissionStats> {
    let params = new HttpParams();
    
    if (problemId) {
      params = params.set('problemId', problemId.toString());
    }

    return this.http.get<SubmissionStats>(`${this.apiUrl}/stats`, { params });
  }

  getUserSubmissionsForProblem(problemId: number): Observable<Submission[]> {
    return this.http.get<Submission[]>(`${this.apiUrl}/problem/${problemId}`);
  }
}