import { useQuery } from 'react-query'
import { useAuth } from '../contexts/AuthContext'
import { contentApi } from '../api/client'
import { BookOpen, Users, Eye, TrendingUp } from 'lucide-react'

const Dashboard = () => {
  const { user } = useAuth()

  const { data: recommendations } = useQuery(
    'recommendations',
    () => contentApi.getRecommendations(5),
    {
      enabled: !!user,
    }
  )

  const { data: contentData } = useQuery('content-stats', () => contentApi.getContent({ limit: 5 }))

  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">
          Welcome back, {user?.firstName}!
        </h1>
        <p className="mt-2 text-gray-600">
          Here's what's happening with your content today.
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <BookOpen className="h-8 w-8 text-indigo-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total Content</p>
              <p className="text-2xl font-semibold text-gray-900">{contentData?.total || 0}</p>
            </div>
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <Users className="h-8 w-8 text-green-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Active Users</p>
              <p className="text-2xl font-semibold text-gray-900">1,234</p>
            </div>
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <Eye className="h-8 w-8 text-blue-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total Views</p>
              <p className="text-2xl font-semibold text-gray-900">5,678</p>
            </div>
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <TrendingUp className="h-8 w-8 text-purple-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Growth Rate</p>
              <p className="text-2xl font-semibold text-gray-900">+12%</p>
            </div>
          </div>
        </div>
      </div>

      {/* Recommendations */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Recommended for You</h2>
          {recommendations?.data?.length > 0 ? (
            <div className="space-y-4">
              {recommendations.data.slice(0, 3).map((content: any) => (
                <div key={content.id} className="border-l-4 border-indigo-500 pl-4">
                  <h3 className="font-medium text-gray-900">{content.title}</h3>
                  <p className="text-sm text-gray-600 mt-1">{content.description}</p>
                  <div className="flex items-center text-sm text-gray-500 mt-2">
                    <span className="mr-4">By {content.author.firstName} {content.author.lastName}</span>
                    <span>{content.viewCount} views</span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500">No recommendations available at the moment.</p>
          )}
        </div>

        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Recent Activity</h2>
          <div className="space-y-4">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 bg-indigo-100 rounded-full flex items-center justify-center">
                  <BookOpen className="h-4 w-4 text-indigo-600" />
                </div>
              </div>
              <div className="ml-4">
                <p className="text-sm text-gray-900">New content published</p>
                <p className="text-xs text-gray-500">2 hours ago</p>
              </div>
            </div>
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 bg-green-100 rounded-full flex items-center justify-center">
                  <Users className="h-4 w-4 text-green-600" />
                </div>
              </div>
              <div className="ml-4">
                <p className="text-sm text-gray-900">New user registered</p>
                <p className="text-xs text-gray-500">5 hours ago</p>
              </div>
            </div>
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 bg-blue-100 rounded-full flex items-center justify-center">
                  <Eye className="h-4 w-4 text-blue-600" />
                </div>
              </div>
              <div className="ml-4">
                <p className="text-sm text-gray-900">Content view milestone reached</p>
                <p className="text-xs text-gray-500">1 day ago</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Dashboard
