import { useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { authApi } from '../api/client'
import { X, User, Mail, Lock, Snowflake } from 'lucide-react'

interface RegisterModalProps {
  isOpen: boolean
  onClose: () => void
  onSwitchToLogin: () => void
}

const RegisterModal: React.FC<RegisterModalProps> = ({ isOpen, onClose, onSwitchToLogin }) => {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    firstName: '',
    lastName: '',
  })
  const [isLoading, setIsLoading] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})
  const { login } = useAuth()

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value
    }))
    
    // Clear error for this field when user starts typing
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }))
    }
  }

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.username.trim()) {
      newErrors.username = 'Username is required'
    } else if (formData.username.length < 3) {
      newErrors.username = 'Username must be at least 3 characters'
    }

    if (!formData.email.trim()) {
      newErrors.email = 'Email is required'
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = 'Please enter a valid email'
    }

    if (!formData.firstName.trim()) {
      newErrors.firstName = 'First name is required'
    }

    if (!formData.lastName.trim()) {
      newErrors.lastName = 'Last name is required'
    }

    if (!formData.password) {
      newErrors.password = 'Password is required'
    } else if (formData.password.length < 6) {
      newErrors.password = 'Password must be at least 6 characters'
    }

    if (!formData.confirmPassword) {
      newErrors.confirmPassword = 'Please confirm your password'
    } else if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) {
      return
    }

    setIsLoading(true)

    try {
      const response = await authApi.register({
        username: formData.username,
        email: formData.email,
        password: formData.password,
        firstName: formData.firstName,
        lastName: formData.lastName,
      })
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
          <div className="relative bg-gradient-to-br from-purple-100 to-pink-100 p-8 text-center">
            {/* Decorative snowflakes */}
            <Snowflake className="absolute left-4 top-4 h-8 w-8 text-purple-300 opacity-50" />
            <Snowflake className="absolute right-6 top-8 h-6 w-6 text-purple-300 opacity-50" />
            <Snowflake className="absolute left-8 bottom-4 h-5 w-5 text-purple-300 opacity-50" />
            
            {/* Polar Bear SVG with party hat */}
            <div className="mb-4">
              <svg
                width="120"
                height="120"
                viewBox="0 0 120 120"
                className="mx-auto"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                {/* Party Hat */}
                <path d="M60 10 L40 35 L80 35 Z" fill="#ec4899" stroke="#be185d" strokeWidth="2"/>
                <circle cx="60" cy="5" r="4" fill="#fbbf24"/>
                
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
                
                {/* Happy Mouth */}
                <path d="M60 42 Q55 46 52 42" stroke="#1f2937" strokeWidth="1.5" fill="none" strokeLinecap="round"/>
                <path d="M60 42 Q65 46 68 42" stroke="#1f2937" strokeWidth="1.5" fill="none" strokeLinecap="round"/>
                
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
            
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Join Our Community!</h2>
            <p className="text-gray-600">Create your account to start your culinary journey</p>
          </div>

          {/* Register Form */}
          <div className="p-6">
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label htmlFor="firstName" className="block text-sm font-medium text-gray-700 mb-1">
                    First Name
                  </label>
                  <input
                    id="firstName"
                    name="firstName"
                    type="text"
                    required
                    className={`w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent ${
                      errors.firstName ? 'border-red-500' : 'border-gray-300'
                    }`}
                    placeholder="First name"
                    value={formData.firstName}
                    onChange={handleChange}
                  />
                  {errors.firstName && (
                    <p className="mt-1 text-sm text-red-600">{errors.firstName}</p>
                  )}
                </div>

                <div>
                  <label htmlFor="lastName" className="block text-sm font-medium text-gray-700 mb-1">
                    Last Name
                  </label>
                  <input
                    id="lastName"
                    name="lastName"
                    type="text"
                    required
                    className={`w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent ${
                      errors.lastName ? 'border-red-500' : 'border-gray-300'
                    }`}
                    placeholder="Last name"
                    value={formData.lastName}
                    onChange={handleChange}
                  />
                  {errors.lastName && (
                    <p className="mt-1 text-sm text-red-600">{errors.lastName}</p>
                  )}
                </div>
              </div>

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
                    className={`w-full pl-10 pr-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent ${
                      errors.username ? 'border-red-500' : 'border-gray-300'
                    }`}
                    placeholder="Choose a username"
                    value={formData.username}
                    onChange={handleChange}
                  />
                </div>
                {errors.username && (
                  <p className="mt-1 text-sm text-red-600">{errors.username}</p>
                )}
              </div>

              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                  Email
                </label>
                <div className="relative">
                  <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
                  <input
                    id="email"
                    name="email"
                    type="email"
                    required
                    className={`w-full pl-10 pr-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent ${
                      errors.email ? 'border-red-500' : 'border-gray-300'
                    }`}
                    placeholder="your@email.com"
                    value={formData.email}
                    onChange={handleChange}
                  />
                </div>
                {errors.email && (
                  <p className="mt-1 text-sm text-red-600">{errors.email}</p>
                )}
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
                    className={`w-full pl-10 pr-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent ${
                      errors.password ? 'border-red-500' : 'border-gray-300'
                    }`}
                    placeholder="Create a password"
                    value={formData.password}
                    onChange={handleChange}
                  />
                </div>
                {errors.password && (
                  <p className="mt-1 text-sm text-red-600">{errors.password}</p>
                )}
              </div>

              <div>
                <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-1">
                  Confirm Password
                </label>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
                  <input
                    id="confirmPassword"
                    name="confirmPassword"
                    type="password"
                    required
                    className={`w-full pl-10 pr-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent ${
                      errors.confirmPassword ? 'border-red-500' : 'border-gray-300'
                    }`}
                    placeholder="Confirm your password"
                    value={formData.confirmPassword}
                    onChange={handleChange}
                  />
                </div>
                {errors.confirmPassword && (
                  <p className="mt-1 text-sm text-red-600">{errors.confirmPassword}</p>
                )}
              </div>

              <div className="flex items-center">
                <input
                  id="terms"
                  name="terms"
                  type="checkbox"
                  required
                  className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label htmlFor="terms" className="ml-2 block text-sm text-gray-700">
                  I agree to the{' '}
                  <button type="button" className="text-indigo-600 hover:text-indigo-500">
                    Terms and Conditions
                  </button>
                </label>
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full bg-indigo-600 text-white py-2 px-4 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {isLoading ? 'Creating Account...' : 'Create Account'}
              </button>
            </form>

            <div className="mt-6 text-center">
              <p className="text-sm text-gray-600">
                Already have an account?{' '}
                <button
                  onClick={onSwitchToLogin}
                  className="font-medium text-indigo-600 hover:text-indigo-500"
                >
                  Sign in here
                </button>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default RegisterModal