import { Routes, Route, Link, useLocation } from 'react-router-dom'
import { Users, FileText, Tags, Settings, BarChart3, UtensilsCrossed } from 'lucide-react'
import DishManagement from './DishManagement'

const AdminDashboard = () => {
  const location = useLocation()
  
  const navigation = [
    { name: 'Dashboard', href: '/admin', icon: BarChart3 },
    { name: 'Users', href: '/admin/users', icon: Users },
    { name: 'Dishes', href: '/admin/dishes', icon: UtensilsCrossed },
    { name: 'Content', href: '/admin/content', icon: FileText },
    { name: 'Categories', href: '/admin/categories', icon: Tags },
    { name: 'Settings', href: '/admin/settings', icon: Settings },
  ]

  const isActive = (path: string) => location.pathname === path

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="flex">
        {/* Sidebar */}
        <div className="hidden md:flex md:w-64 md:flex-col">
          <div className="flex flex-col flex-grow pt-5 bg-white overflow-y-auto">
            <div className="flex items-center flex-shrink-0 px-4">
              <h2 className="text-xl font-semibold text-gray-900">Admin Panel</h2>
            </div>
            <div className="mt-8 flex-1 flex flex-col">
              <nav className="flex-1 px-2 pb-4 space-y-1">
                {navigation.map((item) => {
                  const Icon = item.icon
                  return (
                    <Link
                      key={item.name}
                      to={item.href}
                      className={`group flex items-center px-2 py-2 text-sm font-medium rounded-md ${
                        isActive(item.href)
                          ? 'bg-indigo-100 text-indigo-700'
                          : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                      }`}
                    >
                      <Icon className="mr-3 h-5 w-5" />
                      {item.name}
                    </Link>
                  )
                })}
              </nav>
            </div>
          </div>
        </div>

        {/* Main content */}
        <div className="flex flex-col flex-1">
          <main className="flex-1">
            <Routes>
              <Route path="/" element={<AdminHome />} />
              <Route path="/users" element={<UserManagement />} />
              <Route path="/dishes" element={<DishManagement />} />
              <Route path="/content" element={<ContentManagement />} />
              <Route path="/categories" element={<CategoryManagement />} />
              <Route path="/settings" element={<SettingsPage />} />
            </Routes>
          </main>
        </div>
      </div>
    </div>
  )
}

const AdminHome = () => (
  <div className="p-6">
    <h1 className="text-2xl font-bold text-gray-900 mb-6">Admin Dashboard</h1>
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div className="bg-white p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900">Total Users</h3>
        <p className="text-3xl font-bold text-indigo-600 mt-2">1,234</p>
      </div>
      <div className="bg-white p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900">Total Content</h3>
        <p className="text-3xl font-bold text-green-600 mt-2">567</p>
      </div>
      <div className="bg-white p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900">Categories</h3>
        <p className="text-3xl font-bold text-blue-600 mt-2">12</p>
      </div>
      <div className="bg-white p-6 rounded-lg shadow">
        <h3 className="text-lg font-medium text-gray-900">Active Sessions</h3>
        <p className="text-3xl font-bold text-purple-600 mt-2">89</p>
      </div>
    </div>
  </div>
)

const UserManagement = () => (
  <div className="p-6">
    <h1 className="text-2xl font-bold text-gray-900 mb-6">User Management</h1>
    <div className="bg-white shadow rounded-lg">
      <div className="p-6">
        <p className="text-gray-600">User management functionality would be implemented here.</p>
      </div>
    </div>
  </div>
)

const ContentManagement = () => (
  <div className="p-6">
    <h1 className="text-2xl font-bold text-gray-900 mb-6">Content Management</h1>
    <div className="bg-white shadow rounded-lg">
      <div className="p-6">
        <p className="text-gray-600">Content management functionality would be implemented here.</p>
      </div>
    </div>
  </div>
)

const CategoryManagement = () => (
  <div className="p-6">
    <h1 className="text-2xl font-bold text-gray-900 mb-6">Category Management</h1>
    <div className="bg-white shadow rounded-lg">
      <div className="p-6">
        <p className="text-gray-600">Category management functionality would be implemented here.</p>
      </div>
    </div>
  </div>
)

const SettingsPage = () => (
  <div className="p-6">
    <h1 className="text-2xl font-bold text-gray-900 mb-6">Settings</h1>
    <div className="bg-white shadow rounded-lg">
      <div className="p-6">
        <p className="text-gray-600">Settings functionality would be implemented here.</p>
      </div>
    </div>
  </div>
)

export default AdminDashboard
