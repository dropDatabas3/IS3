// E2E: protección de ruta /my-courses para usuarios no autenticados

describe("Auth Guard - My Courses", () => {
  it("redirecciona a login cuando no está autenticado", () => {
    cy.log("Intentando acceder a /my-courses sin autenticación");
    cy.visit("/my-courses");

    cy.url().should("include", "/auth/login");
    cy.get('[data-test="login-form"]').should("be.visible");
  });
});
