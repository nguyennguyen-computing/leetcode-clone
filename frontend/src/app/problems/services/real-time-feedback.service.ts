import { Injectable, OnDestroy } from '@angular/core';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { filter, takeUntil } from 'rxjs/operators';

export interface SubmissionUpdate {
  submissionId: number;
  status: string;
  progress?: number;
  message?: string;
  result?: any;
}

@Injectable({
  providedIn: 'root'
})
export class RealTimeFeedbackService implements OnDestroy {
  private ws: WebSocket | null = null;
  private destroy$ = new Subject<void>();
  private connectionStatus$ = new BehaviorSubject<'connecting' | 'connected' | 'disconnected'>('disconnected');
  private submissionUpdates$ = new Subject<SubmissionUpdate>();

  constructor() {}

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
    this.disconnect();
  }

  connect(token: string): Observable<'connecting' | 'connected' | 'disconnected'> {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return this.connectionStatus$.asObservable();
    }

    this.connectionStatus$.next('connecting');

    try {
      this.simulateConnection();
    } catch (error) {
      console.error('WebSocket connection failed:', error);
      this.connectionStatus$.next('disconnected');
    }

    return this.connectionStatus$.asObservable();
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.connectionStatus$.next('disconnected');
  }

  subscribeToSubmissionUpdates(submissionId?: number): Observable<SubmissionUpdate> {
    return this.submissionUpdates$.asObservable().pipe(
      filter(update => !submissionId || update.submissionId === submissionId),
      takeUntil(this.destroy$)
    );
  }

  private simulateConnection() {
    // Simulate WebSocket connection for demo purposes
    setTimeout(() => {
      this.connectionStatus$.next('connected');
      
      // Simulate receiving updates
      this.simulateSubmissionUpdates();
    }, 1000);
  }

  private simulateSubmissionUpdates() {
    // This would be replaced with actual WebSocket message handling
    // For now, we'll just simulate some updates for demo purposes
    
    // Example: Simulate a submission being processed
    setTimeout(() => {
      this.submissionUpdates$.next({
        submissionId: 1,
        status: 'processing',
        progress: 25,
        message: 'Compiling code...'
      });
    }, 2000);

    setTimeout(() => {
      this.submissionUpdates$.next({
        submissionId: 1,
        status: 'processing',
        progress: 50,
        message: 'Running test cases...'
      });
    }, 4000);

    setTimeout(() => {
      this.submissionUpdates$.next({
        submissionId: 1,
        status: 'completed',
        progress: 100,
        message: 'Submission completed',
        result: {
          success: true,
          testCasesPassed: 10,
          totalTestCases: 10,
          runtime: 45,
          memory: 12800
        }
      });
    }, 6000);
  }

  // Method to send submission for real-time processing
  submitForRealTimeProcessing(submissionData: any): number {
    // In a real implementation, this would send the submission to the server
    // and return the submission ID for tracking
    
    const submissionId = Math.floor(Math.random() * 10000);
    
    // Simulate immediate acknowledgment
    setTimeout(() => {
      this.submissionUpdates$.next({
        submissionId,
        status: 'queued',
        progress: 0,
        message: 'Submission queued for processing'
      });
    }, 100);

    return submissionId;
  }

  getConnectionStatus(): Observable<'connecting' | 'connected' | 'disconnected'> {
    return this.connectionStatus$.asObservable();
  }

  isConnected(): boolean {
    return this.connectionStatus$.value === 'connected';
  }
}