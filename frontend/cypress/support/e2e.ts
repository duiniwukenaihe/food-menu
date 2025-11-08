// Import commands.js using ES2015 syntax:
import './commands'

// Alternatively you can use CommonJS syntax:
// require('./commands')

// Custom commands for authentication
Cypress.Commands.add('login', (username: string, password: string) => {
  cy.request({
    method: 'POST',
    url: 'http://localhost:8080/api/v1/auth/login',
    body: {
      username,
      password,
    },
  }).then((response) => {
    window.localStorage.setItem('token', response.body.token)
    window.localStorage.setItem('user', JSON.stringify(response.body.user))
  })
})

Cypress.Commands.add('register', (userData: any) => {
  cy.request({
    method: 'POST',
    url: 'http://localhost:8080/api/v1/auth/register',
    body: userData,
  })
})

// Global beforeEach hook to clear localStorage
beforeEach(() => {
  window.localStorage.clear()
})
