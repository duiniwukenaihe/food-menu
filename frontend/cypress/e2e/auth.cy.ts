describe('Authentication Flow', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('should display login page', () => {
    cy.get('[data-cy=login-link]').click()
    cy.url().should('include', '/login')
    cy.get('[data-cy=login-form]').should('be.visible')
    cy.get('[data-cy=username-input]').should('be.visible')
    cy.get('[data-cy=password-input]').should('be.visible')
    cy.get('[data-cy=login-button]').should('be.visible')
  })

  it('should register a new user', () => {
    const userData = {
      username: `testuser${Date.now()}`,
      email: `test${Date.now()}@example.com`,
      password: 'password123',
      firstName: 'Test',
      lastName: 'User',
    }

    cy.get('[data-cy=register-link]').click()
    cy.url().should('include', '/register')

    cy.get('[data-cy=first-name-input]').type(userData.firstName)
    cy.get('[data-cy=last-name-input]').type(userData.lastName)
    cy.get('[data-cy=username-input]').type(userData.username)
    cy.get('[data-cy=email-input]').type(userData.email)
    cy.get('[data-cy=password-input]').type(userData.password)

    cy.get('[data-cy=register-button]').click()

    // Should redirect to dashboard after successful registration
    cy.url().should('include', '/dashboard')
    cy.get('[data-cy=welcome-message]').should('contain', userData.firstName)
  })

  it('should login with valid credentials', () => {
    // First register a user
    const userData = {
      username: `testuser${Date.now()}`,
      email: `test${Date.now()}@example.com`,
      password: 'password123',
      firstName: 'Test',
      lastName: 'User',
    }

    cy.register(userData)

    // Then logout and login again
    cy.clearLocalStorage()
    cy.visit('/login')

    cy.get('[data-cy=username-input]').type(userData.username)
    cy.get('[data-cy=password-input]').type(userData.password)
    cy.get('[data-cy=login-button]').click()

    // Should redirect to dashboard after successful login
    cy.url().should('include', '/dashboard')
    cy.get('[data-cy=welcome-message]').should('contain', userData.firstName)
  })

  it('should show error for invalid credentials', () => {
    cy.visit('/login')

    cy.get('[data-cy=username-input]').type('invaliduser')
    cy.get('[data-cy=password-input]').type('wrongpassword')
    cy.get('[data-cy=login-button]').click()

    // Should show error message
    cy.get('[data-cy=error-message]').should('be.visible')
  })

  it('should logout successfully', () => {
    const userData = {
      username: `testuser${Date.now()}`,
      email: `test${Date.now()}@example.com`,
      password: 'password123',
      firstName: 'Test',
      lastName: 'User',
    }

    cy.register(userData)
    cy.url().should('include', '/dashboard')

    // Click logout button
    cy.get('[data-cy=logout-button]').click()

    // Should redirect to home and show login option
    cy.url().should('include', '/')
    cy.get('[data-cy=login-link]').should('be.visible')
    cy.get('[data-cy=register-link]').should('be.visible')
  })
})
