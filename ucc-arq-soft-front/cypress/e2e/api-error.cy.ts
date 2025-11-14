// E2E: manejo de errores del API en cursos y login usando cy.intercept

describe("API error handling", () => {
  it("muestra algún mensaje o fallback cuando falla el login", () => {
    cy.log("Simulando error de API en login");

    // Interceptar la llamada de login; ajustar URL según endpoint real
    cy.intercept("POST", /auth\/login/i, {
      statusCode: 401,
      body: { message: "Invalid credentials" },
    }).as("loginError");

    cy.visit("/auth/login");

    cy.get('[data-test="login-email-input"]').type("fake@example.com");
    cy.get('[data-test="login-password-input"]').type("wrongpass");
    cy.get('[data-test="login-submit-button"]').click();

    cy.wait("@loginError");

    // Ajustar según cómo la UI muestre el error (toast o mensaje en pantalla)
    cy.contains(/invalid credentials/i).should("exist");
  });

  it("no rompe la UI cuando falla la carga de cursos", () => {
    cy.log("Simulando error de API en /courses");

    // Interceptar la llamada de cursos; ajustar URL según endpoint real
    cy.intercept("GET", /courses/i, {
      statusCode: 500,
      body: { message: "Server error" },
    }).as("coursesError");

    cy.visit("/courses");

    cy.wait("@coursesError");

    // Aquí validamos al menos que la página no se rompe
    // Si tenés un mensaje de error específico, podemos ajustarlo aquí
    cy.contains(/courses/i).should("be.visible");
  });
});
