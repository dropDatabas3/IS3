// E2E: registro -> logout -> login usando usuario único por ejecución

const uniqueUser = () => {
  const timestamp = Date.now();
  // Ajustar dominio si el backend valida dominios específicos
  const email = `testuser+${timestamp}@example.com`;
  const password = `Pass-${timestamp}`;
  const username = `User-${timestamp}`;
  return { email, password, username };
};

describe("Auth - register, logout and login", () => {
  it("registra un usuario nuevo, hace logout y luego login con ese usuario", () => {
    const { email, password, username } = uniqueUser();

    cy.log("Registro de usuario nuevo", { email, username });

    // 1) Register
    cy.visit("/auth/register");
    cy.get('[data-test="register-form"]').should("be.visible");

    cy.get('[data-test="register-username-input"]').type(username);
    cy.get('[data-test="register-email-input"]').type(email);
    cy.get('[data-test="register-password-input"]').type(password);

    cy.get('[data-test="register-submit-button"]').click();

    // El registro hace auto-login y redirige al home
    cy.url().should("eq", `${Cypress.config().baseUrl}/`);

    // 2) Logout (asumiendo que hay un botón de logout visible cuando está logueado)
    // Ajustar selector si el componente usa otro texto/estructura
    cy.log("Realizando logout");
    cy.contains(/logout/i).click();

    // Después de logout deberíamos ver el botón/link de Login en el navbar
    cy.contains(/login/i).should("be.visible");

    // 3) Login con el mismo usuario
    cy.log("Login con usuario registrado", { email });
    cy.visit("/auth/login");

    cy.get('[data-test="login-form"]').should("be.visible");
    cy.get('[data-test="login-email-input"]').type(email);
    cy.get('[data-test="login-password-input"]').type(password);
    cy.get('[data-test="login-submit-button"]').click();

    // Validar que vuelve al home
    cy.url().should("eq", `${Cypress.config().baseUrl}/`);

    // Y que el navbar refleja estado logueado (por ejemplo que ya no aparece Login)
    cy.contains(/login/i).should("not.exist");
  });
});
