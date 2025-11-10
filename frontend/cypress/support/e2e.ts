// Import commands.js using ES2015 syntax:
import './commands'

// Alternatively you can use CommonJS syntax:
// require('./commands')

const SESSION_KEY = 'gourmet-explorer:session'

// Custom commands for authentication
Cypress.Commands.add('login', (username: string, password: string) => {
  cy.request({
    method: 'POST',
    url: 'http://localhost:8080/api/v1/auth/login',
    body: {
      username,
      password,
    },
  }).then(response => {
    const session = {
      token: response.body.token,
      user: response.body.user,
    }
    window.localStorage.setItem(SESSION_KEY, JSON.stringify(session))
  })
})

Cypress.Commands.add('register', (userData: Record<string, unknown>) => {
  cy.request({
    method: 'POST',
    url: 'http://localhost:8080/api/v1/auth/register',
    body: userData,
  })
})

// Global beforeEach hook to clear localStorage
beforeEach(() => {
  window.localStorage.removeItem(SESSION_KEY)
})
