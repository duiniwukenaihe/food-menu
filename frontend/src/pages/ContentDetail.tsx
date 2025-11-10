import { useParams, useNavigate, useState } from 'react-router-dom'
import { useQuery } from 'react-query'
import { contentApi } from '../api/client'
import { Calendar, User, Eye, ArrowLeft, Tag, Clock, ChefHat, Users, Play, ChevronLeft, ChevronRight } from 'lucide-react'

const ContentDetail = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [currentImageIndex, setCurrentImageIndex] = useState(0)

  const { data: contentData, isLoading } = useQuery(
    ['content', id],
    () => contentApi.getContentById(Number(id)),
    {
      enabled: !!id,
    }
  )

  const content = contentData?.data

  if (isLoading) {
    return (
      <div className="max-w-6xl mx-auto">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-3/4 mb-4"></div>
          <div className="h-4 bg-gray-200 rounded w-1/2 mb-8"></div>
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 space-y-4">
              <div className="h-64 bg-gray-200 rounded-lg"></div>
              {[...Array(5)].map((_, i) => (
                <div key={i} className="h-4 bg-gray-200 rounded"></div>
              ))}
            </div>
            <div className="space-y-4">
              <div className="h-32 bg-gray-200 rounded"></div>
              <div className="h-32 bg-gray-200 rounded"></div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!content) {
    return (
      <div className="max-w-4xl mx-auto text-center py-12">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">Recipe not found</h1>
        <button
          onClick={() => navigate('/content')}
          className="text-indigo-600 hover:text-indigo-800"
        >
          Back to recipes
        </button>
      </div>
    )
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }

  // Parse images from content (assuming multiple images might be stored as JSON array or comma-separated)
  const images = content.imageUrl 
    ? (content.imageUrl.includes('[') ? JSON.parse(content.imageUrl) : [content.imageUrl])
    : []

  // Parse video URL (could be YouTube, Vimeo, or direct video link)
  const getVideoEmbedUrl = (url: string) => {
    if (!url) return null
    
    // YouTube
    const youtubeMatch = url.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&\n?#]+)/)
    if (youtubeMatch) {
      return `https://www.youtube.com/embed/${youtubeMatch[1]}`
    }
    
    // Vimeo
    const vimeoMatch = url.match(/vimeo\.com\/(\d+)/)
    if (vimeoMatch) {
      return `https://player.vimeo.com/video/${vimeoMatch[1]}`
    }
    
    // Direct video URL
    if (url.match(/\.(mp4|webm|ogg)$/i)) {
      return url
    }
    
    return null
  }

  const videoEmbedUrl = content.videoUrl ? getVideoEmbedUrl(content.videoUrl) : null

  // Parse cooking steps (assuming they're in the body with step markers)
  const parseSteps = (body: string) => {
    const steps = body.split(/\n(?=\d+\.)/).filter(step => step.trim())
    return steps.map((step, index) => ({
      step: index + 1,
      instruction: step.replace(/^\d+\.\s*/, '').trim()
    }))
  }

  const cookingSteps = parseSteps(content.body)

  const nextImage = () => {
    setCurrentImageIndex((prev) => (prev + 1) % images.length)
  }

  const prevImage = () => {
    setCurrentImageIndex((prev) => (prev - 1 + images.length) % images.length)
  }

  return (
    <div className="max-w-6xl mx-auto">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center text-gray-600 hover:text-gray-900 mb-6"
      >
        <ArrowLeft className="h-4 w-4 mr-2" />
        Back
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Image Carousel */}
          {images.length > 0 && (
            <div className="relative bg-white rounded-lg shadow-lg overflow-hidden">
              <div className="relative h-96">
                <img
                  src={images[currentImageIndex]}
                  alt={`${content.title} - Image ${currentImageIndex + 1}`}
                  className="w-full h-full object-cover"
                />
                
                {images.length > 1 && (
                  <>
                    <button
                      onClick={prevImage}
                      className="absolute left-4 top-1/2 transform -translate-y-1/2 bg-white/80 hover:bg-white text-gray-800 p-2 rounded-full shadow-lg"
                    >
                      <ChevronLeft className="h-5 w-5" />
                    </button>
                    <button
                      onClick={nextImage}
                      className="absolute right-4 top-1/2 transform -translate-y-1/2 bg-white/80 hover:bg-white text-gray-800 p-2 rounded-full shadow-lg"
                    >
                      <ChevronRight className="h-5 w-5" />
                    </button>
                    
                    <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex space-x-2">
                      {images.map((_, index) => (
                        <button
                          key={index}
                          onClick={() => setCurrentImageIndex(index)}
                          className={`w-2 h-2 rounded-full transition-colors ${
                            index === currentImageIndex ? 'bg-white' : 'bg-white/50'
                          }`}
                        />
                      ))}
                    </div>
                  </>
                )}
              </div>
            </div>
          )}

          {/* Video Embed */}
          {videoEmbedUrl && (
            <div className="bg-white rounded-lg shadow-lg overflow-hidden">
              <div className="aspect-w-16 aspect-h-9">
                {videoEmbedUrl.includes('youtube') || videoEmbedUrl.includes('vimeo') ? (
                  <iframe
                    src={videoEmbedUrl}
                    title={`${content.title} - Video Tutorial`}
                    className="w-full h-96"
                    frameBorder="0"
                    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                    allowFullScreen
                  />
                ) : (
                  <video
                    src={videoEmbedUrl}
                    controls
                    className="w-full h-96"
                    title={`${content.title} - Video Tutorial`}
                  >
                    Your browser does not support the video tag.
                  </video>
                )}
              </div>
            </div>
          )}

          {/* Recipe Header */}
          <div className="bg-white rounded-lg shadow-lg p-8">
            <div className="mb-6">
              <div className="flex items-center justify-between mb-4">
                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800">
                  {content.category.name}
                </span>
                <div className="flex items-center text-gray-500 text-sm">
                  <Eye className="h-4 w-4 mr-1" />
                  {content.viewCount} views
                </div>
              </div>
              
              <h1 className="text-3xl font-bold text-gray-900 mb-4">{content.title}</h1>
              
              {content.description && (
                <p className="text-lg text-gray-600 mb-6">{content.description}</p>
              )}

              <div className="flex items-center justify-between border-t border-b py-4">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <div className="h-10 w-10 bg-gray-300 rounded-full flex items-center justify-center">
                      <User className="h-6 w-6 text-gray-600" />
                    </div>
                  </div>
                  <div className="ml-3">
                    <p className="text-sm font-medium text-gray-900">
                      {content.author.firstName} {content.author.lastName}
                    </p>
                    <p className="text-sm text-gray-500">{content.author.email}</p>
                  </div>
                </div>
                
                <div className="flex items-center text-sm text-gray-500">
                  <Calendar className="h-4 w-4 mr-1" />
                  {formatDate(content.createdAt)}
                </div>
              </div>
            </div>

            {content.tags && (
              <div className="mb-6">
                <div className="flex items-center">
                  <Tag className="h-4 w-4 mr-2 text-gray-500" />
                  <div className="flex flex-wrap gap-2">
                    {content.tags.split(',').map((tag: string, index: number) => (
                      <span
                        key={index}
                        className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-gray-100 text-gray-800"
                      >
                        {tag.trim()}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Step-by-Step Instructions */}
          <div className="bg-white rounded-lg shadow-lg p-8">
            <h2 className="text-2xl font-bold text-gray-900 mb-6 flex items-center">
              <ChefHat className="h-6 w-6 mr-2 text-indigo-600" />
              Step-by-Step Instructions
            </h2>
            
            <div className="space-y-6">
              {cookingSteps.map((step) => (
                <div key={step.step} className="flex">
                  <div className="flex-shrink-0">
                    <div className="flex items-center justify-center w-8 h-8 bg-indigo-600 text-white rounded-full font-bold text-sm">
                      {step.step}
                    </div>
                  </div>
                  <div className="ml-4 flex-1">
                    <p className="text-gray-700 leading-relaxed">{step.instruction}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Recipe Info Card */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Recipe Details</h3>
            
            <div className="space-y-3">
              <div className="flex items-center">
                <Clock className="h-5 w-5 text-gray-400 mr-3" />
                <div>
                  <p className="text-sm text-gray-500">Prep Time</p>
                  <p className="font-medium">30 minutes</p>
                </div>
              </div>
              
              <div className="flex items-center">
                <Clock className="h-5 w-5 text-gray-400 mr-3" />
                <div>
                  <p className="text-sm text-gray-500">Cook Time</p>
                  <p className="font-medium">45 minutes</p>
                </div>
              </div>
              
              <div className="flex items-center">
                <Users className="h-5 w-5 text-gray-400 mr-3" />
                <div>
                  <p className="text-sm text-gray-500">Servings</p>
                  <p className="font-medium">4 people</p>
                </div>
              </div>
            </div>
          </div>

          {/* Ingredients */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Ingredients</h3>
            
            <div className="space-y-2">
              {/* This would ideally come from structured data in the backend */}
              <div className="flex justify-between py-2 border-b">
                <span className="text-gray-700">Flour</span>
                <span className="text-gray-600">2 cups</span>
              </div>
              <div className="flex justify-between py-2 border-b">
                <span className="text-gray-700">Sugar</span>
                <span className="text-gray-600">1 cup</span>
              </div>
              <div className="flex justify-between py-2 border-b">
                <span className="text-gray-700">Eggs</span>
                <span className="text-gray-600">2 large</span>
              </div>
              <div className="flex justify-between py-2 border-b">
                <span className="text-gray-700">Butter</span>
                <span className="text-gray-600">1/2 cup</span>
              </div>
              <div className="flex justify-between py-2 border-b">
                <span className="text-gray-700">Milk</span>
                <span className="text-gray-600">1 cup</span>
              </div>
            </div>
          </div>

          {/* Tips & Notes */}
          <div className="bg-indigo-50 rounded-lg p-6">
            <h3 className="text-lg font-semibold text-indigo-900 mb-3">Chef's Tips</h3>
            <ul className="space-y-2 text-sm text-indigo-700">
              <li>• Make sure all ingredients are at room temperature</li>
              <li>• Don't overmix the batter</li>
              <li>• Preheat your oven for best results</li>
              <li>• Let cool completely before serving</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ContentDetail