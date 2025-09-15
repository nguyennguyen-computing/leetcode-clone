export interface Problem {
  id: number;
  title: string;
  slug: string;
  description: string;
  difficulty: 'Easy' | 'Medium' | 'Hard';
  tags: string[];
  examples: Example[];
  constraints: string;
  templateCode: { [language: string]: string };
  createdAt: string;
  isSolved?: boolean;
  acceptanceRate?: number;
}

export interface Example {
  input: string;
  output: string;
  explanation?: string;
}

export interface ProblemFilters {
  difficulty: string[];
  tags: string[];
  status: 'all' | 'solved' | 'unsolved';
  searchQuery: string;
}

export interface ProblemListResponse {
  problems: Problem[];
  total: number;
  page: number;
  limit: number;
}