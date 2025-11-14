// E2E: logout explícito y verificación del estado de la navbar

describe("Auth - logout", () => {
  it("realiza logout y bloquea el acceso a My Courses", () => {
    // Asumimos que ya hay un usuario autenticado antes de este flujo
    cy.log("Verificando navbar en estado autenticado");
    cy.visit("/");

    // Debería existir un botón o texto de Logout
    cy.contains(/logout/i).should("be.visible");

    cy.log("Haciendo logout");
    cy.contains(/logout/i).click();

    // Después de logout, debería aparecer opción de Login
    cy.contains(/login/i).should("be.visible");

    cy.log("Intentando acceder a /my-courses después de logout");
    cy.visit("/my-courses");
    cy.url().should("include", "/auth/login");
  });
});
