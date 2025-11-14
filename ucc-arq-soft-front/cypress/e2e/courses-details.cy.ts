// E2E: listado de cursos y detalle

describe("Courses - listado y detalle", () => {
  it("muestra al menos un curso y permite ir al detalle", () => {
    cy.log("Visitando /courses");
    cy.visit("/courses");

    // Ajustar selectores según estructura real de las cards
    cy.contains(/courses/i).should("be.visible");

    // Suponemos que hay al menos una card de curso con botón
    cy.get("[data-test='course-card']")
      .first()
      .click({ force: true });

    cy.url().should("include", "/course-info");

    cy.contains(/course/i).should("be.visible");
  });
});
