// E2E: tras inscribirse en un curso, aparece en My Courses

describe("My Courses - enrollment persistence", () => {
  it("muestra un curso en My Courses después de inscribirse", () => {
    // Asumimos que el usuario ya está autenticado para este test
    cy.log("Abriendo un curso desde /courses para inscribirse");
    cy.visit("/courses");

    cy.get("[data-test='course-card']")
      .first()
      .as("selectedCourse");

    cy.get("@selectedCourse")
      .find("[data-test='course-card-title']")
      .invoke("text")
      .then((courseTitle) => {
        cy.log("Curso seleccionado", { courseTitle });

        // Abrir modal y navegar a course-info
        cy.get("@selectedCourse").click();
        cy.contains(/see more/i).click();

        cy.url().should("include", "/course-info");

        cy.log("Haciendo click en Enroll");
        cy.get('[data-test="course-enroll-button"]').click();

        // El botón debería indicar que ya está inscrito
        cy.get('[data-test="course-enroll-button"]').contains(/enrolled/i);

        cy.log("Navegando a My Courses para verificar el curso");
        cy.visit("/my-courses");

        cy.contains("My Courses").should("be.visible");
        cy.contains(courseTitle.trim()).should("be.visible");
      });
  });
});
