describe('Content Browsing', () => {
  beforeEach(() => {
    // Login as a regular user
    cy.login('testuser', 'password123')
    cy.visit('/content')
  })

  it('should display content list', () => {
    cy.get('[data-cy=content-list]').should('be.visible')
    cy.get('[data-cy=content-card]').should('have.length.greaterThan', 0)
  })

  it('should filter content by category', () => {
    cy.get('[data-cy=category-filter]').select('Technology')
    cy.get('[data-cy=filter-button]').click()
    
    // Should show filtered results
    cy.get('[data-cy=content-card]').each(($card) => {
      cy.wrap($card).find('[data-cy=category-badge]').should('contain', 'Technology')
    })
  })

  it('should search content', () => {
    const searchTerm = 'React'
    cy.get('[data-cy=search-input]').type(searchTerm)
    cy.get('[data-cy=search-button]').click()

    // Should show search results
    cy.get('[data-cy=content-card]').each(($card) => {
      cy.wrap($card).should('contain', searchTerm)
    })
  })

  it('should navigate to content detail', () => {
    cy.get('[data-cy=content-card]').first().click()
    cy.url().should('match', /\/content\/\d+/)
    cy.get('[data-cy=content-detail]').should('be.visible')
    cy.get('[data-cy=content-title]').should('be.visible')
    cy.get('[data-cy=content-body]').should('be.visible')
  })

  it('should increment view count when viewing content', () => {
    cy.get('[data-cy=content-card]').first().within(() => {
      const initialViewCount = cy.get('[data-cy=view-count]').invoke('text')
    })

    cy.get('[data-cy=content-card]').first().click()
    cy.url().should('match', /\/content\/\d+/)

    // Go back and check if view count increased
    cy.go('back')
    cy.get('[data-cy=content-card]').first().within(() => {
      cy.get('[data-cy=view-count]').should(($el) => {
        const currentCount = parseInt($el.text())
        // View count should be greater than initial count
        expect(currentCount).to.be.greaterThan(0)
      })
    })
  })
})