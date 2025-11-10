export type UserRole = 'guest' | 'user' | 'admin'

export interface User {
  id: string
  name: string
  email: string
  role: UserRole
}
