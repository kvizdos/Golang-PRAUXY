/* 
QR = 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
*/
import { authenticator } from 'otplib';

describe('Login Tests', () => {
    beforeEach(() => {
      cy.task('db:teardown')
      cy.task('db:seed')
    })

    it('Login fails when invalid username and password is used', () => {
      cy.visit('/')
      cy.intercept('POST', '/login', {
        statusCode: 401,
        body: 'invalid username or password'
      }).as('login-attempt')

      cy.get("form#loginForm input#username").type("bad_username")
      cy.get("form#loginForm input#password").type("random_password_here_123!")

      cy.get("form#loginForm").submit()

      cy.wait('@login-attempt')

      cy.get("form#loginForm p#badlogin").should('have.class', 'showing')
    })

    it('succeeds with a valid username & password (live test)', () => {
      cy.visit('/')

      cy.get("form#loginForm input#username").type("test_user")
      cy.get("form#loginForm input#password").type("abc123")

      cy.get("form#loginForm").submit()

      cy.get("p#badlogin").should('not.have.class', 'showing')

      cy.location('pathname').should('eq',  '/account')
      cy.getCookie("sid").should('not.have.property', 'value', 'undefined')

    })

    it('should show a tfa verification screen when enabled on the user', () => {
      cy.visit('/')
      cy.intercept('POST', '/login', {
        statusCode: 200,
        body: '{"mfaSID": "mfaSIDHere"}'
      }).as('login-attempt')

      // Login
      cy.get("form#loginForm input#username").type("test_user_stubbed_tfa")
      cy.get("form#loginForm input#password").type("stubbed_password")
      cy.get("form#loginForm").submit()

      cy.get("#loginForm").should('have.class', 'hide')
      cy.get("#loginMfaForm").should('have.class', 'show')
    })



    it('fail with an invalid totp code', () => {
      cy.visit('/')

      cy.intercept('POST', '/login', {
        statusCode: 200,
        body: '{"mfaSID": "mfaSIDHere"}'
      }).as('login-attempt')

      cy.intercept('POST', '/user/mfa/verify', {
        statusCode: 403,
        body: 'invalid totp code'
      }).as('totp-attempt')

      // Login
      cy.get("form#loginForm input#username").type("test_user_stubbed_tfa")
      cy.get("form#loginForm input#password").type("stubbed_password")
      cy.get("form#loginForm").submit()

      cy.get("form#loginMfaForm .input input#token").type("123456")

      cy.wait('@totp-attempt')

      cy.get('form#loginMfaForm p#verificationStatus').should('have.text', 'Invalid TOTP code provided.')
    })

    // add an edge test for providing a totp code when the account doesn't have it enabled. while this shouldn't show to really any user, its good to handle.
    it('fail if the account doesn\'t have totp enabled', () => {
      cy.visit('/')

      cy.intercept('POST', '/login', {
        statusCode: 200,
        body: '{"mfaSID": "mfaSIDHere"}'
      }).as('login-attempt')

      cy.intercept('POST', '/user/mfa/verify', {
        statusCode: 406,
        body: 'totp not enabled'
      }).as('totp-attempt')

      // Login
      cy.get("form#loginForm input#username").type("test_user_stubbed_tfa")
      cy.get("form#loginForm input#password").type("stubbed_password")
      cy.get("form#loginForm").submit()

      cy.get("form#loginMfaForm .input input#token").type("123456")
      cy.wait('@totp-attempt')

      cy.get('form#loginMfaForm p#verificationStatus').should('have.text', 'TOTP is not enabled on this account.')

    })

    it('return a token when a valid totp code is provided (live test)', () => {
        cy.visit('/')

        // Login
        cy.get("form#loginForm input#username").type("test_user_with_totp")
        cy.get("form#loginForm input#password").type("abc123")
        cy.get("form#loginForm").submit()

        let token = authenticator.generate(Cypress.env('seed_user_totp_secret'))
        cy.get("form#loginMfaForm .input input#token").type(token)

        cy.location('pathname').should('eq',  '/account')
        cy.getCookie("sid").should('not.have.property', 'value', 'undefined')
    })
  })