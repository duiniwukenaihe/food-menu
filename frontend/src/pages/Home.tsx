import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useQuery } from 'react-query'
import { contentApi } from '../api/client'
import { Content, Category } from '../types'
import { Search, Filter, Calendar, Eye, User, Star, Clock, ChefHat } from 'lucide-react'

const Home = () => {
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<number | null>(null)
  const [currentPage, setCurrentPage] = useState(1)

  const { data: contentData, isLoading } = useQuery(
    ['content', searchTerm, selectedCategory, currentPage],
    () => contentApi.getContent({
      page: currentPage,
      limit: 6,
      search: searchTerm,
      category: selectedCategory || undefined,
    }),
    {
      keepPreviousData: true,
    }
  )

  const { data: categoriesData } = useQuery('categories', contentApi.getCategories)
  const { data: seasonalData } = useQuery('seasonal', () => contentApi.getContent({ limit: 6, category: 1 }))
  const { data: recommendationsData } = useQuery('recommendations', () => contentApi.getRecommendations(6))

  const content = contentData?.data || []
  const categories = categoriesData?.data || []
  const seasonal = seasonalData?.data || []
  const recommendations = recommendationsData?.data || []
  const totalPages = Math.ceil((contentData?.total || 0) / 6)

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    setCurrentPage(1)
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString()
  }

  const ContentCard = ({ item, showReason = false }: { item: any; showReason?: boolean }) => (
    <div className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow">
      {item.imageUrl && (
        <img
          src={item.imageUrl}
          alt={item.title}
          className="w-full h-48 object-cover rounded-t-lg"
        />
      )}
      <div className="p-6">
        <div className="flex items-center justify-between mb-2">
          <span className="text-sm text-indigo-600 font-medium">
            {item.category?.name || 'Dish'}
          </span>
          <div className="flex items-center text-gray-500 text-sm">
            <Eye className="h-4 w-4 mr-1" />
            {item.viewCount || 0}
          </div>
        </div>
        <h3 className="text-xl font-semibold text-gray-900 mb-2">
          <Link
            to={`/content/${item.id}`}
            className="hover:text-indigo-600 transition-colors"
          >
            {item.title}
          </Link>
        </h3>
        <p className="text-gray-600 mb-4 line-clamp-2">
          {item.description}
        </p>
        {showReason && item.reason && (
          <div className="flex items-center text-sm text-indigo-600 mb-2">
            <Star className="h-4 w-4 mr-1" />
            {item.reason}
          </div>
        )}
        <div className="flex items-center justify-between text-sm text-gray-500">
          <div className="flex items-center">
            <ChefHat className="h-4 w-4 mr-1" />
            {item.author?.firstName} {item.author?.lastName}
          </div>
          <div className="flex items-center">
            <Calendar className="h-4 w-4 mr-1" />
            {formatDate(item.createdAt)}
          </div>
        </div>
      </div>
    </div>
  )

  return (
    <div className="max-w-7xl mx-auto">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-indigo-500 to-purple-600 rounded-lg shadow-xl p-8 mb-8 text-white">
        <div className="text-center">
          <h1 className="text-4xl font-bold mb-4">Discover Delicious Recipes</h1>
          <p className="text-xl mb-6">Explore seasonal dishes and personalized recommendations</p>
          
          {/* Search Bar */}
          <form onSubmit={handleSearch} className="max-w-2xl mx-auto">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
              <input
                type="text"
                placeholder="Search for recipes, ingredients, or cuisines..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-3 rounded-lg text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-white"
              />
            </div>
          </form>
        </div>
      </div>

      {/* Seasonal Recommendations */}
      <section className="mb-12">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-900 flex items-center">
            <Clock className="h-6 w-6 mr-2 text-indigo-600" />
            Seasonal Recommendations
          </h2>
          <Link
            to="/content?category=1"
            className="text-indigo-600 hover:text-indigo-800 font-medium"
          >
            View All →
          </Link>
        </div>
        
        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="bg-white rounded-lg shadow-md p-6 animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-4"></div>
                <div className="h-3 bg-gray-200 rounded w-full mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-5/6 mb-4"></div>
                <div className="h-3 bg-gray-200 rounded w-1/2"></div>
              </div>
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {seasonal.slice(0, 3).map((item: Content) => (
              <ContentCard key={item.id} item={item} />
            ))}
          </div>
        )}
      </section>

      {/* "Guess You Like" Recommendations */}
      <section className="mb-12">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-900 flex items-center">
            <Star className="h-6 w-6 mr-2 text-indigo-600" />
            Guess You Like
          </h2>
          <Link
            to="/content"
            className="text-indigo-600 hover:text-indigo-800 font-medium"
          >
            Browse All →
          </Link>
        </div>
        
        {recommendations.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {recommendations.slice(0, 3).map((item: any) => (
              <ContentCard key={item.id} item={item.content || item} showReason={true} />
            ))}
          </div>
        ) : (
          <div className="bg-gray-50 rounded-lg p-8 text-center">
            <p className="text-gray-500">Sign in to get personalized recommendations!</p>
          </div>
        )}
      </section>

      {/* Recent Content */}
      <section className="mb-12">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-900">Recent Recipes</h2>
          <div className="flex gap-2">
            <select
              value={selectedCategory || ''}
              onChange={(e) => setSelectedCategory(e.target.value ? Number(e.target.value) : null)}
              className="px-4 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">All Categories</option>
              {categories.map((category: Category) => (
                <option key={category.id} value={category.id}>
                  {category.name}
                </option>
              ))}
            </select>
          </div>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="bg-white rounded-lg shadow-md p-6 animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-4"></div>
                <div className="h-3 bg-gray-200 rounded w-full mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-5/6 mb-4"></div>
                <div className="h-3 bg-gray-200 rounded w-1/2"></div>
              </div>
            ))}
          </div>
        ) : content.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No content found</p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {content.map((item: Content) => (
                <ContentCard key={item.id} item={item} />
              ))}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="mt-8 flex justify-center">
                <nav className="flex items-center space-x-2">
                  <button
                    onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                    disabled={currentPage === 1}
                    className="px-3 py-1 border border-gray-300 rounded-md disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Previous
                  </button>
                  {[...Array(totalPages)].map((_, i) => (
                    <button
                      key={i}
                      onClick={() => setCurrentPage(i + 1)}
                      className={`px-3 py-1 border rounded-md ${
                        currentPage === i + 1
                          ? 'bg-indigo-600 text-white border-indigo-600'
                          : 'border-gray-300 hover:bg-gray-50'
                      }`}
                    >
                      {i + 1}
                    </button>
                  ))}
                  <button
                    onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                    disabled={currentPage === totalPages}
                    className="px-3 py-1 border border-gray-300 rounded-md disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Next
                  </button>
                </nav>
              </div>
            )}
          </>
        )}
      </section>
    </div>
  )
}

export default Home