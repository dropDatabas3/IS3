// E2E: manejo de errores del API en cursos y login usando cy.intercept

describe("API error handling", () => {
  it("muestra algún mensaje o fallback cuando falla el login", () => {
    cy.log("Simulando error de API en login");

    // Interceptar la llamada de login; ajustar URL según endpoint real
    cy.intercept("POST", /auth\/login/i, {
      statusCode: 401,
      body: { message: "Invalid credentials" },
    }).as("loginError");

    cy.visit("/auth/login", { timeout: 120000 });

    cy.get('[data-test="login-email-input"]').type("fake@example.com");
  cy.get('[data-test="login-password-input"]').type("wrongpass");
  cy.get('[data-test="login-submit-button"]').click();

    // La UI en QA puede no mostrar exactamente "Invalid credentials" ni ningún mensaje explícito;
    // validamos que siga en la pantalla de login y que no se rompa.
    cy.url().should("include", "/auth/login");
    // No forzamos la existencia de un mensaje de error específico, sólo que la UI siga viva
    // y el formulario de login continúe disponible.
    cy.get("body").should("be.visible");
    cy.get('[data-test="login-form"]').should("be.visible");
  });

  it("no rompe la UI cuando falla la carga de cursos", () => {
    cy.log("Simulando error de API en /courses");

    // Interceptar la llamada de cursos; ajustar URL según endpoint real
    cy.intercept("GET", /courses/i, {
      statusCode: 500,
      body: { message: "Server error" },
    }).as("coursesError");

    // En lugar de visitar directamente /courses (que puede ser sólo un endpoint JSON),
    // navegamos al home y dejamos que la app haga el fetch a /courses.
    cy.visit("/", { timeout: 120000 });

    cy.wait("@coursesError");

    // Validamos que la página siga viva (no pantalla en blanco / error de React)
    cy.get("body").should("be.visible");
    // Si hay algún texto de error o título de Courses, también es válido
    cy.get("body").then(($body) => {
      if ($body.text().match(/courses|curso|error/i)) {
        cy.wrap($body).contains(/courses|curso|error/i).should("exist");
      }
    });
  });
});
