describe('Admin CRUD Operations', () => {
  beforeEach(() => {
    // Login as admin
    cy.login('admin', 'admin123')
    cy.visit('/admin')
  })

  it('should display admin dashboard', () => {
    cy.url().should('include', '/admin')
    cy.get('[data-cy=admin-dashboard]').should('be.visible')
    cy.get('[data-cy=admin-nav]').should('be.visible')
  })

  describe('User Management', () => {
    beforeEach(() => {
      cy.get('[data-cy=admin-users-link]').click()
      cy.url().should('include', '/admin/users')
    })

    it('should display users list', () => {
      cy.get('[data-cy=users-table]').should('be.visible')
      cy.get('[data-cy=user-row]').should('have.length.greaterThan', 0)
    })

    it('should create new user', () => {
      const newUser = {
        username: `newuser${Date.now()}`,
        email: `new${Date.now()}@example.com`,
        firstName: 'New',
        lastName: 'User',
        password: 'password123',
      }

      cy.get('[data-cy=create-user-button]').click()
      cy.get('[data-cy=create-user-modal]').should('be.visible')

      cy.get('[data-cy=username-input]').type(newUser.username)
      cy.get('[data-cy=email-input]').type(newUser.email)
      cy.get('[data-cy=first-name-input]').type(newUser.firstName)
      cy.get('[data-cy=last-name-input]').type(newUser.lastName)
      cy.get('[data-cy=password-input]').type(newUser.password)

      cy.get('[data-cy=save-user-button]').click()

      // Should show success message and new user in list
      cy.get('[data-cy=success-message]').should('be.visible')
      cy.get('[data-cy=users-table]').should('contain', newUser.username)
    })

    it('should update existing user', () => {
      cy.get('[data-cy=user-row]').first().within(() => {
        cy.get('[data-cy=edit-user-button]').click()
      })

      cy.get('[data-cy=edit-user-modal]').should('be.visible')
      
      const newFirstName = 'Updated'
      cy.get('[data-cy=first-name-input]').clear().type(newFirstName)
      cy.get('[data-cy=save-user-button]').click()

      // Should show success message and updated name
      cy.get('[data-cy=success-message]').should('be.visible')
      cy.get('[data-cy=users-table]').should('contain', newFirstName)
    })

    it('should delete user', () => {
      cy.get('[data-cy=user-row]').first().within(() => {
        cy.get('[data-cy=delete-user-button]').click()
      })

      cy.get('[data-cy=confirm-delete-modal]').should('be.visible')
      cy.get('[data-cy=confirm-delete-button]').click()

      // Should show success message and user should be removed
      cy.get('[data-cy=success-message]').should('be.visible')
    })
  })

  describe('Content Management', () => {
    beforeEach(() => {
      cy.get('[data-cy=admin-content-link]').click()
      cy.url().should('include', '/admin/content')
    })

    it('should display content list', () => {
      cy.get('[data-cy=content-table]').should('be.visible')
      cy.get('[data-cy=content-row]').should('have.length.greaterThan', 0)
    })

    it('should create new content', () => {
      const newContent = {
        title: `New Article ${Date.now()}`,
        description: 'A new test article',
        body: 'This is the body of the new article',
        categoryId: 1,
      }

      cy.get('[data-cy=create-content-button]').click()
      cy.get('[data-cy=create-content-modal]').should('be.visible')

      cy.get('[data-cy=title-input]').type(newContent.title)
      cy.get('[data-cy=description-input]').type(newContent.description)
      cy.get('[data-cy=body-input]').type(newContent.body)
      cy.get('[data-cy=category-select]').select(newContent.categoryId.toString())

      cy.get('[data-cy=save-content-button]').click()

      // Should show success message and new content in list
      cy.get('[data-cy=success-message]').should('be.visible')
      cy.get('[data-cy=content-table]').should('contain', newContent.title)
    })

    it('should update existing content', () => {
      cy.get('[data-cy=content-row]').first().within(() => {
        cy.get('[data-cy=edit-content-button]').click()
      })

      cy.get('[data-cy=edit-content-modal]').should('be.visible')
      
      const newTitle = 'Updated Article Title'
      cy.get('[data-cy=title-input]').clear().type(newTitle)
      cy.get('[data-cy=save-content-button]').click()

      // Should show success message and updated title
      cy.get('[data-cy=success-message]').should('be.visible')
      cy.get('[data-cy=content-table]').should('contain', newTitle)
    })

    it('should delete content', () => {
      cy.get('[data-cy=content-row]').first().within(() => {
        cy.get('[data-cy=delete-content-button]').click()
      })

      cy.get('[data-cy=confirm-delete-modal]').should('be.visible')
      cy.get('[data-cy=confirm-delete-button]').click()

      // Should show success message and content should be removed
      cy.get('[data-cy=success-message]').should('be.visible')
    })
  })

  describe('Category Management', () => {
    beforeEach(() => {
      cy.get('[data-cy=admin-categories-link]').click()
      cy.url().should('include', '/admin/categories')
    })

    it('should display categories list', () => {
      cy.get('[data-cy=categories-table]').should('be.visible')
      cy.get('[data-cy=category-row]').should('have.length.greaterThan', 0)
    })

    it('should create new category', () => {
      const newCategory = {
        name: `New Category ${Date.now()}`,
        description: 'A new test category',
        color: '#FF5733',
      }

      cy.get('[data-cy=create-category-button]').click()
      cy.get('[data-cy=create-category-modal]').should('be.visible')

      cy.get('[data-cy=name-input]').type(newCategory.name)
      cy.get('[data-cy=description-input]').type(newCategory.description)
      cy.get('[data-cy=color-input]').type(newCategory.color)

      cy.get('[data-cy=save-category-button]').click()

      // Should show success message and new category in list
      cy.get('[data-cy=success-message]').should('be.visible')
      cy.get('[data-cy=categories-table]').should('contain', newCategory.name)
    })

    it('should update existing category', () => {
      cy.get('[data-cy=category-row]').first().within(() => {
        cy.get('[data-cy=edit-category-button]').click()
      })

      cy.get('[data-cy=edit-category-modal]').should('be.visible')
      
      const newName = 'Updated Category Name'
      cy.get('[data-cy=name-input]').clear().type(newName)
      cy.get('[data-cy=save-category-button]').click()

      // Should show success message and updated name
      cy.get('[data-cy=success-message]').should('be.visible')
      cy.get('[data-cy=categories-table]').should('contain', newName)
    })

    it('should delete category', () => {
      cy.get('[data-cy=category-row]').first().within(() => {
        cy.get('[data-cy=delete-category-button]').click()
      })

      cy.get('[data-cy=confirm-delete-modal]').should('be.visible')
      cy.get('[data-cy=confirm-delete-button]').click()

      // Should show success message and category should be removed
      cy.get('[data-cy=success-message']).should('be.visible')
    })
  })
})