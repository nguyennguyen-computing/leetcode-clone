import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { CodeSubmission, CodeExecutionResult } from '../components/code-editor/code-editor.component';

export interface ExecutionRequest {
  code: string;
  language: string;
  problemId: number;
  isSubmission?: boolean;
}

export interface ExecutionResponse {
  success: boolean;
  output?: string;
  error?: string;
  runtime?: number;
  memory?: number;
  testCasesPassed?: number;
  totalTestCases?: number;
  submissionId?: number;
}

@Injectable({
  providedIn: 'root'
})
export class ExecutionService {
  private readonly apiUrl = '/api/v1';

  constructor(private http: HttpClient) {}

  runCode(submission: CodeSubmission): Observable<CodeExecutionResult> {
    const request: ExecutionRequest = {
      code: submission.code,
      language: submission.language,
      problemId: submission.problemId,
      isSubmission: false
    };

    return this.http.post<CodeExecutionResult>(`${this.apiUrl}/execute`, request);
  }

  submitCode(submission: CodeSubmission): Observable<CodeExecutionResult> {
    const request: ExecutionRequest = {
      code: submission.code,
      language: submission.language,
      problemId: submission.problemId,
      isSubmission: true
    };

    return this.http.post<CodeExecutionResult>(`${this.apiUrl}/submissions`, request);
  }

  getSubmissionHistory(problemId?: number, limit: number = 20, offset: number = 0): Observable<any> {
    let params = `limit=${limit}&offset=${offset}`;
    if (problemId) {
      params += `&problemId=${problemId}`;
    }
    
    return this.http.get(`${this.apiUrl}/submissions?${params}`);
  }

  getSubmissionById(submissionId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/submissions/${submissionId}`);
  }
}