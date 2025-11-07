import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from 'react-query'
import { AuthProvider } from '../../contexts/AuthContext'
import Login from '../pages/Login'
import { authApi } from '../../api/client'

// Mock the API
jest.mock('../../api/client')
const mockedAuthApi = authApi as jest.Mocked<typeof authApi>

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

describe('Login Page', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    localStorage.clear()
  })

  it('renders login form', () => {
    renderWithProviders(<Login />)
    
    expect(screen.getByText('Sign in to your account')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('Username')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('Password')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Sign in' })).toBeInTheDocument()
  })

  it('shows register link', () => {
    renderWithProviders(<Login />)
    
    const registerLink = screen.getByText('create a new account')
    expect(registerLink).toBeInTheDocument()
    expect(registerLink.closest('a')).toHaveAttribute('href', '/register')
  })

  it('updates form values on input change', () => {
    renderWithProviders(<Login />)
    
    const usernameInput = screen.getByPlaceholderText('Username')
    const passwordInput = screen.getByPlaceholderText('Password')
    
    fireEvent.change(usernameInput, { target: { value: 'testuser' } })
    fireEvent.change(passwordInput, { target: { value: 'password123' } })
    
    expect(usernameInput).toHaveValue('testuser')
    expect(passwordInput).toHaveValue('password123')
  })

  it('submits form with correct data', async () => {
    const mockResponse = {
      token: 'test-token',
      user: {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        role: 'user',
        isActive: true,
        createdAt: '2023-01-01',
        updatedAt: '2023-01-01'
      }
    }
    
    mockedAuthApi.login.mockResolvedValue(mockResponse)
    
    renderWithProviders(<Login />)
    
    const usernameInput = screen.getByPlaceholderText('Username')
    const passwordInput = screen.getByPlaceholderText('Password')
    const submitButton = screen.getByRole('button', { name: 'Sign in' })
    
    fireEvent.change(usernameInput, { target: { value: 'testuser' } })
    fireEvent.change(passwordInput, { target: { value: 'password123' } })
    fireEvent.click(submitButton)
    
    await waitFor(() => {
      expect(mockedAuthApi.login).toHaveBeenCalledWith({
        username: 'testuser',
        password: 'password123'
      })
    })
  })

  it('shows loading state while submitting', async () => {
    mockedAuthApi.login.mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)))
    
    renderWithProviders(<Login />)
    
    const submitButton = screen.getByRole('button', { name: 'Sign in' })
    fireEvent.click(submitButton)
    
    expect(screen.getByText('Signing in...')).toBeInTheDocument()
  })

  it('validates required fields', () => {
    renderWithProviders(<Login />)
    
    const submitButton = screen.getByRole('button', { name: 'Sign in' })
    fireEvent.click(submitButton)
    
    // HTML5 validation should prevent form submission
    const usernameInput = screen.getByPlaceholderText('Username')
    const passwordInput = screen.getByPlaceholderText('Password')
    
    expect(usernameInput).toBeRequired()
    expect(passwordInput).toBeRequired()
  })
})