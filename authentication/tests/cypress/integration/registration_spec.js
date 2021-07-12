describe('Registration Tests', () => {
    beforeEach(() => {
        cy.task('db:teardown')
        cy.task('db:seed')
    })

    it('fails with pre-existing username', () => {
      cy.visit('/')
      cy.intercept('POST', '/register', {
        statusCode: 409,
        body: 'username taken'
      }).as('registration-attempt')

      cy.get("a#transferRegister").click()

      cy.get("form#registrationForm input#email").type("demo@testing.com")
      cy.get("form#registrationForm input#username").type("cool_username")
      cy.get("form#registrationForm input#password").type("random_password_here_123!")

      cy.get("form#registrationForm").submit()

      cy.wait('@registration-attempt')

      cy.get("form#registrationForm label[for='username']").should('have.text', 'Username taken')
    })

    // I plan to only stub the non-crucial routes. For crucial routes, like registering, I want to make sure everything works flawlessly. 
    it('succeeds without stub', () => {
        cy.visit('/')
        cy.get("a#transferRegister").click()

        cy.get("form#registrationForm input#email").type("demo@testing.com")
        cy.get("form#registrationForm input#username").type("my_awesome_username")
        cy.get("form#registrationForm input#password").type("random_password_here_123!")

        cy.get("form#registrationForm").submit()

        cy.get("form#registrationForm input#registerBtn").should('have.value', "Registration complete.")
        cy.get("form#registrationForm input#registerBtn").should('have.class', "success")

        cy.wait(1500)

        cy.get(".page.register").should('not.have.class', 'showing')
        cy.get('.page.login').should('have.class', 'showing')
    })
})