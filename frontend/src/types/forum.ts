export interface Thread {
  id: string;
  showcase_repo_id: string;
  author_id: string;
  title: string;
  body: string;
  comment_count: number;
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: string;
  thread_id: string;
  author_id: string;
  parent_id: string | null;
  body: string;
  created_at: string;
  updated_at: string;
}

export interface Notification {
  id: string;
  user_id: string;
  type: string;
  reference_id: string;
  message: string;
  is_read: boolean;
  created_at: string;
}

export interface ThreadList {
  threads: Thread[];
  next_cursor: string;
  has_more: boolean;
}
