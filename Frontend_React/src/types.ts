export type Role = "user" | "admin";
export type Status = "active" | "inactive" | "banned" | string;

export interface User {
  id: number;
  username: string;
  email: string;

  full_name?: string;
  phone?: string;
  gender?: "male" | "female" | "other" | string;
  date_of_birth?: string;       // "YYYY-MM-DD..." từ BE
  avatar_url?: string;

  street?: string;
  city?: string;
  state?: string;
  country?: string;
  postal_code?: string;

  role: Role;
  status?: Status;

  created_at?: string;          // ISO string
  updated_at?: string;
}
