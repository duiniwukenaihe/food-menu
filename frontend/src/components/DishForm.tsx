import { useState, useEffect } from 'react'
import { X } from 'lucide-react'
import { Dish, CreateDishRequest, UpdateDishRequest } from '../types'
import MediaUpload from './MediaUpload'

interface DishFormProps {
  dish?: Dish
  onSubmit: (data: CreateDishRequest | UpdateDishRequest) => Promise<void>
  onCancel: () => void
}

const MONTHS = [
  { value: '1', label: 'January' },
  { value: '2', label: 'February' },
  { value: '3', label: 'March' },
  { value: '4', label: 'April' },
  { value: '5', label: 'May' },
  { value: '6', label: 'June' },
  { value: '7', label: 'July' },
  { value: '8', label: 'August' },
  { value: '9', label: 'September' },
  { value: '10', label: 'October' },
  { value: '11', label: 'November' },
  { value: '12', label: 'December' },
]

const DishForm: React.FC<DishFormProps> = ({ dish, onSubmit, onCancel }) => {
  const [formData, setFormData] = useState({
    name: dish?.name || '',
    description: dish?.description || '',
    tags: dish?.tags || '',
    isActive: dish?.isActive ?? true,
    isSeasonal: dish?.isSeasonal || false,
    availableMonths: dish?.availableMonths || '',
    seasonalNote: dish?.seasonalNote || '',
    imageUrl: dish?.imageUrl || '',
    thumbnailUrl: dish?.thumbnailUrl || '',
    galleryUrls: dish?.galleryUrls || '',
  })

  const [errors, setErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)
  const [selectedMonths, setSelectedMonths] = useState<string[]>(
    dish?.availableMonths ? dish.availableMonths.split(',').filter(Boolean) : []
  )

  useEffect(() => {
    setFormData((prev) => ({
      ...prev,
      availableMonths: selectedMonths.join(','),
    }))
  }, [selectedMonths])

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.name.trim()) {
      newErrors.name = 'Name is required'
    } else if (formData.name.length > 255) {
      newErrors.name = 'Name must be less than 255 characters'
    }

    if (formData.description.length > 1000) {
      newErrors.description = 'Description must be less than 1000 characters'
    }

    if (formData.isSeasonal && selectedMonths.length === 0) {
      newErrors.availableMonths = 'Please select at least one month for seasonal dishes'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    try {
      setSubmitting(true)
      await onSubmit(formData)
    } finally {
      setSubmitting(false)
    }
  }

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value, type } = e.target
    const checked = (e.target as HTMLInputElement).checked

    setFormData((prev) => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }))

    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: '' }))
    }
  }

  const toggleMonth = (month: string) => {
    setSelectedMonths((prev) =>
      prev.includes(month)
        ? prev.filter((m) => m !== month)
        : [...prev, month].sort((a, b) => parseInt(a) - parseInt(b))
    )
    if (errors.availableMonths) {
      setErrors((prev) => ({ ...prev, availableMonths: '' }))
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label htmlFor="name" className="block text-sm font-medium text-gray-700">
          Name <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          id="name"
          name="name"
          value={formData.name}
          onChange={handleChange}
          className={`mt-1 block w-full rounded-md shadow-sm sm:text-sm ${
            errors.name
              ? 'border-red-300 focus:border-red-500 focus:ring-red-500'
              : 'border-gray-300 focus:border-indigo-500 focus:ring-indigo-500'
          }`}
          disabled={submitting}
        />
        {errors.name && (
          <p className="mt-1 text-sm text-red-600">{errors.name}</p>
        )}
      </div>

      <div>
        <label htmlFor="description" className="block text-sm font-medium text-gray-700">
          Description
        </label>
        <textarea
          id="description"
          name="description"
          rows={4}
          value={formData.description}
          onChange={handleChange}
          className={`mt-1 block w-full rounded-md shadow-sm sm:text-sm ${
            errors.description
              ? 'border-red-300 focus:border-red-500 focus:ring-red-500'
              : 'border-gray-300 focus:border-indigo-500 focus:ring-indigo-500'
          }`}
          disabled={submitting}
        />
        {errors.description && (
          <p className="mt-1 text-sm text-red-600">{errors.description}</p>
        )}
      </div>

      <div>
        <label htmlFor="tags" className="block text-sm font-medium text-gray-700">
          Tags
        </label>
        <input
          type="text"
          id="tags"
          name="tags"
          value={formData.tags}
          onChange={handleChange}
          placeholder="e.g., vegetarian, spicy, gluten-free"
          className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          disabled={submitting}
        />
        <p className="mt-1 text-sm text-gray-500">Comma-separated tags</p>
      </div>

      <div>
        <label htmlFor="imageUrl" className="block text-sm font-medium text-gray-700 mb-2">
          Main Image
        </label>
        <MediaUpload
          value={formData.imageUrl}
          onChange={(url) => setFormData((prev) => ({ ...prev, imageUrl: url }))}
          label="Upload Main Image"
        />
      </div>

      <div className="flex items-center">
        <input
          type="checkbox"
          id="isSeasonal"
          name="isSeasonal"
          checked={formData.isSeasonal}
          onChange={handleChange}
          className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
          disabled={submitting}
        />
        <label htmlFor="isSeasonal" className="ml-2 block text-sm text-gray-700">
          This is a seasonal dish
        </label>
      </div>

      {formData.isSeasonal && (
        <>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Available Months <span className="text-red-500">*</span>
            </label>
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
              {MONTHS.map((month) => (
                <button
                  key={month.value}
                  type="button"
                  onClick={() => toggleMonth(month.value)}
                  className={`px-3 py-2 text-sm rounded-md transition-colors ${
                    selectedMonths.includes(month.value)
                      ? 'bg-indigo-600 text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                  disabled={submitting}
                >
                  {month.label}
                </button>
              ))}
            </div>
            {errors.availableMonths && (
              <p className="mt-1 text-sm text-red-600">{errors.availableMonths}</p>
            )}
          </div>

          <div>
            <label htmlFor="seasonalNote" className="block text-sm font-medium text-gray-700">
              Seasonal Note
            </label>
            <input
              type="text"
              id="seasonalNote"
              name="seasonalNote"
              value={formData.seasonalNote}
              onChange={handleChange}
              placeholder="e.g., Best in winter months"
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              disabled={submitting}
            />
          </div>
        </>
      )}

      <div className="flex items-center">
        <input
          type="checkbox"
          id="isActive"
          name="isActive"
          checked={formData.isActive}
          onChange={handleChange}
          className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
          disabled={submitting}
        />
        <label htmlFor="isActive" className="ml-2 block text-sm text-gray-700">
          Active (visible to customers)
        </label>
      </div>

      <div className="flex justify-end space-x-3 pt-4 border-t">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          disabled={submitting}
        >
          Cancel
        </button>
        <button
          type="submit"
          className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={submitting}
        >
          {submitting ? 'Saving...' : dish ? 'Update Dish' : 'Create Dish'}
        </button>
      </div>
    </form>
  )
}

export default DishForm
