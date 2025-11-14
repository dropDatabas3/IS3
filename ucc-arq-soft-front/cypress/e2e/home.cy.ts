describe("Home page", () => {
  it("carga y muestra la lista de cursos", () => {
    cy.visit("/");
    cy.contains(/cursos/i).should("be.visible");
    // TODO: ajustar selectores data-test cuando estén disponibles
    // cy.get("[data-test='course-card']").should("have.length.at.least", 1);
  });
});
