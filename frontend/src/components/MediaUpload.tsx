import { useState, useRef } from 'react'
import { X, Upload, Image as ImageIcon } from 'lucide-react'
import { adminApi } from '../api/client'
import toast from 'react-hot-toast'

interface MediaUploadProps {
  value?: string
  onChange: (url: string) => void
  onDelete?: () => void
  label?: string
  accept?: string
}

const MediaUpload: React.FC<MediaUploadProps> = ({
  value,
  onChange,
  onDelete,
  label = 'Upload Image',
  accept = 'image/*',
}) => {
  const [uploading, setUploading] = useState(false)
  const [progress, setProgress] = useState(0)
  const [preview, setPreview] = useState<string | undefined>(value)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    // Validate file type
    if (!file.type.startsWith('image/')) {
      toast.error('Please select an image file')
      return
    }

    // Validate file size (10MB)
    const maxSize = 10 * 1024 * 1024
    if (file.size > maxSize) {
      toast.error('File size must be less than 10MB')
      return
    }

    // Create preview
    const reader = new FileReader()
    reader.onloadend = () => {
      setPreview(reader.result as string)
    }
    reader.readAsDataURL(file)

    // Upload file
    try {
      setUploading(true)
      setProgress(0)

      const response = await adminApi.uploadFile(file, (uploadProgress) => {
        setProgress(uploadProgress)
      })

      if (response.success && response.data?.url) {
        onChange(response.data.url)
        toast.success('Image uploaded successfully')
      } else {
        throw new Error('Upload failed')
      }
    } catch (error) {
      console.error('Upload error:', error)
      toast.error('Failed to upload image')
      setPreview(undefined)
    } finally {
      setUploading(false)
      setProgress(0)
    }
  }

  const handleDelete = () => {
    setPreview(undefined)
    onChange('')
    if (onDelete) {
      onDelete()
    }
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const handleClick = () => {
    fileInputRef.current?.click()
  }

  return (
    <div className="w-full">
      <input
        ref={fileInputRef}
        type="file"
        accept={accept}
        onChange={handleFileSelect}
        className="hidden"
        disabled={uploading}
      />

      {preview ? (
        <div className="relative">
          <div className="relative w-full h-48 bg-gray-100 rounded-lg overflow-hidden">
            <img
              src={preview}
              alt="Preview"
              className="w-full h-full object-cover"
            />
            {uploading && (
              <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center">
                <div className="text-center">
                  <div className="text-white mb-2">Uploading...</div>
                  <div className="w-48 h-2 bg-gray-300 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-indigo-600 transition-all duration-300"
                      style={{ width: `${progress}%` }}
                    />
                  </div>
                  <div className="text-white text-sm mt-2">{progress}%</div>
                </div>
              </div>
            )}
          </div>
          {!uploading && (
            <button
              type="button"
              onClick={handleDelete}
              className="absolute top-2 right-2 p-1 bg-red-600 text-white rounded-full hover:bg-red-700 transition-colors"
            >
              <X className="w-4 h-4" />
            </button>
          )}
        </div>
      ) : (
        <button
          type="button"
          onClick={handleClick}
          disabled={uploading}
          className="w-full h-48 border-2 border-dashed border-gray-300 rounded-lg flex flex-col items-center justify-center hover:border-indigo-500 hover:bg-gray-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {uploading ? (
            <>
              <Upload className="w-12 h-12 text-gray-400 animate-bounce mb-2" />
              <span className="text-gray-600">Uploading...</span>
              <div className="w-48 h-2 bg-gray-300 rounded-full overflow-hidden mt-2">
                <div
                  className="h-full bg-indigo-600 transition-all duration-300"
                  style={{ width: `${progress}%` }}
                />
              </div>
              <span className="text-gray-500 text-sm mt-2">{progress}%</span>
            </>
          ) : (
            <>
              <ImageIcon className="w-12 h-12 text-gray-400 mb-2" />
              <span className="text-gray-600 mb-1">{label}</span>
              <span className="text-gray-400 text-sm">
                PNG, JPG, GIF up to 10MB
              </span>
            </>
          )}
        </button>
      )}
    </div>
  )
}

export default MediaUpload
