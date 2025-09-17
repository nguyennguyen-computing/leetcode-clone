export interface Submission {
  id: number;
  userId: number;
  problemId: number;
  language: string;
  code: string;
  status: SubmissionStatus;
  runtimeMs?: number;
  memoryKb?: number;
  testCasesPassed: number;
  totalTestCases: number;
  errorMessage?: string;
  submittedAt: string;
  problem?: {
    id: number;
    title: string;
    difficulty: 'Easy' | 'Medium' | 'Hard';
  };
}

export type SubmissionStatus = 
  | 'Accepted' 
  | 'Wrong Answer' 
  | 'Time Limit Exceeded' 
  | 'Memory Limit Exceeded' 
  | 'Runtime Error' 
  | 'Compilation Error'
  | 'Internal Error';

export interface SubmissionFilters {
  status: SubmissionStatus | 'all';
  language: string | 'all';
  problemId?: number;
  dateRange?: {
    start: string;
    end: string;
  };
}

export interface SubmissionListResponse {
  submissions: Submission[];
  total: number;
  page: number;
  limit: number;
}

export interface SubmissionStats {
  totalSubmissions: number;
  acceptedSubmissions: number;
  acceptanceRate: number;
  languageStats: { [language: string]: number };
  statusStats: { [status: string]: number };
}

export interface TestCaseResult {
  input: string;
  expectedOutput: string;
  actualOutput?: string;
  passed: boolean;
  runtimeMs?: number;
  memoryKb?: number;
  error?: string;
}

export interface DetailedSubmissionResult {
  submission: Submission;
  testCaseResults: TestCaseResult[];
  overallStats: {
    runtimeMs: number;
    memoryKb: number;
    testCasesPassed: number;
    totalTestCases: number;
  };
}