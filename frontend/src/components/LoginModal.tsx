import { useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { authApi } from '../api/client'
import { X, User, Lock, Snowflake } from 'lucide-react'

interface LoginModalProps {
  isOpen: boolean
  onClose: () => void
  onSwitchToRegister: () => void
}

const LoginModal: React.FC<LoginModalProps> = ({ isOpen, onClose, onSwitchToRegister }) => {
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  })
  const [isLoading, setIsLoading] = useState(false)
  const { login } = useAuth()

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)

    try {
      const response = await authApi.login(formData)
      login(response.token, response.user)
      onClose()
    } catch (error) {
      // Error is handled by the API client interceptor
    } finally {
      setIsLoading(false)
    }
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-screen items-center justify-center p-4">
        {/* Backdrop */}
        <div 
          className="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
          onClick={onClose}
        />
        
        {/* Modal */}
        <div className="relative w-full max-w-md transform overflow-hidden rounded-2xl bg-white shadow-2xl transition-all">
          {/* Close button */}
          <button
            onClick={onClose}
            className="absolute right-4 top-4 z-10 rounded-full bg-white/80 p-2 text-gray-600 hover:bg-white hover:text-gray-900"
          >
            <X className="h-5 w-5" />
          </button>

          {/* Polar Bear Header */}
          <div className="relative bg-gradient-to-br from-blue-100 to-indigo-100 p-8 text-center">
            {/* Decorative snowflakes */}
            <Snowflake className="absolute left-4 top-4 h-8 w-8 text-blue-300 opacity-50" />
            <Snowflake className="absolute right-6 top-8 h-6 w-6 text-blue-300 opacity-50" />
            <Snowflake className="absolute left-8 bottom-4 h-5 w-5 text-blue-300 opacity-50" />
            
            {/* Polar Bear SVG */}
            <div className="mb-4">
              <svg
                width="120"
                height="120"
                viewBox="0 0 120 120"
                className="mx-auto"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                {/* Polar Bear Body */}
                <ellipse cx="60" cy="70" rx="35" ry="40" fill="white" stroke="#e5e7eb" strokeWidth="2"/>
                
                {/* Polar Bear Head */}
                <circle cx="60" cy="35" r="25" fill="white" stroke="#e5e7eb" strokeWidth="2"/>
                
                {/* Ears */}
                <circle cx="45" cy="20" r="8" fill="white" stroke="#e5e7eb" strokeWidth="2"/>
                <circle cx="75" cy="20" r="8" fill="white" stroke="#e5e7eb" strokeWidth="2"/>
                <circle cx="45" cy="20" r="4" fill="#1f2937"/>
                <circle cx="75" cy="20" r="4" fill="#1f2937"/>
                
                {/* Eyes */}
                <circle cx="52" cy="32" r="3" fill="#1f2937"/>
                <circle cx="68" cy="32" r="3" fill="#1f2937"/>
                <circle cx="53" cy="31" r="1" fill="white"/>
                <circle cx="69" cy="31" r="1" fill="white"/>
                
                {/* Nose */}
                <ellipse cx="60" cy="40" rx="3" ry="2" fill="#1f2937"/>
                
                {/* Mouth */}
                <path d="M60 42 Q55 45 52 42" stroke="#1f2937" strokeWidth="1.5" fill="none" strokeLinecap="round"/>
                <path d="M60 42 Q65 45 68 42" stroke="#1f2937" strokeWidth="1.5" fill="none" strokeLinecap="round"/>
                
                {/* Paws */}
                <ellipse cx="45" cy="95" rx="8" ry="6" fill="white" stroke="#e5e7eb" strokeWidth="2"/>
                <ellipse cx="75" cy="95" rx="8" ry="6" fill="white" stroke="#e5e7eb" strokeWidth="2"/>
                
                {/* Feet */}
                <ellipse cx="40" cy="105" rx="10" ry="7" fill="white" stroke="#e5e7eb" strokeWidth="2"/>
                <ellipse cx="80" cy="105" rx="10" ry="7" fill="white" stroke="#e5e7eb" strokeWidth="2"/>
                
                {/* Cute details */}
                <circle cx="58" cy="38" r="1" fill="#ef4444" opacity="0.6"/> {/* Blush */}
                <circle cx="62" cy="38" r="1" fill="#ef4444" opacity="0.6"/> {/* Blush */}
              </svg>
            </div>
            
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Welcome Back!</h2>
            <p className="text-gray-600">Sign in to access your personalized recipes</p>
          </div>

          {/* Login Form */}
          <div className="p-6">
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-1">
                  Username
                </label>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
                  <input
                    id="username"
                    name="username"
                    type="text"
                    required
                    className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                    placeholder="Enter your username"
                    value={formData.username}
                    onChange={handleChange}
                  />
                </div>
              </div>

              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
                  Password
                </label>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
                  <input
                    id="password"
                    name="password"
                    type="password"
                    required
                    className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                    placeholder="Enter your password"
                    value={formData.password}
                    onChange={handleChange}
                  />
                </div>
              </div>

              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <input
                    id="remember"
                    name="remember"
                    type="checkbox"
                    className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label htmlFor="remember" className="ml-2 block text-sm text-gray-700">
                    Remember me
                  </label>
                </div>
                <button
                  type="button"
                  className="text-sm text-indigo-600 hover:text-indigo-500"
                >
                  Forgot password?
                </button>
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full bg-indigo-600 text-white py-2 px-4 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {isLoading ? 'Signing in...' : 'Sign In'}
              </button>
            </form>

            <div className="mt-6 text-center">
              <p className="text-sm text-gray-600">
                Don't have an account?{' '}
                <button
                  onClick={onSwitchToRegister}
                  className="font-medium text-indigo-600 hover:text-indigo-500"
                >
                  Sign up for free
                </button>
              </p>
            </div>

            {/* Social Login Options */}
            <div className="mt-6">
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-300" />
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-2 bg-white text-gray-500">Or continue with</span>
                </div>
              </div>

              <div className="mt-6 grid grid-cols-2 gap-3">
                <button
                  type="button"
                  className="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-lg shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                >
                  <svg className="h-5 w-5" viewBox="0 0 24 24">
                    <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                    <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                    <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                    <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                  </svg>
                  <span className="ml-2">Google</span>
                </button>

                <button
                  type="button"
                  className="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-lg shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                >
                  <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                  </svg>
                  <span className="ml-2">GitHub</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default LoginModal