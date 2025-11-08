import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from 'react-query'
import { AuthProvider } from '../../contexts/AuthContext'
import Layout from '../Layout'

const createTestQueryClient = () => new QueryClient({
  defaultOptions: {
    queries: { retry: false },
    mutations: { retry: false },
  },
})

const renderWithProviders = (ui: React.ReactElement) => {
  const testQueryClient = createTestQueryClient()
  return render(
    <QueryClientProvider client={testQueryClient}>
      <BrowserRouter>
        <AuthProvider>
          {ui}
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>
  )
}

describe('Layout Component', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('renders navigation with login/register links when not authenticated', () => {
    renderWithProviders(<Layout>Test Content</Layout>)
    
    expect(screen.getByText('Full-Stack App')).toBeInTheDocument()
    expect(screen.getByText('Home')).toBeInTheDocument()
    expect(screen.getByText('Browse')).toBeInTheDocument()
    expect(screen.getByText('Login')).toBeInTheDocument()
    expect(screen.getByText('Register')).toBeInTheDocument()
  })

  it('renders user menu when authenticated', () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({
      id: 1,
      username: 'testuser',
      email: 'test@example.com',
      firstName: 'Test',
      lastName: 'User',
      role: 'user',
      isActive: true,
      createdAt: '2023-01-01',
      updatedAt: '2023-01-01'
    }))

    renderWithProviders(<Layout>Test Content</Layout>)
    
    expect(screen.getByText('Welcome, Test')).toBeInTheDocument()
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
  })

  it('renders admin link for admin users', () => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('user', JSON.stringify({
      id: 1,
      username: 'admin',
      email: 'admin@example.com',
      firstName: 'Admin',
      lastName: 'User',
      role: 'admin',
      isActive: true,
      createdAt: '2023-01-01',
      updatedAt: '2023-01-01'
    }))

    renderWithProviders(<Layout>Test Content</Layout>)
    
    expect(screen.getByText('Admin')).toBeInTheDocument()
  })

  it('renders children content', () => {
    renderWithProviders(<Layout>Test Content</Layout>)
    
    expect(screen.getByText('Test Content')).toBeInTheDocument()
  })
})
