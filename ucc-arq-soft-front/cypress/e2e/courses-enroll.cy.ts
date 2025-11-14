// E2E: inscribirse a un curso desde el listado y verificar feedback

describe("Courses - enroll", () => {
  it("abre un curso, hace 'See more' y muestra la página de detalle", () => {
    cy.log("Abriendo modal de curso desde /courses");
    cy.visit("/courses");

    cy.get("[data-test='course-card']").first().click();

    // CourseModal debería abrirse con info del curso
    cy.contains(/price:/i).should("be.visible");
    cy.contains(/duration:/i).should("be.visible");

    cy.log("Navegando a /course-info con 'See more'");
    cy.contains(/see more/i).click();

    cy.url().should("include", "/course-info");
  });

  // Nota: la lógica de 'enrolarse' está en la página de course-info a través de CoursesContext.
  // Si el flujo real incluye un botón 'Enroll' en /course-info, este test debería extenderse para:
  // - Hacer click en 'Enroll'
  // - Verificar un toast o cambio en My Courses. Por ahora dejamos el esqueleto preparado.
});
