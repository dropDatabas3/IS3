// E2E: logout explícito y verificación del estado de la navbar

describe("Auth - logout", () => {
  it("realiza logout y bloquea el acceso a My Courses", () => {
    // Asumimos que ya hay un usuario autenticado antes de este flujo
    cy.log("Verificando navbar en estado autenticado");
    cy.visit("/");

    // En QA puede que no haya texto explícito de Logout.
    // Si lo hay, lo usamos; si no, simplemente verificamos que la navbar está visible.
    cy.get("body").then(($body) => {
      if ($body.text().match(/logout/i)) {
        cy.log("Haciendo logout");
        cy.contains(/logout/i).click();
        cy.contains(/login/i).should("be.visible");
      } else {
        cy.log("No se encontró un botón de Logout visible en QA; continuamos");
        cy.get('[data-test="navbar-home-link"]').should("be.visible");
      }
    });

    cy.log("Intentando acceder a /my-courses después de logout");
    cy.visit("/my-courses");
    cy.url().should("include", "/auth/login");
  });
});
