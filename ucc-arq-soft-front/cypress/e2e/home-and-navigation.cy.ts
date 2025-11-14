// E2E: smoke de home y navegación básica

describe("Home & Navigation", () => {
  it("carga el home y muestra la navbar", () => {
    cy.log("Visitando home");
    cy.visit("/");

    cy.get('[data-test="navbar-home-link"]').should("be.visible");
    cy.contains(/duomingo/i).should("be.visible");
  });

  it("navega a Courses desde la navbar", () => {
    cy.log("Navegando a /courses");
    cy.visit("/");

    cy.contains(/courses/i).first().click();
    cy.url().should("include", "/courses");
    cy.contains(/courses/i).should("be.visible");
  });

  it("navega a My Courses (autenticado o redirige a login)", () => {
    cy.log("Navegando a /my-courses desde navbar");
    cy.visit("/");

    cy.contains(/my courses/i).first().click();

    // Dependiendo del estado de autenticación, validamos dos comportamientos válidos:
    cy.url().then((url) => {
      if (url.includes("/auth/login")) {
        cy.log("Redirigido a login por guard");
        cy.get('[data-test="login-form"]').should("be.visible");
      } else {
        cy.log("Acceso directo a /my-courses (usuario ya autenticado)");
        // usamos assert de Cypress en lugar de Jest
        assert.include(url, "/my-courses");
      }
    });
  });

  it("navega a Login desde la navbar", () => {
    cy.log("Navegando a /auth/login");
    cy.visit("/");

    cy.contains(/login/i).first().click();
    cy.url().should("include", "/auth/login");
    cy.get('[data-test="login-form"]').should("be.visible");
  });
});
