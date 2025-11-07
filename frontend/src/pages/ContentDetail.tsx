import { useParams, useNavigate } from 'react-router-dom'
import { useQuery } from 'react-query'
import { contentApi } from '../api/client'
import { Calendar, User, Eye, ArrowLeft, Tag } from 'lucide-react'

const ContentDetail = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()

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
      <div className="max-w-4xl mx-auto">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-3/4 mb-4"></div>
          <div className="h-4 bg-gray-200 rounded w-1/2 mb-8"></div>
          <div className="space-y-4">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-4 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (!content) {
    return (
      <div className="max-w-4xl mx-auto text-center py-12">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">Content not found</h1>
        <button
          onClick={() => navigate('/content')}
          className="text-indigo-600 hover:text-indigo-800"
        >
          Back to content list
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

  return (
    <div className="max-w-4xl mx-auto">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center text-gray-600 hover:text-gray-900 mb-6"
      >
        <ArrowLeft className="h-4 w-4 mr-2" />
        Back
      </button>

      {content.imageUrl && (
        <img
          src={content.imageUrl}
          alt={content.title}
          className="w-full h-64 object-cover rounded-lg mb-8"
        />
      )}

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

        <div className="prose prose-lg max-w-none">
          <div
            dangerouslySetInnerHTML={{
              __html: content.body.replace(/\n/g, '<br />'),
            }}
          />
        </div>
      </div>
    </div>
  )
}

export default ContentDetail
