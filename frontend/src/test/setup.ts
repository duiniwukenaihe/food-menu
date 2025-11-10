import '@testing-library/jest-dom'
import { vi } from 'vitest'

// Mock IntersectionObserver
class MockIntersectionObserver {
  disconnect() {}
  observe() {}
  unobserve() {}
  takeRecords() {
    return []
  }
}

global.IntersectionObserver = MockIntersectionObserver as unknown as typeof IntersectionObserver

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})
