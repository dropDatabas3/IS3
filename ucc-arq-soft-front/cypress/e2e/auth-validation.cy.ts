// E2E: validaciones de formularios de Login y Register (mensajes de yup)

describe("Auth - form validations", () => {
  it("muestra errores de validación en login", () => {
    cy.log("Probando validaciones en /auth/login");
    cy.visit("/auth/login");

    cy.get('[data-test="login-form"]').should("be.visible");

    // Submit vacío
    cy.get('[data-test="login-submit-button"]').click();

    cy.contains("Email is required").should("be.visible");
    cy.contains("Password is required").should("be.visible");

    // Email inválido
    cy.get('[data-test="login-email-input"]').type("not-an-email");
    cy.get('[data-test="login-password-input"]').type("123");
    cy.get('[data-test="login-submit-button"]').click();

    cy.contains("Email is not valid").should("be.visible");
  });

  it("muestra errores de validación en register", () => {
    cy.log("Probando validaciones en /auth/register");
    cy.visit("/auth/register");

    cy.get('[data-test="register-form"]').should("be.visible");

    // Submit vacío
    cy.get('[data-test="register-submit-button"]').click();

    cy.contains("Name is required").should("be.visible");
    cy.contains("Email is required").should("be.visible");
    cy.contains("Password is required").should("be.visible");

    // Email inválido + password corta
    cy.get('[data-test="register-username-input"]').type("Test User");
    cy.get('[data-test="register-email-input"]').type("invalid-email");
    cy.get('[data-test="register-password-input"]').type("123");
    cy.get('[data-test="register-submit-button"]').click();

    cy.contains("Email is not valid").should("be.visible");
  });
});
