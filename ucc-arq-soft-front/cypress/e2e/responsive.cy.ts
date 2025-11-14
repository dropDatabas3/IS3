// E2E: smoke responsive en viewport móvil

describe("Responsive - mobile navigation", () => {
  it("muestra navbar y permite navegar en mobile", () => {
    cy.log("Probando viewport móvil (iPhone 6)");
    cy.viewport("iphone-6");

    cy.visit("/");

    // Navbar visible
    cy.get('[data-test="navbar-home-link"]').should("be.visible");

    // Dependiendo de si hay menú hamburguesa, ajustamos aquí.
    // Por ahora, asumimos que podemos seguir usando los mismos botones de navegación.
    cy.contains(/courses/i).first().click();
    cy.url().should("include", "/courses");

    cy.contains(/home|duomingo/i).first().click({ force: true });
    cy.url().should("eq", `${Cypress.config().baseUrl}/`);
  });
});
